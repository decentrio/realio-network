package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleDenomPrefix     = "asset"
	MaxSubAssetNameLength = 44
	MaxHrpLength          = 16
	MaxCreatorLength      = 59 + MaxHrpLength
)

// GetTokenDenom constructs a denom string for tokens created by asset module
// based on an input creator address and a subname
// The denom constructed is asset/{creator}/{subdenom}
func GetAssetName(creator, subname string) (string, error) {
	if len(subname) > MaxSubAssetNameLength {
		return "", ErrSubdenomTooLong
	}
	if len(creator) > MaxCreatorLength {
		return "", ErrCreatorTooLong
	}
	if strings.Contains(creator, "/") {
		return "", ErrInvalidCreator
	}
	denom := strings.Join([]string{ModuleDenomPrefix, creator, subname}, "/")
	return denom, sdk.ValidateDenom(denom)
}
