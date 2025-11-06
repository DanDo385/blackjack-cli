package game

import (
	"testing"
)

func TestAceSoftHardTotals(t *testing.T) {
	tests := []struct {
		name       string
		cards      []Card
		wantHard   int
		wantSoft   int
		wantIsSoft bool
		wantValue  int
	}{
		{
			name:       "Ace and 6 (soft 17)",
			cards:      []Card{{Rank: Ace, Suit: Spades}, {Rank: Six, Suit: Hearts}},
			wantHard:   7,
			wantSoft:   17,
			wantIsSoft: true,
			wantValue:  17,
		},
		{
			name:       "Ace and 5 (soft 16)",
			cards:      []Card{{Rank: Ace, Suit: Spades}, {Rank: Five, Suit: Hearts}},
			wantHard:   6,
			wantSoft:   16,
			wantIsSoft: true,
			wantValue:  16,
		},
		{
			name:       "Ace, 5, and 10 (hard 16)",
			cards:      []Card{{Rank: Ace, Suit: Spades}, {Rank: Five, Suit: Hearts}, {Rank: Ten, Suit: Clubs}},
			wantHard:   16,
			wantSoft:   16,
			wantIsSoft: false,
			wantValue:  16,
		},
		{
			name:       "Two Aces (soft 12)",
			cards:      []Card{{Rank: Ace, Suit: Spades}, {Rank: Ace, Suit: Hearts}},
			wantHard:   2,
			wantSoft:   12,
			wantIsSoft: true,
			wantValue:  12,
		},
		{
			name:       "King and Queen (20)",
			cards:      []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}},
			wantHard:   20,
			wantSoft:   20,
			wantIsSoft: false,
			wantValue:  20,
		},
		{
			name:       "Ace and King (blackjack)",
			cards:      []Card{{Rank: Ace, Suit: Spades}, {Rank: King, Suit: Hearts}},
			wantHard:   11,
			wantSoft:   21,
			wantIsSoft: true,
			wantValue:  21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := NewHand(10)
			for _, card := range tt.cards {
				hand.Add(card)
			}

			hard, soft, isSoft := hand.Totals()

			if hard != tt.wantHard {
				t.Errorf("Hard total = %d, want %d", hard, tt.wantHard)
			}
			if soft != tt.wantSoft {
				t.Errorf("Soft total = %d, want %d", soft, tt.wantSoft)
			}
			if isSoft != tt.wantIsSoft {
				t.Errorf("IsSoft = %v, want %v", isSoft, tt.wantIsSoft)
			}
			if hand.Value() != tt.wantValue {
				t.Errorf("Value = %d, want %d", hand.Value(), tt.wantValue)
			}
		})
	}
}

func TestBlackjackDetection(t *testing.T) {
	tests := []struct {
		name            string
		cards           []Card
		isInitialDeal   bool
		wantBlackjack   bool
	}{
		{
			name:          "Ace and King (blackjack)",
			cards:         []Card{{Rank: Ace, Suit: Spades}, {Rank: King, Suit: Hearts}},
			isInitialDeal: true,
			wantBlackjack: true,
		},
		{
			name:          "Ten and Ace (blackjack)",
			cards:         []Card{{Rank: Ten, Suit: Spades}, {Rank: Ace, Suit: Hearts}},
			isInitialDeal: true,
			wantBlackjack: true,
		},
		{
			name:          "21 in three cards (not blackjack)",
			cards:         []Card{{Rank: Seven, Suit: Spades}, {Rank: Seven, Suit: Hearts}, {Rank: Seven, Suit: Clubs}},
			isInitialDeal: true,
			wantBlackjack: false,
		},
		{
			name:          "Ace and King after hit (not blackjack)",
			cards:         []Card{{Rank: Ace, Suit: Spades}, {Rank: King, Suit: Hearts}},
			isInitialDeal: false,
			wantBlackjack: false,
		},
		{
			name:          "Ten and Nine (not blackjack)",
			cards:         []Card{{Rank: Ten, Suit: Spades}, {Rank: Nine, Suit: Hearts}},
			isInitialDeal: true,
			wantBlackjack: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := NewHand(10)
			hand.IsInitialDeal = tt.isInitialDeal
			for _, card := range tt.cards {
				hand.Add(card)
			}

			if hand.IsBlackjack() != tt.wantBlackjack {
				t.Errorf("IsBlackjack() = %v, want %v", hand.IsBlackjack(), tt.wantBlackjack)
			}
		})
	}
}

