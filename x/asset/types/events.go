package types

// staking module event types
const (
	EventTypeTokenCreated          = "create_token"
	EventTypeTokenAuthorizeUpdated = "update_authorize_token"

	AttributeKeyTokenId = "token_id"
	AttributeKeyIndex   = "index"
	AttributeKeyAddress = "address"

	AttributeValueCategory = ModuleName
)
