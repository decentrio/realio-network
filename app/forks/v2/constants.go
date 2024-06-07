package v2

import (
	"github.com/realiotech/realio-network/v2/app"
)

const (
	UpgradeName         = "v2"
	ForkHeight          = 7084400
	NewMinCommisionRate = "0.05"
)

var Fork = app.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  ForkHeight,
	BeginForkLogic: RunForkLogic,
}