package hand

import "github.com/Psiloscop/Blackjack-Go-2025/card"

const ace1 = uint(1)
const ace11 = uint(11)
const score21 = uint(21)

type hand struct {
	Cards []card.Card
}

func (h *hand) AddCard(c card.Card) {
	h.Cards = append(h.Cards, c)
}

func (h *hand) GetCards() []card.Card {
	return h.Cards
}

func (h *hand) GetScore() uint {
	var score uint
	var ace1Count uint

	for _, c := range h.Cards {
		switch c.Number {
		case card.Two:
			score += 2
		case card.Three:
			score += 3
		case card.Four:
			score += 4
		case card.Five:
			score += 5
		case card.Six:
			score += 6
		case card.Seven:
			score += 7
		case card.Eight:
			score += 8
		case card.Nine:
			score += 9
		case card.Ten, card.Jack, card.Queen, card.King:
			score += 10
		case card.Ace:
			ace1Count += ace1
		}
	}

	if ace1Count == 0 {
		return score
	}

	if score >= score21 {
		score += ace1Count
	} else {
		for i := ace1Count; i >= uint(0); i-- {
			nextScore := score + ace1Count - i + i*ace11

			if nextScore <= score21 {
				score = nextScore

				break
			}
		}
	}

	return score
}

func (h *hand) GetCardAmount() uint {
	return uint(len(h.Cards))
}

func (h *hand) IsBust() bool {
	return h.GetScore() > score21
}

func (h *hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.GetScore() == score21
}
