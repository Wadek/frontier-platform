-- | OWASP Top 10 (2021) as the first vulnerability definition set V.
-- English policy: english/SECURITY_POLICY_OWASP.md
module Frontier.OWASP
  ( Severity(..)
  , Rule(..)
  , owaspTop10
  , ruleIds
  ) where

data Severity = Medium | High | Critical
  deriving (Eq, Ord, Show)

data Rule = Rule
  { ruleId      :: String
  , owaspRef    :: String
  , severity    :: Severity
  , description :: String
  } deriving (Eq, Show)

-- | V for this policy version: the supplied OWASP rule ids.
owaspTop10 :: [Rule]
owaspTop10 =
  [ Rule "owasp.a01.hardcoded_auth_bypass" "A01" High
      "Hardcoded auth bypass patterns"
  , Rule "owasp.a02.secret_material" "A02" Critical
      "Private key or secret material markers"
  , Rule "owasp.a03.injection_sink" "A03" High
      "SQL/shell/eval injection sinks"
  , Rule "owasp.a04.insecure_flag" "A04" Medium
      "Explicit insecure security flags"
  , Rule "owasp.a05.misconfig" "A05" Medium
      "Debug/misconfiguration defaults"
  , Rule "owasp.a06.unpinned_latest" "A06" Medium
      "Unpinned :latest base images"
  , Rule "owasp.a07.hardcoded_credential" "A07" High
      "Hardcoded credentials/secrets"
  , Rule "owasp.a08.curl_bash" "A08" High
      "curl | bash install pattern"
  , Rule "owasp.a09.sensitive_log" "A09" Medium
      "Logging of secrets/authorization headers"
  , Rule "owasp.a10.open_url_fetch" "A10" Medium
      "Open URL fetch without allowlist note"
  ]

ruleIds :: [String]
ruleIds = map ruleId owaspTop10
