package player

import (
	"errors"
	"slices"

	"github.com/Psiloscop/Blackjack-Go-2025/contract"
)

type InteractableMind struct {
	requestBet    func(p contract.Player) uint
	requestAction func(ctx *contract.PlayerContext) contract.GameAction
	sendError     func(error)
}

func NewInteractableMind(
	requestBet func(p contract.Player) uint,
	requestAction func(ctx *contract.PlayerContext) contract.GameAction,
	sendError func(error),
) contract.Mind {
	return &InteractableMind{
		requestBet,
		requestAction,
		sendError,
	}
}

func (a InteractableMind) PlaceBet(p contract.Player) uint {
	for {
		response := a.requestBet(p)

		if response == 0 {
			a.sendError(errors.New("bet cannot be 0"))
		} else if response > p.GetPurse() {
			a.sendError(errors.New("not enough money in purse"))
		} else {
			return response
		}
	}
}

func (a InteractableMind) ChooseAction(ctx *contract.PlayerContext) contract.GameAction {
	for {
		action := a.requestAction(ctx)
		if !action.IsValid() {
			a.sendError(errors.New("unsupported action chosen"))

			continue
		}

		if !slices.Contains(ctx.GetAvailableActions(), action) {
			a.sendError(errors.New("chosen action is not available"))

			continue
		}

		return action
	}
}
