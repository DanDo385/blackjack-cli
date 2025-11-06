package game

import (
	"testing"
)

func TestDeckHas52Unique(t *testing.T) {
	deck := NewDeck()

	if len(deck) != 52 {
		t.Errorf("Expected deck to have 52 cards, got %d", len(deck))
	}

	// Check for uniqueness
	seen := make(map[Card]bool)
	for _, card := range deck {
		if seen[card] {
			t.Errorf("Duplicate card found: %s", card)
		}
		seen[card] = true
	}
}

func TestShuffleDeterministicWithFixedSeed(t *testing.T) {
	// Create two decks with the same fixed seed
	rng1 := FixedSeededRand()
	rng2 := FixedSeededRand()

	deck1 := NewDeck()
	deck2 := NewDeck()

	Shuffle(deck1, rng1)
	Shuffle(deck2, rng2)

	// Both decks should be shuffled identically
	if len(deck1) != len(deck2) {
		t.Fatalf("Decks have different lengths")
	}

	for i := range deck1 {
		if deck1[i] != deck2[i] {
			t.Errorf("Card at index %d differs: %s vs %s", i, deck1[i], deck2[i])
		}
	}
}

func TestShuffleChangesOrder(t *testing.T) {
	original := NewDeck()
	deck := NewDeck()

	rng := FixedSeededRand()
	Shuffle(deck, rng)

	// After shuffle, at least some cards should be in different positions
	differences := 0
	for i := range deck {
		if deck[i] != original[i] {
			differences++
		}
	}

	if differences == 0 {
		t.Error("Shuffle did not change any card positions")
	}
}

func TestDraw(t *testing.T) {
	deck := NewDeck()

	// Draw 5 cards
	drawn, remaining := Draw(deck, 5)

	if len(drawn) != 5 {
		t.Errorf("Expected 5 drawn cards, got %d", len(drawn))
	}

	if len(remaining) != 47 {
		t.Errorf("Expected 47 remaining cards, got %d", len(remaining))
	}

	// Verify drawn cards match the first 5 of the original deck
	for i := 0; i < 5; i++ {
		if drawn[i] != deck[i] {
			t.Errorf("Drawn card %d doesn't match: %s vs %s", i, drawn[i], deck[i])
		}
	}
}

func TestDrawMoreThanAvailable(t *testing.T) {
	deck := NewDeck()

	// Try to draw more cards than available
	drawn, remaining := Draw(deck, 100)

	if len(drawn) != 52 {
		t.Errorf("Expected 52 drawn cards, got %d", len(drawn))
	}

	if len(remaining) != 0 {
		t.Errorf("Expected 0 remaining cards, got %d", len(remaining))
	}
}

func TestParseCard(t *testing.T) {
	tests := []struct {
		input    string
		expected Card
		wantErr  bool
	}{
		{"AS", Card{Rank: Ace, Suit: Spades}, false},
		{"KD", Card{Rank: King, Suit: Diamonds}, false},
		{"10H", Card{Rank: Ten, Suit: Hearts}, false},
		{"3C", Card{Rank: Three, Suit: Clubs}, false},
		{"", Card{}, true},
		{"X", Card{}, true},
		{"11S", Card{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			card, err := ParseCard(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %s: %v", tt.input, err)
				}
				if card != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, card)
				}
			}
		})
	}
}
