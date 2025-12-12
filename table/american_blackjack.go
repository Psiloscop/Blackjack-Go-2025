package table

import (
	"slices"

	"github.com/Psiloscop/Blackjack-Go-2025/card"
	"github.com/Psiloscop/Blackjack-Go-2025/contract"
)

const americanBlackjackMaxHands = 4 // max 4 hands allowed per player in American blackjack
const americanBlackjackMaxPlayers = 8

type AmericanBlackjack struct {
	shoe                      contract.Shoe
	players                   []*contract.Player
	playerHands               map[contract.PlayerId][]*contract.PlayerHand
	insuredPlayers            []contract.PlayerId
	dealerHand                *contract.DealerHand
	dealerFirstCardFlipped    bool
	dealerDoesntHaveBlackjack bool
	createDealerHand          contract.DealerHandCreator
	createPlayerHand          contract.PlayerHandCreator
	sendGameContext           contract.GameContextSender
	sendMessage               contract.MessageSender
}

func NewAmericanBlackjack(
	shoe contract.Shoe,
	players []*contract.Player,
	createDealerHand contract.DealerHandCreator,
	createPlayerHand contract.PlayerHandCreator,
	sendGameContext contract.GameContextSender,
	sendMessage contract.MessageSender,
) *AmericanBlackjack {
	playerHands := make(map[contract.PlayerId][]*contract.PlayerHand, len(players))

	return &AmericanBlackjack{
		shoe:                   shoe,
		players:                players,
		playerHands:            playerHands,
		dealerHand:             createDealerHand(),
		dealerFirstCardFlipped: true,
		createDealerHand:       createDealerHand,
		createPlayerHand:       createPlayerHand,
		sendGameContext:        sendGameContext,
		sendMessage:            sendMessage,
	}
}

func (abj *AmericanBlackjack) Play() {
	for {
		if len(abj.players) == 0 {
			abj.emitMessage(contract.Message{Text: "No one left in the table. Game finished."})

			break
		}

		abj.prepareShoe()
		abj.requestBets()
		abj.dealCards()
		abj.servePlayers(false)

		if (*abj.dealerHand).IsSecondCardAce() {
			if (*abj.dealerHand).IsBlackjack() {
				abj.completeRoundEarly()
				abj.removeBankruptPlayers()
				abj.clearHands()

				continue
			} else {
				abj.dealerDoesntHaveBlackjack = true
				abj.emitMessage(contract.Message{Text: "Dealer has no blackjack. Continue playing."})
				abj.servePlayers(true)
			}
		}

		abj.takeCardsUpToScore17()
		abj.completeRound()
		abj.removeBankruptPlayers()
		abj.clearHands()
	}
}

func (abj *AmericanBlackjack) prepareShoe() {
	if abj.shoe.IsCutCardReached() {
		abj.shoe.Refill()

		abj.emitMessage(contract.Message{Text: "Refilling and shuffling the shoe because cut card reached in previous round."})
	}

	abj.shoe.Shuffle()
}

func (abj *AmericanBlackjack) clearHands() {
	abj.dealerFirstCardFlipped = true
	abj.insuredPlayers = []contract.PlayerId{}
	abj.dealerHand = abj.createDealerHand()
	for playerId, _ := range abj.playerHands {
		abj.playerHands[playerId] = nil
	}
}

func (abj *AmericanBlackjack) requestBets() {
	for _, player := range abj.players {
		bet := (*player).PlaceBet()
		hand := abj.createPlayerHand(bet)
		abj.playerHands[(*player).GetId()] = []*contract.PlayerHand{hand}
	}

	abj.emitGameContext(contract.PlayerDetails{})
}

func (abj *AmericanBlackjack) dealCards() {
	abj.dealerDoesntHaveBlackjack = false

	for _, player := range abj.players {
		hand := abj.playerHands[(*player).GetId()][0]
		if hand == nil {
			panic("Player hand is nil")
		}

		for i := 0; i < 2; i++ {
			(*hand).AddCard(abj.getCardFromShoe())
		}
	}

	for i := 0; i < 2; i++ {
		(*abj.dealerHand).AddCard(abj.getCardFromShoe())
	}

	//(*abj.dealerHand).AddCard(card.Card{card.Clubs, card.Ten})
	//(*abj.dealerHand).AddCard(card.Card{card.Clubs, card.Three})
	//(*abj.dealerHand).AddCard(card.Card{card.Spades, card.Ace})

	abj.emitGameContext(contract.PlayerDetails{})
}

