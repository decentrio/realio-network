package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var Forks = []Fork{}

type Fork struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string
	// height the upgrade occurs at
	UpgradeHeight int64

	// Function that runs some custom state transition code at the beginning of a fork.
	BeginForkLogic func(ctx sdk.Context, app *RealioNetwork)
}

func BeginBlockForks(ctx sdk.Context, app *RealioNetwork) {
	for _, fork := range Forks {
		if ctx.BlockHeight() == fork.UpgradeHeight {
			fork.BeginForkLogic(ctx, app)
			return
		}
	}
}
