package player

import (
	"github.com/Psiloscop/Blackjack-Go-2025/contract"
	"github.com/google/uuid"
)

type Player struct {
	id    contract.PlayerId
	name  string
	purse uint
	mind  *contract.Mind
}

func (p *Player) GetId() contract.PlayerId {
	return p.id
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetPurse() uint {
	return p.purse
}

func (p *Player) IncreasePurse(amount uint) {
	p.purse += amount
}

func (p *Player) DecreasePurse(amount uint) {
	if p.purse < amount {
		p.purse = 0
	} else {
		p.purse -= amount
	}
}

func (p *Player) HasMoney() bool {
	return p.purse > 0
}

func (p *Player) PlaceBet() uint {
	bet := (*p.mind).PlaceBet(p)
	p.purse -= bet

	return bet
}

func (p *Player) ChooseAction(ctx *contract.PlayerContext) contract.GameAction {
	return (*p.mind).ChooseAction(ctx)
}

func New(name string, purse uint, mind *contract.Mind) contract.Player {
	return &Player{
		contract.PlayerId(uuid.New().String()),
		name,
		purse,
		mind,
	}
}
