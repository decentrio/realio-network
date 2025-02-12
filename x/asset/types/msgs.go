package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateToken{}

func NewMsgCreateToken(issuer []byte, name string, symbol string, description string, decimal uint32, managers [][]byte, extensionsList []string, allowNewExtensions bool) *MsgCreateToken {
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

func NewMsgAssignRoles(issuer []byte, tokenId string, managers [][]byte) *MsgAssignRoles {
	return &MsgAssignRoles{
		Issuer:   issuer,
		TokenId:  tokenId,
		Managers: managers,
	}
}

func NewMsgUnassignRoles(issuer []byte, tokenId string, managers []byte) *MsgUnassignRoles {
	return &MsgUnassignRoles{
		Issuer:   issuer,
		TokenId:  tokenId,
		Managers: managers,
	}
}
