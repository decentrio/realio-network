package erc20

import (
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/evmos/os/precompiles/testutil"
	"github.com/evmos/os/testutil/integration/os/keyring"
	erc20types "github.com/evmos/os/x/erc20/types"
	"github.com/evmos/os/x/evm/core/vm"
	"github.com/evmos/os/x/evm/statedb"
	utiltx "github.com/realiotech/realio-network/testutil/tx"
	assettypes "github.com/realiotech/realio-network/x/asset/types"
)

var (
	tokenDenom = "xmpl"
	// XMPLCoin is a dummy coin used for testing purposes.
	XMPLCoin = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 1e18))
	// toAddr is a dummy address used for testing purposes.
	toAddr = utiltx.GenerateAddress()
	// toAddr is a dummy address used for testing purposes.
	fromAddr = utiltx.GenerateAddress()
)

func (s *PrecompileTestSuite) TestTransfer() {
	method := s.precompile.Methods[TransferMethod]
	// fromAddr is the address of the keyring account used for testing.
	fromAddr := s.keyring.GetKey(0).Addr
	testcases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func()
		expErr      bool
		errContains string
	}{
		{
			"fail - negative amount",
			func() []interface{} {
				return []interface{}{toAddr, big.NewInt(-1)}
			},
			func() {},
			true,
			"coin -1xmpl amount is not positive",
		},
		{
			"fail - invalid to address",
			func() []interface{} {
				return []interface{}{"", big.NewInt(100)}
			},
			func() {},
			true,
			"invalid to address",
		},
		{
			"fail - invalid amount",
			func() []interface{} {
				return []interface{}{toAddr, ""}
			},
			func() {},
			true,
			"invalid amount",
		},
		{
			"fail - not enough balance",
			func() []interface{} {
				return []interface{}{toAddr, big.NewInt(2e18)}
			},
			func() {},
			true,
			ErrTransferAmountExceedsBalance.Error(),
		},
		{
			"pass",
			func() []interface{} {
				return []interface{}{toAddr, big.NewInt(100)}
			},
			func() {
				toAddrBalance := s.network.App.BankKeeper.GetBalance(s.network.GetContext(), toAddr.Bytes(), tokenDenom)
				s.Require().Equal(big.NewInt(100), toAddrBalance.Amount.BigInt(), "expected toAddr to have 100 XMPL")
			},
			false,
			"",
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()
			stateDB := s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx := testutil.NewPrecompileContract(s.T(), s.network.GetContext(), fromAddr, s.precompile, 0)

			// Mint some coins to the module account and then send to the from address
			err := s.network.App.BankKeeper.MintCoins(s.network.GetContext(), erc20types.ModuleName, XMPLCoin)
			s.Require().NoError(err, "failed to mint coins")
			err = s.network.App.BankKeeper.SendCoinsFromModuleToAccount(s.network.GetContext(), erc20types.ModuleName, fromAddr.Bytes(), XMPLCoin)
			s.Require().NoError(err, "failed to send coins from module to account")

			_, err = s.precompile.Transfer(ctx, contract, stateDB, &method, tc.malleate())
			if tc.expErr {
				s.Require().Error(err, "expected transfer transaction to fail")
				s.Require().Contains(err.Error(), tc.errContains, "expected transfer transaction to fail with specific error")
			} else {
				s.Require().NoError(err, "expected transfer transaction succeeded")
				tc.postCheck()
			}
		})
	}
}

