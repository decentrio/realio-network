package types

const (
	EventTypeWithdraw             = "withdraw"
	EventTypeWithdrawClaimConfirm = "withdraw_claim_confirm"

	AttributeKeySender     = "sender"
	AttributeKeyReceiver   = "receiver"
	AttributeKeyAmount     = "amount"
	AttributeKeyDenom      = "denom"
	AttributeKeyConnection = "connection_id"
	AttributeKeyNonce      = "nonce"
	AttributeKeyApprovals  = "approvals"
	AttributeKeyPackage    = "package"

	AttributeValueCategory = ModuleName
)
