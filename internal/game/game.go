package game

import (
	"fmt"
	"math/rand"
)

// Action represents a player action
type Action int

const (
	ActionHit Action = iota
	ActionStand
	ActionDouble
	ActionSplit
	ActionSurrender
)

func (a Action) String() string {
	switch a {
	case ActionHit:
		return "Hit"
	case ActionStand:
		return "Stand"
	case ActionDouble:
		return "Double"
	case ActionSplit:
		return "Split"
	case ActionSurrender:
		return "Surrender"
	default:
		return "Unknown"
	}
}

// Phase represents the current phase of the game
type Phase int

const (
	PhaseBetting Phase = iota
	PhaseInsurance
	PhasePlayerAction
	PhaseDealerAction
	PhaseResolution
	PhaseGameOver
)

// Game represents the game state
type Game struct {
	Bank               int
	Deck               []Card
	PlayerHands        []*Hand
	DealerHand         *Hand
	CurrentPhase       Phase
	ActiveHandIndex    int
	RNG                *rand.Rand
	DealerHasBlackjack bool
	InsuranceOffered   bool
}

// NewGame creates a new game with the starting bank
func NewGame() *Game {
	return &Game{
		Bank:         StartingBank,
		RNG:          NewRand(),
		CurrentPhase: PhaseBetting,
	}
}

// StartHand initializes a new hand with the given bet
func (g *Game) StartHand(bet int) error {
	if bet < MinBet {
		return fmt.Errorf("minimum bet is %d", MinBet)
	}
	if bet > g.Bank {
		return fmt.Errorf("bet exceeds bank balance")
	}

	// Create new deck and shuffle (unless one is already set for testing)
	if len(g.Deck) == 0 {
		g.Deck = NewDeck()
		Shuffle(g.Deck, g.RNG)
	}

	// Initialize hands
	g.PlayerHands = []*Hand{NewHand(bet)}
	g.DealerHand = NewHand(0)
	g.ActiveHandIndex = 0
	g.DealerHasBlackjack = false
	g.InsuranceOffered = false

	// Deal initial cards: player, dealer, player, dealer
	g.dealCard(g.PlayerHands[0])
	g.dealCard(g.DealerHand)
	g.dealCard(g.PlayerHands[0])
	g.dealCard(g.DealerHand)

	// Check if insurance should be offered (based on visible upcard)
	if len(g.DealerHand.Cards) >= 2 {
		upcard := g.DealerHand.Cards[1]
		if upcard.IsAce() {
			g.CurrentPhase = PhaseInsurance
			g.InsuranceOffered = true
		} else {
			// Peek for blackjack if dealer shows 10
			if g.DealerHand.Cards[0].Rank.Value() == 10 {
				if PeekForBlackjack(g.DealerHand.Cards[0], g.DealerHand.Cards[1]) {
					g.DealerHasBlackjack = true
					g.CurrentPhase = PhaseResolution
					return nil
				}
			}
			g.CurrentPhase = PhasePlayerAction
		}
	} else {
		// Peek for blackjack if dealer shows 10
		// Only check if dealer has at least one card
		if len(g.DealerHand.Cards) > 0 && g.DealerHand.Cards[0].Rank.Value() == 10 {
			// Only peek if dealer has 2 cards
			if len(g.DealerHand.Cards) >= 2 {
				if PeekForBlackjack(g.DealerHand.Cards[0], g.DealerHand.Cards[1]) {
					g.DealerHasBlackjack = true
					g.CurrentPhase = PhaseResolution
					return nil
				}
			}
		}
		g.CurrentPhase = PhasePlayerAction
	}

	return nil
}

// TakeInsurance allows the player to take insurance with the given bet
func (g *Game) TakeInsurance(insuranceBet int) error {
	if g.CurrentPhase != PhaseInsurance {
		return fmt.Errorf("insurance not available")
	}

	maxInsurance := g.PlayerHands[0].Bet / 2
	if insuranceBet > maxInsurance {
		return fmt.Errorf("insurance bet cannot exceed half of original bet (%d)", maxInsurance)
	}

	g.PlayerHands[0].InsuranceBet = insuranceBet

	// Peek for dealer blackjack
	if PeekForBlackjack(g.DealerHand.Cards[0], g.DealerHand.Cards[1]) {
		g.DealerHasBlackjack = true
		g.CurrentPhase = PhaseResolution
	} else {
		g.CurrentPhase = PhasePlayerAction
	}

	return nil
}

// DeclineInsurance declines insurance and proceeds with the game
func (g *Game) DeclineInsurance() error {
	if g.CurrentPhase != PhaseInsurance {
		return fmt.Errorf("insurance not available")
	}

	// Peek for dealer blackjack
	if PeekForBlackjack(g.DealerHand.Cards[0], g.DealerHand.Cards[1]) {
		g.DealerHasBlackjack = true
		g.CurrentPhase = PhaseResolution
	} else {
		g.CurrentPhase = PhasePlayerAction
	}

	return nil
}

