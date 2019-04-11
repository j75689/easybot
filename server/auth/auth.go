package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/j75689/easybot/model"
)

var (
	ErrExpiredToken     error = errors.New("Token is expired")
	ErrNotValidYetToken error = errors.New("Token not active yet")
	ErrMalformedToken   error = errors.New("That's not even a token")
	ErrInvalidToken           = errors.New("Invalid token")
)

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       time.Duration
	tokenType     string
}

// TokenInfo jwt info
type TokenInfo struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expire      int64  `json:"expire"`
}

// ServiceAccountClaims custom claime
type ServiceAccountClaims struct {
	jwt.StandardClaims
	Name     string
	EMail    string
	Domain   string
	Provider string
	Scope    string
	Active   int
}

var (
	defaultKey = []byte("easybot")
	Options    = options{
		tokenType:     "Bearer",
		signingMethod: jwt.SigningMethodHS512,
		signingKey:    defaultKey,
		keyfunc: func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return defaultKey, nil
		},
	}
)

func SetTokenType(tokenType string) {
	Options.tokenType = tokenType
}

func SetSigningMethod(method *jwt.SigningMethodHMAC) {
	Options.signingMethod = method
}

func SetSigningKey(secret string) {
	Options.signingKey = []byte(secret)
}

// GenerateToken create new jwt token
func GenerateToken(info *model.ServiceAccount) (*TokenInfo, error) {
	now := time.Now()
	claim := &ServiceAccountClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
		Name:     info.Name,
		EMail:    info.EMail,
		Domain:   info.Domain,
		Provider: info.Provider,
		Scope:    info.Scope,
		Active:   info.Active,
	}
	var expiresAt int64
	if info.Active > 0 {
		expiresAt = now.Add(time.Duration(info.Active) * time.Second).Unix()
		claim.ExpiresAt = expiresAt
	}

	token := jwt.NewWithClaims(Options.signingMethod, claim)

	tokenString, err := token.SignedString(Options.signingKey)
	if err != nil {
		return nil, err
	}
	return &TokenInfo{
		AccessToken: tokenString,
		TokenType:   Options.tokenType,
		Expire:      expiresAt,
	}, nil
}

// ParseToken resolve TokenString to Info
func ParseToken(tokenString string) (*ServiceAccountClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ServiceAccountClaims{}, Options.keyfunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrMalformedToken
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrExpiredToken
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrNotValidYetToken
			} else {
				return nil, ErrInvalidToken
			}
		}
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token.Claims.(*ServiceAccountClaims), nil
}

func RefreshToken(tokenString string) (*TokenInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ServiceAccountClaims{}, Options.keyfunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ServiceAccountClaims); ok && token.Valid {
		var expiresAt int64
		if claims.Active > 0 {
			now := time.Now()
			expiresAt = now.Add(Options.expired).Unix()
			claims.ExpiresAt = time.Now().Add(Options.expired).Unix()
		}

		token := jwt.NewWithClaims(Options.signingMethod, claims)
		tokenString, err := token.SignedString(Options.signingKey)
		if err != nil {
			return nil, err
		}
		return &TokenInfo{
			AccessToken: tokenString,
			TokenType:   Options.tokenType,
			Expire:      expiresAt,
		}, nil
	}

	return nil, ErrInvalidToken
}

// GetTokenFromRequest get token info from http request
func GetTokenFromRequest(request *http.Request) (*ServiceAccountClaims, error) {
	tokenString := request.Header.Get("Authorization")
	prefix := fmt.Sprintf("%s ", Options.tokenType)
	if tokenString != "" && strings.HasPrefix(tokenString, prefix) {
		tokenString = tokenString[len(prefix):]
	} else {
		return nil, ErrMalformedToken
	}

	return ParseToken(tokenString)
}
