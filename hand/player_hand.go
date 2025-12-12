package hand

import (
	"github.com/Psiloscop/Blackjack-Go-2025/card"
	"github.com/Psiloscop/Blackjack-Go-2025/contract"
	"github.com/google/uuid"
)

type PlayerHand struct {
	hand
	id  contract.PlayerHandId
	bet uint
}

func NewPlayerHand(bet uint) contract.PlayerHand {
	return &PlayerHand{
		hand: hand{[]card.Card{}},
		id:   contract.PlayerHandId(uuid.New().String()),
		bet:  bet,
	}
}

func (h *PlayerHand) GetId() contract.PlayerHandId {
	return h.id
}

func (h *PlayerHand) GetBet() uint {
	return h.bet
}

func (h *PlayerHand) IsSplittable() bool {
	return len(h.Cards) == 2 && h.Cards[0].Number == h.Cards[1].Number && h.Cards[0].Number != card.Ace
}

func (h *PlayerHand) Split() card.Card {
	splitCard := h.Cards[1]

	h.Cards = h.Cards[:1]

	return splitCard
}

func (h *PlayerHand) IncreaseBet(amount uint) {
	h.bet += amount
}
