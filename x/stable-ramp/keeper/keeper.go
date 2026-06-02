package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/log"

	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	"github.com/ethereum/go-ethereum/common"

	"github.com/realiotech/realio-network/x/stable-ramp/types"
)

// Keeper of the stable-ramp store
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	bankKeeper   types.BankKeeper
	erc20Keeper  erc20keeper.Keeper
	evmKeeper    *evmkeeper.Keeper

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
	erc20Keeper erc20keeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	authority string,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		bankKeeper:   bankKeeper,
		erc20Keeper:  erc20Keeper,
		evmKeeper:    evmKeeper,
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

// BurnEscrowedCoins burns coins held in the module escrow account.
func (k Keeper) BurnEscrowedCoins(ctx context.Context, coins sdk.Coins) error {
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
}

// FulfillDeposit mints cosmos coins to the receiver then converts them to ERC20 tokens.
func (k Keeper) FulfillDeposit(ctx context.Context, receiver sdk.AccAddress, coins sdk.Coins) error {
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return err
	}

	ethReceiver := common.BytesToAddress(receiver)
	for _, coin := range coins {
		if _, err := k.erc20Keeper.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
			Coin:     coin,
			Receiver: ethReceiver.Hex(),
			Sender:   receiver.String(),
		}); err != nil {
			return err
		}
	}

	return nil
}

// RefundEscrowedCoins releases escrowed cosmos coins back to the sender and
// converts them to ERC20 tokens so the user ends up where they started.
func (k Keeper) RefundEscrowedCoins(ctx context.Context, recipient sdk.AccAddress, coins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins); err != nil {
		return err
	}

	ethReceiver := common.BytesToAddress(recipient)
	for _, coin := range coins {
		if _, err := k.erc20Keeper.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
			Coin:     coin,
			Receiver: ethReceiver.Hex(),
			Sender:   recipient.String(),
		}); err != nil {
			return err
		}
	}

	return nil
}
