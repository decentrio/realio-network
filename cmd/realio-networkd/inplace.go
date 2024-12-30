package main

import (
	// "encoding/base64"
	// "bytes"
	"fmt"
	"io"
	"time"

	// "path/filepath"
	"strings"

	// "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	// tmd25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	tdmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmos "github.com/cometbft/cometbft/libs/os"

	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	// "github.com/cosmos/cosmos-sdk/snapshots"
	// snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	"cosmossdk.io/simapp/params"
	"cosmossdk.io/store"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	//"github.com/cosmos/cosmos-sdk/x/gov/types"
	// tmtypes "github.com/cometbft/cometbft/types"

	"github.com/realiotech/realio-network/app"
	"github.com/realiotech/realio-network/types"
)

var encoding params.EncodingConfig

const (
	valVotingPower int64 = 9000000000000000000
)

var (
	// flagValidatorPrivKey = "validator-privkey"
	flagAccountsToFund = "accounts-to-fund"
)

type valArgs struct {
	newValAddr         tdmbytes.HexBytes
	newOperatorAddress string
	newValPubKey       crypto.PubKey
	// validatorConsPrivKey crypto.PrivKey
	accountsToFund   []sdk.AccAddress
	upgradeToTrigger string
	homeDir          string
}

func NewInplaceCmd(encodingConfig params.EncodingConfig) *cobra.Command {
	encoding = encodingConfig
	cmd := server.InPlaceTestnetCreator(newTestnetApp)
	cmd.Use = "testnet-inplace [newChainID] [newOperatorAddress]"
	cmd.Short = "Updates chain's application and consensus state with provided validator info and starts the node"
	cmd.Long = `The test command modifies both application and consensus stores within a local mainnet node and starts the node,
with the aim of facilitating testing procedures. This command replaces existing validator data with updated information,
thereby removing the old validator set and introducing a new set suitable for local testing purposes. By altering the state extracted from the mainnet node,
it enables developers to configure their local environments to reflect mainnet conditions more accurately.
Example:
	appd testnet chainID-1 cosmosvaloper1w7f3xx7e75p4l7qdym5msqem9rd4dyc4mq79dm --home $HOME/.appd/validator1 --validator-privkey=6dq+/KHNvyiw2TToCgOpUpQKIzrLs69Rb8Az39xvmxPHNoPxY1Cil8FY+4DhT9YwD6s0tFABMlLcpaylzKKBOg== --accounts-to-fund="cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq,cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x" [other_server_start_flags]
	`

	cmd.Example = `appd testnet chainID-1 cosmosvaloper1w7f3xx7e75p4l7qdym5msqem9rd4dyc4mq79dm --home $HOME/.appd/validator1 --validator-privkey=6dq+/KHNvyiw2TToCgOpUpQKIzrLs69Rb8Az39xvmxPHNoPxY1Cil8FY+4DhT9YwD6s0tFABMlLcpaylzKKBOg== --accounts-to-fund="cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq,cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x"`

	// cmd.Flags().String(flagValidatorPrivKey, "", "Validator cometbft/PrivKeyEd25519 consensus private key from the priv_validato_key.json file")
	cmd.Flags().String(flagAccountsToFund, "", "Comma-separated list of account addresses that will be funded for testing purposes")
	return cmd
}

// newTestnetApp starts by running the normal newApp method. From there, the app interface returned is modified in order
// for a testnet to be created from the provided app.
func newTestnetApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	// Create an app and type cast to an App
	newApp := newApp(logger, db, traceStore, encoding, appOpts)
	testApp, ok := newApp.(*app.RealioNetwork)
	if !ok {
		panic("app created from newApp is not of type App")
	}

	// Get command args
	args, err := getCommandArgs(appOpts)
	if err != nil {
		panic(err)
	}

	return initAppForTestnet(testApp, args)
}

