package server

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestKey(t *testing.T) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKey.Seed())
	fmt.Println(privateKeyBase64)
}
func TestKeyBase64(t *testing.T) {
	b64 := "R7xWhPNWejiOzxHPuiD2SRsvOwF81xWXcbxUJtXlG7A="
	seed, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)
	fmt.Println(publicKeyBase64)
}