// PlayerAction performs a player action on the current active hand
func (g *Game) PlayerAction(action Action) error {
	if g.CurrentPhase != PhasePlayerAction {
		return fmt.Errorf("not in player action phase")
	}

	if g.ActiveHandIndex >= len(g.PlayerHands) {
		return fmt.Errorf("invalid hand index")
	}

	hand := g.PlayerHands[g.ActiveHandIndex]

	switch action {
	case ActionHit:
		return g.hit(hand)
	case ActionStand:
		return g.stand()
	case ActionDouble:
		return g.double(hand)
	case ActionSplit:
		return g.split()
	case ActionSurrender:
		return g.surrender(hand)
	default:
		return fmt.Errorf("invalid action")
	}
}

func (g *Game) hit(hand *Hand) error {
	hand.IsInitialDeal = false
	g.dealCard(hand)

	// Automatically stand on 21
	if hand.Value() == 21 {
		return g.advanceToNextHand()
	}

	if hand.IsBust() {
		// Move to next hand or dealer
		return g.advanceToNextHand()
	}

	// Split aces only get one card - advance after dealing the card
	if hand.IsSplitAces {
		return g.advanceToNextHand()
	}

	// For normal hands, stay on the same hand to allow multiple hits
	return nil
}

func (g *Game) stand() error {
	return g.advanceToNextHand()
}

func (g *Game) double(hand *Hand) error {
	if !hand.CanDouble() {
		return fmt.Errorf("cannot double")
	}

	// Double the bet - need to deduct the additional bet amount
	if hand.Bet > g.Bank {
		return fmt.Errorf("insufficient funds to double")
	}
	g.Bank -= hand.Bet // Deduct the additional bet for doubling
	hand.Bet *= 2
	hand.Doubled = true
	hand.IsInitialDeal = false

	// Check for blackjack on initial deal before doubling
	if hand.IsBlackjack() {
		return g.advanceToNextHand()
	}

	// Deal one card and stand
	g.dealCard(hand)

	return g.advanceToNextHand()
}

func (g *Game) split() error {
	hand := g.PlayerHands[g.ActiveHandIndex]

	if !hand.CanSplit() {
		return fmt.Errorf("cannot split")
	}

	// Cannot resplit aces
	if hand.IsSplitAces {
		return fmt.Errorf("cannot resplit aces")
	}

	// Check if we can afford the split
	if hand.Bet > g.Bank {
		return fmt.Errorf("insufficient funds to split")
	}
	g.Bank -= hand.Bet // Deduct the bet for the new hand

	// Check if we've reached the max of 4 hands
	if len(g.PlayerHands) >= 4 {
		return fmt.Errorf("cannot split more than 4 hands")
	}

	// Create new hand with the second card
	newHand := NewHand(hand.Bet)
	newHand.Add(hand.Cards[1])

	// Keep only the first card in the current hand
	hand.Cards = hand.Cards[:1]

	// Check if we're splitting aces
	isAceSplit := hand.Cards[0].IsAce()
	if isAceSplit {
		hand.IsSplitAces = true
		newHand.IsSplitAces = true
	}

	// Deal one card to each hand
	g.dealCard(hand)
	g.dealCard(newHand)

	// Insert the new hand after the current hand
	g.PlayerHands = append(g.PlayerHands[:g.ActiveHandIndex+1], append([]*Hand{newHand}, g.PlayerHands[g.ActiveHandIndex+1:]...)...)

	// Both hands now have their 2 cards, so they're in "initial" state for actions
	// (can split again if they get a pair, can double, etc.)
	// But they cannot have natural blackjack since they came from a split
	hand.IsInitialDeal = true
	hand.IsFromSplit = true
	newHand.IsInitialDeal = true
	newHand.IsFromSplit = true

	// For split aces, the first hand will be handled by the hit() function
	// which will automatically advance after dealing one card
	// We stay on the current hand so the player can see the result
	// The main loop will handle advancing after split aces

	return nil
}

func (g *Game) surrender(hand *Hand) error {
	if !hand.CanSurrender() {
		return fmt.Errorf("cannot surrender")
	}

	hand.Surrendered = true
	hand.IsInitialDeal = false

	return g.advanceToNextHand()
}

func (g *Game) advanceToNextHand() error {
	g.ActiveHandIndex++

	if g.ActiveHandIndex >= len(g.PlayerHands) {
		// All player hands done, move to dealer
		g.CurrentPhase = PhaseDealerAction
		g.PlayDealer()
		g.CurrentPhase = PhaseResolution
		g.ResolvePayouts()
	}

	return nil
}

