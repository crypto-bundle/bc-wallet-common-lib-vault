package vault

const (
	AuthMethodNameTag          = "auth_method"
	RenewTokenChannelStatusTag = "renew_channel_status"
	//nolint:gosec // it's ok...wtf..why gosec triggered on this conts
	RenewTokenLeaseDurationTag = "renew_lease_duration"
)