func (abj *AmericanBlackjack) servePlayers(insuredOnly bool) {
	for _, player := range abj.players {
		if insuredOnly && !slices.Contains(abj.insuredPlayers, (*player).GetId()) {
			continue
		}

		playerId := (*player).GetId()
		for i := 0; i < len(abj.playerHands[playerId]); i++ {
			playerHand := abj.playerHands[playerId][i]
			playerHandId := (*playerHand).GetId()
			for {
				abj.emitGameContext(contract.PlayerDetails{PlayerId: playerId, HandId: playerHandId})
				abj.emitMessage(contract.Message{Text: "Choose action.", PlayerId: playerId, HandId: playerHandId})

				playerHandAmount := uint(len(abj.playerHands[playerId]))
				availableActions := getAvailableActions(
					player,
					abj.dealerHand,
					playerHand,
					playerHandAmount,
					abj.dealerDoesntHaveBlackjack,
				)
				selectedAction := contract.GameActionStay
				if len(availableActions) > 1 {
					playerCtx := contract.NewPlayerContext(
						availableActions,
						(*abj.dealerHand).GetScore(),
						(*playerHand).GetScore(),
					)
					selectedAction = (*player).ChooseAction(playerCtx)
				}

				continuePlaying := abj.applyAction(selectedAction, player, playerHand)
				if !continuePlaying {
					break
				}
			}
		}
	}
}

func (abj *AmericanBlackjack) takeCardsUpToScore17() {
	abj.dealerFirstCardFlipped = false

	for (*abj.dealerHand).GetScore() < 17 {
		(*abj.dealerHand).AddCard(abj.getCardFromShoe())
	}

	abj.emitGameContext(contract.PlayerDetails{})
}

func (abj *AmericanBlackjack) completeRoundEarly() {
	abj.dealerFirstCardFlipped = false
	abj.emitGameContext(contract.PlayerDetails{})

	for _, player := range abj.players {
		playerId := (*player).GetId()
		if slices.Contains(abj.insuredPlayers, (*player).GetId()) {
			hand := abj.playerHands[(*player).GetId()][0]
			handId := (*hand).GetId()
			(*player).IncreasePurse((*hand).GetBet())
			abj.emitMessage(contract.Message{Text: "Dealer has Blackjack. Your bet is saved 2:1.", PlayerId: playerId, HandId: handId})
		} else {
			abj.emitMessage(contract.Message{Text: "Dealer has Blackjack. You've lost.", PlayerId: playerId})
		}
	}
}

func (abj *AmericanBlackjack) completeRound() {
	dealerScore := (*abj.dealerHand).GetScore()
	dealerBusted := (*abj.dealerHand).IsBust()
	dealerHasBlackjack := (*abj.dealerHand).IsBlackjack()

	for _, player := range abj.players {
		playerId := (*player).GetId()
		playerHasBlackjack :=
			len(abj.playerHands[(*player).GetId()]) == 1 &&
				(*abj.playerHands[(*player).GetId()][0]).IsBlackjack()
		if playerHasBlackjack {
			bet := (*abj.playerHands[(*player).GetId()][0]).GetBet()
			if dealerHasBlackjack {
				abj.emitMessage(contract.Message{Text: (*player).GetName() + ", you have Blackjack, as a dealer. Your bet is back.", PlayerId: playerId})
				(*player).IncreasePurse(bet)
			} else {
				abj.emitMessage(contract.Message{Text: (*player).GetName() + ", you have Blackjack! You win 3:2!", PlayerId: playerId})
				(*player).IncreasePurse(bet*2 + bet/2)
			}

			continue
		}

		for _, hand := range abj.playerHands[(*player).GetId()] {
			handId := (*hand).GetId()
			playerBusted := (*hand).GetScore() > 21
			// todo хорошо бы для каждого из сообщений ниже указывать руку, если их >1. Подумать нал репозиторием сообщений, а из этого кода обращаться к ним по ID
			switch {
			case dealerBusted && !playerBusted:
				abj.emitMessage(contract.Message{Text: (*player).GetName() + ", dealer busted. You're win!", PlayerId: playerId, HandId: handId})
				(*player).IncreasePurse((*hand).GetBet() * 2)
			case playerBusted:
				abj.emitMessage(contract.Message{Text: (*player).GetName() + ", you busted!", PlayerId: playerId, HandId: handId})
			case dealerScore > (*hand).GetScore():
				abj.emitMessage(contract.Message{Text: (*player).GetName() + ", you're lost!", PlayerId: playerId, HandId: handId})
			case dealerScore < (*hand).GetScore():
				abj.emitMessage(contract.Message{Text: (*player).GetName() + ", you're win!", PlayerId: playerId, HandId: handId})
				(*player).IncreasePurse((*hand).GetBet() * 2)
			case dealerScore == (*hand).GetScore():
				abj.emitMessage(contract.Message{Text: (*player).GetName() + ", you pushed.", PlayerId: playerId, HandId: handId})
				(*player).IncreasePurse((*hand).GetBet())
			}
		}
	}
}

