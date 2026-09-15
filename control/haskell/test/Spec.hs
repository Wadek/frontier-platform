-- Control laws — runnable witness tests (build where GHC exists: cabal test).
module Main (main) where

import Data.Time.Calendar (fromGregorian)
import Data.Time.Clock (UTCTime (..), secondsToDiffTime)
import Control.Pricing
import Control.Workflow

-- 2026-09-14 is a Monday; 2026-09-12 is a Saturday.
monday :: Double -> UTCTime
monday h = UTCTime (fromGregorian 2026 9 14) (secondsToDiffTime (truncate (h * 3600)))

check :: Bool -> String -> IO ()
check True _ = return ()
check False m = fail ("law broken: " ++ m)

main :: IO ()
main = do
  -- pricing
  check (isPeak (monday 1)) "mon 01:00 is peak"
  check (isPeak (monday 9.5)) "mon 09:30 is peak"
  check (not (isPeak (monday 4.5))) "mon 04:30 is off-peak"
  check (isPeak (monday 10) == False) "mon 10:00 off-peak"
  check (isPeak (monday 5) == False) "mon 05:00 off-peak"
  let sat = UTCTime (fromGregorian 2026 9 12) (secondsToDiffTime 28800)
  check (isPeak sat == False) "saturday off-peak"
  -- provider sizing
  check (route "teacher" (monday 7) == Decision "deepseek-flash" "deepseek-flash" True "pro forbidden at peak")
    "pro deferred at peak"
  check (route "teacher" (monday 10) == Decision "deepseek-v4-pro" "deepseek-v4-pro" True "pro off-peak")
    "pro runs off-peak"
  check (route "atomic" (monday 7) == Decision "local" "local" False "local-first")
    "atomic stays local"
  -- FSM
  check (nextOk Completed Debriefed) "completed -> debriefed"
  check (not (nextOk Completed Closed)) "no closed without debrief (C4)"
  check (not (nextOk Queued Running)) "no review skip"
  check (nextOk Review Deferred && nextOk Deferred Review) "deferred path"
  check (nextOk Running Transferred && nextOk Transferred Running) "transfer path"
  putStrLn "control laws hold"
