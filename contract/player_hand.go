package contract

import "github.com/Psiloscop/Blackjack-Go-2025/card"

type PlayerHand interface {
	hand
	GetId() PlayerHandId
	GetBet() uint
	IsSplittable() bool
	Split() card.Card
	IncreaseBet(amount uint)
}

type PlayerHandId string

type PlayerHandCreator func(bet uint) *PlayerHand
