package server

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
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
func TestEd25519Sign(t *testing.T) {
	b64 := "R7xWhPNWejiOzxHPuiD2SRsvOwF81xWXcbxUJtXlG7A="
	seed, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	msg := []byte("hello world")
	sig := ed25519.Sign(privateKey, msg)
	if !ed25519.Verify(publicKey, msg, sig) {
		t.Errorf("ed25519 verify failed")
	}
}
func TestEd25519SignJWT(t *testing.T) {
	b64 := "R7xWhPNWejiOzxHPuiD2SRsvOwF81xWXcbxUJtXlG7A="
	seed, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.EdDSA,
		Key:       privateKey,
	}, nil)
	if err != nil {
		panic(err)
	}
	claims := jwt.Claims{
		Subject:  "12132312312",
		Expiry:   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token, err := jwt.Signed(signer).Claims(claims).Serialize()
	if err != nil {
		panic(err)
	}
	t.Log(token)
	tok, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.EdDSA})
	if err != nil {
		panic(err)
	}
	var newClaims jwt.Claims
	if err = tok.Claims(publicKey, &newClaims); err != nil {
		panic(err)
	}
	t.Log(newClaims)
	if newClaims.Subject != claims.Subject {
		t.Errorf("claims subject not equal")
	}
}
