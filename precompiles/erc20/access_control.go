// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package erc20

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/os/x/evm/core/vm"
	assettypes "github.com/realiotech/realio-network/x/asset/types"
)

func (p *Precompile) GrantRole(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	to, role, err := ParseGrantRoleArgs(args)
	if err != nil {
		return nil, err
	}

	return p.grantRole(ctx, contract, stateDB, method, to, role)
}

func (p *Precompile) RevokeRole(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	to, role, err := ParseRevokeRoleArgs(args)
	if err != nil {
		return nil, err
	}

	return p.grantRole(ctx, contract, stateDB, method, to, role)
}

func (p *Precompile) grantRole(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	to common.Address,
	role int32,
) (data []byte, err error) {
	sender := contract.CallerAddress

	if role != int32(assettypes.ManagerRole) {
		return nil, fmt.Errorf("only accept manager role")
	}
	err = p.assetKeep.GrantRole(ctx, p.denom, sender.Bytes(), to.Bytes())
	if err != nil {
		return nil, err
	}

	if err = p.EmitGrantRoleEvent(ctx, stateDB, sender, to, role); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (p *Precompile) revokeRole(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	to common.Address,
	role int32,
) (data []byte, err error) {
	sender := contract.CallerAddress

	if role != int32(assettypes.ManagerRole) {
		return nil, fmt.Errorf("only accept manager role")
	}
	err = p.assetKeep.RevokeRole(ctx, p.denom, sender.Bytes(), to.Bytes())
	if err != nil {
		return nil, err
	}

	if err = p.EmitRevokeRoleEvent(ctx, stateDB, sender, to, role); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}
