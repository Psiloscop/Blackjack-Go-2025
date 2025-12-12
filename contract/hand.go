package contract

import "github.com/Psiloscop/Blackjack-Go-2025/card"

type hand interface {
	AddCard(c card.Card)
	GetCards() []card.Card
	GetScore() uint
	GetCardAmount() uint
	IsBust() bool
	IsBlackjack() bool
}
