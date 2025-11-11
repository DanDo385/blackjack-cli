package game

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// PromptBet prompts the user for a bet amount
func PromptBet(reader io.Reader, bank int) (int, error) {
	scanner := bufio.NewScanner(reader)

	for {
		fmt.Printf("Enter bet (1-%d): ", bank)
		if !scanner.Scan() {
			return 0, fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(scanner.Text())
		bet, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if bet < MinBet {
			fmt.Printf("Minimum bet is %d.\n", MinBet)
			continue
		}

		if bet > bank {
			fmt.Printf("Bet exceeds bank balance (%d).\n", bank)
			continue
		}

		return bet, nil
	}
}

// PromptAction prompts the user for an action
func PromptAction(reader io.Reader, actions []Action, handNum int, totalHands int) (Action, error) {
	scanner := bufio.NewScanner(reader)

	for {
		fmt.Print(RenderAvailableActions(actions, handNum, totalHands) + ": ")
		if !scanner.Scan() {
			return 0, fmt.Errorf("failed to read input")
		}

		input := strings.ToLower(strings.TrimSpace(scanner.Text()))

		// Parse action
		var action Action
		var valid bool

		switch input {
		case "h", "hit":
			action = ActionHit
			valid = containsAction(actions, ActionHit)
		case "s", "stand":
			action = ActionStand
			valid = containsAction(actions, ActionStand)
		case "d", "double":
			action = ActionDouble
			valid = containsAction(actions, ActionDouble)
		case "p", "split":
			action = ActionSplit
			valid = containsAction(actions, ActionSplit)
		case "r", "surrender":
			action = ActionSurrender
			valid = containsAction(actions, ActionSurrender)
		default:
			fmt.Println("Invalid action. Please try again.")
			continue
		}

		if !valid {
			fmt.Println("Action not available. Please choose from available actions.")
			continue
		}

		return action, nil
	}
}

// PromptYesNo prompts the user for a yes/no answer
func PromptYesNo(reader io.Reader, prompt string) (bool, error) {
	scanner := bufio.NewScanner(reader)

	for {
		fmt.Print(prompt + " (y/n): ")
		if !scanner.Scan() {
			return false, fmt.Errorf("failed to read input")
		}

		input := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch input {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
			continue
		}
	}
}

// PromptInsurance prompts the user for an insurance bet
func PromptInsurance(reader io.Reader, maxInsurance int) (int, error) {
	scanner := bufio.NewScanner(reader)

	for {
		fmt.Printf("Insurance bet (0-%d): ", maxInsurance)
		if !scanner.Scan() {
			return 0, fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(scanner.Text())
		bet, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if bet < 0 {
			fmt.Println("Insurance bet cannot be negative.")
			continue
		}

		if bet > maxInsurance {
			fmt.Printf("Insurance bet cannot exceed %d.\n", maxInsurance)
			continue
		}

		return bet, nil
	}
}

func containsAction(actions []Action, action Action) bool {
	for _, a := range actions {
		if a == action {
			return true
		}
	}
	return false
}
