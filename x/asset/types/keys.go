package types

import (
	"cosmossdk.io/collections"
)

var (
	ParamsKey        = collections.NewPrefix(0)
	TokenKeyPrefix   = collections.NewPrefix(1)
	IssuerPrefixKey = "issuer"
)

const (
	// ModuleName defines the module name
	ModuleName = "asset"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// Version defines the current version the IBC module supports
	Version = "asset-1"

	// PortID is the default port id that module binds to
	PortID = "asset"
)

// PortKey defines the key to store the port ID in store
var PortKey = KeyPrefix("asset-port-")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
