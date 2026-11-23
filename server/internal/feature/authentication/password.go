package authentication

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory  = 19 * 1024 // 19 MiB
	argonTime    = 2
	argonThreads = 1
	argonKeyLen  = 32
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// es el algoritmo que recomienda OWASP hoy
	// es lento y gasta memoria a proposito,
	// para que crackear hashes robados con GPUs sea dificil
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	b64 := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads,
		b64.EncodeToString(salt), b64.EncodeToString(key)), nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	var version int
	var memory, time uint32
	var threads uint8
	var saltB64, keyB64 string

	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid hash")
	}
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, err
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, err
	}
	saltB64, keyB64 = parts[4], parts[5]

	b64 := base64.RawStdEncoding
	salt, err := b64.DecodeString(saltB64)
	if err != nil {
		return false, err
	}
	want, err := b64.DecodeString(keyB64)
	if err != nil {
		return false, err
	}

	got := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(want)))
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}
