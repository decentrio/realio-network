<!--
order: 6
-->

# Logic

This file describes the core logics in this module.

## Token namespace

The token id/denom for the token will be derived from the Issuer and the Symbol with the format of asset/{Issuer}/{Symbol-Lowercase}. This allow many different issuers to issue token with the same symbol, differentiate their denom by including in their creator.

## Token creation

The token creation process involves the `Issuer` executing `MsgIssueToken` that defines info fields for the tokens. Amount those fields, these are the important feilds that dictacts how the token is operated:

- List of extensions (imutable):
    Choose what exetensions to enable for the token when creating it. The set of extensions is fixed, meaning that all extensions that are not enabled is permanently disabled.
- Manager (mutable):
    Assign the manager account which could be an user account or smart contract account. If this field is blank then the manager will set to be the `issuer` of the token.
- Symbol (imutable):
    This field will be used to derive the token denom, in the format `asset/{Issuer}/{Symbol-Lowercase}`.
- Initial Supply (imutable):
    Upon creation an amount equal to the initial supply will be minted thus create a circulating supply for the token.
- Distributor (imutable):
    Think of this account as the treasury manager account whose task is to distribute the token to holders. The initial supply will be minted to this account.

## Extension

We'll go into details on how each of the extension works

### Mint extension
Only manager is allowed to execute `Mint`. Then mint the corresponding amount to the recipient. Note, total supply can not exceed `MaxSupply`

### Burn extension
Only manager is allowed to execute `Burn`. Then burn the corresponding amount from the address. Note, address that be freezed can not burn.

### Freeze extension
Only manager is allowed to execute `Freeze`. Its will lock all amount of a  asset of that address. An address be freezed with a token can not transfer out or burn.

### TransferAuth extension
Only manager is allowed to execute `TransferAuth`. Its will update the Token's Issue to new receiver.

### 

## EVM integration

### EVM interface

On token creation, all token will be linked to a erc20-precompiles, which allows it to integrate with the ERC20 standard and have an EVM-compatible contract address. This EVM address acts as an abstract interface layer that bypasses the typical logic within ERC20 or EVM contracts. Instead of executing logic directly in the contract, all actions are reflected to the `asset` module's predefined precompiles, where the token’s core state and extensions are managed.

The token itself exists as a coin within the bank state, maintaining its own logic and extensions independently of any ERC20 or EVM contract logic. The ERC20 contract deployed on the EVM serves purely as an interface, with its logic effectively bypassed. When other EVM contracts interact with this interface, their requests are forwarded via JSON-RPC calls to the `asset` module, which directly handles and executes the necessary operations. This is achieved by creating a `dynamic precompile`, ensuring that the token’s behavior aligns with its internal state while still providing compatibility with the EVM ecosystem.

The precompiles actions will depend on the `AllowExtensionList` when creating the token. Therefore different tokens will have precompiles with different addresses and extensions.

### EVM Precompiles

EVM precompiles are EVM interface contracts with state access. These smart contracts can directly interact with Cosmos SDK modules, enabling their own operations while also interacting with the EVM state and other SDK modules.

In `asset` module, there are 2 evm precompiles contracts: `IAsset.sol` corresponding to `asset` precompile and `IERC20Extensions.sol` corresponding to `erc20` precompile.

The `IAsset.sol` is an interface through which Solidity contracts can interact with Realio asset module. This is convenient for developers as they don’t need to know the implementation details behind the `x/asset` module in the Realio Network. Instead, they can interact with `asset` functions using the Ethereum interface they are familiar with.

`asset` precompile provides several functions:

- `issueToken` enables other contracts or users to create an ERC20 token.
- `updateExtensionsList` allow token manager to interact to `asset` module and update the extensions list.
- `assignRoles` allow token issuer to assign role for token.
- `unassignRoles` allow token issuer to unassign role for token.

The functions is defined as follows:

```solidity
    function issueToken(
        address issuerAddress,
        string memory name, 
        string memory symbol,
        uint8 deciaml,
        bool allowNewExtensions,
        string[] memory extensionsList
    ) external returns (bool success, address contractAddress);

    function updateExtensionsList(string memory tokenId, string[] memory newExtensionsList) public;

    struct Role {
        uint8 role; // 1 represent manager, 2 represent distributor
        address account;
    }

    function assignRoles(string memory tokenId, Role[] roles) public;

    function unassignRoles(string memory tokenId, address[] accounts) public;
```

When this function is called, the token-issuer precompile forwards the request to the asset module, invoking the IssueToken function within the module to handle the token creation process.

On the other hand, the `erc20` precompile acts as the ERC20 interface for all tokens managed by the asset module. It implements all standard ERC20 functions as defined in `IERC20Extensions.sol`, including `transfer`, `transferFrom`, `approve`, `increaseAllowance`, and `decreaseAllowance`.

Additionally, the `IERC20Extensions.sol` contract provides extra methods to support interactions with other extensions, enabling more advanced functionality:

```solidity
    function mint(address to, uint256 amount) public;

    function burn(uint256 value) public;
    function burnFrom(address account, uint256 value) public;

    function freeze(address account) public;
```
