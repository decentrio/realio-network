package stableramp

import (
	"context"
	"encoding/hex"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/realiotech/realio-network/x/stable-ramp/keeper"
	"github.com/realiotech/realio-network/x/stable-ramp/types"
)

// BeginBlocker processes pending stable-ramp packages that have reached their
// approval threshold.
func BeginBlocker(ctx context.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return k.PendingPackages.Walk(ctx, nil, func(key collections.Pair[uint64, string], trackedPackage types.TrackedPackage) (bool, error) {
		pkg := trackedPackage.PackageData
		if pkg == nil || pkg.Status != types.PackageStatus_PACKAGE_STATUS_PENDING {
			return false, nil
		}

		committee, err := k.Committees.Get(ctx, pkg.ConnectionId)
		if err != nil {
			return true, err
		}
		if !packageThresholdReached(committee, trackedPackage.Approvals, trackedPackage.TotalMembers) {
			return false, nil
		}

		switch pkg.Action {
		case types.PackageAction_PACKAGE_ACTION_DEPOSIT:
			receiver, err := sdk.AccAddressFromBech32(pkg.ReceiverAddr)
			if err != nil {
				return true, err
			}
			depositCoin := sdk.NewCoin(pkg.Denom, math.NewIntFromUint64(pkg.Amount))
			if err := k.FulfillDeposit(ctx, receiver, sdk.Coins{depositCoin}); err != nil {
				return true, err
			}
			if err := k.PendingPackages.Remove(ctx, key); err != nil {
				return true, err
			}
			sdkCtx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeDepositClaimConfirm,
					sdk.NewAttribute(types.AttributeKeyReceiver, pkg.ReceiverAddr),
					sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatUint(pkg.Amount, 10)),
					sdk.NewAttribute(types.AttributeKeyDenom, pkg.Denom),
					sdk.NewAttribute(types.AttributeKeyConnection, pkg.ConnectionId),
					sdk.NewAttribute(types.AttributeKeyNonce, strconv.FormatUint(pkg.Nonce, 10)),
					sdk.NewAttribute(types.AttributeKeyApprovals, strconv.FormatUint(uint64(trackedPackage.Approvals), 10)),
				),
			)
		case types.PackageAction_PACKAGE_ACTION_WITHDRAW:
			burnCoin := sdk.NewCoin(pkg.Denom, math.NewIntFromUint64(pkg.Amount))
			if err := k.BurnEscrowedCoins(ctx, sdk.Coins{burnCoin}); err != nil {
				return true, err
			}
			trackedPackageBz, err := types.ModuleCdc.Marshal(&trackedPackage)
			if err != nil {
				return true, err
			}
			if err := k.PendingPackages.Remove(ctx, key); err != nil {
				return true, err
			}
			sdkCtx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeWithdrawClaimConfirm,
					sdk.NewAttribute(types.AttributeKeySender, pkg.SenderAddr),
					sdk.NewAttribute(types.AttributeKeyReceiver, pkg.ReceiverAddr),
					sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatUint(pkg.Amount, 10)),
					sdk.NewAttribute(types.AttributeKeyDenom, pkg.Denom),
					sdk.NewAttribute(types.AttributeKeyConnection, pkg.ConnectionId),
					sdk.NewAttribute(types.AttributeKeyNonce, strconv.FormatUint(pkg.Nonce, 10)),
					sdk.NewAttribute(types.AttributeKeyApprovals, strconv.FormatUint(uint64(trackedPackage.Approvals), 10)),
					sdk.NewAttribute(types.AttributeKeyPackage, hex.EncodeToString(trackedPackageBz)),
				),
			)
		}

		return false, nil
	})
}

func packageThresholdReached(committee types.Committee, approvals, totalMembers uint32) bool {
	threshold := committee.Threshold
	if threshold == nil || threshold.Denominator == 0 || totalMembers == 0 {
		return false
	}

	return uint64(approvals)*uint64(threshold.Denominator) >= uint64(totalMembers)*uint64(threshold.Numerator)
}
