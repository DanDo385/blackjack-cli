package game

// Game rules constants
const (
	DealerStandsSoft17 = true // S17 rule
	BlackjackPayout    = 1.5  // 3:2 payout
	InsurancePayout    = 2.0  // 2:1 payout
	MinBet             = 1
	StartingBank       = 1000
)

// Outcome represents the result of a hand
type Outcome int

const (
	OutcomeWin Outcome = iota
	OutcomeLose
	OutcomePush
	OutcomeBlackjack
	OutcomeSurrender
)

func (o Outcome) String() string {
	switch o {
	case OutcomeWin:
		return "Win"
	case OutcomeLose:
		return "Lose"
	case OutcomePush:
		return "Push"
	case OutcomeBlackjack:
		return "Blackjack"
	case OutcomeSurrender:
		return "Surrender"
	default:
		return "Unknown"
	}
}

// Payout calculates the payout for a given outcome and bet
// Returns the delta to the bank (positive for win, negative for loss)
func Payout(outcome Outcome, bet int, isInsurance bool) int {
	if isInsurance {
		if outcome == OutcomeWin {
			return int(float64(bet) * InsurancePayout)
		}
		return -bet
	}

	switch outcome {
	case OutcomeBlackjack:
		// Natural blackjack pays 3:2 (bet + 1.5x bet)
		return int(float64(bet) * (1 + BlackjackPayout))
	case OutcomeWin:
		// Regular win pays 1:1 (bet + bet)
		return bet + bet
	case OutcomePush:
		// Push returns the original bet
		return bet
	case OutcomeLose:
		// Lose forfeits the bet
		return 0
	case OutcomeSurrender:
		// Surrender returns half the bet
		return bet / 2
	default:
		return 0
	}
}

// DetermineOutcome determines the outcome of a hand vs dealer
func DetermineOutcome(playerHand, dealerHand *Hand) Outcome {
	// Check surrender first
	if playerHand.Surrendered {
		return OutcomeSurrender
	}

	// Player bust always loses
	if playerHand.IsBust() {
		return OutcomeLose
	}

	// Dealer bust, player wins
	if dealerHand.IsBust() {
		return OutcomeWin
	}

	playerValue := playerHand.Value()
	dealerValue := dealerHand.Value()

	// Natural blackjack (only on initial 2-card hand)
	if playerHand.IsBlackjack() && !dealerHand.IsBlackjack() {
		return OutcomeBlackjack
	}

	// Dealer blackjack beats non-blackjack
	if dealerHand.IsBlackjack() && !playerHand.IsBlackjack() {
		return OutcomeLose
	}

	// Both blackjack or same value
	if playerValue == dealerValue {
		return OutcomePush
	}

	// Higher value wins
	if playerValue > dealerValue {
		return OutcomeWin
	}

	return OutcomeLose
}
