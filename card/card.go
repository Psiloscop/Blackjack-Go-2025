package card

type Suit uint8

func (s Suit) isValid() bool {
	return s <= Spades
}

type Number uint8

func (n Number) isValid() bool {
	return n <= Ace
}

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

const (
	Two Number = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Card struct {
	Suit   Suit
	Number Number
}

func New(suit Suit, number Number) Card {
	if !suit.isValid() {
		panic("Invalid suit")
	}
	if !number.isValid() {
		panic("Invalid number")
	}

	return Card{suit, number}
}
