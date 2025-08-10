package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
)

var (
	key = []byte("this is the secret key and stuff")
)

// CreateToken is used for creating a new token, it does not have the ID or Token fields because those are generated
// and will be part of the Token
type CreateToken struct {
	Payload   string         `json:"payload"`
	TokenType string         `json:"token_type"`
	TTL       int64          `json:"ttl"`
	Metadata  map[string]any `json:"metadata"`
}

type Token struct {
	BaseModel
	CreateToken
	Token string `json:"token"`
}

func (t *Token) Encrypt() error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	ciphertext := make([]byte, 0, len(t.Payload))
	ciphertext = gcm.Seal(ciphertext, nonce, []byte(t.Payload), nil)
	enc := hex.EncodeToString(ciphertext)
	t.Payload = enc
	return nil
}

func (t *Token) Decrypt() (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	cipherText, err := hex.DecodeString(t.Payload)
	if err != nil {
		return "", err
	}
	decryptedData := make([]byte, 0, len(cipherText))
	decryptedData, err = gcm.Open(decryptedData, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}

func (t *Token) Tokenize() error {
	h := sha512.New512_256()
	h.Write([]byte(t.Payload))
	t.Token = hex.EncodeToString(h.Sum(nil))
	return nil
}
