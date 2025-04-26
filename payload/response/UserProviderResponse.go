package response

type UserProviderResponse struct {
	//UserID           uint   `json:"user_id"`
	Provider string `json:"provider"`
	//AppleId          string `json:"apple_id"`
	UserType         string `json:"user_type"`
	ProviderIdentify string `json:"provider_identify"`
}
