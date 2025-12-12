package console

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Psiloscop/Blackjack-Go-2025/card"
	"github.com/Psiloscop/Blackjack-Go-2025/contract"
)

const spade rune = '\u2660'   // ♠
const heart rune = '\u2665'   // ♥
const diamond rune = '\u2666' // ♦
const club rune = '\u2663'    // ♣
const right rune = '\u25B6'   // ▶

var cashedTableRender string
var cashedMessageRender string

func drawTable(gameCtx *contract.GameContext) {
	if gameCtx == nil {
		if cashedTableRender != "" {
			fmt.Println(cashedTableRender)
		}

		return
	}

	var output string

	if (*gameCtx).DealerFirstCardFlipped {
		output += "Dealer: "
	} else {
		output += "Dealer (" + strconv.Itoa(int((*gameCtx).DealerHand.GetScore())) + "): "
	}

	for i, c := range (*gameCtx).DealerHand.GetCards() {
		if i == 0 && (*gameCtx).DealerFirstCardFlipped {
			output += "## "
		} else {
			output += getCardString(&c) + " "
		}
	}

	output += "\n"

	rightArrow := string(right) + " "
	rightArrowSpaces := strings.Repeat(" ", utf8.RuneCountInString(rightArrow))
	for i, p := range (*gameCtx).Players {
		if i != 0 {
			output += "\n"
		}

		playerTitle := "\n" + p.GetName() + ": "
		playerTitleSpaces := strings.Repeat(" ", utf8.RuneCountInString(playerTitle)-1)
		for j, h := range (*gameCtx).PlayerHands[p.GetId()] {
			if j == 0 {
				output += playerTitle
			} else {
				output += "\n"
				output += playerTitleSpaces
			}

			output += "(" + strconv.Itoa(int(h.GetScore())) + ") "

			if len((*gameCtx).PlayerHands[p.GetId()]) > 1 {
				if h.GetId() == (*gameCtx).CurrentPlayerHandId {
					output += rightArrow
				} else {
					output += rightArrowSpaces
				}
			}

			for _, c := range h.GetCards() {
				output += getCardString(&c) + " "
			}
		}
	}

	cashedTableRender = output

	fmt.Println(output)
}

func drawPlayerBetInput(p contract.Player) {
	var output = "\n"

	output += p.GetName() + ", your purse is $" + strconv.Itoa(int(p.GetPurse())) + ". Enter your bet: "

	fmt.Print(output)
}

func drawGameActionChooser(playerCtx *contract.PlayerContext) {
	var output = "\n"

	for _, aa := range playerCtx.GetAvailableActions() {
		output += strconv.Itoa(int(aa)) + ". " + aa.String() + "\n"
	}

	output += "\n"
	output += "Choose action (enter number): "

	fmt.Print(output)
}

func getCardString(c *card.Card) string {
	var output string

	switch c.Number {
	case card.Two:
		output += "2"
	case card.Three:
		output += "3"
	case card.Four:
		output += "4"
	case card.Five:
		output += "5"
	case card.Six:
		output += "6"
	case card.Seven:
		output += "7"
	case card.Eight:
		output += "8"
	case card.Nine:
		output += "9"
	case card.Ten:
		output += "10"
	case card.Jack:
		output += "J"
	case card.Queen:
		output += "Q"
	case card.King:
		output += "K"
	case card.Ace:
		output += "A"
	}

	switch c.Suit {
	case card.Clubs:
		output += string(club)
	case card.Diamonds:
		output += string(diamond)
	case card.Hearts:
		output += string(heart)
	case card.Spades:
		output += string(spade)
	}

	return output
}

func cacheMessage(message contract.Message) {
	var output = "\n"
	if message.IsError {
		output += "Error: "
	}
	output += message.Text

	cashedMessageRender += output
}
func drawCachedMessages() {
	if cashedMessageRender != "" {
		fmt.Println(cashedMessageRender)
	}
}
func clearCachedMessages() {
	cashedMessageRender = ""
}
