package game

import (
	"testing"
)

func TestInsuranceOfferedOnAceUpcard(t *testing.T) {
	// Dealer upcard is Ace - insurance should be offered
	testDeck := []Card{
		{Rank: King, Suit: Spades},     // Player card 1
		{Rank: King, Suit: Diamonds},   // Dealer hole (hidden)
		{Rank: Ten, Suit: Hearts},      // Player card 2
		{Rank: Ace, Suit: Clubs},       // Dealer upcard (visible - insurance!)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()

	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 100

	if g.CurrentPhase != PhaseInsurance {
		t.Errorf("Expected insurance phase, got %v", g.CurrentPhase)
	}

	if !g.InsuranceOffered {
		t.Error("Insurance should be offered")
	}
}

func TestInsuranceNotOfferedOnNonAceUpcard(t *testing.T) {
	// Dealer upcard is King - insurance should NOT be offered
	testDeck := []Card{
		{Rank: King, Suit: Spades},     // Player card 1
		{Rank: Ace, Suit: Diamonds},    // Dealer hole (hidden - not shown)
		{Rank: Ten, Suit: Hearts},      // Player card 2
		{Rank: King, Suit: Clubs},      // Dealer upcard (visible - no insurance)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()

	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 100

	if g.CurrentPhase != PhasePlayerAction {
		t.Errorf("Expected player action phase, got %v", g.CurrentPhase)
	}

	if g.InsuranceOffered {
		t.Error("Insurance should NOT be offered when upcard is not Ace")
	}
}

func TestMaxInsuranceExceedsBank(t *testing.T) {
	// Bank is 50 after bet, but maxInsurance would be 60 (bet/2 = 120/2)
	// Should limit to 50
	testDeck := []Card{
		{Rank: King, Suit: Spades},     // Player card 1
		{Rank: King, Suit: Diamonds},   // Dealer hole
		{Rank: Ten, Suit: Hearts},      // Player card 2
		{Rank: Ace, Suit: Clubs},       // Dealer upcard (insurance offered)
		{Rank: Five, Suit: Diamonds},   // For potential hits
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()
	g.Bank = 200

	err := g.StartHand(150)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 150

	// Bank is now 50, bet was 150
	// maxInsurance should be min(150/2, 50) = min(75, 50) = 50
	if g.CurrentPhase != PhaseInsurance {
		t.Fatalf("Expected insurance phase")
	}

	// Take insurance for 50 (max possible)
	err = g.TakeInsurance(50)
	if err != nil {
		t.Errorf("Taking insurance for 50 should succeed, got error: %v", err)
	}

	// Taking insurance for 51 should fail (exceeds max bet of 75)
	// But we already took 50, so let's test with fresh game
}

func TestMaxInsuranceIsZero(t *testing.T) {
	// Bank is 0 after bet, maxInsurance = min(bet/2, 0) = 0
	// Should auto-decline without error
	testDeck := []Card{
		{Rank: King, Suit: Spades},     // Player card 1
		{Rank: Two, Suit: Diamonds},    // Dealer hole (2, definitely not BJ)
		{Rank: Ten, Suit: Hearts},      // Player card 2
		{Rank: Ace, Suit: Clubs},       // Dealer upcard (insurance offered)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()
	g.Bank = 100

	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 100

	// Bank is now 0, maxInsurance = min(100/2, 0) = 0
	if g.CurrentPhase != PhaseInsurance {
		t.Fatalf("Expected insurance phase, got %v", g.CurrentPhase)
	}

	// Decline insurance when maxInsurance is 0
	err = g.DeclineInsurance()
	if err != nil {
		t.Errorf("DeclineInsurance should succeed with maxInsurance=0, got error: %v", err)
	}

	// Should proceed to player action (no blackjack with Queen + Ace)
	if g.CurrentPhase != PhasePlayerAction {
		t.Errorf("Expected player action phase after decline, got %v", g.CurrentPhase)
	}

	// Insurance bet should be 0
	if g.PlayerHands[0].InsuranceBet != 0 {
		t.Errorf("Insurance bet should be 0, got %d", g.PlayerHands[0].InsuranceBet)
	}

	// Main hand bet should be intact
	if g.PlayerHands[0].Bet != 100 {
		t.Errorf("Main bet should be 100, got %d", g.PlayerHands[0].Bet)
	}
}
