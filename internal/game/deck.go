package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// NewDeck creates a standard 52-card deck
func NewDeck() []Card {
	deck := make([]Card, 0, 52)
	suits := []Suit{Clubs, Diamonds, Hearts, Spades}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}

	return deck
}

// Shuffle shuffles a deck using the provided random number generator
func Shuffle(deck []Card, rng *rand.Rand) {
	for i := len(deck) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

// Draw draws n cards from the deck and returns them along with the remaining deck
func Draw(deck []Card, n int) ([]Card, []Card) {
	if n > len(deck) {
		n = len(deck)
	}
	drawn := make([]Card, n)
	copy(drawn, deck[:n])
	remaining := deck[n:]
	return drawn, remaining
}

// LoadShoeFromFile loads a shoe from a file for deterministic testing
// File format: one card per line, e.g., "AS", "KD", "3H"
func LoadShoeFromFile(path string) ([]Card, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cards []Card
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		card, err := ParseCard(line)
		if err != nil {
			return nil, fmt.Errorf("invalid card in file: %s", line)
		}
		cards = append(cards, card)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}
