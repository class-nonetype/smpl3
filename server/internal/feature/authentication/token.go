package authentication

import (
	"crypto/rand"
	"crypto/sha256"
)

func newToken() (token string, hash []byte) {
	// genera un string aleatorio criptograficamente seguro de 26 caracteres
	// con 128 bits de entropia
	token = rand.Text()
	return token, hashToken(token)
}

func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}
