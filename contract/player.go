package contract

type Player interface {
	GetId() PlayerId
	GetName() string
	GetPurse() uint
	HasMoney() bool
	IncreasePurse(amount uint)
	DecreasePurse(amount uint)
	PlaceBet() uint
	ChooseAction(ctx *PlayerContext) GameAction
}

type PlayerId string

type PlayerDetails struct {
	PlayerId PlayerId
	HandId   PlayerHandId
}
