-- Tower.FlightPlan — the workflow state machine as truth (F5 witness).
-- No transition without a sealed journal row (F0). No filed without a debrief (T4).
module Tower.FlightPlan
  ( State (..)
  , nextOk
  ) where

data State
  = Queued | Review | Holding | Taxiing | Airborne | HandedOff
  | Failed | Landed | Debriefed | Filed | Rejected
  deriving (Eq, Show, Enum, Bounded)

-- | Is the transition legal?
nextOk :: State -> State -> Bool
nextOk a b = b `elem` allowed a
  where
    allowed Queued    = [Review]
    allowed Review    = [Taxiing, Holding, Rejected]
    allowed Holding   = [Review]
    allowed Taxiing   = [Airborne, Rejected]
    allowed Airborne  = [Landed, HandedOff, Failed]
    allowed HandedOff = [Airborne, Failed]
    allowed Failed    = [Review]
    allowed Landed    = [Debriefed, Failed]
    allowed Debriefed = [Filed, Failed]
    allowed Filed     = []
    allowed Rejected  = []
