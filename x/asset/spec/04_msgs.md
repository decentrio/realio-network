<!--
order: 4
-->

# Messages

## 1. MsgIssueToken

`MsgIssueToken` allow issuer to create token.

```go
    type MsgIssueToken struct {
        Issuer                     address
        Managers                   [ ]address
        Distributors               [ ]address
        ContractAddress            string
        Symbol                     string   
        Decimal                    uint32   
        Description                string 
        SingleRepresentation       bool
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
      "ContractAddress": "0x...",
      "Symbol": "riel",
      "Decimal": "rielio",
      "Description": "",
      "SingleRepresentation": true,
    }
```

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

## 4. ExecuteFunctionality

After setting the managers, the managers can execute their allowed functionality.

```go
    type MsgExecuteFunctionality struct {
        Manager              address     
        TokenId              string     
        FunctionalityMsg     *types.Any
    }
```

CLI:

```bash
    realio-networkd tx execute-functionality [contract-address] [msg.json] [flags]
```

Example:

```bash
    realio-networkd tx execute-functionality 0x... msg.json --from mykey
```

```json
{
  "FreezeMsg": {}
}
```

### 5. Mint

```go
    type MsgMint struct {
        Distributor          address     
        TokenId              string
        Receiver             address
        Amount               math.Int
    }
```

### 6. UpdateDistributionSetting

Distributor can change the max supply and mint ratelimit of the token.

```go
    type MsgUpdateDistributionSetting struct {
        Distributor          address     
        TokenId              string
        NewSettings          DistributionSettings
    }
```

### 7. UpdateParams
