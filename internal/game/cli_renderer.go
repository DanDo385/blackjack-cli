package game

import (
	"fmt"
	"strings"
)

// RenderState renders the current game state
func RenderState(g *Game, hideDealerHole bool) string {
	var sb strings.Builder

	sb.WriteString("+------------------------------------------+\n")

	// Render dealer hand
	sb.WriteString("| Dealer: ")
	if hideDealerHole && len(g.DealerHand.Cards) >= 2 {
		// Hide hole card
		sb.WriteString("[??, ")
		for i := 1; i < len(g.DealerHand.Cards); i++ {
			sb.WriteString(g.DealerHand.Cards[i].String())
			if i < len(g.DealerHand.Cards)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("]")
	} else {
		// Show all cards
		sb.WriteString(g.DealerHand.String())
		if !hideDealerHole {
			sb.WriteString(fmt.Sprintf(" (%d)", g.DealerHand.Value()))
		}
	}
	// Pad to column width (41 chars total including closing |)
	lineLen := len(sb.String()) - strings.LastIndex(sb.String(), "\n") - 1
	if lineLen < 41 {
		sb.WriteString(strings.Repeat(" ", 41-lineLen))
	}
	sb.WriteString("|\n")

	// Render player hands
	for i, hand := range g.PlayerHands {
		if len(g.PlayerHands) > 1 {
			sb.WriteString(fmt.Sprintf("| You (Hand %d/%d): ", i+1, len(g.PlayerHands)))
		} else {
			sb.WriteString("| You: ")
		}

		sb.WriteString(hand.String())
		if !hand.IsBust() {
			sb.WriteString(fmt.Sprintf(" (%d)", hand.Value()))
		} else {
			sb.WriteString(" (BUST)")
		}

		if hand.Surrendered {
			sb.WriteString(" [SURRENDERED]")
		}

		// Pad to column width
		lineLen := len(sb.String()) - strings.LastIndex(sb.String(), "\n") - 1
		if lineLen < 41 {
			sb.WriteString(strings.Repeat(" ", 41-lineLen))
		}
		sb.WriteString("|\n")
	}

	// Render bank and bet info
	sb.WriteString(fmt.Sprintf("| Bank: %-10d", g.Bank))
	if len(g.PlayerHands) > 0 {
		sb.WriteString(fmt.Sprintf(" Bet: %-10d", g.PlayerHands[0].Bet))
	}
	sb.WriteString("       |\n")

	sb.WriteString("+------------------------------------------+")

	return sb.String()
}

// RenderAvailableActions renders the available actions for the current hand
func RenderAvailableActions(actions []Action) string {
	if len(actions) == 0 {
		return ""
	}

	actionStrs := make([]string, 0, len(actions))
	for _, action := range actions {
		switch action {
		case ActionHit:
			actionStrs = append(actionStrs, "(H)it")
		case ActionStand:
			actionStrs = append(actionStrs, "(S)tand")
		case ActionDouble:
			actionStrs = append(actionStrs, "(D)ouble")
		case ActionSplit:
			actionStrs = append(actionStrs, "(P)split")
		case ActionSurrender:
			actionStrs = append(actionStrs, "(R)surrender")
		}
	}

	return "Action: " + strings.Join(actionStrs, ", ")
}

// RenderResult renders the final result of all hands
func RenderResult(g *Game) string {
	var sb strings.Builder

	sb.WriteString("\n" + RenderState(g, false) + "\n\n")
	sb.WriteString("Results:\n")

	for i, hand := range g.PlayerHands {
		handLabel := ""
		if len(g.PlayerHands) > 1 {
			handLabel = fmt.Sprintf("Hand %d/%d: ", i+1, len(g.PlayerHands))
		}

		// Check insurance first
		if hand.InsuranceBet > 0 {
			if g.DealerHasBlackjack {
				insuranceWin := Payout(OutcomeWin, hand.InsuranceBet, true)
				sb.WriteString(fmt.Sprintf("  %sInsurance pays %d chips\n", handLabel, insuranceWin))
			} else {
				sb.WriteString(fmt.Sprintf("  %sInsurance loses %d chips\n", handLabel, hand.InsuranceBet))
			}
		}

		// Main hand outcome
		outcome := DetermineOutcome(hand, g.DealerHand)
		payout := Payout(outcome, hand.Bet, false)

		switch outcome {
		case OutcomeBlackjack:
			sb.WriteString(fmt.Sprintf("  %sBLACKJACK! Wins %d chips\n", handLabel, payout-hand.Bet))
		case OutcomeWin:
			sb.WriteString(fmt.Sprintf("  %sWin! Pays %d chips\n", handLabel, payout-hand.Bet))
		case OutcomePush:
			sb.WriteString(fmt.Sprintf("  %sPush! Returns %d chips\n", handLabel, payout))
		case OutcomeLose:
			sb.WriteString(fmt.Sprintf("  %sLose! Loses %d chips\n", handLabel, hand.Bet))
		case OutcomeSurrender:
			sb.WriteString(fmt.Sprintf("  %sSurrender! Returns %d chips\n", handLabel, payout))
		}
	}

	sb.WriteString(fmt.Sprintf("\nBank: %d chips\n", g.Bank))

	return sb.String()
}
