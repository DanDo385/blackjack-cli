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

		// Start hand (bet is validated but not yet deducted)
		err = g.StartHand(bet)
		if err != nil {
			fmt.Printf("Error starting hand: %v\n", err)
			continue
		}

		// Show initial state
		fmt.Println()
		fmt.Println(game.RenderState(g, true))
		fmt.Println()

		// Deduct the bet after showing the initial state
		// This allows the doubling check to see the full bank balance
		g.Bank -= bet

	
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

		// Player action phase - play each hand independently
		for g.CurrentPhase == game.PhasePlayerAction {
			currentHand := g.GetCurrentHand()
			if currentHand == nil {
				break
			}

			// Capture the current hand index at the start of this iteration
			currentHandIndex := g.ActiveHandIndex
			handNum := currentHandIndex + 1
			totalHands := len(g.PlayerHands)

			// Display current hand info
			fmt.Println()
			fmt.Println(game.RenderCurrentHand(g))
			fmt.Println()

			// Handle split aces - they only get one card and then move on
			if currentHand.IsSplitAces {
				// Split aces already have their one card dealt during split
				// Just show the result and advance
				fmt.Println("Split aces receive only one card.")
				fmt.Println()
				fmt.Println(game.RenderState(g, true))
				fmt.Println()
				
				// Advance to next hand
				err := g.PlayerAction(game.ActionStand)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					break
				}
				continue
			}

			// Play this hand until it's done (stand, bust, double, or surrender)
			// Continue as long as we're still on the same hand index
			for g.CurrentPhase == game.PhasePlayerAction && g.ActiveHandIndex == currentHandIndex {
				// Refresh current hand reference
				currentHand = g.GetCurrentHand()
				if currentHand == nil {
					break
				}

				// Refresh total hands count in case we split
				totalHands = len(g.PlayerHands)

				// Check if hand is bust
				if currentHand.IsBust() {
					fmt.Println("💥 BUST!")
					fmt.Println()
					fmt.Println(game.RenderState(g, true))
					fmt.Println()
					
					err := g.PlayerAction(game.ActionStand)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						break
					}
					break
				}

				// Check if hand is 21 (auto-stand)
				if currentHand.Value() == 21 {
					fmt.Println()
					fmt.Println(game.RenderState(g, true))
					fmt.Println()
					
					err := g.PlayerAction(game.ActionStand)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						break
					}
					break
				}

				// Get available actions
				actions := g.GetAvailableActions()
				if len(actions) == 0 {
					break
				}

				// Show board state
				fmt.Println(game.RenderState(g, true))
				fmt.Println()

				// Prompt for action
				action, err := game.PromptAction(os.Stdin, actions, handNum, totalHands)
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

				// Show result of action
				currentHand = g.GetCurrentHand()
				if currentHand != nil {
					if action == game.ActionHit && currentHand.IsBust() {
						fmt.Println("💥 BUST!")
					} else if action == game.ActionDouble {
						// Double ends the hand, show result
						if currentHand.IsBust() {
							fmt.Println("💥 BUST!")
						}
					} else if action == game.ActionSurrender {
						fmt.Println("Hand surrendered.")
					}
				}

				// If action was stand, double, or surrender, the hand is done
				// (advanceToNextHand was called, so we break out of inner loop)
				if action == game.ActionStand || action == game.ActionDouble || action == game.ActionSurrender {
					break
				}

				// If we split, we continue playing the current hand (first split hand)
				// Refresh the display for the next iteration
				if action == game.ActionSplit {
					fmt.Println()
					fmt.Println(game.RenderCurrentHand(g))
					fmt.Println()
				}
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
