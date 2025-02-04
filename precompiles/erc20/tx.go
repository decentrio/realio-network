// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package erc20

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	cmn "github.com/evmos/os/precompiles/common"
	"github.com/evmos/os/x/evm/core/vm"
	evmtypes "github.com/evmos/os/x/evm/types"
	assettypes "github.com/realiotech/realio-network/x/asset/types"
)

const (
	// TransferMethod defines the ABI method name for the ERC-20 transfer
	// transaction.
	TransferMethod = "transfer"
	// TransferFromMethod defines the ABI method name for the ERC-20 transferFrom
	// transaction.
	TransferFromMethod = "transferFrom"

	BurnMethod = "burn"

	BurnFromMethod = "burnFrom"

	MintMethod = "mint"
)

// SendMsgURL defines the authorization type for MsgSend
var SendMsgURL = sdk.MsgTypeURL(&banktypes.MsgSend{})

// Transfer executes a direct transfer from the caller address to the
// destination address.
func (p *Precompile) Transfer(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	from := contract.CallerAddress
	to, amount, err := ParseTransferArgs(args)
	if err != nil {
		return nil, err
	}

	return p.transfer(ctx, contract, stateDB, method, from, to, amount)
}

// TransferFrom executes a transfer on behalf of the specified from address in
// the call data to the destination address.
func (p *Precompile) TransferFrom(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	from, to, amount, err := ParseTransferFromArgs(args)
	if err != nil {
		return nil, err
	}

	return p.transfer(ctx, contract, stateDB, method, from, to, amount)
}

