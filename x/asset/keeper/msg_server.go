package keeper

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/ethereum/go-ethereum/common"
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

// CreateToken allow issuer to define new token.
func (ms msgServer) CreateToken(ctx context.Context, msg *types.MsgCreateToken) (*types.MsgCreateTokenResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if !ms.GetWhitelistAddress(ctx, msg.Issuer) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorize, "issuer not in whitelisted addresses")
	}

	lowerCaseName := strings.ToLower(msg.Name)
	lowerCaseSymbol := strings.ToLower(msg.Symbol)
	tokenId := fmt.Sprintf("%s/%s/%s", types.ModuleName, msg.Issuer, lowerCaseSymbol)

	isFound := ms.bk.HasSupply(ctx, tokenId)
	if isFound {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "token with id %s already exists", tokenId)
	}

	// Create a evm addr from tokenId
	evmAddr := common.BytesToAddress([]byte(tokenId))

	token := types.NewToken(tokenId, msg.Name, msg.Decimal, msg.Description, msg.Symbol, msg.Issuer, evmAddr.String())
	err := ms.Token.Set(ctx, tokenId, token)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenSet, err.Error())
	}

	tokenManagement := types.NewTokenManagement(msg.Managers, msg.AllowNewExtensions, msg.ExtensionsList, msg.MaxSupply)
	err = ms.TokenManagement.Set(ctx, tokenId, tokenManagement)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenManagementSet, err.Error())
	}

	ms.bk.SetDenomMetaData(ctx, banktypes.Metadata{
		Base: tokenId, Symbol: lowerCaseSymbol, Name: lowerCaseName,
		DenomUnits: []*banktypes.DenomUnit{{Denom: lowerCaseSymbol, Exponent: msg.Decimal}, {Denom: tokenId, Exponent: 0}},
	})

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenCreated,
			sdk.NewAttribute(types.AttributeKeyTokenId, tokenId),
			sdk.NewAttribute(types.AttributeKeyAddress, sdk.AccAddress(msg.Issuer).String()),
		),
	)

	return &types.MsgCreateTokenResponse{
		TokenId: tokenId,
	}, nil
}

func (ms msgServer) AssignRoles(ctx context.Context, msg *types.MsgAssignRoles) (*types.MsgAssignRolesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	token, err := ms.Token.Get(ctx, msg.TokenId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrTokenGet, err.Error())
	}

	if !bytes.Equal(msg.Issuer, token.Issuer) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorize, "issuer not the creator of the token")
	}

	tokenManagement, err := ms.TokenManagement.Get(ctx, msg.TokenId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrTokenManagementGet, err.Error())
	}
	newManagers := append(tokenManagement.Managers, msg.Managers...)
	tokenManagement.Managers = newManagers

	err = ms.TokenManagement.Set(ctx, msg.TokenId, tokenManagement)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenManagementSet, err.Error())
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenAuthorizeUpdated,
			sdk.NewAttribute(types.AttributeKeyTokenId, msg.TokenId),
		),
	)

	return &types.MsgAssignRolesResponse{}, nil
}

func (ms msgServer) UnassignRoles(ctx context.Context, msg *types.MsgUnassignRoles) (*types.MsgUnassignRolesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	token, err := ms.Token.Get(ctx, msg.TokenId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrTokenGet, err.Error())
	}

	if !bytes.Equal(msg.Issuer, token.Issuer) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorize, "issuer not the creator of the token")
	}

	tokenManagement, err := ms.TokenManagement.Get(ctx, msg.TokenId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrTokenManagementGet, err.Error())
	}
	managers := tokenManagement.Managers
	for i, b := range managers {
		if bytes.Equal(b, msg.Managers) {
			tokenManagement.Managers = append(managers[:i], managers[i+1:]...)
		}
	}

	err = ms.TokenManagement.Set(ctx, msg.TokenId, tokenManagement)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenManagementSet, err.Error())
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenAuthorizeUpdated,
			sdk.NewAttribute(types.AttributeKeyTokenId, msg.TokenId),
		),
	)

	return &types.MsgUnassignRolesResponse{}, nil
}

// UpdateParams updates the params.
func (ms msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := ms.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
