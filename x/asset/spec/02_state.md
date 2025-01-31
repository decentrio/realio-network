<!--
order: 2
-->

# State

## State Objects

The `x/asset` module keeps the following objects in state:

| State Object         | Description                            | Key                                                       | Value                                 | Store |
|----------------------|----------------------------------------|-----------------------------------------------------------|---------------------------------------|-------|
| `Params`             | Params of asset module                 | `[]byte{1}`                                               | `[]byte(params)`                      | KV    |
| `Token`              | Token information                      | `[]byte{2} + []byte(token_id)`                            | `[]byte{token}`                       | KV    |
| `TokenManagement`    | TokenManagement info of a denom        | `[]byte{3} + []byte(token_id)`                            | `[]byte{token_manager}`               | KV    |
| `TokenDistribution`  | TokenDistribution info of a denom      | `[]byte{4} + []byte(token_id)`                            | `[]byte{token_distributor}`           | KV    |
| `WhitelistAddresses` | Whitelist Addresses                    | `[]byte{5} + []byte(address)`                             | `[]byte{bool}`                        | KV    |
| `ExtenstionStore`    | State store for each extensions        | `[]byte{6} + []byte(token_id) + []byte(extension_name)`   |  Depend on extension implementation   | KV    |

### Token

Allows creation of tokens with optional user authorization.  

```go
type Token struct {
    TokenId                string               `protobuf:"bytes,1,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
    Issuer                 string               `protobuf:"bytes,2,opt,name=issuer,proto3" json:"issuer,omitempty"`
    Name                   string               `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
    Symbol                 string               `protobuf:"bytes,4,opt,name=symbol,proto3" json:"symbol,omitempty"`
    Decimals               uint32               `protobuf:"varint,5,opt,name=decimals,proto3" json:"decimals,omitempty"`
    Description            string               `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
    EVMAddress             common.Address       `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
}
```

The token id for the token will be derived from the Issuer and the Symbol with the format of asset/{Issuer}/{Symbol-Lowercase}. This will allow 2 tokens to have the same name with different issuers.

The `issuer` is the address that create token. They can control all informations about the token, define other whitelist roles likes `manager` and `distributor`. `issuer` also can enable the token's single evm representation mode, which is showed in [EVM precompiles](README.md#asset-module-and-erc-20-precompiles).

When create the token, `asset` module auto generate for it a evm address. This address is used as a dynamic precompiles.

### TokenManagement

```go
type TokenManagement struct{
    Managers               []string             `protobuf:"bytes,7,rep,name=managers,proto3" json:"managers,omitempty"`
    AllowNewExtensions bool                 `protobuf:"varint,10,opt,name=allow_new_Extensions,json=allowNewExtensions,proto3" json:"allow_new_Extensions,omitempty"`
    ExtensionsList    []string             `protobuf:"bytes,11,rep,name=extensions_list,json=extensionsList,proto3" json:"extensions_list,omitempty"`
    EvmEnable              bool                 `protobuf:"varint,9,opt,name=evm_enable,json=evmEnable,proto3" json:"evm_enable,omitempty"`
   }
```

By setting `allow_new_extensions`, `issuer` can specify whether they accept new extensions or not when creating a new token. If he permits it, when upgrading the chain, the new features will be automatically added to the `extensions_list`and the `manager` can then modify the `extensions_list` as he sees fit. Otherwise, the `manager` can not chaing the `extensions_list`.

### TokenDistribution

```go
type TokenDistribution struct{
    Distributors           []string
    MaxSupply              math.Int
}
```

`MaxSupply` defines the maximum number of tokens can be minted.

### WhitelistAddresses

`WhitelistAddresses` is a list of the address that's allow to create new asset.

