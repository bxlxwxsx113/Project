package consensus

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(content []byte) string {
	hash := sha256.New()
	hash.Write(content)
	return hex.EncodeToString(hash.Sum(nil))
}
