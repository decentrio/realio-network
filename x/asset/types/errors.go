package types

import (
	fmt "fmt"

	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/asset module sentinel errors
var (
	ErrSample               = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout = errorsmod.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = errorsmod.Register(ModuleName, 1501, "invalid version")
	ErrNotAuthorized        = errorsmod.Register(ModuleName, 1502, "transaction not authorized")
	ErrInvalidCreator       = errorsmod.Register(ModuleName, 1503, "invalid creator")
	ErrSubdenomTooLong      = errorsmod.Register(ModuleName, 1504, fmt.Sprintf("sub asset name too long, max length is %d bytes", MaxSubAssetNameLength))
	ErrCreatorTooLong       = errorsmod.Register(ModuleName, 1505, fmt.Sprintf("creator too long, max length is %d bytes", MaxCreatorLength))
	ErrDenomExists          = errorsmod.Register(ModuleName, 1506, "attempting to create a denom that already exists (has bank metadata)")
)
