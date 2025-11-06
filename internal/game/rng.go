package game

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"os"
)

// CryptoSeededRand creates a new rand.Rand seeded with crypto/rand
func CryptoSeededRand() *rand.Rand {
	var seed int64
	binary.Read(cryptorand.Reader, binary.BigEndian, &seed)
	return rand.New(rand.NewSource(seed))
}

// FixedSeededRand creates a new rand.Rand with a fixed seed for testing
func FixedSeededRand() *rand.Rand {
	return rand.New(rand.NewSource(12345))
}

// NewRand creates a new rand.Rand based on the BLACKJACK_SEEDED environment variable
func NewRand() *rand.Rand {
	if os.Getenv("BLACKJACK_SEEDED") == "1" {
		return FixedSeededRand()
	}
	return CryptoSeededRand()
}
