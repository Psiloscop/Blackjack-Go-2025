package table

import (
	"slices"
	"testing"

	"github.com/Psiloscop/Blackjack-Go-2025/card"
	"github.com/Psiloscop/Blackjack-Go-2025/contract"
	"github.com/Psiloscop/Blackjack-Go-2025/hand"
	"github.com/Psiloscop/Blackjack-Go-2025/player"
	"github.com/Psiloscop/Blackjack-Go-2025/shoe"
)

var playerId contract.PlayerId

var tmi contract.Mind
var tmiPtr *contract.Mind
var testMind *TestMind

var testShoe = shoe.New(1, 75)

var testPlayer *contract.Player
var testPlayers []*contract.Player

var createPlayer = func(purse uint) *contract.Player {
	_player := player.New("Player 1", purse, tmiPtr)

	return &_player
}
var dealerHandCreator = func() *contract.DealerHand {
	dealerHand := hand.NewDealerHand()

	return &dealerHand
}
var playerHandCreator = func(bet uint) *contract.PlayerHand {
	playerHand := hand.NewPlayerHand(bet)

	return &playerHand
}

var americanBlackjack *AmericanBlackjack

func setup() {
	tmi = NewTestMind()
	tmiPtr = &tmi
	testPlayer = createPlayer(0)
	testPlayers = []*contract.Player{
		testPlayer,
	}
	playerId = (*testPlayer).GetId()
	tm, ok := (*tmiPtr).(*TestMind)
	if ok {
		testMind = tm
	} else {
		panic("Test mind is not TestMind type")
	}

	americanBlackjack = NewAmericanBlackjack(
		testShoe,
		testPlayers,
		dealerHandCreator,
		playerHandCreator,
		func(ctx *contract.GameContext) {},
		func(msg contract.Message) {},
	)
}

func Test_AmericanBlackjack_requestBets(t *testing.T) {
	setup()

	var expectedBet = uint(1000)

	testMind.BetToPlace = expectedBet

	if len(americanBlackjack.playerHands[playerId]) != 0 {
		t.Error("Player should not has any hands before bet request.")
	}

	americanBlackjack.requestBets()

	if (*americanBlackjack.playerHands[playerId][0]).GetBet() != expectedBet {
		t.Errorf("Player bet should be %d, got %d", expectedBet, (*americanBlackjack.playerHands[playerId][0]).GetBet())
	}

	if len(americanBlackjack.playerHands[playerId]) != 1 {
		t.Error("Player should has one hand after bet request.")
	}
}

func Test_AmericanBlackjack_dealCards(t *testing.T) {
	setup()

	americanBlackjack.requestBets()

	if (*americanBlackjack.playerHands[playerId][0]).GetCardAmount() != 0 {
		t.Errorf(
			"Player hand should be empty before deal. Got %d",
			(*americanBlackjack.playerHands[playerId][0]).GetCardAmount(),
		)
	}
	if (*americanBlackjack.dealerHand).GetCardAmount() != 0 {
		t.Errorf(
			"Dealer hand should be empty before deal. Got %d",
			(*americanBlackjack.playerHands[playerId][0]).GetCardAmount(),
		)
	}

	americanBlackjack.dealCards()

	if (*americanBlackjack.playerHands[playerId][0]).GetCardAmount() != 2 {
		t.Errorf(
			"Player hand should contain two cards after deal. Got %d",
			(*americanBlackjack.playerHands[playerId][0]).GetCardAmount(),
		)
	}
	if (*americanBlackjack.dealerHand).GetCardAmount() != 2 {
		t.Errorf(
			"Dealer hand should contain two cards after deal. Got %d",
			(*americanBlackjack.dealerHand).GetCardAmount(),
		)
	}
}

func Test_AmericanBlackjack_takeCardsUpToScore17(t *testing.T) {
	setup()

	dealerHand := americanBlackjack.dealerHand
	(*dealerHand).AddCard(card.New(card.Spades, card.Two))
	(*dealerHand).AddCard(card.New(card.Spades, card.Two))

	if (*dealerHand).GetScore() != 4 {
		t.Errorf("Dealer hand should have score 4. Got %d", (*dealerHand).GetScore())
	}

	americanBlackjack.takeCardsUpToScore17()

	if (*dealerHand).GetScore() < 17 {
		t.Errorf("Dealer hand should have score 17 or more. Got %d", (*dealerHand).GetScore())
	}
}