func (s *PrecompileTestSuite) TestTransferFrom() {
	var (
		ctx  sdk.Context
		stDB *statedb.StateDB
	)
	method := s.precompile.Methods[TransferFromMethod]
	// owner of the tokens
	owner := s.keyring.GetKey(0)
	// spender of the tokens
	spender := s.keyring.GetKey(1)

	testcases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func()
		expErr      bool
		errContains string
	}{
		{
			"fail - negative amount",
			func() []interface{} {
				return []interface{}{owner.Addr, toAddr, big.NewInt(-1)}
			},
			func() {},
			true,
			"coin -1xmpl amount is not positive",
		},
		{
			"fail - invalid from address",
			func() []interface{} {
				return []interface{}{"", toAddr, big.NewInt(100)}
			},
			func() {},
			true,
			"invalid from address",
		},
		{
			"fail - invalid to address",
			func() []interface{} {
				return []interface{}{owner.Addr, "", big.NewInt(100)}
			},
			func() {},
			true,
			"invalid to address",
		},
		{
			"fail - invalid amount",
			func() []interface{} {
				return []interface{}{owner.Addr, toAddr, ""}
			},
			func() {},
			true,
			"invalid amount",
		},
		{
			"fail - not enough allowance",
			func() []interface{} {
				return []interface{}{owner.Addr, toAddr, big.NewInt(100)}
			},
			func() {},
			true,
			ErrInsufficientAllowance.Error(),
		},
		{
			"fail - not enough balance",
			func() []interface{} {
				expiration := time.Now().Add(time.Hour)
				err := s.network.App.AuthzKeeper.SaveGrant(
					ctx,
					spender.AccAddr,
					owner.AccAddr,
					&banktypes.SendAuthorization{SpendLimit: sdk.Coins{sdk.Coin{Denom: s.tokenDenom, Amount: math.NewInt(5e18)}}},
					&expiration,
				)
				s.Require().NoError(err, "failed to save grant")

				return []interface{}{owner.Addr, toAddr, big.NewInt(2e18)}
			},
			func() {},
			true,
			ErrTransferAmountExceedsBalance.Error(),
		},
		{
			"pass - spend on behalf of other account",
			func() []interface{} {
				expiration := time.Now().Add(time.Hour)
				err := s.network.App.AuthzKeeper.SaveGrant(
					ctx,
					spender.AccAddr,
					owner.AccAddr,
					&banktypes.SendAuthorization{SpendLimit: sdk.Coins{sdk.Coin{Denom: tokenDenom, Amount: math.NewInt(300)}}},
					&expiration,
				)
				s.Require().NoError(err, "failed to save grant")

				return []interface{}{owner.Addr, toAddr, big.NewInt(100)}
			},
			func() {
				toAddrBalance := s.network.App.BankKeeper.GetBalance(ctx, toAddr.Bytes(), tokenDenom)
				s.Require().Equal(big.NewInt(100), toAddrBalance.Amount.BigInt(), "expected toAddr to have 100 XMPL")
			},
			false,
			"",
		},
		{
			"pass - spend on behalf of own account",
			func() []interface{} {
				// Mint some coins to the module account and then send to the spender address
				err := s.network.App.BankKeeper.MintCoins(ctx, erc20types.ModuleName, XMPLCoin)
				s.Require().NoError(err, "failed to mint coins")
				err = s.network.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, spender.AccAddr, XMPLCoin)
				s.Require().NoError(err, "failed to send coins from module to account")

				// NOTE: no authorization is necessary to spend on behalf of the same account
				return []interface{}{spender.Addr, toAddr, big.NewInt(100)}
			},
			func() {
				toAddrBalance := s.network.App.BankKeeper.GetBalance(ctx, toAddr.Bytes(), tokenDenom)
				s.Require().Equal(big.NewInt(100), toAddrBalance.Amount.BigInt(), "expected toAddr to have 100 XMPL")
			},
			false,
			"",
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx = s.network.GetContext()
			stDB = s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, spender.Addr, s.precompile, 0)

			// Mint some coins to the module account and then send to the from address
			err := s.network.App.BankKeeper.MintCoins(ctx, erc20types.ModuleName, XMPLCoin)
			s.Require().NoError(err, "failed to mint coins")
			err = s.network.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, owner.AccAddr, XMPLCoin)
			s.Require().NoError(err, "failed to send coins from module to account")

			_, err = s.precompile.TransferFrom(ctx, contract, stDB, &method, tc.malleate())
			if tc.expErr {
				s.Require().Error(err, "expected transfer transaction to fail")
				s.Require().Contains(err.Error(), tc.errContains, "expected transfer transaction to fail with specific error")
			} else {
				s.Require().NoError(err, "expected transfer transaction succeeded")
				tc.postCheck()
			}
		})
	}
}

