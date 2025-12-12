package contract

type Communicator interface {
	RequestBet() string
	RequestAction(ctx *PlayerContext) string
	SendError(err error)
}
