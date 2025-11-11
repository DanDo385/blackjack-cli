package game

// PeekForBlackjack checks if the dealer has blackjack with the given up and down cards
func PeekForBlackjack(upCard, downCard Card) bool {
	// Only peek if upcard is Ace or 10-value
	if !upCard.IsAce() && upCard.Rank.Value() != 10 {
		return false
	}

	// Check if the two cards make blackjack
	hand := &Hand{Cards: []Card{upCard, downCard}, IsInitialDeal: true}
	return hand.IsBlackjack()
}

// DealerPlay plays out the dealer's hand according to S17 rules
func DealerPlay(deck *[]Card, hand *Hand) {
	// Dealer draws to 17, stands on soft 17
	for {
		value := hand.Value()
		if value > 17 {
			break
		}
		if value == 17 {
			// Stand on all 17s (including soft 17 per S17 rule)
			break
		}

		// Draw a card
		drawn, remaining := Draw(*deck, 1)
		if len(drawn) == 0 {
			// Shouldn't happen with proper deck management
			break
		}
		hand.Add(drawn[0])
		*deck = remaining
	}
}
