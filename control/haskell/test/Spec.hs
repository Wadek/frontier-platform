-- Tower laws — runnable witness tests (build where GHC exists: cabal test).
module Main (main) where

import Data.Time.Calendar (fromGregorian)
import Data.Time.Clock (UTCTime (..), secondsToDiffTime)
import Tower.FlightPlan
import Tower.SessionBudget
import Tower.Weather

-- 2026-09-14 is a Monday; 2026-09-12 is a Saturday.
monday :: Double -> UTCTime
monday h = UTCTime (fromGregorian 2026 9 14) (secondsToDiffTime (truncate (h * 3600)))

check :: Bool -> String -> IO ()
check True _ = return ()
check False m = fail ("law broken: " ++ m)

main :: IO ()
main = do
  -- weather
  check (isPeak (monday 1)) "mon 01:00 is peak"
  check (isPeak (monday 9.5)) "mon 09:30 is peak"
  check (not (isPeak (monday 4.5))) "mon 04:30 is off-peak"
  check (isPeak (monday 10) == False) "mon 10:00 off-peak"
  check (isPeak (monday 5) == False) "mon 05:00 off-peak"
  let sat = UTCTime (fromGregorian 2026 9 12) (secondsToDiffTime 28800)
  check (isPeak sat == False) "saturday off-peak"
  -- runway sizing
  check (route "teacher" (monday 7) == Decision "deepseek-flash" "deepseek-flash" True "pro forbidden at peak")
    "pro grounded at peak"
  check (route "teacher" (monday 10) == Decision "deepseek-v4-pro" "deepseek-v4-pro" True "pro off-peak")
    "pro flies off-peak"
  check (route "atomic" (monday 7) == Decision "local" "qwen2.5-coder:64k" False "local-first")
    "atomic stays local"
  -- FSM
  check (nextOk Landed Debriefed) "landed -> debriefed"
  check (not (nextOk Landed Filed)) "no filed without debrief (T4)"
  check (not (nextOk Queued Airborne)) "no clearance skip"
  check (nextOk Review Holding && nextOk Holding Review) "holding path"
  check (nextOk Airborne HandedOff && nextOk HandedOff Airborne) "handoff path"
  -- session budget (95% rule)
  check (check 95.0 == Verdict False True True True) "95% soft-refuses with handoff"
  check (check 94.99 == Verdict True False False False) "below 95% stays open"
  putStrLn "tower laws hold"