// transfer is a common function that handles transfers for the ERC-20 Transfer
// and TransferFrom methods. It executes a bank Send message if the spender is
// the sender of the transfer, otherwise it executes an authorization.
func (p *Precompile) transfer(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	from, to common.Address,
	amount *big.Int,
) (data []byte, err error) {
	coins := sdk.Coins{{Denom: p.denom, Amount: math.NewIntFromBigInt(amount)}}

	msg := banktypes.NewMsgSend(from.Bytes(), to.Bytes(), coins)

	if err = msg.Amount.Validate(); err != nil {
		return nil, err
	}

	isTransferFrom := method.Name == TransferFromMethod
	owner := sdk.AccAddress(from.Bytes())
	spenderAddr := contract.CallerAddress
	spender := sdk.AccAddress(spenderAddr.Bytes()) // aka. grantee
	ownerIsSpender := spender.Equals(owner)

	var prevAllowance *big.Int
	if ownerIsSpender {
		msgSrv := bankkeeper.NewMsgServerImpl(p.BankKeeper)
		_, err = msgSrv.Send(ctx, msg)
	} else {
		_, _, prevAllowance, err = GetAuthzExpirationAndAllowance(p.AuthzKeeper, ctx, spenderAddr, from, p.denom)
		if err != nil {
			return nil, ConvertErrToERC20Error(errorsmod.Wrap(err, authz.ErrNoAuthorizationFound.Error()))
		}

		_, err = p.AuthzKeeper.DispatchActions(ctx, spender, []sdk.Msg{msg})
	}

	if err != nil {
		err = ConvertErrToERC20Error(err)
		// This should return an error to avoid the contract from being executed and an event being emitted
		return nil, err
	}

	evmDenom := evmtypes.GetEVMCoinDenom()
	if p.denom == evmDenom {
		convertedAmount := evmtypes.ConvertAmountTo18DecimalsBigInt(amount)
		p.SetBalanceChangeEntries(cmn.NewBalanceChangeEntry(from, convertedAmount, cmn.Sub),
			cmn.NewBalanceChangeEntry(to, convertedAmount, cmn.Add))
	}

	if err = p.EmitTransferEvent(ctx, stateDB, from, to, amount); err != nil {
		return nil, err
	}

	// NOTE: if it's a direct transfer, we return here but if used through transferFrom,
	// we need to emit the approval event with the new allowance.
	if !isTransferFrom {
		return method.Outputs.Pack(true)
	}

	var newAllowance *big.Int
	if ownerIsSpender {
		// NOTE: in case the spender is the owner we emit an approval event with
		// the maxUint256 value.
		newAllowance = abi.MaxUint256
	} else {
		newAllowance = new(big.Int).Sub(prevAllowance, amount)
	}

	if err = p.EmitApprovalEvent(ctx, stateDB, from, spenderAddr, newAllowance); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (p *Precompile) Mint(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	to, amount, err := ParseMintArgs(args)
	if err != nil {
		return nil, err
	}

	return p.mint(ctx, contract, stateDB, method, to, amount)
}

func (p *Precompile) mint(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	to common.Address,
	amount *big.Int,
) (data []byte, err error) {

	minter := contract.CallerAddress
	havePerm, err := p.assetKeep.IsTokenManager(ctx, p.denom, minter)
	fmt.Println("have perm", havePerm, err)
	if err != nil || !havePerm {
		return nil, fmt.Errorf("Sender is not token manager")
	}

	mintToAddr := sdk.AccAddress(to.Bytes())

	coins := sdk.Coins{{Denom: p.denom, Amount: math.NewIntFromBigInt(amount)}}

	// Check if new supply exceed max supply

	tm, err := p.assetKeep.GetTokenManager(ctx, p.denom)
	if err != nil {
		return nil, err
	}
	maxSupply := tm.MaxSupply
	currentSupply := p.BankKeeper.GetSupply(ctx, p.denom).Amount
	newSupply := currentSupply.Add(math.NewIntFromBigInt(amount))
	if newSupply.GT(maxSupply) {
		return nil, ConvertErrToERC20Error(fmt.Errorf("Exceed max supply, expected: %d, got: %d", maxSupply, newSupply))
	}

	// Mint coins to asset module then transfer to minter addr
	err = p.BankKeeper.MintCoins(ctx, assettypes.ModuleName, coins)
	if err != nil {
		return nil, ConvertErrToERC20Error(err)
	}

	err = p.BankKeeper.SendCoinsFromModuleToAccount(ctx, assettypes.ModuleName, mintToAddr, coins)
	if err != nil {
		return nil, ConvertErrToERC20Error(err)
	}

	evmDenom := evmtypes.GetEVMCoinDenom()
	if p.denom == evmDenom {
		convertedAmount := evmtypes.ConvertAmountTo18DecimalsBigInt(amount)
		p.SetBalanceChangeEntries(cmn.NewBalanceChangeEntry(to, convertedAmount, cmn.Add))
	}

	if err = p.EmitMintEvent(ctx, stateDB, to, amount); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (p *Precompile) Burn(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	from := contract.CallerAddress
	amount, err := ParseBurnArgs(args)
	if err != nil {
		return nil, err
	}

	return p.burn(ctx, contract, stateDB, method, from, amount)
}

func (p *Precompile) BurnFrom(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	from, amount, err := ParseBurnFromArgs(args)
	if err != nil {
		return nil, err
	}

	return p.burn(ctx, contract, stateDB, method, from, amount)
}

func (p *Precompile) burn(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	from common.Address,
	amount *big.Int,
) (data []byte, err error) {

	minter := contract.CallerAddress
	havePerm, err := p.assetKeep.IsTokenManager(ctx, p.denom, minter)
	if err != nil || !havePerm {
		return nil, err
	}

	burnFromAddr := sdk.AccAddress(from.Bytes())

	coins := sdk.Coins{{Denom: p.denom, Amount: math.NewIntFromBigInt(amount)}}

	// Transfer to asset module then burn
	err = p.BankKeeper.SendCoinsFromAccountToModule(ctx, burnFromAddr, assettypes.ModuleName, coins)
	if err != nil {
		return nil, ConvertErrToERC20Error(err)
	}

	err = p.BankKeeper.BurnCoins(ctx, assettypes.ModuleName, coins)
	if err != nil {
		return nil, ConvertErrToERC20Error(err)
	}

	evmDenom := evmtypes.GetEVMCoinDenom()
	if p.denom == evmDenom {
		convertedAmount := evmtypes.ConvertAmountTo18DecimalsBigInt(amount)
		p.SetBalanceChangeEntries(cmn.NewBalanceChangeEntry(from, convertedAmount, cmn.Add))
	}

	if err = p.EmitBurnEvent(ctx, stateDB, from, amount); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}
