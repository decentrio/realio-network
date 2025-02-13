package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/realiotech/realio-network/x/asset/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	bk           types.BankKeeper
	Ak           types.AccountKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	Schema             collections.Schema
	Params             collections.Item[types.Params]
	Token              collections.Map[string, types.Token]
	TokenManagement    collections.Map[string, types.TokenManagement]
	WhitelistAddresses collections.Map[[]byte, bool]
	FreezeAddresses    collections.Map[[]byte, bool]
}

// NewKeeper returns a new Keeper object with a given codec, dedicated
// store key, a BankKeeper implementation, an AccountKeeper implementation, and a parameter Subspace used to
// store and fetch module parameters. It also has an allowAddrs map[string]bool to skip restrictions for module addresses.
func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:                cdc,
		storeService:       storeService,
		authority:          authority,
		bk:                 bk,
		Ak:                 ak,
		Params:             collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Token:              collections.NewMap(sb, types.TokenKey, "token", collections.StringKey, codec.CollValue[types.Token](cdc)),
		TokenManagement:    collections.NewMap(sb, types.TokenManagementKey, "token_management", collections.StringKey, codec.CollValue[types.TokenManagement](cdc)),
		WhitelistAddresses: collections.NewMap(sb, types.WhitelistAddressesKey, "whitelist_addresses", collections.BytesKey, collections.BoolValue),
		FreezeAddresses:    collections.NewMap(sb, types.FreezeAddressesKey, "freeze_addresses", collections.BytesKey, collections.BoolValue),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return &k
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetWhitelistAddress(ctx context.Context, address []byte) bool {
	found, err := k.WhitelistAddresses.Get(ctx, address)
	if err != nil {
		return false
	}

	return found
}

func (k Keeper) EVMContractExist(ctx context.Context, address common.Address) (bool, string, error) {
	exist := false
	tokenId := ""
	err := k.Token.Walk(ctx, nil, func(key string, token types.Token) (stop bool, err error) {
		if token.EvmAddress == address.String() {
			exist = true
			tokenId = key
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		return false, "", err
	}

	return exist, tokenId, nil
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.Params.Get(ctx)
}

func (k Keeper) IsTokenManager(ctx context.Context, tokenId string, addr []byte) (bool, error) {
	exist := false
	tm, err := k.TokenManagement.Get(ctx, tokenId)
	if err != nil {
		return false, err
	}
	for _, manager := range tm.Managers {
		if bytes.Equal(addr, manager) {
			exist = true
			break
		}
	}
	return exist, nil
}

func (k Keeper) IsTokenOwner(ctx context.Context, tokenId string, addr []byte) (bool, error) {
	token, err := k.Token.Get(ctx, tokenId)
	if err != nil {
		return false, err
	}
	
	if bytes.Equal(token.Issuer, addr) {
		return true, nil
	}
	return false, nil
}

func (k Keeper) IsFreezed(ctx context.Context, addr common.Address) bool {
	exist, err := k.FreezeAddresses.Get(ctx, addr.Bytes())
	if err != nil {
		return false
	}
	return exist
}

func (k Keeper) SetFreezeAddress(ctx context.Context, addr common.Address)  error {
	err := k.FreezeAddresses.Set(ctx, addr.Bytes(), true)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetToken(ctx context.Context, tokenId string) (types.Token, error) {
	return k.Token.Get(ctx, tokenId)
}

func (k Keeper) GetTokenManager(ctx context.Context, tokenId string) (types.TokenManagement, error) {
	return k.TokenManagement.Get(ctx, tokenId)
}
