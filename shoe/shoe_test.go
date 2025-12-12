package shoe

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Psiloscop/Blackjack-Go-2025/card"
)

func Test_New(t *testing.T) {
	shoe := New(1, 75)

	expectedDeckAmount := uint(1)
	expectedDeck := [deckSize]card.Card{
		card.New(card.Clubs, card.Two),
		card.New(card.Clubs, card.Three),
		card.New(card.Clubs, card.Four),
		card.New(card.Clubs, card.Five),
		card.New(card.Clubs, card.Six),
		card.New(card.Clubs, card.Seven),
		card.New(card.Clubs, card.Eight),
		card.New(card.Clubs, card.Nine),
		card.New(card.Clubs, card.Ten),
		card.New(card.Clubs, card.Jack),
		card.New(card.Clubs, card.Queen),
		card.New(card.Clubs, card.King),
		card.New(card.Clubs, card.Ace),
		card.New(card.Diamonds, card.Two),
		card.New(card.Diamonds, card.Three),
		card.New(card.Diamonds, card.Four),
		card.New(card.Diamonds, card.Five),
		card.New(card.Diamonds, card.Six),
		card.New(card.Diamonds, card.Seven),
		card.New(card.Diamonds, card.Eight),
		card.New(card.Diamonds, card.Nine),
		card.New(card.Diamonds, card.Ten),
		card.New(card.Diamonds, card.Jack),
		card.New(card.Diamonds, card.Queen),
		card.New(card.Diamonds, card.King),
		card.New(card.Diamonds, card.Ace),
		card.New(card.Hearts, card.Two),
		card.New(card.Hearts, card.Three),
		card.New(card.Hearts, card.Four),
		card.New(card.Hearts, card.Five),
		card.New(card.Hearts, card.Six),
		card.New(card.Hearts, card.Seven),
		card.New(card.Hearts, card.Eight),
		card.New(card.Hearts, card.Nine),
		card.New(card.Hearts, card.Ten),
		card.New(card.Hearts, card.Jack),
		card.New(card.Hearts, card.Queen),
		card.New(card.Hearts, card.King),
		card.New(card.Hearts, card.Ace),
		card.New(card.Spades, card.Two),
		card.New(card.Spades, card.Three),
		card.New(card.Spades, card.Four),
		card.New(card.Spades, card.Five),
		card.New(card.Spades, card.Six),
		card.New(card.Spades, card.Seven),
		card.New(card.Spades, card.Eight),
		card.New(card.Spades, card.Nine),
		card.New(card.Spades, card.Ten),
		card.New(card.Spades, card.Jack),
		card.New(card.Spades, card.Queen),
		card.New(card.Spades, card.King),
		card.New(card.Spades, card.Ace),
	}

	if shoe.deck != expectedDeck {
		t.Error("Decks do not match")
	}
	if shoe.deckAmount != expectedDeckAmount {
		t.Errorf("Expected deckAmount to be %d, got %d", expectedDeckAmount, shoe.deckAmount)
	}
	if uint(len(shoe.sequence)) != deckSize {
		t.Errorf("Expected sequence length to be %d, got %d", deckSize, len(shoe.sequence))
	}
}

func Test_Refill(t *testing.T) {
	shoe := New(1, 75)

	for cardsToTake := range [10]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		t.Run(fmt.Sprintf("Taking %d cards", cardsToTake), func(t *testing.T) {
			for i := 0; i < cardsToTake; i++ {
				_ = shoe.GetNextCard()
			}

			expectedCursorPosition := uint(cardsToTake)
			if shoe.cursor != expectedCursorPosition {
				t.Errorf("Cursor should be %d after %d calls to GetNextCard", expectedCursorPosition, cardsToTake)
			}

			shoe.Refill()

			if shoe.cursor != 0 {
				t.Errorf("Cursor should be 0 after Refill")
			}
		})
	}
}

func Test_Shuffle(t *testing.T) {
	shoe := New(1, 75)

	expectedSequence := make([]uint, deckSize)
	for i := uint(0); i < deckSize; i++ {
		expectedSequence[i] = i
	}

	if !slices.Equal(shoe.sequence, expectedSequence) {
		t.Errorf("Expected sequence to be %v, got %v", expectedSequence, shoe.sequence)
	}
}

func Test_ShuffleDiscardTray(t *testing.T) {
	shoe := New(1, 75)

	expectedSequence := make([]uint, deckSize)
	for i := uint(0); i < deckSize; i++ {
		expectedSequence[i] = i
	}

	shoe.Shuffle()

	leftPart := shoe.sequence[:shoe.cutCardIndex]
	rightPart := shoe.sequence[shoe.cutCardIndex:]
	leftPartCopy := make([]uint, len(leftPart))
	rightPartCopy := make([]uint, len(rightPart))
	copy(leftPartCopy, leftPart)
	copy(rightPartCopy, rightPart)

	shoe.ShuffleDiscardTray()

	if slices.Equal(leftPart, leftPartCopy) {
		t.Errorf("Expected left part of sequence to be different from %v, got %v", leftPartCopy, leftPart)
	}
	if !slices.Equal(rightPart, rightPartCopy) {
		t.Errorf("Expected right part of sequence to be equal %v, got %v", rightPartCopy, rightPart)
	}
}

func Test_GetNextCard(t *testing.T) {
	testCases := []struct {
		name                   string
		deckAmount             uint
		deckPenetrationPercent uint
		cardAmountToTake       uint
		cursorReversed         bool
		initCursorPosition     uint
		expectedCursorPosition uint
		expectedCutCardReached bool
	}{
		{
			name:                   "Starting from beginning of shoe and not reaching cut card",
			deckAmount:             1,
			deckPenetrationPercent: 75,
			cardAmountToTake:       20,
			cursorReversed:         false,
			initCursorPosition:     0,
			expectedCursorPosition: 20,
			expectedCutCardReached: false,
		},
		{
			name:                   "Starting from cut card and going back and not reaching beginning of shoe",
			deckAmount:             1,
			deckPenetrationPercent: 75,
			cardAmountToTake:       20,
			cursorReversed:         true,
			initCursorPosition:     30,
			expectedCursorPosition: 10,
			expectedCutCardReached: false,
		},
		{
			name:                   "Staring from the middle of the shoe and taking more that shoe has and reaching cut card",
			deckAmount:             1,
			deckPenetrationPercent: 75,
			cardAmountToTake:       35,
			cursorReversed:         false,
			initCursorPosition:     30,
			expectedCursorPosition: 38,
			expectedCutCardReached: true,
		},
		{
			name:                   "Starting from cut card and going back and taking more that discard tray has and reaching beginning of shoe",
			deckAmount:             1,
			deckPenetrationPercent: 75,
			cardAmountToTake:       35,
			cursorReversed:         true,
			initCursorPosition:     30,
			expectedCursorPosition: 0,
			expectedCutCardReached: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shoe := New(tc.deckAmount, tc.deckPenetrationPercent)
			shoe.cursor = tc.initCursorPosition
			shoe.cursorReversed = tc.cursorReversed

			for i := uint(0); i < tc.cardAmountToTake; i++ {
				_ = shoe.GetNextCard()
			}

			if shoe.cursor != tc.expectedCursorPosition {
				t.Errorf(
					"Cursor should be %d after calling GetNextCard %d times but it's %d",
					tc.expectedCursorPosition, tc.cardAmountToTake, shoe.cursor,
				)
			}
			if shoe.isCutCardReached != tc.expectedCutCardReached {
				t.Errorf("Cut card must be reached %t", tc.expectedCutCardReached)
			}
		})
	}
}
