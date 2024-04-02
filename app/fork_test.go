package app

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var oneEnternityLater = time.Date(9999, 9, 9, 9, 9, 9, 9, time.UTC)

func TestFork(t *testing.T) {
	realio := Setup(false, nil)

	ctx := realio.BaseApp.NewContext(false, tmproto.Header{Height: int64(ForkHeight)})
	stakingKeeper := realio.StakingKeeper

	timeKey := time.Date(2024, 4, 1, 1, 1, 1, 1, time.UTC)

	duplicativeUnbondingDelegation := stakingtypes.UnbondingDelegation{
		DelegatorAddress: "test_del_1",
		ValidatorAddress: "test_val_1",
		Entries: []stakingtypes.UnbondingDelegationEntry{
			stakingtypes.NewUnbondingDelegationEntry(int64(ForkHeight), timeKey, math.OneInt()),
		},
	}

	stakingKeeper.InsertUBDQueue(ctx, duplicativeUnbondingDelegation, timeKey)
	stakingKeeper.InsertUBDQueue(ctx, duplicativeUnbondingDelegation, timeKey)

	duplicativeRedelegation := stakingtypes.Redelegation{
		DelegatorAddress:    "test_del_1",
		ValidatorSrcAddress: "test_val_1",
		ValidatorDstAddress: "test_val_2",
		Entries: []stakingtypes.RedelegationEntry{
			stakingtypes.NewRedelegationEntry(int64(ForkHeight), timeKey, math.OneInt(), sdk.OneDec()),
		},
	}
	stakingKeeper.InsertRedelegationQueue(ctx, duplicativeRedelegation, timeKey)
	stakingKeeper.InsertRedelegationQueue(ctx, duplicativeRedelegation, timeKey)
	stakingKeeper.InsertRedelegationQueue(ctx, duplicativeRedelegation, timeKey)

	duplicativeVal := stakingtypes.Validator{
		OperatorAddress: "test_op",
		UnbondingHeight: int64(ForkHeight),
		UnbondingTime:   timeKey,
	}

	stakingKeeper.InsertUnbondingValidatorQueue(ctx, duplicativeVal)
	stakingKeeper.InsertUnbondingValidatorQueue(ctx, duplicativeVal)

	require.True(t, checkDuplicateUBDQueue(ctx, *realio))
	require.True(t, checkDuplicateRelegationQueue(ctx, *realio))
	require.True(t, checkDuplicateValQueue(ctx, *realio))

	realio.ScheduleForkUpgrade(ctx)

	require.False(t, checkDuplicateUBDQueue(ctx, *realio))
	require.False(t, checkDuplicateRelegationQueue(ctx, *realio))
	require.False(t, checkDuplicateValQueue(ctx, *realio))

	dvPairs := stakingKeeper.GetUBDQueueTimeSlice(ctx, timeKey)
	require.Equal(t, dvPairs[0].DelegatorAddress, duplicativeUnbondingDelegation.DelegatorAddress)
	require.Equal(t, dvPairs[0].ValidatorAddress, duplicativeUnbondingDelegation.ValidatorAddress)

	triplets := stakingKeeper.GetRedelegationQueueTimeSlice(ctx, timeKey)
	require.Equal(t, triplets[0].DelegatorAddress, duplicativeRedelegation.DelegatorAddress)
	require.Equal(t, triplets[0].ValidatorDstAddress, duplicativeRedelegation.ValidatorDstAddress)
	require.Equal(t, triplets[0].ValidatorSrcAddress, duplicativeRedelegation.ValidatorSrcAddress)

	vals := stakingKeeper.GetUnbondingValidators(ctx, timeKey, int64(ForkHeight))
	require.Equal(t, vals[0], duplicativeVal.OperatorAddress)

}

func checkDuplicateUBDQueue(ctx sdk.Context, realio RealioNetwork) bool {
	ubdIter := realio.StakingKeeper.UBDQueueIterator(ctx, oneEnternityLater)
	defer ubdIter.Close()

	for ; ubdIter.Valid(); ubdIter.Next() {
		timeslice := stakingtypes.DVPairs{}
		value := ubdIter.Value()
		realio.appCodec.MustUnmarshal(value, &timeslice)
		if checkDuplicateUBD(timeslice.Pairs) {
			return true
		}
	}
	return false
}

func checkDuplicateUBD(eles []stakingtypes.DVPair) bool {
	unique_eles := map[string]bool{}
	for _, ele := range eles {
		unique_eles[ele.String()] = true
	}
	fmt.Println(eles, "eles")

	return len(unique_eles) != len(eles)
}

func checkDuplicateRelegationQueue(ctx sdk.Context, realio RealioNetwork) bool {
	redeIter := realio.StakingKeeper.RedelegationQueueIterator(ctx, oneEnternityLater)
	defer redeIter.Close()

	for ; redeIter.Valid(); redeIter.Next() {
		timeslice := stakingtypes.DVVTriplets{}
		value := redeIter.Value()
		realio.appCodec.MustUnmarshal(value, &timeslice)
		if checkDuplicateRedelegation(timeslice.Triplets) {
			return true
		}
	}
	return false
}

func checkDuplicateRedelegation(eles []stakingtypes.DVVTriplet) bool {
	unique_eles := map[string]bool{}
	for _, ele := range eles {
		unique_eles[ele.String()] = true
	}

	return len(unique_eles) != len(eles)
}

func checkDuplicateValQueue(ctx sdk.Context, realio RealioNetwork) bool {
	valsIter := realio.StakingKeeper.ValidatorQueueIterator(ctx, oneEnternityLater, 9999)
	defer valsIter.Close()

	for ; valsIter.Valid(); valsIter.Next() {
		timeslice := stakingtypes.ValAddresses{}
		value := valsIter.Value()
		realio.appCodec.MustUnmarshal(value, &timeslice)
		if checkDuplicateValAddr(timeslice.Addresses) {
			return true
		}
	}
	return false
}
func checkDuplicateValAddr(eles []string) bool {
	unique_eles := map[string]bool{}
	for _, ele := range eles {
		unique_eles[ele] = true
	}

	return len(unique_eles) != len(eles)
}