func (abj *AmericanBlackjack) getCardFromShoe() card.Card {
	if abj.shoe.IsCutCardReached() && !abj.shoe.IsDiscardTrayMode() {
		abj.shoe.ToggleGettingCardFromDiscardTrayMode()
		abj.shoe.ShuffleDiscardTray()

		abj.emitMessage(contract.Message{Text: "Cut card reached. Got cards from discard tray and shuffled them."})
	}

	return abj.shoe.GetNextCard()
}

func (abj *AmericanBlackjack) removeBankruptPlayers() {
	abj.players = slices.DeleteFunc(abj.players, func(p *contract.Player) bool {
		playerId := (*p).GetId()
		hasNoMoney := !(*p).HasMoney()

		if hasNoMoney {
			abj.emitMessage(contract.Message{Text: (*p).GetName() + "You lost everything. Bye, looser.", PlayerId: playerId})
		}

		return hasNoMoney
	})
}

func (abj *AmericanBlackjack) emitGameContext(playerDetails contract.PlayerDetails) {
	if gameContext == nil {
		gameContext = contract.NewGameContext(
			&abj.players,
			&abj.playerHands,
			abj.dealerHand,
			abj.dealerFirstCardFlipped,
			playerDetails.PlayerId,
			playerDetails.HandId,
		)
	} else {
		contract.UpdateGameContext(
			gameContext,
			&abj.players,
			&abj.playerHands,
			abj.dealerHand,
			abj.dealerFirstCardFlipped,
			playerDetails.PlayerId,
			playerDetails.HandId,
		)
	}

	abj.sendGameContext(gameContext)
}

func (abj *AmericanBlackjack) emitMessage(msg contract.Message) {
	abj.sendMessage(msg)
}

func (abj *AmericanBlackjack) applyAction(
	action contract.GameAction,
	player *contract.Player,
	playerHand *contract.PlayerHand,
) bool {
	continuePlaying := true

	switch action {
	case contract.GameActionHit:
		(*playerHand).AddCard(abj.getCardFromShoe())
	case contract.GameActionStay:
		continuePlaying = false
	case contract.GameActionDoubleDown:
		(*player).DecreasePurse((*playerHand).GetBet())
		(*playerHand).IncreaseBet((*playerHand).GetBet())
		(*playerHand).AddCard(abj.getCardFromShoe())
		continuePlaying = false
	case contract.GameActionSplit:
		(*player).DecreasePurse((*playerHand).GetBet())
		abj.playerHands[(*player).GetId()] = append(abj.playerHands[(*player).GetId()], abj.createPlayerHand((*playerHand).GetBet()))
		splitCard := (*playerHand).Split()
		splitHandIdx := len(abj.playerHands[(*player).GetId()]) - 1
		(*abj.playerHands[(*player).GetId()][splitHandIdx]).AddCard(splitCard)
	case contract.GameActionSurrender:
		(*player).IncreasePurse((*playerHand).GetBet() / 2)
		continuePlaying = false
	case contract.GameActionInsurance:
		insuranceBet := (*playerHand).GetBet() / 2
		abj.insuredPlayers = append(abj.insuredPlayers, (*player).GetId())
		(*player).DecreasePurse(insuranceBet)
		continuePlaying = false
	}

	return continuePlaying
}