// InitAppForTestnet is broken down into two sections:
// Required Changes: Changes that, if not made, will cause the testnet to halt or panic
// Optional Changes: Changes to customize the testnet to one's liking (lower vote times, fund accounts, etc)
func initAppForTestnet(app *app.RealioNetwork, args valArgs) *app.RealioNetwork {
	//
	// Required Changes:
	//
	ctx := app.BaseApp.NewUncachedContext(true, cmtproto.Header{})

	pubkey := &ed25519.PubKey{Key: args.newValPubKey.Bytes()}
	pubkeyAny, err := codectypes.NewAnyWithValue(pubkey)
	if err != nil {
		tmos.Exit(err.Error())
	}

	// STAKING
	//

	// Create Validator struct for our new validator.
	newVal := stakingtypes.Validator{
		OperatorAddress: args.newOperatorAddress,
		ConsensusPubkey: pubkeyAny,
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(valVotingPower),
		DelegatorShares: math.LegacyMustNewDecFromStr("100000000000000000"),
		Description: stakingtypes.Description{
			Moniker: "Testnet Validator",
		},
		Commission: stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.05"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
			},
		},
		MinSelfDelegation: math.OneInt(),
	}

	app.SlashingKeeper.AddPubkey(ctx, pubkey)

	// Remove all validators from power store
	stakingKey := app.GetKey(stakingtypes.ModuleName)
	stakingStore := ctx.KVStore(stakingKey)
	iterator, _ := app.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	// powSkip := app.StakingKeeper.GetLastValidatorPower(ctx, valSkip)
	// Remove all valdiators from last validators store
	iterator, _ = app.StakingKeeper.LastValidatorsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		stakingStore.Delete(key)
	}
	iterator.Close()
	// app.StakingKeeper.SetLastValidatorPower(ctx, valSkip, powSkip)

	// Remove all validators from validators store

	iterator = stakingStore.Iterator(stakingtypes.ValidatorsKey, storetypes.PrefixEndBytes(stakingtypes.ValidatorsKey))
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		stakingStore.Delete(key)
	}
	iterator.Close()

	timestamp := func(key []byte) time.Time {
		bz := key[len(stakingtypes.UnbondingQueueKey):]
		timestamp, err := sdk.ParseTimeBytes(bz)
		if err != nil {
			panic(err)
		}
		return timestamp
	}

	// Remove all validators from unbonding queue
	iterator = stakingStore.Iterator(stakingtypes.UnbondingQueueKey, storetypes.PrefixEndBytes(stakingtypes.UnbondingQueueKey))
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if timestamp(key).After(time.Now()) {
			continue
		} else {
			stakingStore.Delete(key)
		}
	}
	iterator.Close()

	// Add our validator to power and last validators store
	app.StakingKeeper.SetValidator(ctx, newVal)
	err = app.StakingKeeper.SetValidatorByConsAddr(ctx, newVal)
	if err != nil {
		tmos.Exit(err.Error())
	}
	app.StakingKeeper.SetValidatorByPowerIndex(ctx, newVal)
	newValAddr, _ := app.StakingKeeper.ValidatorAddressCodec().StringToBytes(newVal.GetOperator())
	app.StakingKeeper.SetLastValidatorPower(ctx, newValAddr, 1000000000000000000)

	paramStaking, _ := app.StakingKeeper.GetParams(ctx)
	paramStaking.UnbondingTime = 86400 * time.Second
	app.StakingKeeper.SetParams(ctx, paramStaking)

	// DISTRIBUTION
	//

	// Initialize records for this validator across all distribution stores
	app.DistrKeeper.SetValidatorHistoricalRewards(ctx, newValAddr, 0, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))
	app.DistrKeeper.SetValidatorCurrentRewards(ctx, newValAddr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))
	app.DistrKeeper.SetValidatorAccumulatedCommission(ctx, newValAddr, distrtypes.InitialValidatorAccumulatedCommission())
	app.DistrKeeper.SetValidatorOutstandingRewards(ctx, newValAddr, distrtypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})

	// SLASHING
	//

	// Set validator signing info for our new validator.
	newConsAddr := sdk.ConsAddress(args.newValAddr.Bytes())
	newValidatorSigningInfo := slashingtypes.ValidatorSigningInfo{
		Address:     newConsAddr.String(),
		StartHeight: app.LastBlockHeight() - 1,
		Tombstoned:  false,
	}
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newValidatorSigningInfo)

	// // MINT
	// p := app.MintKeeper.GetMinter(ctx)
	// p.AnnualProvisions = sdk.MustNewDecFromStr("47103769757652020228568450937") //.NewInFromString("75000000000000000000000000")
	// p.Inflation = sdk.MustNewDecFromStr("1")
	// app.MintKeeper.SetMinter(ctx, p)

	//
	// Optional Changes:
	//

	// BANK
	//
	bondDenom, _ := app.StakingKeeper.BondDenom(ctx)

	amount, ok := math.NewIntFromString("17103769757652020228568450937")
	if !ok {
		amount = math.NewInt(1000000000000000000)
	}
	defaultCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000000000000000000)), sdk.NewCoin("ario", amount))

	// Fund local accounts
	for _, account := range args.accountsToFund {
		err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, defaultCoins)
		if err != nil {
			tmos.Exit(err.Error())
		}
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, account, defaultCoins)
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	//
	// Optional Changes:
	//

	// GOV
	//1 minute
	govParams, _ := app.GovKeeper.Params.Get(ctx)
	minute := time.Minute
	govParams.VotingPeriod = &minute

	// set deposit 1h and
	minDeposit := math.NewInt(100000000000)
	minDepositCoin := sdk.NewCoins(sdk.NewCoin(bondDenom, minDeposit))

	govParams.MinDeposit = minDepositCoin
	err = app.GovKeeper.Params.Set(ctx, govParams)
	if err != nil {
		tmos.Exit(err.Error())
	}
	return app
}

