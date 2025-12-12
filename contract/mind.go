package contract

type Mind interface {
	PlaceBet(p Player) uint
	ChooseAction(ctx *PlayerContext) GameAction
}
