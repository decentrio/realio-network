package keeper

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/realiotech/realio-network/x/asset/types"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		storeKey         storetypes.StoreKey
		memKey           storetypes.StoreKey
		paramstore       paramtypes.Subspace
		bankKeeper       types.BankKeeper
		ak               types.AccountKeeper
		PrivilegeManager map[string]types.PrivilegeI
	}
)

// NewKeeper returns a new Keeper object with a given codec, dedicated
// store key, a BankKeeper implementation, an AccountKeeper implementation, and a parameter Subspace used to
// store and fetch module parameters. It also has an allowAddrs map[string]bool to skip restrictions for module addresses.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	ak types.AccountKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		paramstore:       ps,
		bankKeeper:       bankKeeper,
		ak:               ak,
		PrivilegeManager: map[string]types.PrivilegeI{},
	}
}

func (k Keeper) AddPrivilege(priv types.PrivilegeI) error {
	if _, ok := k.PrivilegeManager[priv.Name()]; ok {
		return fmt.Errorf("privilege %s already exists", priv.Name())
	}

	k.PrivilegeManager[priv.Name()] = priv
	// regiester the privilege's interfaces
	priv.RegisterInterfaces()
	return nil
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
