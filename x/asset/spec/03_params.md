<!--
order: 3
-->

# Parameters

The asset module contains the following parameters:

| Key                  | Type          | Example                |
|----------------------|---------------|------------------------|
| AllowExtensions      | []string      | ["burn","freeze"]      |
| RatelimitDuration    | time.Duration | "86400s"               |
| WhitelistAddresses   | []address     | ["realio1..."]         |

## Details

- AllowExtensions: list of extensions that the module provides. They can be update after the chain upgrade to enable new extension add-on to the module.
- RatelimitDuration: duration of ratelimit for `mint` extension.
- WhitelistAddresses: list of the address that's allow to create new token.
