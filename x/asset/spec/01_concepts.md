<!--
order: 1
-->

# Concepts

## The Realio Asset Token Model

The Realio Asset module is centered aroumd a token model. It contains the following fields:

### Token

```protobuf
message Token {
  string token_id = 1;
  string issuer = 2;
  string name = 3;
  string symbol = 4;
  uint32 decimals = 5;
  string description = 6;
  []address manager_addresses = 7;
  []address distributor_addresses = 8;
  bool evm_enable = 9;
  bool allow_new_fuctionalities = 10;
  []string functionalities_list = 11;
  DistributionSettings distribution_settings = 12;
}
```

The `issuer` is the address that create token. They can control all informations about the token, define other whitelist roles likes `manager` and `distributor`. `issuer` also can enable the token's single evm representation mode, which is showed in [EVM precompiles](README.md#asset-module-and-erc-20-precompiles).

When creating a new token, `issuer` can define that they accept new functionalities or not by setting `allow_new_fuctionalities`. If he allow it, the new functionality when upgrading chain will be automatically added to the `functionalities_list` and the `issuer` can adjust the `functionalities_list` by his will.

By setting `allow_new_fuctionalities`, `issuer` can specify whether they accept new functionalities or not when creating a new token. If he permits it, when upgrading the chain, the new features will be automatically added to the `functionalities_list`and the `manager` can then modify the `functionalities_list` as he sees fit. Otherwise, the `manager` can not chaing the `functionalities_list`.

### Role

In token model, each token has 2 roles which can execute different functionality. They are whitelisted address that is defined by the issuer of the token. While the `manager` can execute the `freeze` and `burn` functionality, `distributor` can control the `mint` functionality and custom the `DistributionSettings`.

- "ROLE_UNSPECIFIED": 0
- "ROLE_MANAGER": 1
- "ROLE_DISTRIBUTOR": 2

### DistributionSettings

```protobuf
message DistributionSettings {
  string max_supply = 1[(gogoproto.customtype) = "cosmossdk.io/math.Int"]; 
  string max_ratelimit = 2[(gogoproto.customtype) = "cosmossdk.io/math.Int"];
}
```