func Test_AmericanBlackjack_serveInsuredPlayers(t *testing.T) {
	setup()

	player1 := createPlayer(0)
	player2 := createPlayer(0)
	player3 := createPlayer(0)

	americanBlackjack.players = []*contract.Player{
		player1,
		player2,
		player3,
	}
	americanBlackjack.playerHands[(*player1).GetId()] = []*contract.PlayerHand{playerHandCreator(100)}
	americanBlackjack.playerHands[(*player2).GetId()] = []*contract.PlayerHand{playerHandCreator(100)}
	americanBlackjack.playerHands[(*player3).GetId()] = []*contract.PlayerHand{playerHandCreator(100)}

	americanBlackjack.insuredPlayers = []contract.PlayerId{
		(*player1).GetId(),
		(*player3).GetId(),
	}

	americanBlackjack.completeRoundEarly()

	if (*player1).GetPurse() != 100 {
		t.Errorf("Player 1 purse should be 100 because its insured. Got %d", (*player1).GetPurse())
	}
	if (*player2).GetPurse() != 0 {
		t.Errorf("Player 2 purse should be 0 because its not insured. Got %d", (*player2).GetPurse())
	}
	if (*player3).GetPurse() != 100 {
		t.Errorf("Player 3 purse should be 100 because its insured. Got %d", (*player3).GetPurse())
	}
}

func Test_AmericanBlackjack_removeBankruptPlayers(t *testing.T) {
	setup()

	player1 := createPlayer(100)
	player2 := createPlayer(0)
	player3 := createPlayer(100)

	americanBlackjack.players = []*contract.Player{
		player1,
		player2,
		player3,
	}
	expectedPlayers := []*contract.Player{player1, player3}

	americanBlackjack.removeBankruptPlayers()

	if !slices.Equal(americanBlackjack.players, expectedPlayers) {
		t.Errorf("Not all bankrupt players removed %T. Got %T", americanBlackjack.players, expectedPlayers)
	}
}

func Test_AmericanBlackjack_applyActon_Hit(t *testing.T) {
	setup()

	americanBlackjack.requestBets()
	americanBlackjack.dealCards()

	playerHand := *americanBlackjack.playerHands[playerId][0]

	if playerHand.GetCardAmount() != 2 {
		t.Errorf("hand should contain two cards before hit. Got %d", playerHand.GetCardAmount())
	}

	americanBlackjack.applyAction(contract.GameActionHit, testPlayer, &playerHand)

	if playerHand.GetCardAmount() != 3 {
		t.Errorf("hand should contain three cards after one hit. Got %d", playerHand.GetCardAmount())
	}
}

func Test_AmericanBlackjack_applyActon_DoubleDown(t *testing.T) {
	setup()
	testMind.BetToPlace = 500
	(*testPlayer).IncreasePurse(2000)

	americanBlackjack.requestBets()
	americanBlackjack.dealCards()

	playerHand := *americanBlackjack.playerHands[playerId][0]

	if playerHand.GetBet() != 500 {
		t.Errorf("Player bet should be 500 before double down. Got %d", playerHand.GetBet())
	}
	if playerHand.GetCardAmount() != 2 {
		t.Errorf("hand should contain two cards before double down. Got %d", playerHand.GetCardAmount())
	}

	americanBlackjack.applyAction(contract.GameActionDoubleDown, testPlayer, &playerHand)

	if (*testPlayer).GetPurse() != 1000 {
		t.Errorf("Player purse should be 1000 after double down. Got %d", (*testPlayer).GetPurse())
	}
	if playerHand.GetBet() != 1000 {
		t.Errorf("Player bet should be 1000 after double down. Got %d", playerHand.GetBet())
	}
	if playerHand.GetCardAmount() != 3 {
		t.Errorf("hand should contain three cards after double down. Got %d", playerHand.GetCardAmount())
	}
}

