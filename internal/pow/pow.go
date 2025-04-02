package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Challenge struct {
	ID         string    `json:"id"`
	Data       string    `json:"data"`
	Difficulty int       `json:"difficulty"`
	Timestamp  time.Time `json:"timestamp"`
	TTL        int       `json:"ttl"` // seconds
}

func GenerateChallenge(id string, difficulty, ttl int) Challenge {
	data := fmt.Sprintf("%d", time.Now().UnixNano())
	return Challenge{
		ID:         id,
		Data:       data,
		Difficulty: difficulty,
		Timestamp:  time.Now(),
		TTL:        ttl,
	}
}

func (c Challenge) Verify(nonce string) bool {
	hash := sha256.Sum256([]byte(c.Data + nonce))
	hashStr := hex.EncodeToString(hash[:])
	prefix := ""
	for i := 0; i < c.Difficulty; i++ {
		prefix += "0"
	}
	return hashStr[:c.Difficulty] == prefix
}
