package game

import (
	"testing"
)

func TestPayouts(t *testing.T) {
	tests := []struct {
		name        string
		outcome     Outcome
		bet         int
		isInsurance bool
		want        int
	}{
		{
			name:        "Blackjack pays 3:2",
			outcome:     OutcomeBlackjack,
			bet:         100,
			isInsurance: false,
			want:        250, // 100 + 150
		},
		{
			name:        "Regular win pays 1:1",
			outcome:     OutcomeWin,
			bet:         100,
			isInsurance: false,
			want:        200, // 100 + 100
		},
		{
			name:        "Push returns bet",
			outcome:     OutcomePush,
			bet:         100,
			isInsurance: false,
			want:        100,
		},
		{
			name:        "Lose forfeits bet",
			outcome:     OutcomeLose,
			bet:         100,
			isInsurance: false,
			want:        0,
		},
		{
			name:        "Surrender returns half bet",
			outcome:     OutcomeSurrender,
			bet:         100,
			isInsurance: false,
			want:        50,
		},
		{
			name:        "Insurance win pays 2:1",
			outcome:     OutcomeWin,
			bet:         50,
			isInsurance: true,
			want:        100, // 50 * 2
		},
		{
			name:        "Insurance lose forfeits bet",
			outcome:     OutcomeLose,
			bet:         50,
			isInsurance: true,
			want:        -50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Payout(tt.outcome, tt.bet, tt.isInsurance)
			if got != tt.want {
				t.Errorf("Payout() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDetermineOutcome(t *testing.T) {
	tests := []struct {
		name        string
		playerCards []Card
		dealerCards []Card
		playerFlags map[string]bool
		want        Outcome
	}{
		{
			name:        "Player blackjack vs dealer 20",
			playerCards: []Card{{Rank: Ace, Suit: Spades}, {Rank: King, Suit: Hearts}},
			dealerCards: []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}},
			playerFlags: map[string]bool{"isInitialDeal": true},
			want:        OutcomeBlackjack,
		},
		{
			name:        "Player blackjack vs dealer blackjack",
			playerCards: []Card{{Rank: Ace, Suit: Spades}, {Rank: King, Suit: Hearts}},
			dealerCards: []Card{{Rank: Ace, Suit: Clubs}, {Rank: Queen, Suit: Hearts}},
			playerFlags: map[string]bool{"isInitialDeal": true},
			want:        OutcomePush,
		},
		{
			name:        "Player 20 vs dealer blackjack",
			playerCards: []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}},
			dealerCards: []Card{{Rank: Ace, Suit: Clubs}, {Rank: King, Suit: Clubs}},
			playerFlags: map[string]bool{"isInitialDeal": true},
			want:        OutcomeLose,
		},
		{
			name:        "Player 20 vs dealer 19",
			playerCards: []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}},
			dealerCards: []Card{{Rank: King, Suit: Clubs}, {Rank: Nine, Suit: Clubs}},
			playerFlags: map[string]bool{"isInitialDeal": true},
			want:        OutcomeWin,
		},
		{
			name:        "Player 20 vs dealer 20",
			playerCards: []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}},
			dealerCards: []Card{{Rank: King, Suit: Clubs}, {Rank: Queen, Suit: Clubs}},
			playerFlags: map[string]bool{"isInitialDeal": true},
			want:        OutcomePush,
		},
		{
			name:        "Player bust vs dealer 20",
			playerCards: []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}, {Rank: Five, Suit: Clubs}},
			dealerCards: []Card{{Rank: King, Suit: Clubs}, {Rank: Queen, Suit: Clubs}},
			playerFlags: map[string]bool{"isInitialDeal": false},
			want:        OutcomeLose,
		},
		{
			name:        "Player 20 vs dealer bust",
			playerCards: []Card{{Rank: King, Suit: Spades}, {Rank: Queen, Suit: Hearts}},
			dealerCards: []Card{{Rank: King, Suit: Clubs}, {Rank: Queen, Suit: Clubs}, {Rank: Five, Suit: Clubs}},
			playerFlags: map[string]bool{"isInitialDeal": true},
			want:        OutcomeWin,
		},
		{
			name:        "Player surrendered",
			playerCards: []Card{{Rank: Ten, Suit: Spades}, {Rank: Six, Suit: Hearts}},
			dealerCards: []Card{{Rank: King, Suit: Clubs}, {Rank: Queen, Suit: Clubs}},
			playerFlags: map[string]bool{"surrendered": true},
			want:        OutcomeSurrender,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playerHand := NewHand(100)
			for _, card := range tt.playerCards {
				playerHand.Add(card)
			}
			if v, ok := tt.playerFlags["isInitialDeal"]; ok {
				playerHand.IsInitialDeal = v
			}
			if v, ok := tt.playerFlags["surrendered"]; ok {
				playerHand.Surrendered = v
			}

			dealerHand := NewHand(0)
			dealerHand.IsInitialDeal = true
			for _, card := range tt.dealerCards {
				dealerHand.Add(card)
			}

			got := DetermineOutcome(playerHand, dealerHand)
			if got != tt.want {
				t.Errorf("DetermineOutcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeekForBlackjack(t *testing.T) {
	tests := []struct {
		name    string
		upCard  Card
		downCard Card
		want    bool
	}{
		{
			name:    "Ace and King",
			upCard:  Card{Rank: Ace, Suit: Spades},
			downCard: Card{Rank: King, Suit: Hearts},
			want:    true,
		},
		{
			name:    "King and Ace",
			upCard:  Card{Rank: King, Suit: Spades},
			downCard: Card{Rank: Ace, Suit: Hearts},
			want:    true,
		},
		{
			name:    "Ten and Ace",
			upCard:  Card{Rank: Ten, Suit: Spades},
			downCard: Card{Rank: Ace, Suit: Hearts},
			want:    true,
		},
		{
			name:    "Ace and Nine (not blackjack)",
			upCard:  Card{Rank: Ace, Suit: Spades},
			downCard: Card{Rank: Nine, Suit: Hearts},
			want:    false,
		},
		{
			name:    "King and Queen (not blackjack)",
			upCard:  Card{Rank: King, Suit: Spades},
			downCard: Card{Rank: Queen, Suit: Hearts},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PeekForBlackjack(tt.upCard, tt.downCard)
			if got != tt.want {
				t.Errorf("PeekForBlackjack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealerS17(t *testing.T) {
	tests := []struct {
		name        string
		initialCards []Card
		deckCards   []Card
		wantFinalValue int
		wantFinalCards int
	}{
		{
			name:        "Dealer stands on hard 17",
			initialCards: []Card{{Rank: Ten, Suit: Spades}, {Rank: Seven, Suit: Hearts}},
			deckCards:   []Card{{Rank: Five, Suit: Clubs}},
			wantFinalValue: 17,
			wantFinalCards: 2,
		},
		{
			name:        "Dealer stands on soft 17 (A+6)",
			initialCards: []Card{{Rank: Ace, Suit: Spades}, {Rank: Six, Suit: Hearts}},
			deckCards:   []Card{{Rank: Five, Suit: Clubs}},
			wantFinalValue: 17,
			wantFinalCards: 2,
		},
		{
			name:        "Dealer hits on 16",
			initialCards: []Card{{Rank: Ten, Suit: Spades}, {Rank: Six, Suit: Hearts}},
			deckCards:   []Card{{Rank: Five, Suit: Clubs}},
			wantFinalValue: 21,
			wantFinalCards: 3,
		},
		{
			name:        "Dealer hits on soft 16",
			initialCards: []Card{{Rank: Ace, Suit: Spades}, {Rank: Five, Suit: Hearts}},
			deckCards:   []Card{{Rank: Five, Suit: Clubs}},
			wantFinalValue: 21,
			wantFinalCards: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := NewHand(0)
			for _, card := range tt.initialCards {
				hand.Add(card)
			}

			deck := tt.deckCards
			DealerPlay(&deck, hand)

			if hand.Value() != tt.wantFinalValue {
				t.Errorf("Final value = %d, want %d", hand.Value(), tt.wantFinalValue)
			}

			if len(hand.Cards) != tt.wantFinalCards {
				t.Errorf("Final card count = %d, want %d", len(hand.Cards), tt.wantFinalCards)
			}
		})
	}
}
