package token

// TokenInfo jwt info
type TokenInfo struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expire      int64  `json:"expire"`
}
