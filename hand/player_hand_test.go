package hand

import (
	"fmt"
	"testing"

	"github.com/Psiloscop/Blackjack-Go-2025/card"
)

func Test_NewPlayerHand(t *testing.T) {
	expectedBet := uint(777)
	_hand := NewPlayerHand(expectedBet)
	hand := _hand.(*PlayerHand)

	if hand.bet != expectedBet {
		t.Errorf("bet should be %d, got %d", expectedBet, hand.bet)
	}
}

func Test_PlayerHand_IsSplittable(t *testing.T) {
	testCases := []struct {
		cards              []card.Card
		expectedSplittable bool
	}{
		{
			[]card.Card{
				card.New(card.Hearts, card.Three),
				card.New(card.Hearts, card.Three),
			},
			true,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Ten),
				card.New(card.Hearts, card.Ten),
			},
			true,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Jack),
				card.New(card.Hearts, card.Queen),
			},
			false,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Queen),
				card.New(card.Hearts, card.King),
			},
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Expected splittable = %t", testCase.expectedSplittable), func(t *testing.T) {
			hand := NewPlayerHand(777)
			for _, c := range testCase.cards {
				hand.AddCard(c)
			}

			if hand.IsSplittable() != testCase.expectedSplittable {
				t.Errorf("Expected splittable, but it's not. Score is %d", hand.GetScore())
			}
		})
	}
}

func Test_PlayerHand_Split(t *testing.T) {
	_hand := NewPlayerHand(777)
	hand := _hand.(*PlayerHand)
	hand.AddCard(card.New(card.Hearts, card.Jack))
	hand.AddCard(card.New(card.Hearts, card.Jack))

	if !hand.IsSplittable() {
		t.Error("hand should be splittable, but it's not")
	}

	splitCard := hand.Split()

	if len(hand.Cards) != 1 {
		t.Error("hand should contain only one card after split")
	}

	if splitCard.Number != hand.Cards[0].Number {
		t.Errorf("Split cards should be the same, but got %d and %d", splitCard.Number, hand.Cards[0].Number)
	}

	hand.AddCard(card.New(card.Hearts, card.King))
	if len(hand.Cards) != 2 || hand.Cards[0].Number != card.Jack || hand.Cards[1].Number != card.King {
		t.Error("Something went wrong with hand splitting.")
	}
}
