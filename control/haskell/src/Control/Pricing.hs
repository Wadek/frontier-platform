-- Control.Pricing — the pricing policy as truth (F5 witness).
-- Official source: https://api-docs.deepseek.com/quick_start/pricing
-- Peak Mon-Fri 01:00-04:00 and 06:00-10:00 UTC; everything else off-peak.
module Control.Pricing
  ( Decision (..)
  , isPeak
  , route
  ) where

import Data.Time.Calendar (DayOfWeek (..), dayOfWeek)
import Data.Time.Clock (UTCTime (..))

peakWindows :: [(Double, Double)]
peakWindows = [(1, 4), (6, 10)]

-- | True inside a weekday peak window (UTC).
isPeak :: UTCTime -> Bool
isPeak t = weekday && any inWindow peakWindows
  where
    weekday = dayOfWeek (utctDay t) `notElem` [Saturday, Sunday]
    h = realToFrac (utctDayTime t) / 3600
    inWindow (lo, hi) = h >= lo && h < hi

data Decision = Decision
  { target :: String
  , model  :: String
  , egress :: Bool
  , reason :: String
  } deriving (Eq, Show)

localKinds, flashKinds, proKinds :: [String]
localKinds = ["atomic", "audit", "ops", "ingest", "train"]
flashKinds = ["plan", "label", "census-draft", "tips", "vision", "coach"]
proKinds   = ["census-deep", "verify", "teacher", "plan-pro", "label-pro"]

-- | Deterministic provider sizing. A model never decides which model to use.
route :: String -> UTCTime -> Decision
route kind now
  | kind `elem` localKinds = Decision "local" "local" False "local-first"
  | kind `elem` flashKinds = Decision "deepseek-flash" "deepseek-flash" True (kind ++ " after egress; flash")
  | kind `elem` proKinds
  = if isPeak now
      then Decision "deepseek-flash" "deepseek-flash" True "pro forbidden at peak"
      else Decision "deepseek-v4-pro" "deepseek-v4-pro" True "pro off-peak"
  | otherwise = Decision "local" "local" False "default local"
