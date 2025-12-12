package table

import (
	"slices"

	"github.com/Psiloscop/Blackjack-Go-2025/contract"
)

var gameContext *contract.GameContext

func getAvailableActions(
	player *contract.Player,
	dealerHand *contract.DealerHand,
	playerHand *contract.PlayerHand,
	playerHandAmount uint,
	dealerDoesntHaveBlackjack bool,
) []contract.GameAction {
	actions := make([]contract.GameAction, 0, 6)
	actions = append(actions, contract.GameActionStay)

	if (*playerHand).GetScore() <= 21 {
		actions = append(actions, contract.GameActionHit)
	}
	if (*playerHand).GetCardAmount() <= 2 && (*player).GetPurse() >= (*playerHand).GetBet() {
		actions = append(actions, contract.GameActionDoubleDown)
	}
	if (*playerHand).IsSplittable() && (*player).GetPurse() >= (*playerHand).GetBet() && playerHandAmount < americanBlackjackMaxHands {
		actions = append(actions, contract.GameActionSplit)
	}
	if (*playerHand).GetCardAmount() == 2 && !(*dealerHand).IsSecondCardAce() && playerHandAmount == 1 {
		actions = append(actions, contract.GameActionSurrender)
	}
	if (*dealerHand).IsSecondCardAce() && playerHandAmount == 1 && !dealerDoesntHaveBlackjack {
		actions = append(actions, contract.GameActionInsurance)
	}

	slices.Sort(actions)

	return actions
}
