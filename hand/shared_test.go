package hand

import (
	"fmt"
	"testing"

	"github.com/Psiloscop/Blackjack-Go-2025/card"
)

func Test_Hand_AddCard(t *testing.T) {
	hand := hand{[]card.Card{}}
	if len(hand.Cards) != 0 {
		t.Error("hand should be empty")
	}

	hand.AddCard(card.New(card.Clubs, card.Ace))
	if len(hand.Cards) != 1 {
		t.Error("hand should contain one card")
	}

	hand.AddCard(card.New(card.Diamonds, card.Ace))
	if len(hand.Cards) != 2 {
		t.Error("hand should contain two cards")
	}
}

func Test_Hand_GetScore(t *testing.T) {
	testCases := []struct {
		cards         []card.Card
		expectedScore uint
	}{
		{
			[]card.Card{
				card.New(card.Hearts, card.Two),
				card.New(card.Hearts, card.Five),
				card.New(card.Hearts, card.Eight),
			},
			15,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Two),
				card.New(card.Hearts, card.Five),
				card.New(card.Hearts, card.Eight),
				card.New(card.Hearts, card.Ace),
			},
			16,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Ace),
			},
			11,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Ace),
				card.New(card.Hearts, card.Ace),
			},
			12,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Ten),
				card.New(card.Hearts, card.Ace),
				card.New(card.Hearts, card.Ace),
			},
			12,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Ace),
				card.New(card.Hearts, card.Ace),
				card.New(card.Hearts, card.Ace),
			},
			13,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Seven),
				card.New(card.Hearts, card.Four),
				card.New(card.Hearts, card.Ace),
			},
			22,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Seven),
				card.New(card.Hearts, card.Six),
				card.New(card.Hearts, card.Ace),
			},
			24,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Ten),
				card.New(card.Hearts, card.Ace),
			},
			score21,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Ace),
			},
			score21,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Queen),
				card.New(card.Hearts, card.Ace),
			},
			score21,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.King),
				card.New(card.Hearts, card.Ace),
			},
			score21,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Expecting %d score", testCase.expectedScore), func(t *testing.T) {
			hand := hand{testCase.cards}

			if hand.GetScore() != testCase.expectedScore {
				t.Errorf("Expected score to be %d, got %d", testCase.expectedScore, hand.GetScore())
			}
		})
	}
}

func Test_Hand_IsBust(t *testing.T) {
	testCases := []struct {
		cards        []card.Card
		expectedBust bool
	}{
		{
			[]card.Card{
				card.New(card.Hearts, card.Two),
				card.New(card.Hearts, card.Five),
				card.New(card.Hearts, card.Eight),
			},
			false,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Queen),
				card.New(card.Hearts, card.King),
			},
			true,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Ace),
			},
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Expected bust = %t", testCase.expectedBust), func(t *testing.T) {
			hand := hand{testCase.cards}

			if hand.IsBust() != testCase.expectedBust {
				t.Errorf("Expected bust, but it's not. Score is %d", hand.GetScore())
			}
		})
	}
}

func Test_Hand_IsBlackjack(t *testing.T) {
	testCases := []struct {
		cards             []card.Card
		expectedBlackjack bool
	}{
		{
			[]card.Card{
				card.New(card.Hearts, card.Ten),
				card.New(card.Hearts, card.Eight),
				card.New(card.Hearts, card.Two),
			},
			false,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Queen),
				card.New(card.Hearts, card.Ace),
			},
			false,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Ace),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Expected blackjack = %t", testCase.expectedBlackjack), func(t *testing.T) {
			hand := hand{testCase.cards}

			if hand.IsBlackjack() != testCase.expectedBlackjack {
				t.Errorf("Expected blackjack, but it's not. Score is %d", hand.GetScore())
			}
		})
	}
}
