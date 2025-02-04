package erc20

import (
	"fmt"
	"math/big"
	"slices"
	"time"

	errorsmod "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/os/precompiles/erc20"
	commonfactory "github.com/evmos/os/testutil/integration/common/factory"
	commonnetwork "github.com/evmos/os/testutil/integration/common/network"
	"github.com/evmos/os/testutil/integration/os/factory"
	testutils "github.com/evmos/os/testutil/integration/os/utils"
	erc20types "github.com/evmos/os/x/erc20/types"
	network "github.com/realiotech/realio-network/testutil/integration/os/network"
	utiltx "github.com/realiotech/realio-network/testutil/tx"
)

// setupSendAuthz is a helper function to set up a SendAuthorization for
// a given grantee and granter combination for a given amount.
//
// NOTE: A default expiration of 1 hour after the current block time is used.
func (s *PrecompileTestSuite) setupSendAuthz(
	grantee sdk.AccAddress, granterPriv cryptotypes.PrivKey, amount sdk.Coins,
) {
	err := setupSendAuthz(
		s.network,
		s.factory,
		grantee,
		granterPriv,
		amount,
	)
	s.Require().NoError(err, "failed to set up send authorization")
}

func setupSendAuthz(
	network commonnetwork.Network,
	factory commonfactory.BaseTxFactory,
	grantee sdk.AccAddress,
	granterPriv cryptotypes.PrivKey,
	amount sdk.Coins,
) error {
	granter := sdk.AccAddress(granterPriv.PubKey().Address())
	expiration := network.GetContext().BlockHeader().Time.Add(time.Hour)
	sendAuthz := banktypes.NewSendAuthorization(
		amount,
		[]sdk.AccAddress{},
	)

	msgGrant, err := authz.NewMsgGrant(
		granter,
		grantee,
		sendAuthz,
		&expiration,
	)
	if err != nil {
		return errorsmod.Wrap(err, "failed to create MsgGrant")
	}

	// Create an authorization
	txArgs := commonfactory.CosmosTxArgs{Msgs: []sdk.Msg{msgGrant}}
	_, err = factory.ExecuteCosmosTx(granterPriv, txArgs)
	if err != nil {
		return errorsmod.Wrap(err, "failed to execute MsgGrant")
	}

	return nil
}

// requireOut is a helper utility to reduce the amount of boilerplate code in the query tests.
//
// It requires the output bytes and error to match the expected values. Additionally, the method outputs
// are unpacked and the first value is compared to the expected value.
//
// NOTE: It's sufficient to only check the first value because all methods in the ERC20 precompile only
// return a single value.
func (s *PrecompileTestSuite) requireOut(
	bz []byte,
	err error,
	method abi.Method,
	expPass bool,
	errContains string,
	expValue interface{},
) {
	if expPass {
		s.Require().NoError(err, "expected no error")
		s.Require().NotEmpty(bz, "expected bytes not to be empty")

		// Unpack the name into a string
		out, err := method.Outputs.Unpack(bz)
		s.Require().NoError(err, "expected no error unpacking")

		// Check if expValue is a big.Int. Because of a difference in uninitialized/empty values for big.Ints,
		// this comparison is often not working as expected, so we convert to Int64 here and compare those values.
		bigExp, ok := expValue.(*big.Int)
		if ok {
			bigOut, ok := out[0].(*big.Int)
			s.Require().True(ok, "expected output to be a big.Int")
			s.Require().Equal(bigExp.Int64(), bigOut.Int64(), "expected different value")
		} else {
			s.Require().Equal(expValue, out[0], "expected different value")
		}
	} else {
		s.Require().Error(err, "expected error")
		s.Require().Contains(err.Error(), errContains, "expected different error")
	}
}

// requireSendAuthz is a helper function to check that a SendAuthorization
// exists for a given grantee and granter combination for a given amount.
//
// NOTE: This helper expects only one authorization to exist.
func (s *PrecompileTestSuite) requireSendAuthz(grantee, granter sdk.AccAddress, amount sdk.Coins, allowList []string) {
	grants, err := s.grpcHandler.GetGrantsByGrantee(grantee.String())
	s.Require().NoError(err, "expected no error querying the grants")
	s.Require().Len(grants, 1, "expected one grant")
	s.Require().Equal(grantee.String(), grants[0].Grantee, "expected different grantee")
	s.Require().Equal(granter.String(), grants[0].Granter, "expected different granter")

	authzs, err := s.grpcHandler.GetAuthorizationsByGrantee(grantee.String())
	s.Require().NoError(err, "expected no error unpacking the authorization")
	s.Require().Len(authzs, 1, "expected one authorization")

	sendAuthz, ok := authzs[0].(*banktypes.SendAuthorization)
	s.Require().True(ok, "expected send authorization")

	s.Require().Equal(amount, sendAuthz.SpendLimit, "expected different spend limit amount")
	if len(allowList) == 0 {
		s.Require().Empty(sendAuthz.AllowList, "expected empty allow list")
	} else {
		s.Require().Equal(allowList, sendAuthz.AllowList, "expected different allow list")
	}
}

