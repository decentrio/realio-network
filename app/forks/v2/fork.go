package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	baseapp "github.com/realiotech/realio-network/v2/app"
)

func RunForkLogic(ctx sdk.Context, app *baseapp.RealioNetwork) {
	ctx.Logger().Info("Starting upgrade for Realio-network v2...")
	fixMinCommisionRate(ctx, app.StakingKeeper, app.GetSubspace(stakingtypes.ModuleName))
	migrateParamSubspace(ctx, app.ConsensusParamsKeeper, app.ParamsKeeper)

	if _, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, app.AppCodec(), app.IBCKeeper.ClientKeeper); err != nil {
		panic(err)
	}
}
