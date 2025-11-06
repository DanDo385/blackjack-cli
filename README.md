# Blackjack CLI

A fully-featured command-line Blackjack game written in Go.

## Features

- Classic Blackjack gameplay with complete rule implementation
- Starting bank of 1000 chips
- Single 52-card deck, reshuffled after every hand
- Full player actions: Hit, Stand, Double, Split, Surrender
- Insurance when dealer shows Ace
- Advanced rules:
  - Dealer stands on soft 17 (S17)
  - Late surrender (before first action)
  - Double after split (except split aces)
  - Split up to 4 hands (resplit any pair except aces)
  - Split aces receive one card only
  - Blackjack after split counts as 21 (not natural)
- Payouts:
  - Natural Blackjack: 3:2
  - Insurance: 2:1
  - Regular Win: 1:1
  - Push: returns bet
  - Surrender: returns half bet

## Requirements

- Go 1.22 or higher

## Installation

### Using Make

```bash
make build
```

This creates the executable at `./bin/blackjack`.

### Using Scripts

```bash
./scripts/build.sh
```

### Manual Build

```bash
go build -o ./bin/blackjack ./cmd/blackjack
```

## Running the Game

### Using Make

```bash
make run
```

### Using Scripts

```bash
./scripts/run.sh
```

### Direct Execution

```bash
./bin/blackjack
```

## How to Play

1. The game starts with a bank of 1000 chips
2. Enter your bet amount (minimum 1 chip, maximum your current bank)
3. Cards are dealt: 2 to you, 2 to dealer (one face down)
4. If dealer shows an Ace, you'll be offered insurance
5. Choose from available actions:
   - **(H)it**: Take another card
   - **(S)tand**: Keep your current hand
   - **(D)ouble**: Double your bet and receive exactly one more card
   - **(P)split**: Split a pair into two separate hands (up to 4 total hands)
   - **(R)surrender**: Forfeit the hand and recover half your bet
6. After all your hands are played, the dealer reveals their hole card and plays
7. Winnings are calculated and added to your bank
8. Continue playing until you run out of chips or choose to quit

## Game Rules Summary

- **Dealer**: Stands on all 17s (including soft 17)
- **Blackjack**: Natural 21 with first two cards pays 3:2
- **Insurance**: Offered when dealer shows Ace; costs up to half your bet; pays 2:1 if dealer has blackjack
- **Double Down**: Available on first action only (except split aces)
- **Splitting**:
  - Can split any pair of same rank
  - Can resplit to a maximum of 4 hands
  - Aces can only be split once and receive one card each
  - No doubling on split aces
  - 21 after split is not blackjack
- **Surrender**: Late surrender only (before first action); recovers half your bet

## Testing

### Run All Tests

```bash
make test
```

This runs all tests with race detection.

### Run Specific Test Suites

```bash
go test ./internal/game -v
```

### Deterministic Testing

The game supports deterministic testing using a seeded shoe. Set the `BLACKJACK_SEEDED=1` environment variable to use a predetermined card sequence from `internal/game/testdata/seeded_shoe.txt`:

```bash
BLACKJACK_SEEDED=1 go test ./internal/game -run Engine -v
```

This is useful for testing specific scenarios like:
- Player/Dealer blackjack combinations
- Insurance outcomes
- Split mechanics (including split aces)
- Double down scenarios
- Surrender situations
- Dealer S17 behavior

## Test Coverage

The test suite includes:

- **Deck Tests**: Verify deck creation, shuffling, and card drawing
- **Hand Tests**: Validate hand value calculations, soft/hard totals, blackjack detection
- **Rules Tests**: Confirm payouts, outcomes, and dealer behavior
- **Engine Tests**: End-to-end gameplay scenarios including:
  - Blackjack payouts (3:2)
  - Insurance mechanics (2:1 payout)
  - Split hands (including aces-only-one-card rule)
  - Double down
  - Late surrender
  - Dealer stands on soft 17
  - Maximum 4 hands after splitting

## Code Formatting

```bash
make fmt
```

This runs `gofmt` and `go vet` on all Go files.

## Project Structure

```
blackjack-cli/
├── cmd/
│   └── blackjack/
│       └── main.go           # CLI entry point
├── internal/
│   └── game/
│       ├── card.go           # Card, Suit, Rank types
│       ├── deck.go           # Deck creation and shuffling
│       ├── hand.go           # Hand logic and calculations
│       ├── rules.go          # Game rules and payouts
│       ├── dealer.go         # Dealer behavior
│       ├── game.go           # Main game engine
│       ├── cli_renderer.go   # ASCII rendering
│       ├── input.go          # User input handling
│       ├── rng.go            # Random number generation
│       ├── *_test.go         # Test files
│       └── testdata/
│           └── seeded_shoe.txt
├── scripts/
│   ├── build.sh
│   └── run.sh
├── go.mod
├── Makefile
└── README.md
```

## Architecture

The game is built with clean separation of concerns:

- **Pure Game Engine** (`internal/game`): Deterministic, testable logic with no I/O dependencies
- **Thin CLI Layer** (`cmd/blackjack`): Handles user interaction and rendering
- **Testability**: All game logic can be tested without console I/O
- **Deterministic Testing**: Supports seeded RNG for reproducible test scenarios

## License

This project is provided as-is for educational and entertainment purposes.
