package game

import "fmt"

// Suit represents a card suit
type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

func (s Suit) String() string {
	switch s {
	case Clubs:
		return "♣"
	case Diamonds:
		return "♦"
	case Hearts:
		return "♥"
	case Spades:
		return "♠"
	default:
		return "?"
	}
}

// Rank represents a card rank
type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return "?"
	}
}

// RankValue returns the blackjack value of the rank (Ace is 11, face cards are 10)
func (r Rank) Value() int {
	switch r {
	case Ace:
		return 11
	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten, Jack, Queen, King:
		return 10
	default:
		return 0
	}
}

// Card represents a playing card
type Card struct {
	Rank Rank
	Suit Suit
}

// String returns the string representation of a card (e.g., "K♣")
func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

// IsAce returns true if the card is an Ace
func (c Card) IsAce() bool {
	return c.Rank == Ace
}

// ParseCard parses a card from a string like "AS", "KD", "10H"
func ParseCard(s string) (Card, error) {
	if len(s) < 2 {
		return Card{}, fmt.Errorf("invalid card string: %s", s)
	}

	var rank Rank
	var suitChar byte

	if len(s) == 2 {
		// Single character rank
		suitChar = s[1]
		switch s[0] {
		case 'A':
			rank = Ace
		case '2':
			rank = Two
		case '3':
			rank = Three
		case '4':
			rank = Four
		case '5':
			rank = Five
		case '6':
			rank = Six
		case '7':
			rank = Seven
		case '8':
			rank = Eight
		case '9':
			rank = Nine
		case 'J':
			rank = Jack
		case 'Q':
			rank = Queen
		case 'K':
			rank = King
		default:
			return Card{}, fmt.Errorf("invalid rank: %c", s[0])
		}
	} else if len(s) == 3 && s[0:2] == "10" {
		rank = Ten
		suitChar = s[2]
	} else {
		return Card{}, fmt.Errorf("invalid card string: %s", s)
	}

	var suit Suit
	switch suitChar {
	case 'C':
		suit = Clubs
	case 'D':
		suit = Diamonds
	case 'H':
		suit = Hearts
	case 'S':
		suit = Spades
	default:
		return Card{}, fmt.Errorf("invalid suit: %c", suitChar)
	}

	return Card{Rank: rank, Suit: suit}, nil
}
