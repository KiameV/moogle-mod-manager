package config

type GameName string

const (
	FfPrI   GameName = "FF PR I"
	FfPrII  GameName = "FF PR II"
	FfPrIII GameName = "FF PR III"
	FfPrIV  GameName = "FF PR IV"
	FfPrV   GameName = "FF PR V"
	FfPrVI  GameName = "FF PR VI"
)

type Game int

const (
	I Game = iota
	II
	III
	IV
	V
	VI
)

func String(game Game) (name string) {
	switch game {
	case I:
		name = "I"
	case II:
		name = "II"
	case III:
		name = "III"
	case IV:
		name = "IV"
	case V:
		name = "V"
	case VI:
		name = "VI"
	}
	return
}

func GameNameString(game Game) string {
	return "Final Fantasy " + String(game)
}

func FromString(s string) (game Game) {
	switch s {
	case "Final Fantasy I":
		return I
	case "Final Fantasy II":
		return II
	case "Final Fantasy III":
		return III
	case "Final Fantasy IV":
		return IV
	case "Final Fantasy V":
		return V
	case "Final Fantasy VI":
		return VI
	}
	return VI
}

func NameToGame(n GameName) Game {
	switch n {
	case FfPrI:
		return I
	case FfPrII:
		return II
	case FfPrIII:
		return III
	case FfPrIV:
		return IV
	case FfPrV:
		return V
	}
	return VI
}

func GameToName(game Game) GameName {
	switch game {
	case I:
		return FfPrI
	case II:
		return FfPrII
	case III:
		return FfPrIII
	case IV:
		return FfPrIV
	case V:
		return FfPrV
	}
	return FfPrVI
}