func (s *PrecompileTestSuite) TestMint() {
	method := s.precompile.Methods[MintMethod]
	// fromAddr is the address of the keyring account used for testing.
	sender := s.keyring.GetKey(0)
	invalidSender := s.keyring.GetKey(1)
	maxSupply := math.NewInt(200)
	testcases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func()
		expErr      bool
		errContains string
		sender      keyring.Key
	}{
		{
			"fail - negative amount",
			func() []interface{} {
				return []interface{}{toAddr, big.NewInt(-1)}
			},
			func() {},
			true,
			"coin -1xmpl amount is not positive",
			sender,
		},
		{
			"fail - invalid to address",
			func() []interface{} {
				return []interface{}{"", big.NewInt(100)}
			},
			func() {},
			true,
			"invalid to address",
			sender,
		},
		{
			"fail - invalid amount",
			func() []interface{} {
				return []interface{}{toAddr, ""}
			},
			func() {},
			true,
			"invalid amount",
			sender,
		},
		{
			"fail - sender is not manager",
			func() []interface{} {
				return []interface{}{toAddr, big.NewInt(2e18)}
			},
			func() {},
			true,
			ErrTransferAmountExceedsBalance.Error(),
			invalidSender,
		},
		{
			"fail - exceed max supply",
			func() []interface{} {
				return []interface{}{toAddr, big.NewInt(300)}
			},
			func() {},
			true,
			"Exceed max supply",
			invalidSender,
		},
		{
			"pass",
			func() []interface{} {
				return []interface{}{toAddr, big.NewInt(100)}
			},
			func() {
				toAddrBalance := s.network.App.BankKeeper.GetBalance(s.network.GetContext(), toAddr.Bytes(), tokenDenom)
				s.Require().Equal(big.NewInt(100), toAddrBalance.Amount.BigInt(), "expected toAddr to have 100 XMPL")
			},
			false,
			"",
			sender,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()
			stateDB := s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx := testutil.NewPrecompileContract(s.T(), s.network.GetContext(), tc.sender.Addr, s.precompile, 0)

			// Set up manager role for valid sender
			err := s.precompile.assetKeep.TokenManagement.Set(
				ctx, 
				s.tokenDenom, 
				assettypes.TokenManagement{
					Managers: []string{sender.AccAddr.String()},
					ExtensionsList: []string{"mint"},
					MaxSupply: maxSupply,
				},
			)
			s.Require().NoError(err)

			_, err = s.precompile.Mint(ctx, contract, stateDB, &method, tc.malleate())
			fmt.Println("errrrr", err)
			if tc.expErr {
				s.Require().Error(err, "expected mint transaction to fail")
				// s.Require().Contains(err.Error(), tc.errContains, "expected transfer transaction to fail with specific error")
			} else {
				s.Require().NoError(err, "expected transfer transaction succeeded")
				tc.postCheck()
			}
		})
	}
}

