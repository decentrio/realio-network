package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/log"

	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/realiotech/realio-network/x/stable-ramp/types"
)

// Keeper of the stable-ramp store
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	bankKeeper   types.BankKeeper

	// the address capable of executing a MsgUpdateParams and MsgUpdateCommittee message. Typically, this
	// should be the x/gov module account.
	authority string

	Schema collections.Schema

	// Committees maps a connection ID to committee data.
	Committees collections.Map[string, types.Committee]

	// PendingPackages maps pending packages by nonce and connection ID.
	PendingPackages collections.Map[collections.Pair[uint64, string], types.TrackedPackage]

	// CommitteeUpdates maps a connection ID to a pending new committee.
	CommitteeUpdates collections.Map[string, types.Committee]

	// EscrowAccounts stores escrow account addresses.
	EscrowAccounts collections.KeySet[sdk.AccAddress]

	// NextPackageNonce tracks the next package nonce.
	NextPackageNonce collections.Sequence
}

// NewKeeper creates a new stable-ramp Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		bankKeeper:   bankKeeper,
		authority:    authority,
		Committees: collections.NewMap(
			sb,
			types.CommitteesPrefix,
			"committees",
			collections.StringKey,
			codec.CollValue[types.Committee](cdc),
		),
		PendingPackages: collections.NewMap(
			sb,
			types.PendingPackagesPrefix,
			"pending_packages",
			collections.PairKeyCodec(collections.Uint64Key, collections.StringKey),
			codec.CollValue[types.TrackedPackage](cdc),
		),
		CommitteeUpdates: collections.NewMap(
			sb,
			types.CommitteeUpdatesPrefix,
			"committee_updates",
			collections.StringKey,
			codec.CollValue[types.Committee](cdc),
		),
		EscrowAccounts: collections.NewKeySet(
			sb,
			types.EscrowAccountsPrefix,
			"escrow_accounts",
			sdk.AccAddressKey,
		),
		NextPackageNonce: collections.NewSequence(
			sb,
			types.NextPackageNoncePrefix,
			"next_package_nonce",
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

// GetAuthority returns the x/stable-ramp module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}
