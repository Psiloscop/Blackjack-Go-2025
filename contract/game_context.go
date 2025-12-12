package contract

type GameContext struct {
	Players                []Player
	PlayerHands            map[PlayerId][]PlayerHand
	DealerHand             DealerHand
	DealerFirstCardFlipped bool
	CurrentPlayerId        PlayerId
	CurrentPlayerHandId    PlayerHandId
}

func NewGameContext(
	players *[]*Player,
	playerHands *map[PlayerId][]*PlayerHand,
	dealerHand *DealerHand,
	dealerFirstCardFlipped bool,
	currentPlayerId PlayerId,
	currentPlayerHandId PlayerHandId,
) *GameContext {
	_players := convertPlayers(players)
	_playerHands := convertPlayerHands(playerHands)

	return &GameContext{
		*_players,
		*_playerHands,
		*dealerHand,
		dealerFirstCardFlipped,
		currentPlayerId,
		currentPlayerHandId,
	}
}

func UpdateGameContext(
	gameContext *GameContext,
	players *[]*Player,
	playerHands *map[PlayerId][]*PlayerHand,
	dealerHand *DealerHand,
	dealerFirstCardFlipped bool,
	currentPlayerId PlayerId,
	currentPlayerHandId PlayerHandId,
) {
	_players := convertPlayers(players)
	_playerHands := convertPlayerHands(playerHands)

	gameContext.Players = *_players
	gameContext.PlayerHands = *_playerHands
	gameContext.DealerHand = *dealerHand
	gameContext.DealerFirstCardFlipped = dealerFirstCardFlipped
	gameContext.CurrentPlayerId = currentPlayerId
	gameContext.CurrentPlayerHandId = currentPlayerHandId
}

func convertPlayers(players *[]*Player) *[]Player {
	_players := make([]Player, 0, len(*players))
	for _, player := range *players {
		_players = append(_players, *player)
	}
	return &_players
}

func convertPlayerHands(playerHands *map[PlayerId][]*PlayerHand) *map[PlayerId][]PlayerHand {
	_playerHands := make(map[PlayerId][]PlayerHand, len(*playerHands))
	for playerId, hands := range *playerHands {
		_playerHands[playerId] = make([]PlayerHand, 0, len(hands))
		for _, hand := range hands {
			_playerHands[playerId] = append(_playerHands[playerId], *hand)
		}
	}

	return &_playerHands
}

type GameContextSender func(ctx *GameContext)
