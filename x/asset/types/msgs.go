package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgAuthorizeAddress   = "authorize_address"
	TypeMsgUnAuthorizeAddress = "un_authorize_address"
	TypeMsgCreateToken        = "create_token"
	TypeMsgUpdateToken        = "update_token"
)

var (
	_ sdk.Msg = &MsgAuthorizeAddress{}
	_ sdk.Msg = &MsgUnAuthorizeAddress{}
	_ sdk.Msg = &MsgCreateToken{}
	_ sdk.Msg = &MsgUpdateToken{}
)
func NewMsgAuthorizeAddress(manager string, symbol string, address string) *MsgAuthorizeAddress {
	return &MsgAuthorizeAddress{
		Manager: manager,
		Symbol:  symbol,
		Address: address,
	}
}

func (msg *MsgAuthorizeAddress) Route() string {
	return RouterKey
}

func (msg *MsgAuthorizeAddress) Type() string {
	return TypeMsgAuthorizeAddress
}

func (msg *MsgAuthorizeAddress) GetSigners() []sdk.AccAddress {
	manager, err := sdk.AccAddressFromBech32(msg.Manager)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{manager}
}

func (msg *MsgAuthorizeAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAuthorizeAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Manager)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid manager address (%s)", err)
	}
	return nil
}

func NewMsgUnAuthorizeAddress(manager string, symbol string, address string) *MsgUnAuthorizeAddress {
	return &MsgUnAuthorizeAddress{
		Manager: manager,
		Symbol:  symbol,
		Address: address,
	}
}

func (msg *MsgUnAuthorizeAddress) Route() string {
	return RouterKey
}

func (msg *MsgUnAuthorizeAddress) Type() string {
	return TypeMsgUnAuthorizeAddress
}

func (msg *MsgUnAuthorizeAddress) GetSigners() []sdk.AccAddress {
	manager, err := sdk.AccAddressFromBech32(msg.Manager)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{manager}
}

func (msg *MsgUnAuthorizeAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnAuthorizeAddress) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Manager); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid manager address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", err)
	}

	return nil
}

func NewMsgCreateToken(manager string, name string, symbol string, total string, authorizationRequired bool) *MsgCreateToken {
	return &MsgCreateToken{
		Manager:               manager,
		Name:                  name,
		Symbol:                symbol,
		Total:                 total,
		AuthorizationRequired: authorizationRequired,
	}
}

func (msg *MsgCreateToken) Route() string {
	return RouterKey
}

func (msg *MsgCreateToken) Type() string {
	return TypeMsgCreateToken
}

func (msg *MsgCreateToken) GetSigners() []sdk.AccAddress {
	manager, err := sdk.AccAddressFromBech32(msg.Manager)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{manager}
}

func (msg *MsgCreateToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Manager)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid manager address (%s)", err)
	}
	return nil
}

func NewMsgUpdateToken(manager string, symbol string, authorizationRequired bool) *MsgUpdateToken {
	return &MsgUpdateToken{
		Manager:               manager,
		Symbol:                symbol,
		AuthorizationRequired: authorizationRequired,
	}
}

func (msg *MsgUpdateToken) Route() string {
	return RouterKey
}

func (msg *MsgUpdateToken) Type() string {
	return TypeMsgUpdateToken
}

func (msg *MsgUpdateToken) GetSigners() []sdk.AccAddress {
	manager, err := sdk.AccAddressFromBech32(msg.Manager)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{manager}
}

func (msg *MsgUpdateToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateToken) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Manager)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid manager address (%s)", err)
	}
	return nil
}
