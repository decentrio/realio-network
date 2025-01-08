<!--
order: 3
-->

# Parameters

The asset module contains the following parameters:

| Key                  | Type     | Example                |
|----------------------|----------|------------------------|
| AllowFunctionalities | []string | ["burn","freeze"]      |
| WhitelistAddresses   | []address| ["realio1..."]         |

## Details

- AllowFunctionalities: list of functionalities that the module provides. They can be update after the chain upgrade to enable new functionality add-on to the module.
- WhitelistAddresses: list of the address that's allow to create new token.
