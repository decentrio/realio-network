package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/realiotech/realio-network/x/asset/types"
)

func (k Keeper) GrantRole(ctx context.Context, tokenId string, issuer []byte, manager []byte) error {
	token, err := k.Token.Get(ctx, tokenId)
	if err != nil {
		return errorsmod.Wrapf(types.ErrTokenGet, err.Error())
	}

	if !bytes.Equal(issuer, token.Issuer) {
		return errorsmod.Wrapf(types.ErrUnauthorize, "issuer not the creator of the token")
	}

	tokenManagement, err := k.TokenManagement.Get(ctx, tokenId)
	if err != nil {
		return errorsmod.Wrapf(types.ErrTokenManagementGet, err.Error())
	}
	newManagers := append(tokenManagement.Managers, manager)
	tokenManagement.Managers = newManagers

	return k.TokenManagement.Set(ctx, tokenId, tokenManagement)
}

func (k Keeper) RevokeRole(ctx context.Context, tokenId string, issuer []byte, manager []byte) error {
	token, err := k.Token.Get(ctx, tokenId)
	if err != nil {
		return errorsmod.Wrapf(types.ErrTokenGet, err.Error())
	}

	if !bytes.Equal(issuer, token.Issuer) {
		return errorsmod.Wrapf(types.ErrUnauthorize, "issuer not the creator of the token")
	}

	tokenManagement, err := k.TokenManagement.Get(ctx, tokenId)
	if err != nil {
		return errorsmod.Wrapf(types.ErrTokenManagementGet, err.Error())
	}
	managers := tokenManagement.Managers
	for i, b := range managers {
		if bytes.Equal(b, manager) {
			tokenManagement.Managers = append(managers[:i], managers[i+1:]...)
		}
	}

	return k.TokenManagement.Set(ctx, tokenId, tokenManagement)
}
