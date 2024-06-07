package v1

import (
	"github.com/realiotech/realio-network/v2/app"
	"time"
)

var (
	UpgradeName       = "v1"
	ForkHeight        = int64(5989487)
	oneEnternityLater = time.Date(9999, 9, 9, 9, 9, 9, 9, time.UTC)
)

var Fork = app.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  ForkHeight,
	BeginForkLogic: RunForkLogic,
}
