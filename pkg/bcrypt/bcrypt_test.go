//
// Copyright 2011 ZooWar.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package bcrypt

import (
	"reflect"
	"testing"
)

func checkError(t *testing.T, i int, err, terr error) {
	if err != terr {
		t.Errorf("test(%d): err: %v expected: %v", i, err, terr)
	}
}

type HashTest struct {
	pw   []byte
	hash []byte
	err  error
}

var hashTests []HashTest = []HashTest{
	{[]byte(nil), []byte(nil), InvalidSalt},
	{[]byte(""), []byte(nil), InvalidSalt},
	{[]byte(" "), []byte(nil), InvalidSalt},
	{[]byte("testing"), []byte(nil), InvalidSalt},
	{[]byte("testing"), []byte(""), InvalidSalt},
}

var testBCrypt = NewDefaultBcrypt()

func TestHashBytes(t *testing.T) {
	bad := []byte("saltine")
	for i, test := range hashTests {
		pw := test.pw
		hash, err := hashBytes(pw, test.hash)
		checkError(t, i, err, test.err)
		if err != nil {
			continue
		}
		if !testBCrypt.MatchesBytes(pw, hash) {
			t.Errorf("test(%d): b", i)
		}
		if testBCrypt.MatchesBytes(bad, hash) {
			t.Errorf("test(%d): c", i)
		}
	}

	salt, _ := genSalt(testBCrypt.version, 10)
	h1, _ := hashBytes(nil, salt)
	h2, _ := hashBytes(nil, salt)
	if !reflect.DeepEqual(h1, h2) {
		t.Errorf("salt random")
	}
}

type SaltTest struct {
	rounds int
	err    error
}

var saltTests []SaltTest = []SaltTest{
	{-1, nil},
	{MinRounds - 1, InvalidRounds},
	{MinRounds, nil},
	{MaxRounds, nil},
	{MaxRounds + 1, InvalidRounds},
}

func TestSalt(t *testing.T) {
	for i, test := range saltTests {
		_, err := genSalt(testBCrypt.version, test.rounds)
		checkError(t, i, err, test.err)
	}
	_, err := genSalt(testBCrypt.version, testBCrypt.strength)
	if err != nil {
		t.Errorf("salt variadic: %v\n", err)
	}
	_, err = genSalt(testBCrypt.version, DefaultRounds)
	if err != nil {
		t.Errorf("salt variadic: %v\n", err)
	}
}

type TestString struct {
	plain string
	salt  string
	hash  string
}