// parse the input flags and returns valArgs
func getCommandArgs(appOpts servertypes.AppOptions) (valArgs, error) {
	args := valArgs{}

	newValAddr, ok := appOpts.Get(server.KeyNewValAddr).(tdmbytes.HexBytes)
	if !ok {
		panic("newValAddr is not of type bytes.HexBytes")
	}
	args.newValAddr = newValAddr
	newValPubKey, ok := appOpts.Get(server.KeyUserPubKey).(crypto.PubKey)
	if !ok {
		panic("newValPubKey is not of type crypto.PubKey")
	}
	args.newValPubKey = newValPubKey
	newOperatorAddress, ok := appOpts.Get(server.KeyNewOpAddr).(string)
	if !ok {
		panic("newOperatorAddress is not of type string")
	}
	args.newOperatorAddress = newOperatorAddress
	upgradeToTrigger, ok := appOpts.Get(server.KeyTriggerTestnetUpgrade).(string)
	if !ok {
		panic("upgradeToTrigger is not of type string")
	}
	args.upgradeToTrigger = upgradeToTrigger

	// validate  and set validator privkey
	// validatorPrivKey := cast.ToString(appOpts.Get(flagValidatorPrivKey))
	// if validatorPrivKey == "" {
	// 	return args, fmt.Errorf("invalid validator private key")
	// }
	// decPrivKey, err := base64.StdEncoding.DecodeString(validatorPrivKey)
	// if err != nil {
	// 	return args, fmt.Errorf("cannot decode validator private key %w", err)
	// }
	// args.validatorConsPrivKey = tmd25519.PrivKey([]byte(decPrivKey))

	// validate  and set accounts to fund
	accountsString := cast.ToString(appOpts.Get(flagAccountsToFund))

	for _, account := range strings.Split(accountsString, ",") {
		if account != "" {
			addr, err := sdk.AccAddressFromBech32(account)
			if err != nil {
				return args, fmt.Errorf("invalid bech32 address format %w", err)
			}
			args.accountsToFund = append(args.accountsToFund, addr)
		}
	}

	// home dir
	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	if homeDir == "" {
		return args, fmt.Errorf("invalid home dir")
	}
	args.homeDir = homeDir

	return args, nil
}

// newApp creates the application
func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	encodingConfig params.EncodingConfig,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	var cache storetypes.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	// homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	// chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	// if chainID == "" {
	// 	// fallback to genesis chain-id
	// 	appGenesis, err := tmtypes.GenesisDocFromFile(filepath.Join(homeDir, "config", "genesis.json"))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	chainID = appGenesis.ChainID
	// }

	baseappOptions := []func(*baseapp.BaseApp){
		baseapp.SetPruning(pruningOpts),
		baseapp.SetChainID(types.MainnetChainID + "-1"),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(server.FlagIAVLCacheSize))),
		baseapp.SetIAVLDisableFastNode(cast.ToBool(appOpts.Get(server.FlagDisableIAVLFastNode))),
		// baseapp.SetChainID(chainID),
	}

	// If this is an in place testnet, set any new stores that may exist

	return app.New(
		logger, db, traceStore, true,
		map[int64]bool{},
		app.DefaultNodeHome,
		0,
		encodingConfig,
		appOpts,
		baseappOptions...,
	)
}
