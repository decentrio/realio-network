<!--
order: 4
-->

# Messages

## 1. MsgIssueToken

`MsgIssueToken` allow issuer to create token. The issuer must be in param's whitelist addresses to be able to execute this msg.

```go
    type MsgIssueToken struct {
        Issuer                     address
        Managers                   [ ]address
        Distributors               [ ]address
        Symbol                     string   
        Decimal                    uint32   
        Description                string 
        EvmEnable                  bool
        AllowNewExtensions     bool
        ExtensionsList             [ ]string
    }
```

```go
    type MsgIssueTokenResponse struct {
    }
```

CLI:

```bash
    realio-networkd tx issue-token [token.json] [flags]
```

Example token.json:

```json
   {
      "Manager": ["realioabc..."],
      "Distributor": ["realioabc2..."],
      "Symbol": "riel",
      "Decimal": "rielio",
      "Description": "",
      "EvmEnable": true,
      "AllowNewExtensions": true,
      "ExtensionsList": [],
    }
```

Validation:

- Check if Creator is whitelisted. We only allow some certain accounts to create tokens, these accounts is determined via gov proposal.
- Check if token has been created or not by iterating through all denom existing.
- Sanity check on token info like decimal, description

Flow:

1. The denom for the token will be derived from Creator and Symbol with the format of asset/{Issuer}/{Symbol-Lowercase}
2. If `EvmEnable` is true, create a dynamic precompiles for the token.
3. Save the token basic information (name, symbol, decimal and description) in the x/bank metadata store
4. Save the token management info and distribution info in the x/asset store.

## 2. AssignRoles

`MsgAssignRoles` allow issue to set role likes manager or distributor for the token.

```go
    type MsgAssignRoles struct {
        TokenId         string
        Issuer          address
        Addresses       mapping[Role]([]addresses)
    }
```

```go
    type MsgAssignRolesResponse struct {
    }
```

CLI:

```bash
    realio-networkd tx assign-roles [privilege.json] [flags]
```

Example privilege.json:

```json
    {
        "TokenId": "asset/realio1.../tokena",
        "Issuer": "realio1...",
        "Assign": [
            {
                "role": 1 (manager),
                "addresses": ["realio2..."],
            },
            {
                "role": 2 (distributor),
                "addresses": ["realio3..."],
            }
        ]
    }
```

Validation:

- Check if token exists
- Check if caller is issuer of the token
- Check if addresses is valid
- Check if manager doesn't exist in the current managers list of token
- Check if distributor doesn't exist in the current distributor list of token

Flow:

- Get `TokenManager` and `TokenDistributor` from store by token_id
- Loop through addresses and append manager addresses to `TokenManager.Managers`, distributor addresses to `TokenDistributor.Distributors`

## 3. UnassignRole

```go
    type MsgUnassignRole struct {
        TokenId         string
        Issuer          address
        Assigners       []address
    }
```

```go
    type MsgUnassignRoleResponse struct {
    }
```

Validation:

- Check if token exists
- Check if caller is issuer of the token
- Check if addresses is valid
- Check if addresses is in `TokenManager.Managers` or `TokenDistributor.Distributors`

Flow:

- Get `TokenManager` and `TokenDistributor` from store by token_id
- Loop through addresses and remove manager addresses from `TokenManager.Managers`, distributor addresses to `TokenDistributor.Distributors`

## 4. ExecuteExtension

After setting the managers, the managers can execute their allowed extension.

```go
    type MsgExecuteExtension struct {
        Manager              address     
        TokenId              string     
        ExtensionMsg     *types.Any
    }
```

Validation:

- Checks if the token specified in the msg exists.
- Checks if the extension is supported.
- Checks if the `Msg.Address` has the corresponding `Extension` specified by `ExtensionMsg.NeedExtension()`

Flow:

- Prepare store for the extension of the token via `MakeExtensionStore(extension name, token denom)`. That store is the only store accessable by the extension's `MsgHandler`.
- `ExtensionMsgRouting` routes the `ExtensionMsg` to the its `MsgHandler`.
- `MsgHandler` now handles the `ExtensionMsg`.

### 5. Mint

This msg only can be executed when the token's `ExtensionsList` has `mint` extension.

```go
    type MsgMint struct {
        Distributor          address     
        TokenId              string
        Receiver             address
        Amount               math.Int
    }
```

Validation:

- Checks if the token specified in the msg exists.
- Checks if the extension is supported.
- Check if addresses is valid
- Checks if the distributor address is in `TokenDistributor.Distributors`
- Checks if mint amount exceed `MaxSupply` or `MaxRatelimit`.

Flow:

- Get `TokenDistributor` from store by token_id
- Mint the asset for corresponding reciever
- Increase the ratelimit

### 6. UpdateDistributionSetting

Distributor can change the max supply and mint ratelimit of the token.

```go
    type MsgUpdateDistributionSetting struct {
        Distributor          address     
        TokenId              string
        NewSettings          DistributionSettings
    }
```

Validation:

- Checks if the token specified in the msg exists.
- Checks if the extension is supported.
- Checks if the distributor address is in `TokenDistributor.Distributors`
- Checks if current supply exceed new settings `MaxSupply`

### 7. UpdateExtensionsList

Manager can update the `ExtensionsList` of the token. This only can be executed when the token's `AllowNewExtensions` is enable.

```go
    type ExtensionsList struct {
        Manager              address     
        TokenId              string
        NewExtensions   []string
    }
```

Validation:

- Checks if the token specified in the msg exists.
- Checks if manager addresses is in `TokenManager.Managers`
- Checks if the new extension is supported.

### 8. UpdateParams
