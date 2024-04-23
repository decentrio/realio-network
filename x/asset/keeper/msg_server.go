package keeper

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	"cosmossdk.io/math"

	realionetworktypes "github.com/realiotech/realio-network/types"
	"github.com/realiotech/realio-network/x/asset/types"
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

func (k msgServer) AuthorizeAddress(goCtx context.Context, msg *types.MsgAuthorizeAddress) (*types.MsgAuthorizeAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	token, isFound := k.GetToken(ctx, msg.Symbol)
	if !isFound {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "symbol %s does not exists", msg.Symbol)
	}

	// assert that the manager account is the only signer of the message
	if msg.Manager != token.Manager {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "caller not authorized")
	}

	accAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}

	token.AuthorizeAddress(accAddress)
	k.SetToken(ctx, token)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenAuthorized,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
		),
	)

	return &types.MsgAuthorizeAddressResponse{}, nil
}

func (k msgServer) UnAuthorizeAddress(goCtx context.Context, msg *types.MsgUnAuthorizeAddress) (*types.MsgUnAuthorizeAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	token, isFound := k.GetToken(ctx, msg.Symbol)
	if !isFound {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "symbol %s does not exists", msg.Symbol)
	}

	// assert that the manager account is the only signer of the message
	if msg.Manager != token.Manager {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "caller not authorized")
	}

	accAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}

	token.UnAuthorizeAddress(accAddress)
	k.SetToken(ctx, token)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenUnAuthorized,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
		),
	)

	return &types.MsgUnAuthorizeAddressResponse{}, nil
}

func (k msgServer) CreateToken(goCtx context.Context, msg *types.MsgCreateToken) (*types.MsgCreateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	lowerCaseSymbol := strings.ToLower(msg.Symbol)
	lowerCaseName := strings.ToLower(msg.Name)
	baseDenom := fmt.Sprintf("a%s", lowerCaseSymbol)

	isFound := k.bankKeeper.HasSupply(ctx, baseDenom)
	if isFound {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "token with denom %s already exists", baseDenom)
	}

	managerAccAddress, err := sdk.AccAddressFromBech32(msg.Manager)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid manager address")
	}

	token := types.NewToken(lowerCaseName, lowerCaseSymbol, msg.Total, msg.Manager, msg.AuthorizationRequired)

	if msg.AuthorizationRequired {
		// create authorization for module account and manager
		assetModuleAddress := k.ak.GetModuleAddress(types.ModuleName)
		moduleAuthorization := types.NewAuthorization(assetModuleAddress)
		newAuthorizationManager := types.NewAuthorization(managerAccAddress)
		token.Authorized = append(token.Authorized, moduleAuthorization, newAuthorizationManager)
	}

	k.SetToken(ctx, token)

	k.bankKeeper.SetDenomMetaData(ctx, bank.Metadata{
		Base: baseDenom, Symbol: lowerCaseSymbol, Name: lowerCaseName,
		DenomUnits: []*bank.DenomUnit{{Denom: lowerCaseSymbol, Exponent: 18}, {Denom: baseDenom, Exponent: 0}},
	})

	// mint coins for the current module
	// normalize into chains 10^18 denomination
	totalInt, _ := math.NewIntFromString(msg.Total)
	canonicalAmount := totalInt.Mul(realionetworktypes.PowerReduction)
	coin := sdk.Coins{{Denom: baseDenom, Amount: canonicalAmount}}

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coin)
	if err != nil {
		panic(err)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, managerAccAddress, coin)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenCreated,
			sdk.NewAttribute(sdk.AttributeKeyAmount, fmt.Sprint(msg.Total)),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
		),
	)

	return &types.MsgCreateTokenResponse{}, nil
}

func (k msgServer) UpdateToken(goCtx context.Context, msg *types.MsgUpdateToken) (*types.MsgUpdateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	existing, isFound := k.GetToken(ctx, msg.Symbol)
	if !isFound {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "symbol %s does not exists", msg.Symbol)
	}

	// assert that the manager account is the only signer of the message
	if msg.Manager != existing.Manager {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "caller not authorized")
	}

	// only Authorization Flag is updatable at this time
	token := types.Token{
		Name:                  existing.Name,
		Symbol:                existing.Symbol,
		Total:                 existing.Total,
		Manager:               existing.Manager,
		AuthorizationRequired: msg.AuthorizationRequired,
	}

	k.SetToken(ctx, token)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenUpdated,
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
		),
	)

	return &types.MsgUpdateTokenResponse{}, nil
}