func Test_AmericanBlackjack_applyActon_Split(t *testing.T) {
	setup()

	testCases := []struct {
		name          string
		expectedPurse uint
	}{
		{
			"Player purse should be 1750 first split",
			1500,
		},
		{
			"Player purse should be 1500 second split",
			1250,
		},
		{
			"Player purse should be 1250 third split",
			1000,
		},
	}
	expectedBedPerHand := uint(250)

	testMind.BetToPlace = expectedBedPerHand
	(*testPlayer).IncreasePurse(2000)

	americanBlackjack.requestBets()

	firstPlayerHand := *americanBlackjack.playerHands[playerId][0]
	firstPlayerHand.AddCard(card.New(card.Spades, card.Eight))

	if firstPlayerHand.GetBet() != expectedBedPerHand {
		t.Errorf("Player bet should be 250 before split. Got %d", firstPlayerHand.GetBet())
	}

	for idx, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			splitNumber := idx + 1
			expectedHands := splitNumber + 1

			firstPlayerHand.AddCard(card.New(card.Spades, card.Eight))

			americanBlackjack.applyAction(contract.GameActionSplit, testPlayer, &firstPlayerHand)

			splitPlayerHand := *americanBlackjack.playerHands[playerId][splitNumber]

			if len(americanBlackjack.playerHands[playerId]) != expectedHands {
				t.Errorf(
					"Player should has %d hands after %d split. Got %d",
					expectedHands,
					splitNumber,
					len(americanBlackjack.playerHands[playerId]),
				)
			}
			if firstPlayerHand.GetBet() != expectedBedPerHand {
				t.Errorf(
					"First player bet should stay %d after split. Got %d",
					expectedBedPerHand,
					firstPlayerHand.GetBet(),
				)
			}
			if firstPlayerHand.GetCardAmount() != 1 {
				t.Errorf(
					"First player hand should contain one card after %d split. Got %d",
					splitNumber,
					splitPlayerHand.GetCardAmount(),
				)
			}
			if splitPlayerHand.GetBet() != expectedBedPerHand {
				t.Errorf(
					"Player's split bet should be %d after %d split. Got %d",
					expectedBedPerHand,
					splitNumber,
					firstPlayerHand.GetBet(),
				)
			}
			if splitPlayerHand.GetCardAmount() != 1 {
				t.Errorf(
					"Player's split hand should contain one card after %d split. Got %d",
					splitNumber,
					splitPlayerHand.GetCardAmount(),
				)
			}
			if (*testPlayer).GetPurse() != tc.expectedPurse {
				t.Errorf(
					"Player purse should be %d after %d split. Got %d",
					tc.expectedPurse,
					splitNumber,
					(*testPlayer).GetPurse(),
				)
			}
		})
	}
}

func Test_AmericanBlackjack_applyActon_Surrender(t *testing.T) {
	setup()
	testMind.BetToPlace = 500
	(*testPlayer).IncreasePurse(2000)

	americanBlackjack.requestBets()
	americanBlackjack.dealCards()

	playerHand := *americanBlackjack.playerHands[playerId][0]

	if (*testPlayer).GetPurse() != 1500 {
		t.Errorf("Player purse should be 1500 before surrender. Got %d", (*testPlayer).GetPurse())
	}
	if playerHand.GetBet() != 500 {
		t.Errorf("Player bet should be 500 before surrender. Got %d", playerHand.GetBet())
	}
	if playerHand.GetCardAmount() != 2 {
		t.Errorf("hand should contain two cards before surrender. Got %d", playerHand.GetCardAmount())
	}

	americanBlackjack.applyAction(contract.GameActionSurrender, testPlayer, &playerHand)

	if (*testPlayer).GetPurse() != 1750 {
		t.Errorf("Player purse should be 1750 after surrender. Got %d", (*testPlayer).GetPurse())
	}
}

func Test_AmericanBlackjack_applyActon_Insurance(t *testing.T) {
	setup()
	testMind.BetToPlace = 500
	(*testPlayer).IncreasePurse(2000)

	americanBlackjack.requestBets()
	americanBlackjack.dealCards()

	playerHand := *americanBlackjack.playerHands[playerId][0]

	if (*testPlayer).GetPurse() != 1500 {
		t.Errorf("Player purse should be 1500 before insurance. Got %d", (*testPlayer).GetPurse())
	}
	if playerHand.GetBet() != 500 {
		t.Errorf("Player bet should be 500 before insurance. Got %d", playerHand.GetBet())
	}

	americanBlackjack.applyAction(contract.GameActionInsurance, testPlayer, &playerHand)

	if (*testPlayer).GetPurse() != 1250 {
		t.Errorf("Player purse should be 1250 after surrender. Got %d", (*testPlayer).GetPurse())
	}
	if !slices.Contains(americanBlackjack.insuredPlayers, (*testPlayer).GetId()) {
		t.Errorf("Player ID  should be in insured slice")
	}
}

