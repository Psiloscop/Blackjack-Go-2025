package contract

type PlayerContext struct {
	availableActions []GameAction
	dealerScore      uint
	playerScore      uint
}

func NewPlayerContext(availableActions []GameAction, dealerScore uint, playerScore uint) *PlayerContext {
	return &PlayerContext{availableActions, dealerScore, playerScore}
}

func (ctx PlayerContext) GetAvailableActions() []GameAction {
	return ctx.availableActions
}

func (ctx PlayerContext) GetDealerScore() uint {
	return ctx.dealerScore
}

func (ctx PlayerContext) GetPlayerScore() uint {
	return ctx.playerScore
}
