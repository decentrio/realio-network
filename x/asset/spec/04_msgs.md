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
        Name                       string   
        Symbol                     string   
        Decimal                    uint32   
        Description                string 
        ExtensionsList             [ ]string
        Distributor                [ ]string
        InitialSupply              [ ]math.Int
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
      "Symbol": "riel",
      "Decimal": 18,
      "Description": "",
      "AllowNewExtensions": true,
      "ExtensionsList": [],
    }
```

Validation:

- Check if Creator is whitelisted. We only allow some certain accounts to create tokens, these accounts is determined via gov proposal.
- Check if token has been created or not by iterating through all denom existing.
- Sanity check on token info like decimal, description

Flow:

1. The token-id for the token will be derived from Creator and Symbol with the format of asset/{Issuer}/{Symbol-Lowercase}
2. Create a evm address for the asset.
3. Create a dynamic precompiles linking to the newly created evm address.
4. Save the token basic information (name, symbol, decimal and description) in the x/bank metadata store
5. Save the token management info and distribution info in the x/asset store.

## 2. AssignManagers

`MsgAssignManagers` allow issue to set managers for the token.

```go
    type MsgAssignManagers struct {
        TokenId         string
        Issuer          address
        Addresses       []addresses
    }
```

```go
    type MsgAssignManagersResponse struct {
    }
```

CLI:

```bash
    realio-networkd tx assign-managers [privilege.json] [flags]
```

Example privilege.json:

```json
    {
        "TokenId": "asset/realio1.../tokena",
        "Issuer": "realio1...",
        "Assign": [
            "realio2...",
            "realio3..."
        ]
    }
```

Validation:

- Check if token exists
- Check if caller is issuer of the token
- Check if addresses is valid
- Check if manager doesn't exist in the current managers list of token

Flow:

- Get `TokenManager` from store by token_id
- Loop through addresses and append manager addresses to `TokenManager.Managers`

## 3. UnassignManager

```go
    type MsgUnassignRoles struct {
        TokenId         string
        Issuer          address
        Assigners       []address
    }
```

```go
    type MsgUnassignRolesResponse struct {
    }
```

Validation:

- Check if token exists
- Check if caller is issuer of the token
- Check if addresses is valid
- Check if addresses is in `TokenExtensions.Managers` 

Flow:

- Get `TokenManager` from store by token_id
- Loop through addresses and remove manager addresses from `TokenManager.Managers`

## 4. Burn

This msg only can be executed when the token's `ExtensionsList` has `burn` extension.

```go
    type MsgBurn struct {
        Manager              address     
        TokenId              string     
        BurnFromAddr         address
        Amount               math.Int
    }
```

Validation:

- Checks if the token specified in the msg exists.
- Checks if the extension is supported.
- Check if addresses is valid
- Checks if the address is in `TokenManager.Managers`
- Checks if address is freezed in `FreezeAddresses`

Flow:

- Get `TokenManager` from store by token_id
- Check if `BurnFromAddr` has enough token to burn
- Burn the asset from `BurnFromAddr`

### 5. Mint

This msg only can be executed when the token's `ExtensionsList` has `mint` extension.

```go
    type MsgMint struct {
        Manager              address     
        TokenId              string
        Receiver             address
        Amount               math.Int
    }
```

Validation:

- Checks if the token specified in the msg exists.
- Checks if the extension is supported.
- Check if addresses is valid
- Checks if the address is in `TokenManager.Managers`
- Checks if mint amount exceed `MaxSupply`.

Flow:

- Get `TokenManager` from store by token_id
- Mint the asset for corresponding receiver
- Increase the supply.

### 6. Freeze

This msg only can be executed when the token's `ExtensionsList` has `freeze` extension.

```go
    type MsgFreeze struct {
        Manager              address     
        TokenId              string
        Receiver             address
    }
```

Validation:

- Checks if the token specified in the msg exists.
- Checks if the extension is supported.
- Check if addresses is valid
- Checks if the address is in `TokenManager.Managers`

Flow:

- Get `TokenManager` from store by token_id
- Set address into `FreezeAddresses`
- All account in `FreezeAddresses` can not be transfer token out or burned.

### 7. Set max supply

```go
    type MsgSetMaxSupply struct {
        Manager                    address
        TokenId                    string
        MaxSupply                  int64
    }
```

This message can only executed once, it will set the maximum supply for the token



