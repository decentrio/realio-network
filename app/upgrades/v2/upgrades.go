package v2

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/ethereum/go-ethereum/common"
	ostypes "github.com/evmos/os/types"
	evmkeeper "github.com/evmos/os/x/evm/keeper"
	evmtypes "github.com/evmos/os/x/evm/types"
	assettypes "github.com/realiotech/realio-network/x/asset/types"
	bridgetypes "github.com/realiotech/realio-network/x/bridge/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	accountkeeper authkeeper.AccountKeeper,
	evmkeeper *evmkeeper.Keeper,
	paramskeeper paramskeeper.Keeper,
	consensuskeeper consensusparamkeeper.Keeper,
	IBCKeeper ibckeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		for _, subspace := range paramskeeper.GetSubspaces() {
			subspace := subspace
			fmt.Println("subspace.Name(): ", subspace.Name())

			var keyTable paramstypes.KeyTable
			switch subspace.Name() {
			case evmtypes.ModuleName:
				keyTable = evmtypes.ParamKeyTable() //nolint:staticcheck
			case assettypes.ModuleName:
				keyTable = assettypes.ParamKeyTable() //nolint:staticcheck
			case bridgetypes.ModuleName:
				keyTable = bridgetypes.ParamKeyTable() //nolint:staticcheck
			case banktypes.ModuleName:
				keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
			case stakingtypes.ModuleName:
				keyTable = stakingtypes.ParamKeyTable() //nolint:staticcheck
			case minttypes.ModuleName:
				keyTable = minttypes.ParamKeyTable() //nolint:staticcheck
			case distrtypes.ModuleName:
				keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck
			case slashingtypes.ModuleName:
				keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
			case govtypes.ModuleName:
				keyTable = govv1.ParamKeyTable() //nolint:staticcheck
			case crisistypes.ModuleName:
				keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck
			case authtypes.ModuleName:
				keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
			}

			if !subspace.HasKeyTable() {
				subspace.WithKeyTable(keyTable)
			}
		}

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		legacyBaseAppSubspace := paramskeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(sdkCtx, legacyBaseAppSubspace, consensuskeeper.ParamsStore)

		legacyClientSubspace, _ := paramskeeper.GetSubspace(exported.ModuleName)
		var params clienttypes.Params
		legacyClientSubspace.GetParamSet(sdkCtx, &params)

		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		IBCKeeper.ClientKeeper.SetParams(sdkCtx, params)

		MigrateEthAccountsToBaseAccounts(sdkCtx, accountkeeper, evmkeeper)
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// MigrateEthAccountsToBaseAccounts is used to store the code hash of the associated
// smart contracts in the dedicated store in the EVM module and convert the former
// EthAccounts to standard Cosmos SDK accounts.
func MigrateEthAccountsToBaseAccounts(ctx sdk.Context, ak authkeeper.AccountKeeper, ek *evmkeeper.Keeper) {
	ak.IterateAccounts(ctx, func(account sdk.AccountI) (stop bool) {
		ethAcc, ok := account.(*ostypes.EthAccount)
		if !ok {
			return false
		}

		// NOTE: we only need to add store entries for smart contracts
		codeHashBytes := common.HexToHash(ethAcc.CodeHash).Bytes()
		if !evmtypes.IsEmptyCodeHash(codeHashBytes) {
			ek.SetCodeHash(ctx, ethAcc.EthAddress().Bytes(), codeHashBytes)
		}

		// Set the base account in the account keeper instead of the EthAccount
		ak.SetAccount(ctx, ethAcc.BaseAccount)

		return false
	})
}
