package claim

import jwt "github.com/dgrijalva/jwt-go"

// ServiceAccountClaims custom claime
type ServiceAccountClaims struct {
	jwt.StandardClaims
	Name     string
	EMail    string
	Domain   string
	Provider string
	Scope    string
	Active   int64
}
