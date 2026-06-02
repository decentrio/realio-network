package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"strconv"
	"strings"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/realiotech/realio-network/crypto/ossecp256k1"
	"github.com/realiotech/realio-network/x/stable-ramp/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (ms msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "withdraw amount must be positive")
	}

	sender, err := sdk.AccAddressFromBech32(msg.SenderAddr)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !common.IsHexAddress(msg.ReceiverAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "receiver_addr %s is not a valid Ethereum address", msg.ReceiverAddr)
	}

	if !common.IsHexAddress(msg.TokenAddress) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "token_address %s is not a valid Ethereum address", msg.TokenAddress)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	tokenAddr := common.HexToAddress(msg.TokenAddress)

	// Check erc20 token has presentation
	denom, err := ms.erc20Keeper.GetTokenDenom(ctx, tokenAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "token_address %s is not a registered ERC20 token: %s", msg.TokenAddress, err)
	}

	amount := math.NewIntFromUint64(msg.Amount)
	coin := sdk.NewCoin(denom, amount)

	committee, err := ms.Committees.Get(goCtx, msg.ConnectionId)
	if err != nil {
		return nil, err
	}

	nonce, err := ms.NextPackageNonce.Next(goCtx)
	if err != nil {
		return nil, err
	}

	pkg := types.Package{
		SenderAddr:   msg.SenderAddr,
		ReceiverAddr: msg.ReceiverAddr,
		ConnectionId: msg.ConnectionId,
		Amount:       msg.Amount,
		Denom:        coin.Denom,
		Nonce:        nonce,
		Status:       types.PackageStatus_PACKAGE_STATUS_PENDING,
		Action:       types.PackageAction_PACKAGE_ACTION_WITHDRAW,
	}

	trackedPackage := types.TrackedPackage{
		PackageData:  &pkg,
		Approvals:    0,
		TotalMembers: uint32(len(committee.Members)),
	}
	if err := ms.PendingPackages.Set(goCtx, collections.Join(nonce, msg.ConnectionId), trackedPackage); err != nil {
		return nil, err
	}

	// Convert ERC20 token into cosmos coin representation, minted to sender's account.
	_, err = ms.erc20Keeper.ConvertERC20(goCtx, &erc20types.MsgConvertERC20{
		Sender:          common.BytesToAddress(sender).Hex(),
		ContractAddress: msg.TokenAddress,
		Amount:          amount,
		Receiver:        msg.SenderAddr,
	})
	if err != nil {
		return nil, err
	}

	// Lock into module account for burning later
	if err := ms.bankKeeper.SendCoinsFromAccountToModule(goCtx, sender, types.ModuleName, sdk.Coins{coin}); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdraw,
			sdk.NewAttribute(types.AttributeKeySender, msg.SenderAddr),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.ReceiverAddr),
			sdk.NewAttribute(types.AttributeKeyAmount, coin.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, coin.Denom),
			sdk.NewAttribute(types.AttributeKeyConnection, msg.ConnectionId),
			sdk.NewAttribute(types.AttributeKeyNonce, strconv.FormatUint(nonce, 10)),
		),
	)

	return &types.MsgWithdrawResponse{}, nil
}

func (ms msgServer) DepositClaim(goCtx context.Context, msg *types.MsgDepositClaim) (*types.MsgDepositClaimResponse, error) {
	memberAddr, err := sdk.AccAddressFromBech32(msg.MemberAddr)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	pendingPackageKey := collections.Join(msg.PackageNonce, msg.ConnectionId)
	trackedPackage, err := ms.PendingPackages.Get(goCtx, pendingPackageKey)
	if err != nil {
		return nil, err
	}
	if trackedPackage.PackageData == nil {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "cannot find package with nounce %v for connection %s",
			msg.PackageNonce,
			msg.ConnectionId,
		)
	}
	if trackedPackage.PackageData.ConnectionId != msg.ConnectionId {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"claim connection_id %s does not match package connection_id %s",
			msg.ConnectionId,
			trackedPackage.PackageData.ConnectionId,
		)
	}
	if trackedPackage.PackageData.Action != types.PackageAction_PACKAGE_ACTION_DEPOSIT {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "package %d is not a deposit package", msg.PackageNonce)
	}

	committee, err := ms.Committees.Get(goCtx, msg.ConnectionId)
	if err != nil {
		return nil, err
	}
	if !committeeContainsMemberAddress(committee, memberAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "member %s is not in the committee", msg.MemberAddr)
	}

	if hasVoted(trackedPackage.VotedMembers, msg.MemberAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "member %s has already voted on package %d", msg.MemberAddr, msg.PackageNonce)
	}

	switch msg.Action {
	case types.ClaimAction_CLAIM_ACTION_APPROVE:
		if trackedPackage.Approvals < trackedPackage.TotalMembers {
			trackedPackage.Approvals++
		}
	case types.ClaimAction_CLAIM_ACTION_REJECT:
		trackedPackage.PackageData.Status = types.PackageStatus_PACKAGE_STATUS_FAILED
	default:
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid claim action: %s", msg.Action.String())
	}

	trackedPackage.VotedMembers = append(trackedPackage.VotedMembers, msg.MemberAddr)

	if err := ms.PendingPackages.Set(goCtx, pendingPackageKey, trackedPackage); err != nil {
		return nil, err
	}

	return &types.MsgDepositClaimResponse{}, nil
}

