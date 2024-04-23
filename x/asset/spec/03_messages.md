<!--
order: 3
-->

# Messages

## MsgCreateAsset
Asset can be created by address that allowed via `MsgCreateAsset`. First, it'll check if the `asset_name` already exists. If not, the asset will be created. By default, creator is admin of the asset after created.
```
message MsgCreateAsset {
    string creator; (sender)
    string asset_name;
    string total_supply;
    cosmos.bank.v1beta1.Metadata metadata;
}
```
The message handling should fail if:
* `creator` not allowed by gov
* `asset_name` already exist
## MsgSetAssetMetadata
Admin of the asset can change asset's metadata with the `MsgSetAssetMetadata`.
```
message MsgCreateAsset {
    string admin; (sender)
    string asset_name;
    cosmos.bank.v1beta1.Metadata metadata;
}
```
The message handling should fail if:
* `asset_name` not exist
* mess`sender` isn't asset's admin
* `metadata` validate fail
## MsgChangeAdmin
`MsgChangeAdmin` allow to change admin of the asset to the different account address.
```
message MsgCreateAsset {
    string old_admin; (sender)
    string new_admin;
    string asset_name;
}
```
The message handling should fail if:
* `asset_name` not exist
* mess`sender` isn't asset's admin
## MsgMint
```
message MsgMint {
    string minter; (sender)
    string to_address;
    cosmos.base.v1beta1.Coin amount;
}
```
- to_address => can be in config (admin address, distributor address)
## MsgBurn
```
message MsgBurn {
    string burner; (sender)
    string from_address;
    cosmos.base.v1beta1.Coin amount;
}
```
## MsgFreezingAccount
```
message MsgFreezingAccount {
    string admin; (sender)
    string asset_name;
    string address;
}
```
## MsgUnFreezingAccount
```
message MsgUnFreezingAccount {
    string admin; (sender)
    string asset_name;
    string address;
}
```
## MsgRevokeAsset
```
message MsgRevokeAsset {
    string admin; (sender)
    string address;
    cosmos.base.v1beta1.Coin amount;
}
```
## MsgPromoteRole
```
message MsgRevokeAsset {
    string admin; (sender)
    string address;
    string role;
}
```
## MsgBlockAsset