func (s *PrecompileTestSuite) TestBurn() {
	method := s.precompile.Methods[BurnMethod]
	// fromAddr is the address of the keyring account used for testing.
	sender := s.keyring.GetKey(0)
	invalidSender := s.keyring.GetKey(1)
	maxSupply := math.NewInt(200)
	testcases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func()
		expErr      bool
		errContains string
		sender      keyring.Key
	}{
		{
			"fail - negative amount",
			func() []interface{} {
				return []interface{}{big.NewInt(-1)}
			},
			func() {},
			true,
			"coin -1xmpl amount is not positive",
			sender,
		},
		{
			"fail - invalid amount",
			func() []interface{} {
				return []interface{}{""}
			},
			func() {},
			true,
			"invalid amount",
			sender,
		},
		{
			"fail - sender is not manager",
			func() []interface{} {
				return []interface{}{big.NewInt(100)}
			},
			func() {},
			true,
			"sender is not token manager",
			invalidSender,
		},
		{
			"fail - not enough balance",
			func() []interface{} {
				return []interface{}{big.NewInt(300)}
			},
			func() {},
			true,
			"not enough balance",
			sender,
		},
		{
			"pass",
			func() []interface{} {
				return []interface{}{big.NewInt(100)}
			},
			func() {
				addrBalance := s.network.App.BankKeeper.GetBalance(s.network.GetContext(), sender.AccAddr.Bytes(), tokenDenom)
				s.Require().Equal(big.NewInt(100), addrBalance.Amount.BigInt(), "expected toAddr to have 100 XMPL")
			},
			false,
			"",
			sender,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()
			stateDB := s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx := testutil.NewPrecompileContract(s.T(), s.network.GetContext(), tc.sender.Addr, s.precompile, 0)

			// Set up manager role for valid sender
			err := s.precompile.assetKeep.TokenManagement.Set(
				ctx, 
				s.tokenDenom, 
				assettypes.TokenManagement{
					Managers: []string{sender.AccAddr.String()},
					ExtensionsList: []string{"mint"},
					MaxSupply: maxSupply,
				},
			)
			s.Require().NoError(err)

			// Mint amount to fromAddr to burn lately
			// _, err = s.precompile.Mint(ctx, contract, stateDB, &method, []interface{}{sender.Addr, big.NewInt(maxSupply.Int64())})
			// fmt.Println("errrrr", err)
			// s.Require().NoError(err)

			// Mint some coins to the module account and then send to the from address
			err = s.network.App.BankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(sdk.NewCoin(tokenDenom, maxSupply)))
			s.Require().NoError(err, "failed to mint coins")
			err = s.network.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, sender.AccAddr, sdk.NewCoins(sdk.NewCoin(tokenDenom, maxSupply)))
			s.Require().NoError(err, "failed to send coins from module to account")

			_, err = s.precompile.Burn(ctx, contract, stateDB, &method, tc.malleate())
			if tc.expErr {
				s.Require().Error(err, "expected burn transaction to fail")
				// s.Require().Contains(err.Error(), tc.errContains, "expected transfer transaction to fail with specific error")
			} else {
				s.Require().NoError(err, "expected burn transaction succeeded")
				tc.postCheck()
			}
		})
	}
}

