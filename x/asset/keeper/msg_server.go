package keeper

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

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

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	issuerAddr, err := ms.ak.AddressCodec().StringToBytes(msg.Issuer)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrAccAddress, err.Error())
	}

	if !ms.GetWhitelistAddress(ctx, issuerAddr) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorize, "issuer not in whitelisted addresses")
	}

	lowerCaseName := strings.ToLower(msg.Name)
	lowerCaseSymbol := strings.ToLower(msg.Symbol)
	tokenId := fmt.Sprintf("%s/%s/%s", types.ModuleName, msg.Issuer, lowerCaseSymbol)

	isFound := ms.bk.HasSupply(ctx, tokenId)
	if isFound {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "token with id %s already exists", tokenId)
	}

	// TODO: create evm precompile here

	token := types.NewToken(tokenId, msg.Name, msg.Decimal, msg.Description, msg.Symbol, msg.Issuer)
	err = ms.Token.Set(ctx, tokenId, token)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenSet, err.Error())
	}

	tokenManagement := types.NewTokenManagement(msg.Managers, msg.AllowNewExtensions, msg.ExtensionsList)
	err = ms.TokenManagement.Set(ctx, tokenId, tokenManagement)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenManagementSet, err.Error())
	}

	tokenDistribution := types.NewTokenDistribution(msg.Distributors, msg.MaxSupply)
	err = ms.TokenDistribution.Set(ctx, tokenId, tokenDistribution)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenDistributionSet, err.Error())
	}

	ms.bk.SetDenomMetaData(ctx, banktypes.Metadata{
		Base: tokenId, Symbol: lowerCaseSymbol, Name: lowerCaseName,
		DenomUnits: []*banktypes.DenomUnit{{Denom: lowerCaseSymbol, Exponent: msg.Decimal}, {Denom: tokenId, Exponent: 0}},
	})

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenCreated,
			sdk.NewAttribute(types.AttributeKeyTokenId, tokenId),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Issuer),
		),
	)

	return &types.MsgCreateTokenResponse{
		TokenId: tokenId,
	}, nil
}

func (ms msgServer) AssignRoles(ctx context.Context, msg *types.MsgAssignRoles) (*types.MsgAssignRolesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	token, err := ms.Token.Get(ctx, msg.TokenId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrTokenGet, err.Error())
	}

	if msg.Issuer != token.Issuer {
		return nil, errorsmod.Wrapf(types.ErrUnauthorize, "issuer not the creator of the token")
	}

	tokenManagement, err := ms.TokenManagement.Get(ctx, msg.TokenId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrTokenManagementGet, err.Error())
	}
	tokenManagement.Managers = msg.Managers
	err = ms.TokenManagement.Set(ctx, msg.TokenId, tokenManagement)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenManagementSet, err.Error())
	}

	tokenDistribution, err := ms.TokenDistribution.Get(ctx, msg.TokenId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrTokenDistributionGet, err.Error())
	}
	tokenDistribution.Distributors = msg.Distributors
	err = ms.TokenDistribution.Set(ctx, msg.TokenId, tokenDistribution)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrTokenDistributionSet, err.Error())
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenAuthorizeUpdated,
			sdk.NewAttribute(types.AttributeKeyTokenId, msg.TokenId),
		),
	)

	return &types.MsgAssignRolesResponse{}, nil
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
