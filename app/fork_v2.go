package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	evmtypes "github.com/evmos/evmos/v18/x/evm/types"
	feemarkettypes "github.com/evmos/evmos/v18/x/feemarket/types"
	minttypes "github.com/realiotech/realio-network/v2/x/mint/types"
)

const (
	V2UpgradeName       = "v2"
	V2ForkHeight        = 7084410
	NewMinCommisionRate = "0.05"
)

var v2Fork = Fork{
	UpgradeName: V2UpgradeName,
	UpgradeHeight: V2ForkHeight,
	BeginForkLogic: V2RunForkLogic,
}

func V2RunForkLogic(ctx sdk.Context, app *RealioNetwork) {
	ctx.Logger().Info("Starting upgrade for Realio-network v2...")
	fixMinCommisionRate(ctx, app.StakingKeeper, app.GetSubspace(stakingtypes.ModuleName))
	migrateParamSubspace(ctx, app.ConsensusParamsKeeper, app.ParamsKeeper)

	if _, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, app.AppCodec(), app.IBCKeeper.ClientKeeper); err != nil {
		panic(err)
	}

	storeUpgrades := &storetypes.StoreUpgrades{
		Added: []string{crisistypes.ModuleName, consensusparamtypes.ModuleName},
	}
	app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(V2ForkHeight, storeUpgrades))
}

func fixMinCommisionRate(ctx sdk.Context, staking *stakingkeeper.Keeper, stakingLegacySubspace paramstypes.Subspace) {
	// Upgrade every validators min-commission rate
	validators := staking.GetAllValidators(ctx)
	minComm := sdk.MustNewDecFromStr(NewMinCommisionRate)
	if stakingLegacySubspace.HasKeyTable() {
		stakingLegacySubspace.Set(ctx, stakingtypes.KeyMinCommissionRate, minComm)
	} else {
		stakingLegacySubspace.WithKeyTable(stakingtypes.ParamKeyTable())
		stakingLegacySubspace.Set(ctx, stakingtypes.KeyMinCommissionRate, minComm)
	}

	for _, v := range validators {
		//nolint
		if v.Commission.Rate.LT(minComm) {
			v.Commission = updateValidatorCommission(ctx, v, minComm)

			// call the before-modification hook since we're about to update the commission
			staking.Hooks().BeforeValidatorModified(ctx, v.GetOperator())
			staking.SetValidator(ctx, v)
		}
	}
}

func updateValidatorCommission(ctx sdk.Context,
	validator stakingtypes.Validator, newRate sdk.Dec,
) stakingtypes.Commission {
	commission := validator.Commission
	blockTime := ctx.BlockHeader().Time

	commission.Rate = newRate
	if commission.MaxRate.LT(newRate) {
		commission.MaxRate = newRate
	}

	commission.UpdateTime = blockTime

	return commission
}

func migrateParamSubspace(ctx sdk.Context, ck consensuskeeper.Keeper, pk paramskeeper.Keeper) {
	for _, subspace := range pk.GetSubspaces() {
		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case authtypes.ModuleName:
			keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
		case banktypes.ModuleName:
			keyTable = banktypes.ParamKeyTable() //nolint:staticcheck,nolintlint
		case stakingtypes.ModuleName:
			keyTable = stakingtypes.ParamKeyTable()
		case minttypes.ModuleName:
			keyTable = minttypes.ParamKeyTable()
		case distrtypes.ModuleName:
			keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck,nolintlint
		case slashingtypes.ModuleName:
			keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
		case govtypes.ModuleName:
			keyTable = govv1.ParamKeyTable() //nolint:staticcheck
		case crisistypes.ModuleName:
			keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck
		case ibctransfertypes.ModuleName:
			keyTable = ibctransfertypes.ParamKeyTable()
		case evmtypes.ModuleName:
			keyTable = evmtypes.ParamKeyTable() //nolint:staticcheck
		case feemarkettypes.ModuleName:
			keyTable = feemarkettypes.ParamKeyTable()
		default:
			continue
		}
		if !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}

	baseAppLegacySS := pk.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
	baseapp.MigrateParams(ctx, baseAppLegacySS, &ck)
}