// setupERC20Precompile is a helper function to set up an instance of the ERC20 precompile for
// a given token denomination, set the token pair in the ERC20 keeper and adds the precompile
// to the available and active precompiles.
func (s *PrecompileTestSuite) setupERC20Precompile(denom string) *Precompile {
	tokenPair := erc20types.NewTokenPair(utiltx.GenerateAddress(), denom, erc20types.OWNER_MODULE)
	s.network.App.Erc20Keeper.SetTokenPair(s.network.GetContext(), tokenPair)

	precompile, err := setupERC20PrecompileForTokenPair(*s.network, tokenPair)
	s.Require().NoError(err, "failed to set up %q erc20 precompile", tokenPair.Denom)

	return precompile
}

// setupERC20PrecompileForTokenPair is a helper function to set up an instance of the ERC20 precompile for
// a given token pair and adds the precompile to the available and active precompiles.
// Do not use this function for integration tests.
func setupERC20PrecompileForTokenPair(
	unitNetwork network.UnitTestNetwork, tokenPair erc20types.TokenPair,
) (*Precompile, error) {
	precompile, err := NewPrecompile(
		tokenPair.Denom,
		tokenPair.GetERC20Contract(),
		unitNetwork.App.BankKeeper,
		unitNetwork.App.AuthzKeeper,
		unitNetwork.App.TransferKeeper,
		unitNetwork.App.AssetKeeper,
	)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to create %q erc20 precompile", tokenPair.Denom)
	}

	err = unitNetwork.App.Erc20Keeper.EnableDynamicPrecompiles(
		unitNetwork.GetContext(),
		precompile.Address(),
	)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to add %q erc20 precompile to EVM extensions", tokenPair.Denom)
	}

	return precompile, nil
}

// setupNewERC20PrecompileForTokenPair is a helper function to set up an instance of the ERC20 precompile for
// a given token pair and adds the precompile to the available and active precompiles.
// This function should be used for integration tests
func setupNewERC20PrecompileForTokenPair(
	privKey cryptotypes.PrivKey,
	unitNetwork *network.UnitTestNetwork,
	tf factory.TxFactory, tokenPair erc20types.TokenPair,
) (*erc20.Precompile, error) {
	precompile, err := erc20.NewPrecompile(
		tokenPair,
		unitNetwork.App.BankKeeper,
		unitNetwork.App.AuthzKeeper,
		unitNetwork.App.TransferKeeper,
	)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to create %q erc20 precompile", tokenPair.Denom)
	}

	// Update the params via gov proposal
	params := unitNetwork.App.Erc20Keeper.GetParams(unitNetwork.GetContext())
	params.DynamicPrecompiles = append(params.DynamicPrecompiles, precompile.Address().Hex())
	slices.Sort(params.DynamicPrecompiles)

	if err := params.Validate(); err != nil {
		return nil, err
	}

	if err := testutils.UpdateERC20Params(testutils.UpdateParamsInput{
		Pk:      privKey,
		Tf:      tf,
		Network: unitNetwork,
		Params:  params,
	}); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to add %q erc20 precompile to EVM extensions", tokenPair.Denom)
	}

	return precompile, nil
}

// CallType indicates which type of contract call is made during the integration tests.
type CallType int

// callType constants to differentiate between direct calls and calls through a contract.
const (
	directCall CallType = iota + 1
	directCallToken2
	contractCall
	contractCallToken2
	erc20Call
	erc20CallerCall
	erc20V5Call
	erc20V5CallerCall
)

var (
	nativeCallTypes = []CallType{directCall, directCallToken2, contractCall, contractCallToken2}
	erc20CallTypes  = []CallType{erc20Call, erc20CallerCall, erc20V5Call, erc20V5CallerCall}
)

// ExpectedBalance is a helper struct to check the balances of accounts.
type ExpectedBalance struct {
	address  sdk.AccAddress
	expCoins sdk.Coins
}

// ContractsData is a helper struct to hold the addresses and ABIs for the
// different contract instances that are subject to testing here.
type ContractsData struct {
	contractData map[CallType]ContractData
	ownerPriv    cryptotypes.PrivKey
}

// ContractData is a helper struct to hold the address and ABI for a given contract.
type ContractData struct {
	Address common.Address
	ABI     abi.ABI
}

// GetContractData is a helper function to return the contract data for a given call type.
func (cd ContractsData) GetContractData(callType CallType) ContractData {
	data, found := cd.contractData[callType]
	if !found {
		panic(fmt.Sprintf("no contract data found for call type: %d", callType))
	}
	return data
}