func TestSplitRules(t *testing.T) {
	tests := []struct {
		name          string
		cards         []Card
		isInitialDeal bool
		wantCanSplit  bool
	}{
		{
			name:          "Pair of 8s",
			cards:         []Card{{Rank: Eight, Suit: Spades}, {Rank: Eight, Suit: Hearts}},
			isInitialDeal: true,
			wantCanSplit:  true,
		},
		{
			name:          "Pair of Aces",
			cards:         []Card{{Rank: Ace, Suit: Spades}, {Rank: Ace, Suit: Hearts}},
			isInitialDeal: true,
			wantCanSplit:  true,
		},
		{
			name:          "King and Queen (both 10 value but different ranks)",
			cards:         []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}},
			isInitialDeal: true,
			wantCanSplit:  false,
		},
		{
			name:          "Three cards",
			cards:         []Card{{Rank: Eight, Suit: Spades}, {Rank: Eight, Suit: Hearts}, {Rank: Five, Suit: Clubs}},
			isInitialDeal: true,
			wantCanSplit:  false,
		},
		{
			name:          "Pair of 8s after hit",
			cards:         []Card{{Rank: Eight, Suit: Spades}, {Rank: Eight, Suit: Hearts}},
			isInitialDeal: false,
			wantCanSplit:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := NewHand(10)
			hand.IsInitialDeal = tt.isInitialDeal
			for _, card := range tt.cards {
				hand.Add(card)
			}

			if hand.CanSplit() != tt.wantCanSplit {
				t.Errorf("CanSplit() = %v, want %v", hand.CanSplit(), tt.wantCanSplit)
			}
		})
	}
}

func TestDoubleRules(t *testing.T) {
	tests := []struct {
		name          string
		isInitialDeal bool
		isSplitAces   bool
		wantCanDouble bool
	}{
		{
			name:          "Initial deal, not split aces",
			isInitialDeal: true,
			isSplitAces:   false,
			wantCanDouble: true,
		},
		{
			name:          "After hit",
			isInitialDeal: false,
			isSplitAces:   false,
			wantCanDouble: false,
		},
		{
			name:          "Split aces",
			isInitialDeal: true,
			isSplitAces:   true,
			wantCanDouble: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := NewHand(10)
			hand.IsInitialDeal = tt.isInitialDeal
			hand.IsSplitAces = tt.isSplitAces

			if hand.CanDouble() != tt.wantCanDouble {
				t.Errorf("CanDouble() = %v, want %v", hand.CanDouble(), tt.wantCanDouble)
			}
		})
	}
}

func TestSurrenderLate(t *testing.T) {
	tests := []struct {
		name             string
		cards            []Card
		isInitialDeal    bool
		wantCanSurrender bool
	}{
		{
			name:             "Initial two cards",
			cards:            []Card{{Rank: Ten, Suit: Spades}, {Rank: Six, Suit: Hearts}},
			isInitialDeal:    true,
			wantCanSurrender: true,
		},
		{
			name:             "After hit",
			cards:            []Card{{Rank: Ten, Suit: Spades}, {Rank: Six, Suit: Hearts}},
			isInitialDeal:    false,
			wantCanSurrender: false,
		},
		{
			name:             "Three cards",
			cards:            []Card{{Rank: Ten, Suit: Spades}, {Rank: Six, Suit: Hearts}, {Rank: Five, Suit: Clubs}},
			isInitialDeal:    true,
			wantCanSurrender: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := NewHand(10)
			hand.IsInitialDeal = tt.isInitialDeal
			for _, card := range tt.cards {
				hand.Add(card)
			}

			if hand.CanSurrender() != tt.wantCanSurrender {
				t.Errorf("CanSurrender() = %v, want %v", hand.CanSurrender(), tt.wantCanSurrender)
			}
		})
	}
}

func TestBust(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		wantBust bool
	}{
		{
			name:     "20 (not bust)",
			cards:    []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}},
			wantBust: false,
		},
		{
			name:     "21 (not bust)",
			cards:    []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}, {Rank: Ace, Suit: Clubs}},
			wantBust: false,
		},
		{
			name:     "22 (bust)",
			cards:    []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}, {Rank: Two, Suit: Clubs}},
			wantBust: true,
		},
		{
			name:     "26 (bust)",
			cards:    []Card{{Rank: King, Suit: Spades}, {Rank: King, Suit: Hearts}, {Rank: Six, Suit: Clubs}},
			wantBust: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := NewHand(10)
			for _, card := range tt.cards {
				hand.Add(card)
			}

			if hand.IsBust() != tt.wantBust {
				t.Errorf("IsBust() = %v, want %v (value: %d)", hand.IsBust(), tt.wantBust, hand.Value())
			}
		})
	}
}
