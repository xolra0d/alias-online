package main

import (
	crand "crypto/rand"
	"encoding/base32"
	"encoding/base64"
	mrand "math/rand"
	"strconv"
	"strings"

	"go.pact.im/x/option"
	"go.pact.im/x/phcformat"
	"go.pact.im/x/phcformat/encode"
	"golang.org/x/crypto/argon2"
)

type Secrets struct {
	logger *PrefixLogger

	Argon2idTime    uint32
	Argon2idMemory  uint32
	Argon2idThreads uint8
	Argon2idOutLen  uint32
}

// hashSecret hashes any password with random salt and returns result in phcformat string.
// More: https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md
func (s *Secrets) hashSecret(secret string) string {
	salt := s.GenerateSecretBase32()
	output := argon2.IDKey([]byte(secret), []byte(salt), s.Argon2idTime, s.Argon2idMemory, s.Argon2idThreads, s.Argon2idOutLen)

	newParam := func(k string, v uint) encode.KV[encode.String, encode.Byte, encode.Uint] {
		return encode.NewKV(encode.NewByte('='), encode.NewString(k), encode.NewUint(v))
	}

	out := string(phcformat.Append(nil,
		encode.String("argon2id"),
		option.Value(encode.NewUint(uint(argon2.Version))),
		option.Value(encode.NewList(
			encode.NewByte(','),
			newParam("t", uint(s.Argon2idTime)),
			newParam("m", uint(s.Argon2idMemory)),
			newParam("p", uint(s.Argon2idThreads)),
			newParam("l", uint(s.Argon2idOutLen)),
		)),
		option.Value(encode.NewBase64(salt)),
		option.Value(encode.NewBase64(output)),
	))

	return out
}

// verifyArgon2id checks that secret is the same as encoded one in phcformat hash.
func (s *Secrets) verifyArgon2id(secret string, hash phcformat.Hash) bool {
	paramsStr, ok := hash.Params.Unwrap()
	if !ok {
		s.logger.Error("params for hashing algorithm are not set", "hash", hash.String())
		return false
	}
	params := make(map[string]uint32)
	for part := range strings.SplitSeq(paramsStr, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			s.logger.Error("params for hashing algorithm are not in key-value", "hash", hash.String())
			return false
		}

		v, err := strconv.ParseUint(kv[1], 10, 32)
		if err != nil {
			s.logger.Error("could not parse params for hashing algorithm", "hash", hash.String(), "name", kv[0], "val", kv[1])
			return false
		}
		params[kv[0]] = uint32(v)
	}

	t, ok := params["t"]
	if !ok {
		s.logger.Error("time is not set", "hash", hash.String())
		return false
	}
	m, ok := params["m"]
	if !ok {
		s.logger.Error("memory is not set", "hash", hash.String())
		return false
	}
	p, ok := params["p"]
	if !ok {
		s.logger.Error("processors is not set", "hash", hash.String())
		return false
	}
	l, ok := params["l"]
	if !ok {
		s.logger.Error("length is not set", "hash", hash.String())
		return false
	}
	salt, ok := hash.Salt.Unwrap()
	if !ok {
		s.logger.Error("salt is not set", "hash", hash.String())
		return false
	}
	saltDecoded, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		s.logger.Error("salt decode error", "hash", hash.String(), "err", err)
		return false
	}
	expected, ok := hash.Output.Unwrap()
	if !ok {
		s.logger.Error("output is not set", "hash", hash.String(), "err", err)
		return false
	}
	received := argon2.IDKey([]byte(secret), saltDecoded, t, m, uint8(p), l)

	if expected != base64.RawStdEncoding.EncodeToString(received) {
		return false
	}
	return true
}

// VerifyPassword checks if secret is equal to hash's secret.
func (s *Secrets) VerifyPassword(secret, hash string) bool {
	h, ok := phcformat.Parse(hash)
	if !ok {
		s.logger.Error("Could not decode phcformat hash", "hash", hash)
		return false
	}

	switch h.ID {
	case "argon2id":
		return s.verifyArgon2id(secret, h)
	default:
		s.logger.Warn("Invalid hash ID", "hash", hash)
		return false
	}
}

// GenerateRoomId creates new 40 bit base32 roomId
func (s *Secrets) GenerateRoomId() string {
	data := [5]byte{}
	_, _ = crand.Read(data[:]) // possible collision at ~1 million games.
	return base32.StdEncoding.EncodeToString(data[:])
}

// GenerateName creates a new name for account in form `AdjectiveNoun(0-99)`.
func (s *Secrets) GenerateName() string {
	adjectives := []string{"Grumpy", "Sleepy", "Chaotic", "Spicy", "Wobbly", "Fluffy", "Sneaky"}
	nouns := []string{"Waffle", "Penguin", "Muffin", "Wizard", "Noodle", "Taco", "Biscuit"}
	return adjectives[mrand.Intn(len(adjectives))] + nouns[mrand.Intn(len(nouns))] + strconv.Itoa(mrand.Intn(100))
}

// GenerateSecretBase32 creates secure base32 secret.
func (s *Secrets) GenerateSecretBase32() string {
	return crand.Text()
}
