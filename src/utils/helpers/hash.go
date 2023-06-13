package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func PasswordHash256(username string, password string) string {
	combine := strings.ToUpper(username) + password
	hash := []byte(combine)
	hash_byte := sha256.Sum256(hash)
	hash_str := hex.EncodeToString(hash_byte[:])
	return hash_str
}
