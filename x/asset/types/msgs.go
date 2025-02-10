package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateToken{}

func NewMsgCreateToken(issuer string, name string, symbol string, description string, decimal uint32, managers, distributors, extensionsList []string, allowNewExtensions bool) *MsgCreateToken {
	return &MsgCreateToken{
		Issuer:             issuer,
		Name:               name,
		Symbol:             symbol,
		Description:        description,
		Decimal:            decimal,
		Managers:           managers,
		ExtensionsList:     extensionsList,
		AllowNewExtensions: allowNewExtensions,
	}
}

func (msg *MsgCreateToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Issuer)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid issuer address (%s)", err)
	}
	for _, manager := range msg.Managers {
		_, err := sdk.AccAddressFromBech32(manager)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid manager address (%s): %s", manager, err)
		}
	}

	return nil
}

func NewMsgAssignRoles(issuer string, tokenId string, managers []string) *MsgAssignRoles {
	return &MsgAssignRoles{
		Issuer:       issuer,
		TokenId:      tokenId,
		Managers:     managers,
	}
}

func (msg *MsgAssignRoles) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Issuer)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid issuer address (%s)", err)
	}

	for _, manager := range msg.Managers {
		_, err := sdk.AccAddressFromBech32(manager)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid manager address (%s): %s", manager, err)
		}
	}

	return ValidateTokenId(msg.TokenId)
}

func NewMsgUnassignRoles(issuer string, tokenId string, managers []string) *MsgUnassignRoles {
	return &MsgUnassignRoles{
		Issuer:       issuer,
		TokenId:      tokenId,
		Managers:     managers,
	}
}

func (msg *MsgUnassignRoles) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Issuer)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid issuer address (%s)", err)
	}

	for _, manager := range msg.Managers {
		_, err := sdk.AccAddressFromBech32(manager)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid manager address (%s): %s", manager, err)
		}
	}

	return ValidateTokenId(msg.TokenId)
}
