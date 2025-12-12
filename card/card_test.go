package card

import "testing"

func Test_SuitId_isValid(t *testing.T) {
	var suit Suit

	suit = Spades
	if !suit.isValid() {
		t.Error("Valid suitId should return true, got false")
	}

	suit = 4
	if suit.isValid() {
		t.Error("Invalid suitId should return false, got true")
	}
}

func Test_ValueId_isValid(t *testing.T) {
	var number Number

	number = Ace
	if !number.isValid() {
		t.Error("Valid number should return true, got false")
	}

	number = 13
	if number.isValid() {
		t.Error("Invalid number should return false, got true")
	}
}

func Test_New_ExpectPanicBecauseOfInvalidSuit(t *testing.T) {
	assertPanic(t, func() {
		New(4, Seven)
	})
}
func Test_New_ExpectPanicBecauseOfInvalidNumber(t *testing.T) {
	assertPanic(t, func() {
		New(Clubs, 13)
	})
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected function to panic, but it didn't")
		}
	}()
	f()
}
