package hand

import (
	"fmt"
	"testing"

	"github.com/Psiloscop/Blackjack-Go-2025/card"
)

func Test_NewDealerHand(t *testing.T) {
	_hand := NewDealerHand()
	hand := _hand.(*DealerHand)

	if len(hand.hand.Cards) != 0 {
		t.Error("Dealer hand should be empty after creation")
	}
}

func Test_DealerHand_IsSecondCardAce(t *testing.T) {
	testCases := []struct {
		cards                 []card.Card
		expectedSecondCardAce bool
	}{
		{
			[]card.Card{
				card.New(card.Hearts, card.Queen),
				card.New(card.Hearts, card.Ace),
			},
			true,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Ace),
				card.New(card.Hearts, card.King),
			},
			false,
		},
		{
			[]card.Card{
				card.New(card.Hearts, card.Ace),
				card.New(card.Hearts, card.Ace),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Expected second card ace = %t", testCase.expectedSecondCardAce), func(t *testing.T) {
			hand := NewDealerHand()
			for _, c := range testCase.cards {
				hand.AddCard(c)
			}

			if hand.IsSecondCardAce() != testCase.expectedSecondCardAce {
				t.Errorf("Expected second card ace, but it's not. Score is %d", hand.GetScore())
			}
		})
	}
}
