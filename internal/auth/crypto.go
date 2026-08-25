package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/box"
)

func SealTokenForCLI(publicKeyB64, token string) (string, error) {
	pubBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return "", err
	}
	if len(pubBytes) != 32 {
		return "", fmt.Errorf("public key harus 32 byte, dapat %d byte", len(pubBytes))
	}
	var recipientPub [32]byte
	copy(recipientPub[:], pubBytes)
	ephPub, ephPriv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}
	nonce, err := sealNonce(ephPub[:], recipientPub[:])
	if err != nil {
		return "", err
	}
	out := make([]byte, 0, 32+len(token)+box.Overhead)
	out = append(out, ephPub[:]...)
	out = box.Seal(out, []byte(token), &nonce, &recipientPub, ephPriv)
	return base64.StdEncoding.EncodeToString(out), nil
}

func sealNonce(ephPub, recipientPub []byte) ([24]byte, error) {
	var nonce [24]byte
	h, err := blake2b.New(24, nil)
	if err != nil {
		return nonce, err
	}
	h.Write(ephPub)
	h.Write(recipientPub)
	copy(nonce[:], h.Sum(nil))
	return nonce, nil
}
