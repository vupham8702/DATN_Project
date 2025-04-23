package response

type UserToken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	Error        error  `json:"-"`
}
