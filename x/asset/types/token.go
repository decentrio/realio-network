package types

import (
	fmt "fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewToken creates a new Token instance
func NewToken(id, name string, decimal uint32, description string, symbol string, issuer []byte, evmAddress string) Token {
	return Token{
		TokenId:     id,
		Name:        name,
		Decimal:     decimal,
		Description: description,
		Symbol:      symbol,
		Issuer:      issuer,
		EvmAddress:  evmAddress,
	}
}

func NewTokenManagement(managers [][]byte, allowNewExtension bool, extensionList []string, maxSupply math.Int) TokenManagement {
	return TokenManagement{
		Managers:           managers,
		AllowNewExtensions: allowNewExtension,
		ExtensionsList:     extensionList,
		MaxSupply: maxSupply,
	}
}

func ValidateTokenId(tokenId string) error {
	tokenParts := strings.Split(tokenId, "/")
	if len(tokenParts) < 3 {
		return fmt.Errorf("invalid token id format, should be asset/IssuerAddress/lowercaseTokenName")
	}

	if tokenParts[0] != ModuleName {
		return fmt.Errorf("invalid token id format, should be asset/IssuerAddress/lowercaseTokenName")
	}

	_, err := sdk.AccAddressFromBech32(tokenParts[1])
	if err != nil {
		return fmt.Errorf("invalid issuer address")
	}

	tokenName := strings.Join(tokenParts[2:], "/")
	if strings.ToLower(tokenName) != tokenName {
		return fmt.Errorf("token name should be in lower case")
	}

	return nil
}
