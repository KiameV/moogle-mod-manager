package config

const GameCount = 9

type GameName string

const (
	FfPrI           GameName = "FF PR I"
	FfPrII          GameName = "FF PR II"
	FfPrIII         GameName = "FF PR III"
	FfPrIV          GameName = "FF PR IV"
	FfPrV           GameName = "FF PR V"
	FfPrVI          GameName = "FF PR VI"
	ChronoCrossName          = "Chrono Cross"
	BofIIIName      GameName = "BoF III"
	BofIVName       GameName = "BoF IV"
)

type Game int

const (
	I Game = iota
	II
	III
	IV
	V
	VI
	ChronoCross
	BofIII
	BofIV
)

var GameNames = []string{
	GameNameString(I),
	GameNameString(II),
	GameNameString(III),
	GameNameString(IV),
	GameNameString(V),
	GameNameString(VI),
	GameNameString(ChronoCross),
	// TODO BOF
	//GameNameString(BofIII),
	//GameNameString(BofIV),
}

func String(game Game) string {
	switch game {
	case I:
		return "I"
	case II:
		return "II"
	case III:
		return "III"
	case IV:
		return "IV"
	case V:
		return "V"
	case VI:
		return "VI"
	case ChronoCross:
		return "Chrono Cross"
	case BofIII:
		return "BoF III"
	case BofIV:
		return "BoF IV"
	}
	panic("invalid game " + string(game))
}

func GameNameString(game Game) string {
	if game <= VI {
		return "Final Fantasy " + String(game)
	} else if game == ChronoCross {
		return "Chrono Cross"
	}
	return "Breath of Fire " + String(game)
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
	case "Chrono Cross":
		return ChronoCross
	case "Breath of Fire III":
		return BofIII
	case "Breath of Fire IV":
		return BofIV
	}
	panic("invalid game name " + s)
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
	case FfPrVI:
		return VI
	case ChronoCrossName:
		return ChronoCross
	case BofIIIName:
		return BofIII
	case BofIVName:
		return BofIV
	}
	return I
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
	case VI:
		return FfPrVI
	case ChronoCross:
		return ChronoCrossName
	case BofIII:
		return BofIIIName
	case BofIV:
		return BofIVName
	}
	panic("invalid game " + string(game))
}
