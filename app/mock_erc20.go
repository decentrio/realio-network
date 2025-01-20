package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/os/x/evm/core/vm"
	transferkeeper "github.com/evmos/os/x/ibc/transfer/keeper"
	"github.com/realiotech/realio-network/precompiles/erc20"
	assetkeeper "github.com/realiotech/realio-network/x/asset/keeper"
)

type MockErc20Keeper struct {
	ContractAddrs []string

	bankKeeper     bankkeeper.Keeper
	authzKeeper    authzkeeper.Keeper
	transferKeeper *transferkeeper.Keeper
	assetKeeper    assetkeeper.Keeper
}

func NewMockErc20KeeperWithAddrs(addrs []string) MockErc20Keeper {
	return MockErc20Keeper{ContractAddrs: addrs}
}

func NewEmptyMockErc20Keeper() MockErc20Keeper {
	return MockErc20Keeper{}
}

func (k MockErc20Keeper) GetERC20PrecompileInstance(ctx sdk.Context, addr common.Address) (contract vm.PrecompiledContract, found bool, err error) {
	// Check if contract address in list
	exist, err := k.assetKeeper.EVMContractExist(ctx, addr)
	if err != nil || !exist {
		return nil, false, nil
	}

	precompile, err := erc20.NewPrecompile(addr, k.bankKeeper, k.authzKeeper, *k.transferKeeper, k.assetKeeper)
	if err != nil {
		return nil, false, nil
	}

	return precompile, true, nil
}

func (k MockErc20Keeper) AddContractAddress(_ sdk.Context, addr common.Address) {
	k.ContractAddrs = append(k.ContractAddrs, addr.String())
}
