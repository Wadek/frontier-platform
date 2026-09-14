-- Tower.SessionBudget — the 95% rule as truth (F5 witness).
module Tower.SessionBudget
  ( Verdict (..)
  , check
  , softLimit
  ) where

softLimit :: Double
softLimit = 95.0

data Verdict = Verdict
  { ok                :: Bool
  , refuseNewRuns     :: Bool
  , requireHandoff    :: Bool
  , requireNewSession :: Bool
  } deriving (Eq, Show)

-- | At 95% session context/cache: soft-refuse new runs, insist on a handoff
-- and a new session. Soft = the artifacts are required; it is never silent.
check :: Double -> Verdict
check pct
  | pct >= softLimit = Verdict False True True True
  | otherwise        = Verdict True False False False
