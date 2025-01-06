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
  bool single_representation = 9;
  bool is_freeze = 10;
  DistributionSettings distribution_settings = 11;
}
```

The `issuer` is the address that create token. They can control all informations about the token, define other whitelist roles likes `manager` and `distributor`. `issuer` also can enable the token's single evm representation mode, which is showed in [EVM precompiles](README.md#asset-module-and-erc-20-precompiles).

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
