package main

import (
	"github.com/Psiloscop/Blackjack-Go-2025/console"
	"github.com/Psiloscop/Blackjack-Go-2025/contract"
	"github.com/Psiloscop/Blackjack-Go-2025/hand"
	"github.com/Psiloscop/Blackjack-Go-2025/player"
	"github.com/Psiloscop/Blackjack-Go-2025/shoe"
	"github.com/Psiloscop/Blackjack-Go-2025/table"
)

func main() {
	interactableMind := player.NewInteractableMind(
		console.RequestBet,
		console.RequestAction,
		console.SendError,
	)
	pl1 := player.New("Pierce Brosnan", 25000, &interactableMind)
	//pl2 := player.New("Don Johnson", 25000, &interactableMind)
	players := []*contract.Player{
		&pl1,
		//&pl2,
	}

	americanBlackjack := table.NewAmericanBlackjack(
		shoe.New(6, 85),
		players,
		func() *contract.DealerHand {
			dealerHand := hand.NewDealerHand()

			return &dealerHand
		},
		func(bet uint) *contract.PlayerHand {
			playerHand := hand.NewPlayerHand(bet)

			return &playerHand
		},
		console.DisplayTable,
		console.DisplayMessage,
	)

	americanBlackjack.Play()
}
