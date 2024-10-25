package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realiotech/realio-network/x/asset/types"
)

func (k msgServer) AuthorizeAddress(goCtx context.Context, msg *types.MsgAuthorizeAddress) (*types.MsgAuthorizeAddressResponse, error) {
	// Check if the value exists
	token, err := k.Token.Get(goCtx, types.TokenKey(msg.Symbol))
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "symbol %s does not exists : %s", msg.Symbol, err.Error())
	}

	// Checks if the token manager signed
	signers, _, err := k.cdc.GetMsgV1Signers(msg)
	if err != nil {
		return nil, err
	}

	if len(signers) != 1 {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "invalid signers")
	}

	// assert that the manager account is the only signer of the message
	if sdk.AccAddress(signers[0]).String() != token.Manager {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "caller not authorized")
	}

	accAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}

	token.AuthorizeAddress(accAddress)
	err = k.Token.Set(goCtx, types.TokenKey(msg.Symbol), token)
	if err != nil {
		return nil, types.ErrSetTokenUnable
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenAuthorized,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
		),
	)

	return &types.MsgAuthorizeAddressResponse{}, nil
}
