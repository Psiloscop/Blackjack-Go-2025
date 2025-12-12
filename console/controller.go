package console

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/Psiloscop/Blackjack-Go-2025/contract"
)

func DisplayTable(ctx *contract.GameContext) {
	clearConsole()
	drawTable(ctx)
}

func DisplayMessage(msg contract.Message) {
	clearConsole()
	drawTable(nil)
	cacheMessage(msg)
	drawCachedMessages()
}

func RequestBet(p contract.Player) uint {
	for {
		clearConsole()
		drawTable(nil)
		drawCachedMessages()
		clearCachedMessages()
		drawPlayerBetInput(p)

		var input string
		_, scanErr := fmt.Scan(&input)
		if scanErr == nil {
			value, parseErr := strconv.ParseUint(input, 10, 32)

			if parseErr == nil {
				return uint(value)
			}
		}

		cacheMessage(contract.Message{Text: "Invalid input. Enter unsigned value.", IsError: true})
	}
}

func RequestAction(ctx *contract.PlayerContext) contract.GameAction {
	for {
		clearConsole()
		drawTable(nil)
		drawCachedMessages()
		clearCachedMessages()
		drawGameActionChooser(ctx)

		var input string
		_, scanErr := fmt.Scan(&input)
		if scanErr == nil {
			value, parseErr := strconv.ParseUint(input, 10, 8)

			if parseErr == nil {
				return contract.GameAction(value)
			}
		}

		cacheMessage(contract.Message{Text: "Invalid input. Enter unsigned acton number.", IsError: true})
	}
}

func SendError(err error) {
	cacheMessage(contract.Message{Text: err.Error(), IsError: true})
}

func clearConsole() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else { // linux, darwin (macOS)
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
