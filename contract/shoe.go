package contract

import "github.com/Psiloscop/Blackjack-Go-2025/card"

type Shoe interface {
	Refill()
	Shuffle()
	ShuffleDiscardTray()
	GetNextCard() card.Card
	ToggleGettingCardFromDiscardTrayMode()
	IsCutCardReached() bool
	IsDiscardTrayMode() bool
}
