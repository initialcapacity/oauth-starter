package pkcesupport

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

// BASE64URL-ENCODE(SHA256(ASCII(code_verifier))) == code_challenge (pkce ietf 4.6)

func CodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	return safeEncode(h.Sum(nil))
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func CodeVerifier() string {
	var characters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, 64)
	for i := range b {
		b[i] = characters[seededRand.Intn(len(characters))]
	}
	return safeEncode(b)
}

func safeEncode(b []byte) string {
	encoded := base64.StdEncoding.EncodeToString(b)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}
