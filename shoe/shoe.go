package shoe

import (
	"math/rand"

	"github.com/Psiloscop/Blackjack-Go-2025/card"
)

const deckSize = uint(52) // todo некоторые вариации блекджека имеют неполную колоду

type Shoe struct {
	deckAmount       uint
	deck             [deckSize]card.Card
	sequence         []uint
	cursor           uint
	cursorReversed   bool
	cutCardIndex     uint
	isShuffled       bool
	isCutCardReached bool
}

func New(deckAmount, deckPenetrationPercent uint) *Shoe {
	if deckAmount == 0 {
		panic("Deck amount cannot be zero")
	}
	if deckPenetrationPercent < 10 || deckPenetrationPercent > 100 {
		panic("Deck penetration percentage cannot be greater than 100 or less than 30")
	}

	cardAmount := deckSize * deckAmount
	penetration := float64(deckPenetrationPercent) / 100.0
	remainingCards := uint(float64(cardAmount) * (1 - penetration))

	var shoe Shoe
	shoe.deckAmount = deckAmount
	shoe.cutCardIndex = cardAmount - remainingCards - 1

	var deckIndex uint8
	for suit := range [4]card.Suit{
		card.Clubs,
		card.Diamonds,
		card.Hearts,
		card.Spades,
	} {
		for number := range [13]card.Number{
			card.Two,
			card.Three,
			card.Four,
			card.Five,
			card.Six,
			card.Seven,
			card.Eight,
			card.Nine,
			card.Ten,
			card.Jack,
			card.Queen,
			card.King,
			card.Ace,
		} {
			shoe.deck[deckIndex] = card.New(card.Suit(suit), card.Number(number))
			deckIndex++
		}
	}

	shoe.Refill()

	return &shoe
}

func (s *Shoe) Refill() {
	s.isCutCardReached = false
	s.cursorReversed = false
	s.cursor = 0

	if s.isShuffled || len(s.sequence) == 0 {
		cardAmount := deckSize * s.deckAmount
		s.sequence = make([]uint, cardAmount)
		for i := uint(0); i < s.deckAmount; i++ {
			for j := uint(0); j < deckSize; j++ {
				deckNumber := deckSize * i
				s.sequence[deckNumber+j] = j
			}
		}

		s.isShuffled = false
	}
}

func (s *Shoe) Shuffle() {
	for i := len(s.sequence) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s.sequence[i], s.sequence[j] = s.sequence[j], s.sequence[i]
	}

	s.isShuffled = true
}

func (s *Shoe) ShuffleDiscardTray() {
	for i := len(s.sequence) - int(s.cutCardIndex) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s.sequence[i], s.sequence[j] = s.sequence[j], s.sequence[i]
	}
}

func (s *Shoe) ToggleGettingCardFromDiscardTrayMode() {
	s.cursorReversed = !s.cursorReversed

	if s.cursor == s.cutCardIndex {
		s.cursor-- // move the cursor back away from the cut card
	}
}

func (s *Shoe) IsCutCardReached() bool {
	return s.isCutCardReached
}

func (s *Shoe) IsDiscardTrayMode() bool {
	return s.cursorReversed
}

func (s *Shoe) IsDiscardTrayEnded() bool {
	return s.cursorReversed && s.cursor == 0
}

func (s *Shoe) GetNextCard() card.Card {
	c := s.deck[s.sequence[s.cursor]]

	//if s.cursorReversed {
	//	if s.cursor == 0 {
	//		panic("Discard tray ended.")
	//	}
	//
	//	s.cursor--
	//} else {
	//	if s.cursor == s.cutCardIndex {
	//		panic("Cut card reached.")
	//	}
	//
	//	s.cursor++
	//
	//	if s.cursor == s.cutCardIndex {
	//		s.isCutCardReached = true
	//	}
	//}

	if s.cursor == s.cutCardIndex {
		s.isCutCardReached = true
	} else {
		if s.cursorReversed {
			if s.cursor != 0 {
				s.cursor--
			}
		} else {
			s.cursor++
		}
	}

	return c
}
