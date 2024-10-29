package app

import (
	"testing"
	"time"

	testifysuite "github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"strconv"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
)

var (
	globalStartTime = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
)

type TransferTestSuite struct {
	testifysuite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func TestTransferTestSuite(t *testing.T) {
	testifysuite.Run(t, new(TransferTestSuite))
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(GetChainID(2))
}

func NewTransferPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	pathAtoB := ibctesting.NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(pathAtoB)

	timeoutHeight := clienttypes.NewHeight(1, 110)

	amount, ok := sdkmath.NewIntFromString("9223372036854775808") // 2^63 (one above int64)
	suite.Require().True(ok)
	coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, amount)

	// send from chainA to chainB
	msg := types.NewMsgTransfer(pathAtoB.EndpointA.ChannelConfig.PortID, pathAtoB.EndpointA.ChannelID, coinToSendToB, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
	res, err := suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err) // message committed

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	suite.Require().NoError(err)

	err = pathAtoB.RelayPacket(packet)
	suite.Require().NoError(err) // relay committed

	// check that module account escrow address has locked the tokens
	escrowAddress := types.GetEscrowAddress(packet.GetSourcePort(), packet.GetSourceChannel())
	balance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), escrowAddress, sdk.DefaultBondDenom)
	suite.Require().Equal(coinToSendToB, balance)
}

func NewCoordinator(t *testing.T, n int) *ibctesting.Coordinator {
	t.Helper()
	chains := make(map[string]*ibctesting.TestChain)
	coord := &ibctesting.Coordinator{
		T:           t,
		CurrentTime: globalStartTime,
	}

	for i := 1; i <= n; i++ {
		chainID := GetChainID(i)
		chains[chainID] = NewTestChain(t, coord, chainID)
	}
	coord.Chains = chains

	return coord
}

func NewTestChain(t *testing.T, coord *ibctesting.Coordinator, chainID string) *ibctesting.TestChain {
	t.Helper()
	var (
		validatorsPerChain = 4
	)

	valSet, signersByAddress := GenValSetForTestingApp(validatorsPerChain)

	return NewTestChainWithValSet(t, coord, chainID, valSet, signersByAddress)
}

func NewTestChainWithValSet(tb testing.TB, coord *ibctesting.Coordinator, chainID string, valSet *tmtypes.ValidatorSet, signers map[string]tmtypes.PrivValidator) *ibctesting.TestChain {
	tb.Helper()
	app, senderAccs := SetupForTestingApp(chainID, 4, valSet)

	// create current header and call begin block
	header := tmproto.Header{
		ChainID: chainID,
		Height:  1,
		Time:    coord.CurrentTime.UTC(),
	}

	txConfig := app.GetTxConfig()

	// create an account to send transactions from
	chain := &ibctesting.TestChain{
		TB:             tb,
		Coordinator:    coord,
		ChainID:        chainID,
		App:            app,
		CurrentHeader:  header,
		QueryServer:    app.GetIBCKeeper(),
		TxConfig:       txConfig,
		Codec:          app.AppCodec(),
		Vals:           valSet,
		NextVals:       valSet,
		Signers:        signers,
		SenderPrivKey:  senderAccs[0].SenderPrivKey,
		SenderAccount:  senderAccs[0].SenderAccount,
		SenderAccounts: senderAccs,
	}

	// commit genesis block
	chain.NextBlock()

	return chain
}

func GetChainID(index int) string {
	return "realionetworklocal_7777-" + strconv.Itoa(index)
}
