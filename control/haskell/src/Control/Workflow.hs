-- Control.Workflow — the workflow state machine as truth (F5 witness).
-- No transition without a sealed journal row (F0). No closed without a debrief (T4).
module Control.Workflow
  ( State (..)
  , nextOk
  ) where

data State
  = Queued | Review | Deferred | Preparing | Running | Transferred
  | Failed | Completed | Debriefed | Closed | Rejected
  deriving (Eq, Show, Enum, Bounded)

-- | Is the transition legal?
nextOk :: State -> State -> Bool
nextOk a b = b `elem` allowed a
  where
    allowed Queued      = [Review]
    allowed Review      = [Preparing, Deferred, Rejected]
    allowed Deferred    = [Review]
    allowed Preparing   = [Running, Rejected]
    allowed Running     = [Completed, Transferred, Failed]
    allowed Transferred = [Running, Failed]
    allowed Failed      = [Review]
    allowed Completed   = [Debriefed, Failed]
    allowed Debriefed   = [Closed, Failed]
    allowed Closed      = []
    allowed Rejected    = []
