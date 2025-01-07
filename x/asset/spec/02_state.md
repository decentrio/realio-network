<!--
order: 2
-->

# State

## State Objects

The `x/asset` module keeps the following objects in state:

| State Object         | Description                    | Key                      | Value                | Store |
|----------------------|--------------------------------|--------------------------| ---------------------|-------|
| `Token`              | Token bytecode                 | `[]byte{1} + []byte(id)` | `[]byte{token}`      | KV    |
| `FreezedAddresses`   | Addresses bytecode             | `[]byte{21} + []byte(id)`| `[]byte{[]address}`  | KV    |
| `Params`             | Params bytecode                | `[]byte{3} + []byte(id)` | `[]byte(id)`         | KV    |

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
    Managers               []string             `protobuf:"bytes,7,rep,name=managers,proto3" json:"managers,omitempty"`
    Distributors           []string             `protobuf:"bytes,8,rep,name=distributors,proto3" json:"distributors,omitempty"`
    EvmEnable              bool                 `protobuf:"varint,9,opt,name=evm_enable,json=evmEnable,proto3" json:"evm_enable,omitempty"`
    AllowNewFuctionalities bool                 `protobuf:"varint,10,opt,name=allow_new_fuctionalities,json=allowNewFuctionalities,proto3" json:"allow_new_fuctionalities,omitempty"`
    FunctionalitiesList    []string             `protobuf:"bytes,11,rep,name=functionalities_list,json=functionalitiesList,proto3" json:"functionalities_list,omitempty"`
    DistributionSettings   DistributionSettings `protobuf:"bytes,12,opt,name=distribution_settings,json=distributionSettings,proto3" json:"distribution_settings"`
}
```

### FreezedAddress

List of addresses that is freezed by the manager. This only exists when the Token enable the `freeze` functionality. The addresses in list will not be able to execute any msg about the token.

### Params

```go
type Params struct {
    AllowFunctionalities []string `protobuf:"bytes,1,rep,name=allow_functionalities,json=allowFunctionalities,proto3" json:"allow_functionalities,omitempty"`
}
```

## Genesis State

The `x/asset` module's `GenesisState` defines the state necessary for initializing the chain from a previous exported height. It contains the module parameters and the registered token pairs :

```go
// GenesisState defines the module's genesis state.
type GenesisState struct {
    Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
    PortId string `protobuf:"bytes,2,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
}
```
