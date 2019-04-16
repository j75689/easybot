package model

import (
	"github.com/j75689/easybot/auth/claim"
	"github.com/j75689/easybot/auth/token"
)

// ServiceAccount api access token account
type ServiceAccount struct {
	ID       string `json:"-" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	EMail    string `json:"email" bson:"email"`
	Domain   string `json:"domain" bson:"domain"`
	Provider string `json:"provider" bson:"provider"`
	Scope    string `json:"scope" bson:"scope"`
	Active   int64  `json:"active" bson:"active"`

	Generate int64  `json:"generate" bson:"generate"`
	Expired  int64  `json:"expired" bson:"expired"`
	Token    string `json:"token" bson:"token"`
}

// ValidInfo token info match account
func (account *ServiceAccount) ValidInfo(token *token.TokenInfo, claim *claim.ServiceAccountClaims) bool {
	return account.ValidName(claim.Name) &&
		account.ValidDomain(claim.Domain) &&
		account.ValidEmail(claim.EMail) &&
		account.ValidProvider(claim.Provider) &&
		account.ValidScope(claim.Scope) &&
		account.ValidActive(claim.Active) &&
		account.ValidExpired(claim.ExpiresAt) &&
		account.ValidToken(token.AccessToken)
}

// ValidName field
func (account *ServiceAccount) ValidName(name string) bool {
	return account.Name == name
}

// ValidEmail field
func (account *ServiceAccount) ValidEmail(email string) bool {
	return account.EMail == email
}

// ValidDomain field
func (account *ServiceAccount) ValidDomain(domain string) bool {
	return account.Domain == domain
}

// ValidProvider field
func (account *ServiceAccount) ValidProvider(provider string) bool {
	return account.Provider == provider
}

// ValidScope field
func (account *ServiceAccount) ValidScope(scope string) bool {
	return account.Scope == scope
}

// ValidActive field
func (account *ServiceAccount) ValidActive(active int64) bool {
	return account.Active == active
}

// ValidExpired field
func (account *ServiceAccount) ValidExpired(expired int64) bool {
	return account.Expired == expired
}

// ValidToken field
func (account *ServiceAccount) ValidToken(token string) bool {
	return account.Token == token
}
