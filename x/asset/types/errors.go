package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/asset module sentinel errors
var (
	ErrUnauthorize          = errorsmod.Register(ModuleName, 1501, "unauthorized address")
	ErrTokenSet             = errorsmod.Register(ModuleName, 1502, "token is unable to be set")
	ErrTokenManagementSet   = errorsmod.Register(ModuleName, 1503, "token management is unable to be set")
	ErrTokenDistributionSet = errorsmod.Register(ModuleName, 1504, "token distribution is unable to be set")
	ErrTokenGet             = errorsmod.Register(ModuleName, 1505, "token is unable to be get")
	ErrTokenManagementGet   = errorsmod.Register(ModuleName, 1506, "token management is unable to be get")
	ErrTokenDistributionGet = errorsmod.Register(ModuleName, 1507, "token distribution is unable to be get")
	ErrAccAddress           = errorsmod.Register(ModuleName, 1508, "unable to convert string to acc address")
)