var testStrings []TestString = []TestString{
	{"", "$2a$06$DCq7YPn5Rq63x1Lad4cll.",
		"$2a$06$DCq7YPn5Rq63x1Lad4cll.TV4S6ytwfsfvkgY8jIucDrjc8deX1s."},
	{"", "$2a$08$HqWuK6/Ng6sg9gQzbLrgb.",
		"$2a$08$HqWuK6/Ng6sg9gQzbLrgb.Tl.ZHfXLhvt/SgVyWhQqgqcZ7ZuUtye"},
	{"", "$2a$10$k1wbIrmNyFAPwPVPSVa/ze",
		"$2a$10$k1wbIrmNyFAPwPVPSVa/zecw2BCEnBwVS2GbrmgzxFUOqW9dk4TCW"},
	{"", "$2a$12$k42ZFHFWqBp3vWli.nIn8u",
		"$2a$12$k42ZFHFWqBp3vWli.nIn8uYyIkbvYRvodzbfbK18SSsY.CsIQPlxO"},
	{"a", "$2a$06$m0CrhHm10qJ3lXRY.5zDGO",
		"$2a$06$m0CrhHm10qJ3lXRY.5zDGO3rS2KdeeWLuGmsfGlMfOxih58VYVfxe"},
	{"a", "$2a$08$cfcvVd2aQ8CMvoMpP2EBfe",
		"$2a$08$cfcvVd2aQ8CMvoMpP2EBfeodLEkkFJ9umNEfPD18.hUF62qqlC/V."},
	{"a", "$2a$10$k87L/MF28Q673VKh8/cPi.",
		"$2a$10$k87L/MF28Q673VKh8/cPi.SUl7MU/rWuSiIDDFayrKk/1tBsSQu4u"},
	{"a", "$2a$12$8NJH3LsPrANStV6XtBakCe",
		"$2a$12$8NJH3LsPrANStV6XtBakCez0cKHXVxmvxIlcz785vxAIZrihHZpeS"},
	{"abc", "$2a$06$If6bvum7DFjUnE9p2uDeDu",
		"$2a$06$If6bvum7DFjUnE9p2uDeDu0YHzrHM6tf.iqN8.yx.jNN1ILEf7h0i"},
	{"abc", "$2a$08$Ro0CUfOqk6cXEKf3dyaM7O",
		"$2a$08$Ro0CUfOqk6cXEKf3dyaM7OhSCvnwM9s4wIX9JeLapehKK5YdLxKcm"},
	{"abc", "$2a$10$WvvTPHKwdBJ3uk0Z37EMR.",
		"$2a$10$WvvTPHKwdBJ3uk0Z37EMR.hLA2W6N9AEBhEgrAOljy2Ae5MtaSIUi"},
	{"abc", "$2a$12$EXRkfkdmXn2gzds2SSitu.",
		"$2a$12$EXRkfkdmXn2gzds2SSitu.MW9.gAVqa9eLS1//RYtYCmB1eLHg.9q"},
	{"abcdefghijklmnopqrstuvwxyz", "$2a$06$.rCVZVOThsIa97pEDOxvGu",
		"$2a$06$.rCVZVOThsIa97pEDOxvGuRRgzG64bvtJ0938xuqzv18d3ZpQhstC"},
	{"abcdefghijklmnopqrstuvwxyz", "$2a$08$aTsUwsyowQuzRrDqFflhge",
		"$2a$08$aTsUwsyowQuzRrDqFflhgekJ8d9/7Z3GV3UcgvzQW3J5zMyrTvlz."},
	{"abcdefghijklmnopqrstuvwxyz", "$2a$10$fVH8e28OQRj9tqiDXs1e1u",
		"$2a$10$fVH8e28OQRj9tqiDXs1e1uxpsjN0c7II7YPKXua2NAKYvM6iQk7dq"},
	{"abcdefghijklmnopqrstuvwxyz", "$2a$12$D4G5f18o7aMMfwasBL7Gpu",
		"$2a$12$D4G5f18o7aMMfwasBL7GpuQWuP3pkrZrOAnqP.bmezbMng.QwJ/pG"},
	{"~!@#$%^&*()      ~!@#$%^&*()PNBFRD", "$2a$06$fPIsBO8qRqkjj273rfaOI.",
		"$2a$06$fPIsBO8qRqkjj273rfaOI.HtSV9jLDpTbZn782DC6/t7qT67P6FfO"},
	{"~!@#$%^&*()      ~!@#$%^&*()PNBFRD", "$2a$08$Eq2r4G/76Wv39MzSX262hu",
		"$2a$08$Eq2r4G/76Wv39MzSX262huzPz612MZiYHVUJe/OcOql2jo4.9UxTW"},
	{"~!@#$%^&*()      ~!@#$%^&*()PNBFRD", "$2a$10$LgfYWkbzEvQ4JakH7rOvHe",
		"$2a$10$LgfYWkbzEvQ4JakH7rOvHe0y8pHKF9OaFgwUZ2q7W2FFZmZzJYlfS"},
	{"~!@#$%^&*()      ~!@#$%^&*()PNBFRD", "$2a$12$WApznUOJfkEGSmYRfnkrPO",
		"$2a$12$WApznUOJfkEGSmYRfnkrPOr466oFDCaj4b6HY3EXGvfxm43seyhgC"},
}

func TestStrings(t *testing.T) {
	for i, test := range testStrings {
		hash, err := hash(test.plain, test.salt)
		if err != nil {
			t.Errorf("hash(%d): %v", i, err)
			continue
		}
		if hash != test.hash {
			t.Errorf("test(%d): equal: %v", i, hash)
			t.Errorf("test(%d): equal: %v", i, test.hash)
		}
	}
	for _, r := range []int{4, 8, 14} {
		test := testStrings[r]
		salt, err := genSalt(testBCrypt.version, r)
		hash, err := hash(test.plain, string(salt))
		if !testBCrypt.Matches(test.plain, hash) {
			t.Errorf("Rounds(%d): %v", r, err)
		}
	}
}

func TestEncode(t *testing.T) {
	for _, r := range []int{4, 8, 14} {
		test := testStrings[r]
		testBC := NewBcrypt(BCryptVersion2A, r)
		hash := testBC.Encode(test.plain)
		if !testBC.Matches(test.plain, hash) {
			t.Errorf("Rounds(%d)", r)
		}
	}
}

type BytesType struct {
	plain []byte
	salt  []byte
	hash  []byte
}

func TestBytes(t *testing.T) {
	testBytes := make([]BytesType, len(testStrings))
	for i, v := range testStrings {
		testBytes[i] = BytesType{[]byte(v.plain), []byte(v.salt), []byte(v.hash)}
	}
	for i, test := range testBytes {
		hash, err := hashBytes(test.plain, test.salt)
		if err != nil {
			t.Errorf("hash(%d):%v", i, err)
			continue
		}
		if !reflect.DeepEqual(hash, test.hash) {
			t.Errorf("test(%d): equal: %v", i, hash)
			t.Errorf("test(%d): equal: %v", i, test.hash)
		}
	}
	for _, r := range []int{4, 8, 14} {
		test := testBytes[r]
		salt, err := genSalt(testBCrypt.version, r)
		hash, err := hashBytes(test.plain, salt)
		if !testBCrypt.MatchesBytes(test.plain, hash) {
			t.Errorf("Rounds(%d): %v", r, err)
		}
	}
}