func Test_AmericanBlackjack_completeRound(t *testing.T) {
	setup()

	testCases := []struct {
		name              string
		playerCreator     func() *contract.Player
		dealerHandCreator func() *contract.DealerHand
		playerHandCreator func() *contract.PlayerHand
		expectedPurse     uint
	}{
		{
			"Player win.",
			func() *contract.Player {
				_player := player.New("Player 1", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.Ten))
				(*dealerHand).AddCard(card.New(card.Spades, card.Five))
				(*dealerHand).AddCard(card.New(card.Spades, card.Two))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(100)
				(*playerHand).AddCard(card.New(card.Spades, card.King))
				(*playerHand).AddCard(card.New(card.Spades, card.Three))
				(*playerHand).AddCard(card.New(card.Spades, card.Eight))

				return playerHand
			},
			200,
		},
		{
			"Dealer win.",
			func() *contract.Player {
				_player := player.New("Player 1", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.Ten))
				(*dealerHand).AddCard(card.New(card.Spades, card.Three))
				(*dealerHand).AddCard(card.New(card.Spades, card.Two))
				(*dealerHand).AddCard(card.New(card.Spades, card.Five))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(100)
				(*playerHand).AddCard(card.New(card.Spades, card.Jack))
				(*playerHand).AddCard(card.New(card.Spades, card.Eight))

				return playerHand
			},
			0,
		},
		{
			"Player busted. Dealer not.",
			func() *contract.Player {
				_player := player.New("Player 1", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.King))
				(*dealerHand).AddCard(card.New(card.Spades, card.Nine))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(100)
				(*playerHand).AddCard(card.New(card.Spades, card.King))
				(*playerHand).AddCard(card.New(card.Spades, card.Nine))
				(*playerHand).AddCard(card.New(card.Spades, card.Eight))

				return playerHand
			},
			0,
		},
		{
			"Dealer busted. Player not.",
			func() *contract.Player {
				_player := player.New("Player 1", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.King))
				(*dealerHand).AddCard(card.New(card.Spades, card.Nine))
				(*dealerHand).AddCard(card.New(card.Spades, card.Eight))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(100)
				(*playerHand).AddCard(card.New(card.Spades, card.King))
				(*playerHand).AddCard(card.New(card.Spades, card.Nine))

				return playerHand
			},
			200,
		},
		{
			"Dealer busted. Player busted.",
			func() *contract.Player {
				_player := player.New("Player 1", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.King))
				(*dealerHand).AddCard(card.New(card.Spades, card.Nine))
				(*dealerHand).AddCard(card.New(card.Spades, card.Eight))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(100)
				(*playerHand).AddCard(card.New(card.Spades, card.King))
				(*playerHand).AddCard(card.New(card.Spades, card.Nine))
				(*playerHand).AddCard(card.New(card.Spades, card.Eight))

				return playerHand
			},
			0,
		},
		{
			"Dealer nor Player win.",
			func() *contract.Player {
				_player := player.New("Player 1", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.King))
				(*dealerHand).AddCard(card.New(card.Spades, card.Nine))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(100)
				(*playerHand).AddCard(card.New(card.Spades, card.King))
				(*playerHand).AddCard(card.New(card.Spades, card.Nine))

				return playerHand
			},
			100,
		},
		{
			"Dealer and Player have blackjack.",
			func() *contract.Player {
				_player := player.New("Player 1", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.King))
				(*dealerHand).AddCard(card.New(card.Spades, card.Ace))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(100)
				(*playerHand).AddCard(card.New(card.Spades, card.King))
				(*playerHand).AddCard(card.New(card.Spades, card.Ace))

				return playerHand
			},
			100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testPlayer := tc.playerCreator()
			americanBlackjack.players = []*contract.Player{testPlayer}
			americanBlackjack.playerHands[(*testPlayer).GetId()] = []*contract.PlayerHand{tc.playerHandCreator()}
			americanBlackjack.dealerHand = tc.dealerHandCreator()

			americanBlackjack.completeRound()

			if (*testPlayer).GetPurse() != tc.expectedPurse {
				t.Errorf(
					"Player purse should be %d after round completion. Got %d",
					tc.expectedPurse,
					(*testPlayer).GetPurse(),
				)
			}
		})
	}
}