func (ms msgServer) WithdrawClaim(goCtx context.Context, msg *types.MsgWithdrawClaim) (*types.MsgWithdrawClaimResponse, error) {
	memberAddr, err := sdk.AccAddressFromBech32(msg.MemberAddr)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	pendingPackageKey := collections.Join(msg.PackageNonce, msg.ConnectionId)
	trackedPackage, err := ms.PendingPackages.Get(goCtx, pendingPackageKey)
	if err != nil {
		return nil, err
	}
	if trackedPackage.PackageData == nil {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "cannot find package with nonce %v for connection %s",
			msg.PackageNonce,
			msg.ConnectionId,
		)
	}
	if trackedPackage.PackageData.ConnectionId != msg.ConnectionId {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"claim connection_id %s does not match package connection_id %s",
			msg.ConnectionId,
			trackedPackage.PackageData.ConnectionId,
		)
	}
	if trackedPackage.PackageData.Action != types.PackageAction_PACKAGE_ACTION_WITHDRAW {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "package %d is not a withdraw package", msg.PackageNonce)
	}

	committee, err := ms.Committees.Get(goCtx, msg.ConnectionId)
	if err != nil {
		return nil, err
	}
	if !committeeContainsMemberAddress(committee, memberAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "member %s is not in the committee", msg.MemberAddr)
	}

	if hasVoted(trackedPackage.VotedMembers, msg.MemberAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "member %s has already voted on package %d", msg.MemberAddr, msg.PackageNonce)
	}

	switch msg.Action {
	case types.ClaimAction_CLAIM_ACTION_APPROVE:
		if trackedPackage.Approvals < trackedPackage.TotalMembers {
			trackedPackage.Approvals++
		}
		trackedPackage.Signatures = append(trackedPackage.Signatures, append([]byte(nil), msg.Signature...))
		trackedPackage.VotedMembers = append(trackedPackage.VotedMembers, msg.MemberAddr)
	case types.ClaimAction_CLAIM_ACTION_REJECT:
		senderAddr, err := sdk.AccAddressFromBech32(trackedPackage.PackageData.SenderAddr)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
		}
		refundCoin := sdk.NewCoin(trackedPackage.PackageData.Denom, math.NewIntFromUint64(trackedPackage.PackageData.Amount))
		if err := ms.RefundEscrowedCoins(goCtx, senderAddr, sdk.Coins{refundCoin}); err != nil {
			return nil, err
		}
		return &types.MsgWithdrawClaimResponse{}, ms.PendingPackages.Remove(goCtx, pendingPackageKey)
	default:
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid claim action: %s", msg.Action.String())
	}

	if err := ms.PendingPackages.Set(goCtx, pendingPackageKey, trackedPackage); err != nil {
		return nil, err
	}

	return &types.MsgWithdrawClaimResponse{}, nil
}

func hasVoted(votedMembers []string, memberAddr string) bool {
	for _, v := range votedMembers {
		if v == memberAddr {
			return true
		}
	}
	return false
}

func committeeContainsMemberAddress(committee types.Committee, memberAddr sdk.AccAddress) bool {
	for _, member := range committee.Members {
		if memberAddressMatches(member, memberAddr) {
			return true
		}
	}
	return false
}

func memberAddressMatches(member *types.Member, memberAddr sdk.AccAddress) bool {
	if member == nil {
		return false
	}

	pubkeyBytes, err := hex.DecodeString(strings.TrimPrefix(member.RealioPubkey, "0x"))
	if err != nil {
		return false
	}

	pubkey := ossecp256k1.PubKey{Key: pubkeyBytes}
	return bytes.Equal(sdk.AccAddress(pubkey.Address()), memberAddr)
}
