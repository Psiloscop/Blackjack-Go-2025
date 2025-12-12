package table

import (
	"github.com/Psiloscop/Blackjack-Go-2025/contract"
)

type TestMind struct {
	BetToPlace     uint
	ActionToChoose contract.GameAction
}

func NewTestMind() contract.Mind {
	return &TestMind{}
}

func (a TestMind) PlaceBet(p contract.Player) uint {
	return a.BetToPlace
}

func (a TestMind) ChooseAction(ctx *contract.PlayerContext) contract.GameAction {
	return a.ActionToChoose
}
