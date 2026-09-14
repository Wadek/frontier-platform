-- | Push gate. Pure. No IO. No models.
-- English: before push, branch must be a feature branch, HEAD set, tree clean.
module Frontier.Gate
  ( GateInput(..)
  , GateResult(..)
  , evaluatePushGate
  ) where

data GateInput = GateInput
  { branch     :: String
  , headSha    :: String
  , dirty      :: Bool
  , allowDirty :: Bool
  } deriving (Eq, Show)

data GateResult = GateResult
  { ok      :: Bool
  , reasons :: [String]
  } deriving (Eq, Show)

evaluatePushGate :: GateInput -> GateResult
evaluatePushGate g =
  let rs =
        [ msg
        | (cond, msg) <- checks
        , cond
        ]
  in GateResult { ok = null rs, reasons = rs }
  where
    checks =
      [ (null (branch g), "detached HEAD or empty branch")
      , (null (headSha g), "missing HEAD")
      , (dirty g && not (allowDirty g), "working tree dirty; commit or clean before push")
      , (isMain (branch g), "refusing direct push to main/master (use a feature branch)")
      ]

isMain :: String -> Bool
isMain b =
  let x = map toLower' b
  in x == "main" || x == "master"
  where
    toLower' c
      | c >= 'A' && c <= 'Z' = toEnum (fromEnum c + 32)
      | otherwise            = c
