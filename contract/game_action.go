package contract

type GameAction uint8

func (a GameAction) IsValid() bool {
	return a <= GameActionInsurance
}
func (a GameAction) String() string {
	switch a {
	case GameActionHit:
		return "Hit"
	case GameActionStay:
		return "Stay"
	case GameActionDoubleDown:
		return "Double Down"
	case GameActionSplit:
		return "Split"
	case GameActionSurrender:
		return "Surrender"
	case GameActionInsurance:
		return "Insurance"
	default:
		return "Undefined"
	}
}

const (
	GameActionHit GameAction = iota
	GameActionStay
	GameActionDoubleDown
	GameActionSplit
	GameActionSurrender
	GameActionInsurance
)
