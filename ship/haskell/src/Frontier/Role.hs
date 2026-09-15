-- | Frontier roles. Pure. No IO.
-- English: Observer looks, Analyst advises, Operator prepares, Executor pushes.
module Frontier.Role
  ( Role(..)
  , can
  , elevate
  , roleName
  ) where

data Role
  = Observer
  | Analyst
  | Operator
  | Executor
  deriving (Eq, Ord, Show, Enum, Bounded)

-- | May this role use a tool that requires at least 'minRole'?
can :: Role -> Role -> Bool
can have minRole = have >= minRole

-- | Step up exactly one rung. Nothing above Executor.
elevate :: Role -> Either String Role
elevate Executor = Left "already at executor"
elevate r        = Right (succ r)

roleName :: Role -> String
roleName Observer = "observer"
roleName Analyst  = "analyst"
roleName Operator = "operator"
roleName Executor = "executor"
