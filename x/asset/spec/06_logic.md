<!--
order: 6
-->

# Logic

This file describes the core logics in this module.

## Functionality

### Register new functionality

To intergrate with the `asset module` Each type of functionality has to implement this interface

```go
type Functionality interface {
    RegisterInterfaces()
    MsgHandler() MsgHandler
    QueryHandler() QueryHandler
    CLI() *cobra.Command
}

type MsgHandler func(context Context, funcMsg FunctionalityMsg) error

type QueryHandler func(context Context, funcQuery FunctionalityMsg) error
```

This interface provides all the functionality necessary for a functionality, including a message handler, query handler and cli

All the `FunctionalityMsg` of a functionality should return the name of that functionality when called `FunctionalityName()`. A message handler should handle all the `FunctionalityMsg` of that functionality.

### Upgrade Priviliges

All functionalities are located in a seperate packages, for example: asset/functionalities, therefore the exentsion or upgrade of functionalities is unrelated to core logic of `Asset` module, all the modification and addition happen in only asset/functionalities package. The core Asset module does not need to be aware of the specifics of functionality handling. It interacts with functionalitys through defined interfaces or protocols.

Each Priviliges has their own proto to define their state and execute/query messages. By assigning a distinct proto to each functionality, you ensure that the logic, messages, and data associated with one functionality do not interfere with or complicate the others. This also makes the design easier to understand, maintain, and scale.

When adding a `Functionality`, we calls `FunctionalityManager.AddFunctionality()` in `app.go` which inturn maps all the `FunctionalityMsg` of that functionality to its `MsgHandler`. This mapping logic will later be used when running a `MsgExecuteFunctionality`.

### Message/Query routing

The message we pass in `ExecuteFunctionality` is `msg.Any` type. This type refered that it could be any type of message.
After receive this message, we will unpack the `msg.Any` type to an interface which implements what we want:

```go
type FunctionalityMsg interface {
    FunctionalityName() string
}
```

As we defined the `AllowFunctionalities` in Params. If the message name is in the list, they will also exist there own interface in `FunctionalityRouting`.

`FunctionalityRouting` acts as a centralized hub for routing messages, making it easy to manage and audit the flow of messages in the system.
It will route the `FunctionalityMsg` to its `MsgHandler` - where the msg is executed. In the `MsgHandler`, the message is further routed based on its type to the correct execution logic. This additional layer of routing within the MsgHandler ensures precise handling through message types, enabling fine-grained control and execution workflows.

This flow will also work with `QueryHandler`, as long as we can unpack the `msg.Any` and get the name of the functionality.

## EVM interaction

### Enable EVM interface

The token includes a field named `evm_enable`, which allows it to integrate with the ERC20 standard and have an EVM-compatible contract address. This EVM address acts as an abstract interface layer that bypasses the typical logic within ERC20 or EVM contracts. Instead of executing logic directly in the contract, all actions are reflected to the `asset` module's predefined precompiles, where the token’s core state and functionalities are managed.

The token itself exists as a coin within the bank state, maintaining its own logic and functionalities independently of any ERC20 or EVM contract logic. The ERC20 contract deployed on the EVM serves purely as an interface, with its logic effectively bypassed. When other EVM contracts interact with this interface, their requests are forwarded via JSON-RPC calls to the `asset` module, which directly handles and executes the necessary operations. This is achieved by creating a `dynamic precompile` when `evm_enable` is activated, ensuring that the token’s behavior aligns with its internal state while still providing compatibility with the EVM ecosystem.

### ERC20 Precompiles