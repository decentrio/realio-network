package types

import "cosmossdk.io/collections"

var (
	CommitteesPrefix       = collections.NewPrefix(0)
	PendingPackagesPrefix  = collections.NewPrefix(1)
	CommitteeUpdatesPrefix = collections.NewPrefix(2)
	EscrowAccountsPrefix   = collections.NewPrefix(3)
	NextPackageNoncePrefix = collections.NewPrefix(4)
)

const (
	// ModuleName defines the module name.
	ModuleName = "stable-ramp"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for stable-ramp.
	RouterKey = ModuleName
)