func Test_getAvailableActions(t *testing.T) {
	setup()

	testCases := []struct {
		name              string
		playerCreator     func() *contract.Player
		dealerHandCreator func() *contract.DealerHand
		playerHandCreator func() *contract.PlayerHand
		playerHandAmount  uint
		expectedActions   []contract.GameAction
	}{
		{
			"Available actions: Hit, Stay.",
			func() *contract.Player {
				_player := player.New("Player 1", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				return dealerHandCreator()
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(1000)
				(*playerHand).AddCard(card.New(card.Spades, card.Two))
				(*playerHand).AddCard(card.New(card.Spades, card.Three))
				(*playerHand).AddCard(card.New(card.Spades, card.Four))

				return playerHand
			},
			1,
			[]contract.GameAction{
				contract.GameActionHit,
				contract.GameActionStay,
			},
		},
		{
			"Available actions: Hit, Stay, Surrender.",
			func() *contract.Player {
				_player := player.New("Player", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				return dealerHandCreator()
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(1000)
				(*playerHand).AddCard(card.New(card.Spades, card.Two))
				(*playerHand).AddCard(card.New(card.Spades, card.Two))

				return playerHand
			},
			1,
			[]contract.GameAction{
				contract.GameActionHit,
				contract.GameActionStay,
				contract.GameActionSurrender,
			},
		},
		{
			"Available actions: Hit, Stay, Double Down, Split, Surrender.",
			func() *contract.Player {
				_player := player.New("Player", 1000, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.Two))
				(*dealerHand).AddCard(card.New(card.Spades, card.Three))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(1000)
				(*playerHand).AddCard(card.New(card.Spades, card.Five))
				(*playerHand).AddCard(card.New(card.Spades, card.Five))

				return playerHand
			},
			1,
			[]contract.GameAction{
				contract.GameActionHit,
				contract.GameActionStay,
				contract.GameActionDoubleDown,
				contract.GameActionSplit,
				contract.GameActionSurrender,
			},
		},
		{
			"Available actions: Hit, Stay, Insurance.",
			func() *contract.Player {
				_player := player.New("Player", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.Two))
				(*dealerHand).AddCard(card.New(card.Spades, card.Ace))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(1000)
				(*playerHand).AddCard(card.New(card.Spades, card.Two))
				(*playerHand).AddCard(card.New(card.Spades, card.Three))

				return playerHand
			},
			1,
			[]contract.GameAction{
				contract.GameActionHit,
				contract.GameActionStay,
				contract.GameActionInsurance,
			},
		},
		{
			"Available actions: Hit, Stay.",
			func() *contract.Player {
				_player := player.New("Player", 0, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.Two))
				(*dealerHand).AddCard(card.New(card.Spades, card.Ace))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(1000)
				(*playerHand).AddCard(card.New(card.Spades, card.Two))
				(*playerHand).AddCard(card.New(card.Spades, card.Three))
				(*playerHand).AddCard(card.New(card.Spades, card.Four))

				return playerHand
			},
			1,
			[]contract.GameAction{
				contract.GameActionHit,
				contract.GameActionStay,
			},
		},
		{
			"Available actions: Hit, Stay, Double Down, Split, Insurance.",
			func() *contract.Player {
				_player := player.New("Player", 1000, tmiPtr)

				return &_player
			},
			func() *contract.DealerHand {
				dealerHand := dealerHandCreator()
				(*dealerHand).AddCard(card.New(card.Spades, card.Two))
				(*dealerHand).AddCard(card.New(card.Spades, card.Ace))

				return dealerHand
			},
			func() *contract.PlayerHand {
				playerHand := playerHandCreator(1000)
				(*playerHand).AddCard(card.New(card.Spades, card.Five))
				(*playerHand).AddCard(card.New(card.Spades, card.Five))

				return playerHand
			},
			1,
			[]contract.GameAction{
				contract.GameActionHit,
				contract.GameActionStay,
				contract.GameActionDoubleDown,
				contract.GameActionSplit,
				contract.GameActionInsurance,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			availableActions := getAvailableActions(
				tc.playerCreator(),
				tc.dealerHandCreator(),
				tc.playerHandCreator(),
				tc.playerHandAmount,
			)

			if !slices.Equal(tc.expectedActions, availableActions) {
				t.Errorf("Expected actions: %v, got %v", tc.expectedActions, availableActions)
			}
		})
	}
}
