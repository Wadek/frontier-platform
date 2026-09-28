# YouTube source (listen / dual)

YouTube is an input adapter for the learning skill. It is not a project skill.

## Fetch

1. If `scripts/fetch_transcript.py` exists in the runtime, run it on the URL.
2. Else search the public web for a published transcript (`"<title>" transcript` or video id).
3. Else ask the human to paste YouTube → More → Show transcript.

Need metadata only: `https://noembed.com/embed?url=https://www.youtube.com/watch?v=ID`.

## After fetch

- Digest requested → TL;DR + takeaways + timestamped claims.
- Practice requested → `listen`: 8–12 speech segments, one at a time, text hidden.
- Shadowing → same segments; human speaks; agent only checks wording from captions.

## Do not

- Download the media file unless the human explicitly orders it and review policy allows.
- Create `youtube-listening` under a personal skill folder.
- Invent captions when fetch fails.
