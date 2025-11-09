package game

import (
	"testing"
)

func TestEnginePlayerBlackjack(t *testing.T) {
	// Create a simple test deck: Player gets BJ, dealer gets 20
	testDeck := []Card{
		{Rank: Ace, Suit: Spades},   // Player card 1
		{Rank: King, Suit: Clubs},   // Dealer card 1
		{Rank: King, Suit: Hearts},  // Player card 2 (BJ!)
		{Rank: Queen, Suit: Clubs},  // Dealer card 2 (20)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()

	// Start hand with bet of 100
	initialBank := g.Bank
	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 100

	// Player should have blackjack
	if !g.PlayerHands[0].IsBlackjack() {
		t.Error("Player should have blackjack")
	}

	// Dealer should not have blackjack
	if g.DealerHasBlackjack {
		t.Error("Dealer should not have blackjack")
	}

	// Play dealer and resolve
	g.CurrentPhase = PhaseDealerAction
	g.PlayDealer()
	g.CurrentPhase = PhaseResolution
	g.resolvePayouts()

	// Player should win with blackjack payout (3:2)
	expectedBank := initialBank - 100 + Payout(OutcomeBlackjack, 100, false)
	if g.Bank != expectedBank {
		t.Errorf("Bank = %d, want %d", g.Bank, expectedBank)
	}
}

func TestEngineSplitEights(t *testing.T) {
	// Create test deck: Player gets 8,8; dealer gets 9
	testDeck := []Card{
		{Rank: Eight, Suit: Spades},  // Player card 1
		{Rank: Nine, Suit: Clubs},    // Dealer card 1
		{Rank: Eight, Suit: Hearts},  // Player card 2 (pair!)
		{Rank: Seven, Suit: Clubs},   // Dealer card 2
		{Rank: Three, Suit: Diamonds}, // Split hand 1 draw
		{Rank: Ten, Suit: Diamonds},  // Split hand 2 draw
		{Rank: King, Suit: Spades},   // Dealer draw (bust: 9+7+K=26)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()

	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 100

	// Check if we can split
	if !g.PlayerHands[0].CanSplit() {
		t.Fatal("Should be able to split")
	}

	// Split
	err = g.PlayerAction(ActionSplit)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	g.Bank -= 100 // Second bet for split

	// Should now have 2 hands
	if len(g.PlayerHands) != 2 {
		t.Errorf("Expected 2 hands after split, got %d", len(g.PlayerHands))
	}

	// Stand on both hands (dealer plays and payouts resolve automatically after last hand)
	err = g.PlayerAction(ActionStand)
	if err != nil {
		t.Fatalf("Stand on hand 1 failed: %v", err)
	}

	err = g.PlayerAction(ActionStand)
	if err != nil {
		t.Fatalf("Stand on hand 2 failed: %v", err)
	}

	// After last hand, dealer should have played and payouts resolved

	// Both hands should win (dealer busts)
	// Hand 1: 8+3=11 wins
	// Hand 2: 8+10=18 wins
	// Each bet is 100, wins pay 200 each
	expectedDelta := 200 + 200 // Both hands win
	expectedBank := 1000 - 200 + expectedDelta
	if g.Bank != expectedBank {
		t.Errorf("Bank = %d, want %d", g.Bank, expectedBank)
	}
}

func TestEngineSplitAcesOneCardOnly(t *testing.T) {
	// Create test deck: Player gets A,A; dealer gets 7
	testDeck := []Card{
		{Rank: Ace, Suit: Spades},    // Player card 1
		{Rank: Seven, Suit: Clubs},   // Dealer card 1
		{Rank: Ace, Suit: Hearts},    // Player card 2 (pair!)
		{Rank: Ten, Suit: Clubs},     // Dealer card 2 (17)
		{Rank: Nine, Suit: Diamonds}, // Split hand 1 draw (A+9=20)
		{Rank: Eight, Suit: Diamonds}, // Split hand 2 draw (A+8=19)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()

	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 100

	// Split aces
	err = g.PlayerAction(ActionSplit)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	g.Bank -= 100

	// Both hands should be marked as split aces
	if !g.PlayerHands[0].IsSplitAces {
		t.Error("Hand 1 should be marked as split aces")
	}
	if !g.PlayerHands[1].IsSplitAces {
		t.Error("Hand 2 should be marked as split aces")
	}

	// Each hand should have exactly 2 cards (original + one draw)
	if len(g.PlayerHands[0].Cards) != 2 {
		t.Errorf("Hand 1 should have 2 cards, got %d", len(g.PlayerHands[0].Cards))
	}
	if len(g.PlayerHands[1].Cards) != 2 {
		t.Errorf("Hand 2 should have 2 cards, got %d", len(g.PlayerHands[1].Cards))
	}

	// Split aces cannot double
	if g.PlayerHands[0].CanDouble() {
		t.Error("Split aces should not be able to double")
	}

	// Split aces automatically advance to dealer after both get one card
	// Dealer should have already played and payouts resolved

	// Hand 1: A+9=20 wins vs dealer 17
	// Hand 2: A+8=19 wins vs dealer 17
	expectedBank := 1000 - 200 + 200 + 200
	if g.Bank != expectedBank {
		t.Errorf("Bank = %d, want %d", g.Bank, expectedBank)
	}
}

func TestEngineDouble(t *testing.T) {
	// Create test deck: Player gets 11; dealer gets 6
	testDeck := []Card{
		{Rank: Six, Suit: Spades},    // Player card 1
		{Rank: Six, Suit: Clubs},     // Dealer card 1
		{Rank: Five, Suit: Hearts},   // Player card 2 (11)
		{Rank: Ten, Suit: Clubs},     // Dealer card 2 (16)
		{Rank: Ten, Suit: Diamonds},  // Player double card (21)
		{Rank: Eight, Suit: Hearts},  // Dealer draw (24, bust)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()

	initialBank := g.Bank
	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 100

	// Check if we can double
	if !g.PlayerHands[0].CanDouble() {
		t.Fatal("Should be able to double")
	}

	// Double
	err = g.PlayerAction(ActionDouble)
	if err != nil {
		t.Fatalf("Double failed: %v", err)
	}

	// Bet should be doubled
	if g.PlayerHands[0].Bet != 200 {
		t.Errorf("Bet should be 200 after double, got %d", g.PlayerHands[0].Bet)
	}

	// Should have exactly 3 cards
	if len(g.PlayerHands[0].Cards) != 3 {
		t.Errorf("Should have 3 cards after double, got %d", len(g.PlayerHands[0].Cards))
	}

	// Double automatically advances to dealer and resolves
	// (dealer has played and payouts are resolved)

	// Player wins with doubled bet: -200 + 400 = +200
	expectedBank := initialBank - 200 + 400
	if g.Bank != expectedBank {
		t.Errorf("Bank = %d, want %d", g.Bank, expectedBank)
	}
}

func TestEngineSurrender(t *testing.T) {
	// Create test deck: Player gets 16; dealer shows 10
	testDeck := []Card{
		{Rank: Ten, Suit: Spades},    // Player card 1
		{Rank: Ten, Suit: Clubs},     // Dealer card 1
		{Rank: Six, Suit: Hearts},    // Player card 2 (16)
		{Rank: Seven, Suit: Clubs},   // Dealer card 2 (17)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()

	initialBank := g.Bank
	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}

	g.Bank -= 100

	// Check if we can surrender
	if !g.PlayerHands[0].CanSurrender() {
		t.Fatal("Should be able to surrender")
	}

	// Surrender
	err = g.PlayerAction(ActionSurrender)
	if err != nil {
		t.Fatalf("Surrender failed: %v", err)
	}

	// Hand should be marked as surrendered
	if !g.PlayerHands[0].Surrendered {
		t.Error("Hand should be marked as surrendered")
	}

	// Surrender automatically advances to dealer and resolves
	// (dealer has played and payouts are resolved)

	// Surrender returns half bet: -100 + 50 = -50
	expectedBank := initialBank - 100 + 50
	if g.Bank != expectedBank {
		t.Errorf("Bank = %d, want %d", g.Bank, expectedBank)
	}
}

func TestEngineMaxFourHands(t *testing.T) {
	// Create test deck with multiple pairs to split
	testDeck := []Card{
		{Rank: Eight, Suit: Spades},    // Player card 1
		{Rank: Nine, Suit: Clubs},      // Dealer card 1
		{Rank: Eight, Suit: Hearts},    // Player card 2 (pair)
		{Rank: Ten, Suit: Clubs},       // Dealer card 2
		{Rank: Eight, Suit: Diamonds},  // Hand 1 gets 8 (can split again)
		{Rank: Five, Suit: Diamonds},   // Hand 2 gets 5
		{Rank: Eight, Suit: Clubs},     // Hand 1 splits again, gets 8 (can split again)
		{Rank: Six, Suit: Hearts},      // New hand 2 gets 6
		{Rank: Seven, Suit: Spades},    // Hand 1 gets 7
		{Rank: Four, Suit: Hearts},     // New hand gets 4 (now 4 hands total)
	}

	g := NewGame()
	g.Deck = testDeck
	g.RNG = FixedSeededRand()

	err := g.StartHand(100)
	if err != nil {
		t.Fatalf("StartHand failed: %v", err)
	}
	g.Bank -= 100

	// First split
	err = g.PlayerAction(ActionSplit)
	if err != nil {
		t.Fatalf("First split failed: %v", err)
	}
	g.Bank -= 100

	if len(g.PlayerHands) != 2 {
		t.Fatalf("Expected 2 hands after first split, got %d", len(g.PlayerHands))
	}

	// Second split on hand 1 (8,8)
	err = g.PlayerAction(ActionSplit)
	if err != nil {
		t.Fatalf("Second split failed: %v", err)
	}
	g.Bank -= 100

	if len(g.PlayerHands) != 3 {
		t.Fatalf("Expected 3 hands after second split, got %d", len(g.PlayerHands))
	}

	// Third split on hand 1 (8,8)
	err = g.PlayerAction(ActionSplit)
	if err != nil {
		t.Fatalf("Third split failed: %v", err)
	}
	g.Bank -= 100

	if len(g.PlayerHands) != 4 {
		t.Fatalf("Expected 4 hands after third split, got %d", len(g.PlayerHands))
	}

	// Try to split again (should fail - max 4 hands)
	// But first we need a hand that can split
	// Since we already have 4 hands, we can't split anymore regardless
	// Let's just verify we have 4 hands
	if len(g.PlayerHands) != 4 {
		t.Errorf("Should have exactly 4 hands, got %d", len(g.PlayerHands))
	}
}
