<!--
order: 1
-->

# Concepts

## The Realio Asset Token Model

The Realio Asset module is centered around a token model where certain whitelisted accounts can issue their own token. A token issued by this module will be managed by `manager` accounts assigned by the issuer of the asset.

### Token extensions

Token extensions are additional features that can be flug-in for each token. There're are four types of extensions `Mint`, `Burn`, `Transfer Auth` and `Freeze`. The `Issuer` can choose what extensions to be included for his token at creation time, and only the `manager` can trigger the extension's logic.

### EVM enable

While it is the asset token in represented in the bank module, enabling the token interface in evm environment is very convenient and open up the possibility of integrating new features into the ecosystem.

Each token is automatically enabled to work in the evm environment when created, which means user can interact with the token through evm side like metamask or anyother evm wallet and more other protocol integrated in the future.
