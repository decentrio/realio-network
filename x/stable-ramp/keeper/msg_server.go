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

	"github.com/realiotech/realio-network/crypto/ossecp256k1"
	realionetworktypes "github.com/realiotech/realio-network/types"
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

	amount := math.NewIntFromUint64(msg.Amount)
	coin := sdk.NewCoin(realionetworktypes.BaseDenom, amount)
	balance := ms.bankKeeper.GetBalance(goCtx, sender, coin.Denom)
	if balance.Amount.LT(coin.Amount) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"sender %s has %s, required %s",
			msg.SenderAddr,
			balance.String(),
			coin.String(),
		)
	}

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

	ctx := sdk.UnwrapSDKContext(goCtx)
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

	committee, err := ms.Committees.Get(goCtx, msg.ConnectionId)
	if err != nil {
		return nil, err
	}
	if !committeeContainsMemberAddress(committee, memberAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "member %s is not in the committee", msg.MemberAddr)
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

	committee, err := ms.Committees.Get(goCtx, msg.ConnectionId)
	if err != nil {
		return nil, err
	}
	if !committeeContainsMemberAddress(committee, memberAddr) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "member %s is not in the committee", msg.MemberAddr)
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

	if err := ms.PendingPackages.Set(goCtx, pendingPackageKey, trackedPackage); err != nil {
		return nil, err
	}

	return &types.MsgWithdrawClaimResponse{}, nil
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
