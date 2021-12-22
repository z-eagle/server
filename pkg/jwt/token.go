package jwt

import (
	goJWT "github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

const (
	Header = "x-token"
	Param  = "token"
	Cookie = "token"
	Prefix = "Bearer "
	Issuer = "ZhouQiaoKeJi"
	secret = "uidupQNPG1sBZZNA34U9eTgECs6BVfhAIOZtWi/BR0Y="
)

const (
	IPADDRESS = "ipaddr"
	TERMINAL  = "term"
)

type User struct {
	Id       uint64
	UserName string
	IpAddr   string
	Terminal string
}

func GenerateToken(user User) string {
	claims := newStandardClaims()
	claims.SetIssuer(Issuer)
	claims.SetIpAddr(user.IpAddr)
	claims.SetTerminal(user.Terminal)
	claims.RegisteredClaims = goJWT.RegisteredClaims{
		// A usual scenario is to set the expiration time relative to the current time
		ExpiresAt: goJWT.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  goJWT.NewNumericDate(time.Now()),
		NotBefore: goJWT.NewNumericDate(time.Now()),
		Subject:   user.UserName,
		Issuer:    Issuer,
		ID:        strconv.FormatUint(user.Id, 10),
		Audience:  []string{IPADDRESS, TERMINAL},
	}
	token := goJWT.NewWithClaims(SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func ValidateToken(t string, user User) (*User, error) {
	var (
		err    error
		claims = newStandardClaims()
	)
	err = ParseWithStandardClaims(t, claims, []byte(secret))
	if err != nil {
		return nil, err
	}

	if claims.Issuer != Issuer {
		return nil, NewValidationError("issuer not match",
			ValidationErrorIssuer)
	}

	if claims.Terminal != user.Terminal {
		return nil, NewValidationError("subject not match",
			ValidationErrorTerminal)
	}

	if err = claims.Valid(); err != nil {
		return nil, err
	}
	id, _ := strconv.Atoi(claims.ID)
	return &User{Id: uint64(id), UserName: claims.Subject}, nil
}
