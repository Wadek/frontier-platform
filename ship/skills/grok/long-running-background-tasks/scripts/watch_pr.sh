#!/usr/bin/env bash
# watch_pr.sh: deterministic PR monitor for the long-running-background-tasks skill.
#
# Usage: watch_pr.sh <pr-url-or-number> [poll-seconds]
# Needs: authenticated `gh`. Default poll is 30s; lower values hit GitHub rate
# limits when several watchers share a token.
#
# Stdout carries only wake-up tokens, one line each:
#   DONE: merged                       terminal, exit 0
#   CANCELLED: closed without merge    terminal, exit 0
#   FAILED: <reason>                   terminal, exit 1 (PR not found, auth lost)
#   ACTION_REQUIRED: <kind> ...        wake-up; the watcher keeps running
# Kinds: ci_failed, conflict, comment, review, review_comment, mergeable.
# Each kind fires at most once per head SHA (comment counts fire on increase).
# Commit comments are not watched. Diagnostics go to a log file keyed on the PR.
set -uo pipefail

PR=${1:?usage: watch_pr.sh <pr-url-or-number> [poll-seconds]}
POLL=${2:-30}
DIR=${WATCH_DIR:-$HOME/.grok/long-running-background-tasks}
mkdir -p "$DIR"
KEY=$(printf '%s' "$PR" | tr -c 'A-Za-z0-9' '_')
LOG=$DIR/watch_pr_$KEY.log
log() { printf '%s %s\n' "$(date -u +%FT%TZ)" "$*" >>"$LOG"; }

declare -A fired=()
once() { # once <kind> <sha> <message>
  local k="$1@$2"
  [[ -n ${fired[$k]:-} ]] && return
  fired[$k]=1
  echo "ACTION_REQUIRED: $1 $3"
}

prev_comments=-1 prev_reviews=-1 prev_inline=-1 misses=0
log "start pr=$PR poll=$POLL"
while :; do
  if ! view=$(gh pr view "$PR" --json state,headRefOid,mergeable,mergeStateStatus,reviewDecision,url,number,comments,reviews,statusCheckRollup \
      --jq '[.state, .headRefOid, .mergeable, .mergeStateStatus, (.reviewDecision // ""), .url, (.number|tostring),
             (.comments|length|tostring), (.reviews|length|tostring),
             ([.statusCheckRollup[]? | select((.conclusion // .state // "") as $c |
                 ["FAILURE","TIMED_OUT","CANCELLED","STARTUP_FAILURE","ERROR","ACTION_REQUIRED"] | index($c))] | length | tostring),
             ([.statusCheckRollup[]? | select((.status // "COMPLETED") != "COMPLETED" or (.state // "") == "PENDING")] | length | tostring)
            ] | join("\u001f")' 2>>"$LOG"); then
    misses=$((misses + 1)); log "gh pr view failed ($misses)"
    if (( misses >= 5 )); then echo "FAILED: cannot read $PR after 5 attempts; see $LOG"; exit 1; fi
    sleep "$POLL"; continue
  fi
  misses=0
  IFS=$'\x1f' read -r state sha mergeable mstate review url num ncomments nreviews nfail npending <<<"$view"
  log "state=$state sha=${sha:0:7} mergeable=$mergeable mstate=$mstate review=$review fail=$nfail pending=$npending comments=$ncomments reviews=$nreviews"

  case $state in
    MERGED) echo "DONE: merged $url"; exit 0 ;;
    CLOSED) echo "CANCELLED: closed without merge $url"; exit 0 ;;
  esac

  repo=$(sed -E 's#https://github.com/([^/]+/[^/]+)/pull/.*#\1#' <<<"$url")
  ninline=$(gh api "repos/$repo/pulls/$num/comments" --paginate --jq 'length' 2>>"$LOG" | awk '{s+=$1} END {print s+0}')

  (( nfail > 0 )) && once ci_failed "$sha" "$nfail failing check(s) on ${sha:0:7} $url"
  [[ $mergeable == CONFLICTING ]] && once conflict "$sha" "merge conflict on ${sha:0:7} $url"
  (( prev_comments >= 0 && ncomments > prev_comments )) && echo "ACTION_REQUIRED: comment $((ncomments - prev_comments)) new conversation comment(s) $url"
  (( prev_reviews >= 0 && nreviews > prev_reviews )) && echo "ACTION_REQUIRED: review $((nreviews - prev_reviews)) new review(s) $url"
  (( prev_inline >= 0 && ninline > prev_inline )) && echo "ACTION_REQUIRED: review_comment $((ninline - prev_inline)) new inline comment(s) $url"
  prev_comments=$ncomments prev_reviews=$nreviews prev_inline=$ninline

  if [[ $mergeable == MERGEABLE && $nfail == 0 && $npending == 0 && $review != CHANGES_REQUESTED && $review != REVIEW_REQUIRED \
        && ( $mstate == CLEAN || $mstate == HAS_HOOKS || $mstate == UNSTABLE ) ]]; then
    once mergeable "$sha" "ready to merge ${sha:0:7} $url"
  fi
  sleep "$POLL"
done
