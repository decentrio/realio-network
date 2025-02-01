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
| `WhitelistAddresses` | Whitelist Addresses                    | `[]byte{4} + []byte(address)`                             | `[]byte{bool}`                        | KV    |
| `FreezeAddresses` | Whitelist Addresses                    | `[]byte{5} + []byte(address)`                             | `[]byte{bool}`                        | KV    |

### Token

Allows creation of tokens with optional user authorization.  

```go
type Token struct {
	TokenId     string `protobuf:"bytes,1,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
	Issuer      string `protobuf:"bytes,2,opt,name=issuer,proto3" json:"issuer,omitempty"`
	Name        string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Symbol      string `protobuf:"bytes,4,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Decimal     uint32 `protobuf:"varint,5,opt,name=decimal,proto3" json:"decimal,omitempty"`
	Description string `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	EvmAddress  string `protobuf:"bytes,9,opt,name=evm_address,json=evmAddress,proto3" json:"evm_address,omitempty"`
}
```

When create the token, `asset` module auto generate for it a evm address. This address is used as a dynamic precompiles.

### TokenManagement

```go
type TokenManagement struct {
	Managers           []string              `protobuf:"bytes,1,rep,name=managers,proto3" json:"managers,omitempty"`
	ExtensionsList     []string              `protobuf:"bytes,3,rep,name=extensions_list,json=extensionsList,proto3" json:"extensions_list,omitempty"`
	MaxSupply          cosmossdk_io_math.Int `protobuf:"bytes,4,opt,name=max_supply,json=maxSupply,proto3,customtype=cosmossdk.io/math.Int" json:"max_supply"`
}
```

`extensions_list` is the list of actions that the manager can execute.

`MaxSupply` defines the maximum number of tokens can be minted.

### WhitelistAddresses

`WhitelistAddresses` is a list of the address that's allow to create new asset.

