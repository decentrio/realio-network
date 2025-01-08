<!--
order: 5
-->

# Queries

## 1. QueryParams

The `QueryParams` allows users to query asset module params

```go
    type QueryParamsRequest struct {
    }
```

```go
    type QueryParamsResponse struct {
        Params Params
    }
```

CLI:

```bash
    realio-networkd q params
```

## 2. QueryToken

The `QueryToken` allows users to query a token related information

```go
    type QueryTokenRequest struct {
        TokenId string
    }
```

```go
    type QueryTokenResponse struct {
        Token               Token
        TokenManager        TokenManager
        TokenDistributor    TokenDistributor
    }
```

CLI:

```bash
    realio-networkd q token [token-id]
```

## 3. QueryAllTokens

The `QueryAllTokens` allows users to query all tokens related information

```go
    type QueryAllTokensRequest struct {
    }
```

```go
    type QueryAllTokensResponse struct {
        TokenInfo []TokenInfo

    }
```

```go
    type TokenInfo struct {
        Token               Token
        TokenManager        TokenManager
        TokenDistributor    TokenDistributor

    }
```

CLI:

```bash
    realio-networkd q all-tokens
```

## 4. QueryFrozenAddresses

The `QueryFrozenAddresses` allows users to query all frozen addresses

```go
    type QueryFrozenAddressesRequest struct {
    }
```

```go
    type QueryFrozenAddressesResponse struct {
        FrozenAddresses []address
    }
```

CLI:

```bash
    realio-networkd q frozen-addresses
```