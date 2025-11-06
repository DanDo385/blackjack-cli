package game

// Hand represents a blackjack hand
type Hand struct {
	Cards         []Card
	Bet           int
	IsSplitAces   bool
	Doubled       bool
	Surrendered   bool
	IsInitialDeal bool // True if this hand has had no actions yet
	IsFromSplit   bool // True if this hand came from a split (cannot have natural blackjack)
	InsuranceBet  int
}

// NewHand creates a new hand with the given bet
func NewHand(bet int) *Hand {
	return &Hand{
		Cards:         make([]Card, 0),
		Bet:           bet,
		IsInitialDeal: true,
	}
}

// Add adds a card to the hand
func (h *Hand) Add(c Card) {
	h.Cards = append(h.Cards, c)
}

// Totals returns the hard total, soft total, and whether the hand is soft
func (h *Hand) Totals() (hard int, soft int, isSoft bool) {
	hard = 0
	aces := 0

	for _, card := range h.Cards {
		if card.IsAce() {
			aces++
			hard += 1
		} else {
			hard += card.Rank.Value()
		}
	}

	soft = hard
	// Try to use one ace as 11 if it doesn't bust
	if aces > 0 && hard+10 <= 21 {
		soft = hard + 10
		isSoft = true
	}

	return hard, soft, isSoft
}

// Value returns the best total for the hand (soft if it doesn't bust, otherwise hard)
func (h *Hand) Value() int {
	hard, soft, isSoft := h.Totals()
	if isSoft {
		return soft
	}
	return hard
}

// IsBlackjack returns true if the hand is a natural blackjack (2 cards totaling 21)
// Hands from splits cannot have natural blackjack
func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21 && h.IsInitialDeal && !h.IsFromSplit
}

// IsBust returns true if the hand is bust (over 21)
func (h *Hand) IsBust() bool {
	return h.Value() > 21
}

// IsSoft returns true if the hand is soft (contains an ace counted as 11)
func (h *Hand) IsSoft() bool {
	_, _, isSoft := h.Totals()
	return isSoft
}

// CanSplit returns true if the hand can be split
func (h *Hand) CanSplit() bool {
	// Can split on initial deal (or after a split) with exactly 2 cards of the same rank
	// Cannot split after hitting (IsInitialDeal becomes false after first action)
	// Special case: split aces cannot be resplit (handled in game logic)
	return h.IsInitialDeal && len(h.Cards) == 2 && h.Cards[0].Rank == h.Cards[1].Rank
}

// CanDouble returns true if the hand can be doubled
func (h *Hand) CanDouble() bool {
	// Can only double on first action
	// Cannot double on split aces (they only get one card)
	if h.IsSplitAces {
		return false
	}
	return h.IsInitialDeal
}

// CanSurrender returns true if the hand can surrender (late surrender only)
func (h *Hand) CanSurrender() bool {
	// Can only surrender on first action
	return h.IsInitialDeal && len(h.Cards) == 2
}

// String returns a string representation of the hand
func (h *Hand) String() string {
	s := "["
	for i, card := range h.Cards {
		if i > 0 {
			s += ", "
		}
		s += card.String()
	}
	s += "]"
	return s
}
