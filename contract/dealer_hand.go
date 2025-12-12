package contract

type DealerHand interface {
	hand
	IsSecondCardAce() bool
}

type DealerHandCreator func() *DealerHand
