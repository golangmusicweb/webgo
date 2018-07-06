package utils

import (
	"hash"
	"golang.org/x/crypto/pbkdf2"
	"encoding/base64"
	"bytes"
	"math/rand"
)

type HashKey struct {
	Diggest func() hash.Hash
	SaltSize int
	KeyLen int
	Iterations int
}

type HashedResult struct {

CipherText string
Salt       string
}

func NewHashKey(diggest func() hash.Hash, saltSize int, keyLen int, iter int) *HashKey {
	return &HashKey{
		Diggest: diggest,
		SaltSize: saltSize,
		KeyLen: keyLen,
		Iterations: iter,
	}
}

func (p *HashKey) HashPassword(password string) HashedResult {
	saltBytes := make([]byte, p.SaltSize)
	rand.Read(saltBytes)
	saltString := base64.StdEncoding.EncodeToString(saltBytes)
	salt := bytes.NewBufferString(saltString).Bytes()
	df := pbkdf2.Key([]byte(password), salt, p.Iterations, p.KeyLen, p.Diggest)
	cipherText := base64.StdEncoding.EncodeToString(df)
	return HashedResult{CipherText: cipherText, Salt: saltString}
}

func (p *HashKey) VerifyPassword(password, cipherText, salt string) bool {
	saltBytes := bytes.NewBufferString(salt).Bytes()
	df := pbkdf2.Key([]byte(password), saltBytes, p.Iterations, p.KeyLen, p.Diggest)
	newCipherText := base64.StdEncoding.EncodeToString(df)
	valid := newCipherText == cipherText
	return valid
}