// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package jwt

// This file is a re-implementation of the original code here with some
// additional allocation tweaks reproduced using GODEBUG=allocfreetrace=1
// original file https://github.com/golang-jwt/jwt/blob/main/parser.go
// borrowed under MIT License https://github.com/golang-jwt/jwt/blob/main/LICENSE

import (
	"encoding/base64"
	"time"

	goJWT "github.com/golang-jwt/jwt/v4"
)

// Specific instances for HS256, HS384, HS512
var (
	SigningMethodHS256 = goJWT.SigningMethodHS256
	SigningMethodHS384 = goJWT.SigningMethodHS384
	SigningMethodHS512 = goJWT.SigningMethodHS512
)

// ECDSA - EC256 and company.
var (
	SigningMethodES256 = goJWT.SigningMethodES256
	SigningMethodES384 = goJWT.SigningMethodES384
	SigningMethodES512 = goJWT.SigningMethodES512
)

// StandardClaims are basically standard claims with "accessKey"
type StandardClaims struct {
	AccessKey string `json:"accessKey,omitempty"`
	IpAddr    string `json:"IpAddr,omitempty"`
	Terminal  string `json:"Terminal,omitempty"`
	goJWT.RegisteredClaims
}

// MapClaims - implements custom unmarshaller
type MapClaims struct {
	AccessKey string `json:"accessKey,omitempty"`
	IpAddr    string `json:"IpAddr,omitempty"`
	Terminal  string `json:"Terminal,omitempty"`
	goJWT.MapClaims
}

// GetAccessKey will return the access key.
// If nil an empty string will be returned.
func (c *MapClaims) GetAccessKey() string {
	if c == nil {
		return ""
	}
	return c.AccessKey
}

// newStandardClaims - initializes standard claims
func newStandardClaims() *StandardClaims {
	return &StandardClaims{}
}

// SetIssuer sets issuer for these claims
func (c *StandardClaims) SetIssuer(issuer string) {
	c.Issuer = issuer
}

// SetAudience sets audience for these claims
func (c *StandardClaims) SetAudience(aud goJWT.ClaimStrings) {
	c.Audience = aud
}

// SetExpiry sets expiry in unix epoch secs
func (c *StandardClaims) SetExpiry(t goJWT.NumericDate) {
	c.ExpiresAt = &t
}

// SetAccessKey sets access key as jwt subject and custom
// "accessKey" field.
func (c *StandardClaims) SetAccessKey(accessKey string) {
	c.AccessKey = accessKey
}

// SetIpAddr sets IpAddr for these claims
func (c *StandardClaims) SetIpAddr(ipaddr string) {
	c.IpAddr = ipaddr
}

// SetTerminal sets Terminal for these claims
func (c *StandardClaims) SetTerminal(terminal string) {
	c.Terminal = terminal
}

// Valid - implements https://godoc.org/github.com/golang-jwt/jwt#Claims compatible
// claims interface, additionally validates "accessKey" fields.
func (c *StandardClaims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}

	if c.AccessKey == "" && c.Subject == "" {
		return NewValidationError("accessKey/sub missing",
			ValidationErrorClaimsInvalid)
	}

	return c.RegisteredClaims.Valid()
}

// newMapClaims - Initializes a new map claims
func newMapClaims() *MapClaims {
	return &MapClaims{MapClaims: goJWT.MapClaims{}}
}

// Lookup returns the value and if the key is found.
func (c *MapClaims) Lookup(key string) (value string, ok bool) {
	if c == nil {
		return "", false
	}
	var vinterface interface{}
	vinterface, ok = c.MapClaims[key]
	if ok {
		value, ok = vinterface.(string)
	}
	return
}

// SetExpiry sets expiry in unix epoch secs
func (c *MapClaims) SetExpiry(t time.Time) {
	c.MapClaims["exp"] = t.Unix()
}

// SetAccessKey sets access key as jwt subject and custom
// "accessKey" field.
func (c *MapClaims) SetAccessKey(accessKey string) {
	c.MapClaims["sub"] = accessKey
	c.MapClaims["accessKey"] = accessKey
}

// Valid - implements https://godoc.org/github.com/golang-jwt/jwt#Claims compatible
// claims interface, additionally validates "accessKey" fields.
func (c *MapClaims) Valid() error {
	if err := c.MapClaims.Valid(); err != nil {
		return err
	}

	if c.AccessKey == "" {
		return NewValidationError("accessKey/sub missing", ValidationErrorClaimsInvalid)
	}

	return nil
}

// Map returns underlying low-level map claims.
func (c *MapClaims) Map() map[string]interface{} {
	if c == nil {
		return nil
	}
	return c.MapClaims
}

// ParseWithStandardClaims - parse the token string, valid methods.
func ParseWithStandardClaims(tokenStr string, claims *StandardClaims, key []byte) error {
	var (
		err error
	)
	// Key is not provided.
	if key == nil {
		// keyFunc was not provided, return error.
		return NewValidationError("no key was provided.", ValidationErrorUnverifiable)
	}

	if _, err = goJWT.ParseWithClaims(tokenStr, claims, func(token *goJWT.Token) (interface{}, error) {
		return key, nil
	}); err != nil {
		return NewValidationError("token parse fail.", ValidationErrorUnverifiable)
	}

	// Signature is valid, lets validate the claims for
	// other fields such as expiry etc.
	return nil
}

// ParseWithClaims - parse the token string, valid methods.
func ParseWithClaims(tokenStr string, claims *MapClaims, fn func(*MapClaims) ([]byte, error)) (*goJWT.Token, error) {
	var (
		token *goJWT.Token
		err   error
	)
	// Key lookup function has to be provided.
	if fn == nil {
		// keyFunc was not provided, return error.
		return nil, NewValidationError("no Keyfunc was provided.", ValidationErrorUnverifiable)
	}

	if token, err = goJWT.ParseWithClaims(tokenStr, claims, func(token *goJWT.Token) (interface{}, error) {
		return fn(claims)
	}); err != nil {
		return nil, NewValidationError("token parse fail.", ValidationErrorUnverifiable)
	}

	// Signature is valid, lets validate the claims for
	// other fields such as expiry etc.
	return token, nil
}

// base64DecodeBytes returns the bytes represented by the base64 string s.
func base64DecodeBytes(b []byte, buf []byte) (int, error) {
	return base64.RawURLEncoding.Decode(buf, b)
}
