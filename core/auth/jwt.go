/*
	jwt.go
	Purpose: Methods for jwt.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package auth

import (
	"time"

	"app/core/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/xid"
)

var Secret = []byte(util.RandStr(32))

/*

	type RegisteredClaims struct {
		// the `iss` (Issuer) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
		Issuer string `json:"iss,omitempty"`

		// the `sub` (Subject) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
		Subject string `json:"sub,omitempty"`

		// the `aud` (Audience) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
		Audience ClaimStrings `json:"aud,omitempty"`

		// the `exp` (Expiration Time) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
		ExpiresAt *NumericDate `json:"exp,omitempty"`

		// the `nbf` (Not Before) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
		NotBefore *NumericDate `json:"nbf,omitempty"`

		// the `iat` (Issued At) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
		IssuedAt *NumericDate `json:"iat,omitempty"`

		// the `jti` (JWT ID) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
		ID string `json:"jti,omitempty"`
	}

*/

type Claims struct {
	User  string `json:"usr,omitempty"`
	Group int    `json:"grp,omitempty"`
	Dept  string `json:"dept,omitempty"`
	jwt.RegisteredClaims
}

type UserInfo struct {
	Username string
	Group    Group
	Dept     string
}

func (c *Claims) Token() (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(Secret)
}

func (c *Claims) MustToken() string {
	s, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(Secret)
	return s
}

func NewClaim(user *UserInfo) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        xid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		User:  user.Username,
		Group: int(user.Group),
		Dept:  user.Dept,
	}
}

func ParseToken(tok string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tok, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
