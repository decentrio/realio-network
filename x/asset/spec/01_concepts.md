<!--
order: 1
-->

# Concepts

## The Realio Asset Token Model

The Realio Asset module is centered around a token model where certain whitelisted accounts can issue their own token. A token issued by this module will be managed by a set of `manager` and `distributor` accounts. These accounts are assigned role by the issuer of the asset.

### System of privileged accounts

Privileged accounts of a token are accounts that can execute certain actions for that token. There're are several types of extensions, each has its own logic to define the actions which accounts of said type can execute. We wanna decouple the logic of these extensions from the `Asset module` logic, meaning that extensions will be defined in separate packages/modules, thus, developers can customize their type of extension without modifying the `Asset Module`. Doing this allows our extensions system to be extensible while keeping the core logic of `Asset Module` untouched and simple, avoiding complicated migration when we expand our extensions system.

In order for a extension to integrate into the `Asset Module`. It has to implement the `Extension` interface and has its implementation registered via the method `AddExtension`. Once that is done, we can make said extension available onchain by executing `SoftwareUpgradeProposal` like a regular chain upgrade process.

Currently, there are 2 type of privileged accounts: `manager` and `distributor`. Each can execute different extensions. While `distributor` can control the `mint` extension and custom the `DistributionSettings`, the `manager` can execute the other extensions like `burn` or `freeze` and could modify the  `extensions_list`. It's important to note that the `manager` can choose what extensions it wants to disable for its token.

### EVM enable

While it is useful to represent the token in bank module, enabling the token to be in action in evm environment is very convenient and pave the possibility of integrating new features into the ecosystem.

Each token can be enabled to work in the evm environment by the token manager, which means user can interact with the token through evm side like metamask or anyother evm wallet and more other protocol integrated in the future. Note that, token manager can disable or enable the token to be used in the evm environment.
