<!--
order: 6
-->

# Logic

This file describes the core logics in this module.

## Extension

### Register new extension

To intergrate with the `asset module` Each type of extension has to implement this interface

```go
type Extension interface {
    RegisterInterfaces()
    MsgHandler() MsgHandler
    QueryHandler() QueryHandler
    CLI() *cobra.Command
}

type MsgHandler func(context Context, funcMsg ExtensionMsg) error

type QueryHandler func(context Context, funcQuery ExtensionMsg) error
```

This interface provides all the extension necessary for a extension, including a message handler, query handler and cli.

All the `ExtensionMsg` of a extension should return the name of that extension when called `ExtensionName()`. A message handler should handle all the `ExtensionMsg` of that extension.

### Upgrade Extensions

All extensions are located in a seperate packages, for example: asset/extensions, therefore the exentsion or upgrade of extensions is unrelated to core logic of `Asset` module, all the modification and addition happen in only asset/extensions package. The core Asset module does not need to be aware of the specifics of extension handling. It interacts with extensions through defined interfaces or protocols.

Each Extensions has their own proto to define their state and execute/query messages. By assigning a distinct proto to each extension, you ensure that the logic, messages, and data associated with one extension do not interfere with or complicate the others. This also makes the design easier to understand, maintain, and scale.

When adding a `Extension`, we calls `ExtensionRouting.AddExtension()` in `app.go` which inturn maps all the `ExtensionMsg` of that extension to its `MsgHandler`. This mapping logic will later be used when running a `MsgExecuteExtension`.

### Message/Query routing

The message we pass in `ExecuteExtension` is `msg.Any` type. This type refered that it could be any type of message.
After receive this message, we will unpack the `msg.Any` type to an interface which implements what we want:

```go
type ExtensionMsg interface {
    ExtensionName() string
}
```

As we defined the `AllowExtensions` in Params. If the message name is in the list, they will also exist there own interface in `ExtensionRouting`.

`ExtensionRouting` acts as a centralized hub for routing messages, making it easy to manage and audit the flow of messages in the system.
It will route the `ExtensionMsg` to its `MsgHandler` - where the msg is executed. In the `MsgHandler`, the message is further routed based on its type to the correct execution logic. This additional layer of routing within the MsgHandler ensures precise handling through message types, enabling fine-grained control and execution workflows.

This flow will also work with `QueryHandler`, as long as we can unpack the `msg.Any` and get the name of the extension.

### Extension Store

## EVM interaction

### Enable EVM interface

The token includes a field named `evm_enable`, which allows it to integrate with the ERC20 standard and have an EVM-compatible contract address. This EVM address acts as an abstract interface layer that bypasses the typical logic within ERC20 or EVM contracts. Instead of executing logic directly in the contract, all actions are reflected to the `asset` module's predefined precompiles, where the token’s core state and extensions are managed.

The token itself exists as a coin within the bank state, maintaining its own logic and extensions independently of any ERC20 or EVM contract logic. The ERC20 contract deployed on the EVM serves purely as an interface, with its logic effectively bypassed. When other EVM contracts interact with this interface, their requests are forwarded via JSON-RPC calls to the `asset` module, which directly handles and executes the necessary operations. This is achieved by creating a `dynamic precompile` when `evm_enable` is activated, ensuring that the token’s behavior aligns with its internal state while still providing compatibility with the EVM ecosystem.

### ERC20 Precompiles
