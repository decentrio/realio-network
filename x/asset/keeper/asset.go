package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realiotech/realio-network/x/asset/types"
)

func (k Keeper) validateAsset(ctx sdk.Context, creator string, subName string) (string, error) {
	if k.bankKeeper.HasSupply(ctx, subName) {
		return "", fmt.Errorf("can't create subdenoms that are the same as a native denom")
	}

	assetName, err := types.GetAssetName(creator, subName)
	if err != nil {
		return "", err
	}

	_, found := k.bankKeeper.GetDenomMetaData(ctx, assetName)
	if found {
		return "", types.ErrDenomExists
	}

	return assetName, nil
}

func (k Keeper) createAsset(ctx sdk.Context, creator string, denom string) error {

}

func (k Keeper) SetAsset(ctx sdk.Context, asset types.Asset) error {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.AssetKeyPrefix)

	key := []byte(asset.AssetName)
	bz, err := k.cdc.Marshal(&asset)
	if err != nil {
		return err
	}

	prefixStore.Set(key, bz)
	return nil
}

func (k Keeper) GetAsset(ctx sdk.Context, assetName string) (types.Asset, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.AssetKeyPrefix)

	key := []byte(assetName)
	bz := prefixStore.Get(key)

	var asset types.Asset
	err := k.cdc.Unmarshal(bz, &asset)
	if err != nil {
		return types.Asset{}, err
	}

	return asset, nil
}
