package hand

import (
	"github.com/Psiloscop/Blackjack-Go-2025/card"
	"github.com/Psiloscop/Blackjack-Go-2025/contract"
)

type DealerHand struct {
	hand
}

func NewDealerHand() contract.DealerHand {
	return &DealerHand{hand{[]card.Card{}}}
}

func (h *DealerHand) IsSecondCardAce() bool {
	if len(h.Cards) < 2 {
		return false
	}

	return h.Cards[1].Number == card.Ace
}