// PlayDealer plays out the dealer's hand according to house rules
func (g *Game) PlayDealer() {
	// Check if all player hands are bust or surrendered
	allBustOrSurrendered := true
	for _, hand := range g.PlayerHands {
		if !hand.IsBust() && !hand.Surrendered {
			allBustOrSurrendered = false
			break
		}
	}

	// Dealer doesn't play if all player hands are bust or surrendered
	if allBustOrSurrendered {
		return
	}

	// Dealer already has blackjack
	if g.DealerHasBlackjack {
		return
	}

	DealerPlay(&g.Deck, g.DealerHand)
}

func (g *Game) ResolvePayouts() {
	// Start with the current bank, which no longer includes the bets (already deducted)
	finalBank := g.Bank

	for _, hand := range g.PlayerHands {
		// Resolve insurance bet
		if hand.InsuranceBet > 0 {
			if g.DealerHasBlackjack {
				// Insurance pays 2:1
				finalBank += Payout(OutcomeWin, hand.InsuranceBet, true)
			} else {
				// Insurance loses
				finalBank += Payout(OutcomeLose, hand.InsuranceBet, true)
			}
		}

		// If dealer has blackjack
		if g.DealerHasBlackjack {
			if hand.IsBlackjack() {
				// Push - bet is returned
				finalBank += hand.Bet
			} else {
				// Player loses bet (already deducted, so nothing to add back)
			}
			continue
		}

		// Determine outcome
		outcome := DetermineOutcome(hand, g.DealerHand)
		payout := Payout(outcome, hand.Bet, false)

		// The payout function returns the total amount given to the player.
		// Since the bet was already deducted from the bank, we add back the full payout.
		finalBank += payout
	}

	// Update the bank with the final calculated value
	g.Bank = finalBank

	if g.Bank <= 0 {
		g.CurrentPhase = PhaseGameOver
	} else {
		g.CurrentPhase = PhaseBetting
	}
}

func (g *Game) dealCard(hand *Hand) {
	drawn, remaining := Draw(g.Deck, 1)
	if len(drawn) > 0 {
		hand.Add(drawn[0])
		g.Deck = remaining
	}
}

// GetAvailableActions returns the available actions for the current active hand
func (g *Game) GetAvailableActions() []Action {
	if g.CurrentPhase != PhasePlayerAction {
		return nil
	}

	if g.ActiveHandIndex >= len(g.PlayerHands) {
		return nil
	}

	hand := g.PlayerHands[g.ActiveHandIndex]

	// If hand is bust, no actions available
	if hand.IsBust() {
		return nil
	}

	// If hand is 21, no actions available (auto-stand)
	if hand.Value() == 21 {
		return nil
	}

	// Split aces only get one card - no actions available except stand
	// (handled in main loop, but this prevents any other actions)
	if hand.IsSplitAces && len(hand.Cards) >= 2 {
		return []Action{ActionStand}
	}

	actions := []Action{ActionHit, ActionStand}

	// For doubling, we need enough total chips to cover the doubled bet
	// Since the bet was already deducted in main.go, we check if g.Bank + hand.Bet >= hand.Bet * 2
	// This simplifies to: g.Bank >= hand.Bet
	// Example: Bank=1006, Bet=1000 -> After deduction: Bank=6, need 6>=1000? No, can't double
	// Example: Bank=2000, Bet=1000 -> After deduction: Bank=1000, need 1000>=1000? Yes, can double
	// However, if we want to allow doubling when total chips are sufficient, we should check:
	// g.Bank + hand.Bet >= hand.Bet * 2, which means g.Bank >= hand.Bet
	// But the user expects to be able to double with 1006 chips and 1000 bet, which means
	// they want: g.Bank + hand.Bet >= hand.Bet * 2 -> 6 + 1000 >= 2000 -> 1006 >= 2000? No
	// So the current check is correct - you need at least 2000 chips total to double a 1000 bet
	if hand.CanDouble() && g.Bank >= hand.Bet {
		actions = append(actions, ActionDouble)
	}

	// For splitting, we need enough remaining bank to cover the bet for the new hand
	// Since the bet was already deducted, we check if g.Bank >= hand.Bet
	if hand.CanSplit() && g.Bank >= hand.Bet && len(g.PlayerHands) < 4 {
		actions = append(actions, ActionSplit)
	}

	if hand.CanSurrender() {
		actions = append(actions, ActionSurrender)
	}

	return actions
}

// GetCurrentHand returns the current active hand
func (g *Game) GetCurrentHand() *Hand {
	if g.ActiveHandIndex >= len(g.PlayerHands) {
		return nil
	}
	return g.PlayerHands[g.ActiveHandIndex]
}
