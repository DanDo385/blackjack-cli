package main

import (
	"fmt"
	"os"

	"github.com/DanDo385/blackjack-cli/internal/game"
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║         BLACKJACK CLI GAME             ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	g := game.NewGame()

	for g.Bank > 0 {
		// Betting phase
		fmt.Printf("\n🎰 Current Bank: %d chips\n", g.Bank)
		bet, err := game.PromptBet(os.Stdin, g.Bank)
		if err != nil {
			fmt.Printf("Error reading bet: %v\n", err)
			continue
		}

		// Start hand
		err = g.StartHand(bet)
		if err != nil {
			fmt.Printf("Error starting hand: %v\n", err)
			continue
		}

		// Deduct the bet
		g.Bank -= bet

		// Show initial state
		fmt.Println()
		fmt.Println(game.RenderState(g, true))
		fmt.Println()

	
		// Insurance phase
		if g.CurrentPhase == game.PhaseInsurance {
			dealerCard := g.DealerHand.Cards[1]
			maxInsurance := bet / 2
			if maxInsurance > g.Bank {
				maxInsurance = g.Bank
			}

			if maxInsurance == 0 {
				err = g.DeclineInsurance()
				if err != nil {
					fmt.Printf("Error declining insurance: %v\n", err)
					continue
				}
			} else {
				prompt := fmt.Sprintf("Dealer shows %s. Take insurance?", dealerCard.String())
				takeInsurance, err := game.PromptYesNo(os.Stdin, prompt)
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					continue
				}

				if takeInsurance {
					insuranceBet, err := game.PromptInsurance(os.Stdin, maxInsurance)
					if err != nil {
						fmt.Printf("Error reading insurance bet: %v\n", err)
						continue
					}

					if insuranceBet > 0 {
						g.Bank -= insuranceBet
						err = g.TakeInsurance(insuranceBet)
						if err != nil {
							fmt.Printf("Error taking insurance: %v\n", err)
							g.Bank += insuranceBet // Refund
							continue
						}
					} else {
						err = g.DeclineInsurance()
						if err != nil {
							fmt.Printf("Error declining insurance: %v\n", err)
							continue
						}
					}
				} else {
					err = g.DeclineInsurance()
					if err != nil {
						fmt.Printf("Error declining insurance: %v\n", err)
						continue
					}
				}
			}
		}

		// Check if dealer has blackjack after insurance
		if g.CurrentPhase == game.PhaseResolution {
			// Resolve payouts (insurance + main hand)
			g.ResolvePayouts()

			fmt.Println("\n🃏 Dealer has Blackjack!")
			fmt.Println(game.RenderResult(g))

			// Continue to next hand
			if !promptContinue() {
				break
			}
			continue
		}

		// Check for player blackjack and skip to dealer
		if g.PlayerHands[0].IsBlackjack() {
			fmt.Println("\n🃏 Blackjack!")
			// Skip player action and go straight to dealer
			for g.CurrentPhase == game.PhasePlayerAction {
				g.PlayerAction(game.ActionStand)
			}
		}

		// Player action phase
		for g.CurrentPhase == game.PhasePlayerAction {
			currentHand := g.GetCurrentHand()
			if currentHand == nil {
				break
			}

			// Check if hand is automatically done (split aces with one card dealt)
			if currentHand.IsSplitAces && len(currentHand.Cards) > 1 {
				fmt.Println("Split aces receive only one card.")
				err := g.PlayerAction(game.ActionStand)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					break
				}
				fmt.Println()
				fmt.Println(game.RenderState(g, true))
				fmt.Println()
				continue
			}

			// Check if hand is bust
			if currentHand.IsBust() {
				fmt.Println("BUST!")
				err := g.PlayerAction(game.ActionStand)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					break
				}
				fmt.Println()
				fmt.Println(game.RenderState(g, true))
				fmt.Println()
				continue
			}

			// Get available actions
			actions := g.GetAvailableActions()
			if len(actions) == 0 {
				break
			}

			// Prompt for action
			action, err := game.PromptAction(os.Stdin, actions)
			if err != nil {
				fmt.Printf("Error reading action: %v\n", err)
				continue
			}

			// Perform action
			err = g.PlayerAction(action)
			if err != nil {
				fmt.Printf("Error performing action: %v\n", err)
				continue
			}

			// Show board after action
			fmt.Println()
			fmt.Println(game.RenderState(g, true))
			fmt.Println()

			// Show result of action (bust)
			if action == game.ActionHit && currentHand.IsBust() {
				fmt.Println("💥 BUST!")
			}
		}

		// Show final result
		fmt.Println(game.RenderResult(g))

		// Check if game is over
		if g.Bank <= 0 {
			fmt.Println("\n💸 You're busted. Thanks for playing!")
			break
		}

		// Continue prompt
		if !promptContinue() {
			break
		}
	}

	// Final bank
	fmt.Printf("\n🏦 Final Bank: %d chips\n", g.Bank)
	fmt.Println("\nThanks for playing!")
}

func promptContinue() bool {
	cont, err := game.PromptYesNo(os.Stdin, "\nPlay another hand?")
	if err != nil {
		return false
	}
	return cont
}
