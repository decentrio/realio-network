package types

import "cosmossdk.io/collections"

var (
	CommitteesPrefix                = collections.NewPrefix(0)
	CommitteeUpdatesPrefix          = collections.NewPrefix(2)
	EscrowAccountsPrefix            = collections.NewPrefix(3)
	DepositPackagesPrefix           = collections.NewPrefix(5)
	WithdrawPackagesPrefix          = collections.NewPrefix(6)
	WithdrawNoncesPrefix            = collections.NewPrefix(7)
	LastObservedDepositNoncePrefix  = collections.NewPrefix(8)
)

const (
	// ModuleName defines the module name.
	ModuleName = "stable-ramp"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for stable-ramp.
	RouterKey = ModuleName
)
