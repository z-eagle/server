package bcrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

type BCrypt struct {
	version  string
	strength int
}

const (
	BCryptVersion2A = "$2a"
	BCryptVersion2Y = "$2y"
	BCryptVersion2B = "$2b"
	base64_code     = "./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	MaxRounds       = 31
	MinRounds       = 4
	DefaultRounds   = 12
	SaltLen         = 16
	BlowfishRounds  = 16
)

var (
	InvalidRounds = errors.New("bcrypt: Invalid rounds parameter")
	InvalidSalt   = errors.New("bcrypt: Invalid salt supplied")
	NotFoundSalt  = errors.New("bcrypt: Salt not found")
)

var enc = base64.NewEncoding(base64_code)

// Helper function to build the bcrypt hash string
// payload takes :
//		* []byte -> which it base64 encodes it (trims padding "=") and writes it to the buffer
//		* string -> which it writes straight to the buffer
func buildBcryptStr(minor byte, rounds uint, payload ...interface{}) []byte {
	rs := bytes.NewBuffer(make([]byte, 0, 61))
	rs.WriteString("$2")
	if minor >= 'a' {
		rs.WriteByte(minor)
	}

	rs.WriteByte('$')
	if rounds < 10 {
		rs.WriteByte('0')
	}

	rs.WriteString(strconv.FormatUint(uint64(rounds), 10))
	rs.WriteByte('$')
	for _, p := range payload {
		if pb, ok := p.([]byte); ok {
			rs.WriteString(strings.TrimRight(enc.EncodeToString(pb), "="))
		} else if ps, ok := p.(string); ok {
			rs.WriteString(ps)
		}
	}
	return rs.Bytes()
}

// Salt generation
func genSalt(version string, strength int) (salt []byte, err error) {
	var r int
	r = strength
	if strength == -1 {
		r = DefaultRounds
	}

	if r < MinRounds || r > MaxRounds {
		return nil, InvalidRounds
	}

	rnd := make([]byte, SaltLen)
	read, err := rand.Read(rnd)
	if read != SaltLen || err != nil {
		return nil, err
	}
	return buildBcryptStr(version[2], uint(r), rnd), nil
}

func consume(r *bytes.Buffer, b byte) bool {
	got, err := r.ReadByte()
	if err != nil {
		return false
	}
	if got != b {
		_ = r.UnreadByte()
		return false
	}

	return true
}

func hash(password string, salt ...string) (ps string, err error) {
	var s []byte
	var pb []byte

	if len(salt) == 0 {
		return "", NotFoundSalt
	} else if len(salt) > 0 {
		s = []byte(salt[0])
	}

	pb, err = hashBytes([]byte(password), s)
	return string(pb), err
}

func hashBytes(password []byte, salt ...[]byte) (hash []byte, err error) {
	var s []byte

	if len(salt) == 0 {
		return nil, NotFoundSalt
	} else if len(salt) > 0 {
		s = salt[0]
	}

	// Ok, extract the required information
	minor := byte(0)
	sr := bytes.NewBuffer(s)

	if !consume(sr, '$') || !consume(sr, '2') {
		return nil, InvalidSalt
	}

	if !consume(sr, '$') {
		minor, _ = sr.ReadByte()
		if minor != 'a'  && minor != 'x' && minor != 'y' && minor != 'b'  || !consume(sr, '$') {
			return nil, InvalidSalt
		}
	}

	roundsBytes := make([]byte, 2)
	read, err := sr.Read(roundsBytes)
	if err != nil || read != 2 {
		return nil, InvalidSalt
	}

	if !consume(sr, '$') {
		return nil, InvalidSalt
	}

	var rounds64 uint64
	rounds64, err = strconv.ParseUint(string(roundsBytes), 10, 0)
	if err != nil {
		return nil, InvalidSalt
	}

	rounds := uint(rounds64)

	// TODO: can't we use base64.NewDecoder(enc, sr) ?
	saltBytes := make([]byte, 22)
	read, err = sr.Read(saltBytes)
	if err != nil || read != 22 {
		return nil, InvalidSalt
	}

	var saltB []byte
	// encoding/base64 expects 4 byte blocks padded, since bcrypt uses only 22 bytes we need to go up
	saltB, err = enc.DecodeString(string(saltBytes) + "==")
	if err != nil {
		return nil, err
	}

	// cipher expects null terminated input (go initializes everything with zero values so this works)
	passwordTerm := make([]byte, len(password)+1)
	copy(passwordTerm, password)

	hashed := cryptRaw(passwordTerm, saltB[:SaltLen], rounds)
	return buildBcryptStr(minor, rounds, string(saltBytes), hashed[:len(bfCryptCiphertext)*4-1]), nil
}

func (crypt *BCrypt) Matches(password, hash string) bool {
	return crypt.MatchesBytes([]byte(password), []byte(hash))
}

func (crypt *BCrypt) MatchesBytes(password []byte, hash []byte) bool {
	h, err := hashBytes(password, hash)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(h, hash) == 1
}

func (crypt *BCrypt) Encode(password string) string {
	salt, _ := genSalt(crypt.version, crypt.strength)
	ps, _ := hash(password, string(salt))
	return ps
}

func NewDefaultBcrypt() *BCrypt {
	return &BCrypt{BCryptVersion2A, DefaultRounds}
}

func NewBcrypt(version string, strength int) *BCrypt {
	return &BCrypt{version, strength}
}