func (s *PrecompileTestSuite) TestBurnFrom() {
	method := s.precompile.Methods[BurnFromMethod]
	// fromAddr is the address of the keyring account used for testing.
	sender := s.keyring.GetKey(0)
	invalidSender := s.keyring.GetKey(1)
	maxSupply := math.NewInt(200)
	testcases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func()
		expErr      bool
		errContains string
		sender      keyring.Key
	}{
		{
			"fail - negative amount",
			func() []interface{} {
				return []interface{}{fromAddr, big.NewInt(-1)}
			},
			func() {},
			true,
			"coin -1xmpl amount is not positive",
			sender,
		},
		{
			"fail - invalid from address",
			func() []interface{} {
				return []interface{}{"", big.NewInt(100)}
			},
			func() {},
			true,
			"invalid from address",
			sender,
		},
		{
			"fail - invalid amount",
			func() []interface{} {
				return []interface{}{fromAddr, ""}
			},
			func() {},
			true,
			"invalid amount",
			sender,
		},
		{
			"fail - sender is not manager",
			func() []interface{} {
				return []interface{}{fromAddr, big.NewInt(100)}
			},
			func() {},
			true,
			"sender is not token manager",
			invalidSender,
		},
		{
			"fail - not enough balance",
			func() []interface{} {
				return []interface{}{fromAddr, big.NewInt(300)}
			},
			func() {},
			true,
			"not enough balance",
			sender,
		},
		{
			"pass",
			func() []interface{} {
				return []interface{}{fromAddr, big.NewInt(100)}
			},
			func() {
				addrBalance := s.network.App.BankKeeper.GetBalance(s.network.GetContext(), fromAddr.Bytes(), tokenDenom)
				s.Require().Equal(big.NewInt(100), addrBalance.Amount.BigInt(), "expected toAddr to have 100 XMPL")
			},
			false,
			"",
			sender,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()
			stateDB := s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx := testutil.NewPrecompileContract(s.T(), s.network.GetContext(), tc.sender.Addr, s.precompile, 0)

			// Set up manager role for valid sender
			err := s.precompile.assetKeep.TokenManagement.Set(
				ctx, 
				s.tokenDenom, 
				assettypes.TokenManagement{
					Managers: []string{sender.AccAddr.String()},
					ExtensionsList: []string{"mint"},
					MaxSupply: maxSupply,
				},
			)
			s.Require().NoError(err)

			// Mint amount to fromAddr to burn lately
			// _, err = s.precompile.Mint(ctx, contract, stateDB, &method, []interface{}{sender.Addr, big.NewInt(maxSupply.Int64())})
			// fmt.Println("errrrr", err)
			// s.Require().NoError(err)

			// Mint some coins to the module account and then send to the from address
			err = s.network.App.BankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(sdk.NewCoin(tokenDenom, maxSupply)))
			s.Require().NoError(err, "failed to mint coins")
			err = s.network.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, fromAddr.Bytes(), sdk.NewCoins(sdk.NewCoin(tokenDenom, maxSupply)))
			s.Require().NoError(err, "failed to send coins from module to account")

			_, err = s.precompile.BurnFrom(ctx, contract, stateDB, &method, tc.malleate())
			if tc.expErr {
				s.Require().Error(err, "expected burn transaction to fail")
				// s.Require().Contains(err.Error(), tc.errContains, "expected transfer transaction to fail with specific error")
			} else {
				s.Require().NoError(err, "expected burn transaction succeeded")
				tc.postCheck()
			}
		})
	}
}

func (s *PrecompileTestSuite) TestFreeze() {
	method := s.precompile.Methods[FreezeMethod]
	// fromAddr is the address of the keyring account used for testing.
	sender := s.keyring.GetKey(0)
	invalidSender := s.keyring.GetKey(1)
	testcases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func()
		expErr      bool
		errContains string
		sender      keyring.Key
	}{
		{
			"fail - invalid to address",
			func() []interface{} {
				return []interface{}{""}
			},
			func() {},
			true,
			"invalid from address",
			sender,
		},
		{
			"fail - sender is not manager",
			func() []interface{} {
				return []interface{}{toAddr}
			},
			func() {},
			true,
			"sender is not token manager",
			invalidSender,
		},
		{
			"pass",
			func() []interface{} {
				return []interface{}{toAddr}
			},
			func() {
				exist := s.network.App.AssetKeeper.IsFreezed(s.network.GetContext(), toAddr)
				s.Require().True(exist)
			},
			false,
			"",
			sender,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()
			stateDB := s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx := testutil.NewPrecompileContract(s.T(), s.network.GetContext(), tc.sender.Addr, s.precompile, 0)

			// Set up manager role for valid sender
			err := s.precompile.assetKeep.TokenManagement.Set(
				ctx, 
				s.tokenDenom, 
				assettypes.TokenManagement{
					Managers: []string{sender.AccAddr.String()},
					ExtensionsList: []string{"mint"},
				},
			)
			s.Require().NoError(err)

			// Mint amount to fromAddr to burn lately
			// _, err = s.precompile.Mint(ctx, contract, stateDB, &method, []interface{}{sender.Addr, big.NewInt(maxSupply.Int64())})
			// fmt.Println("errrrr", err)
			// s.Require().NoError(err)

			_, err = s.precompile.Freeze(ctx, contract, stateDB, &method, tc.malleate())
			if tc.expErr {
				s.Require().Error(err, "expected burn transaction to fail")
				// s.Require().Contains(err.Error(), tc.errContains, "expected transfer transaction to fail with specific error")
			} else {
				s.Require().NoError(err, "expected burn transaction succeeded")
				tc.postCheck()
			}
		})
	}
}
