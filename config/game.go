package config

import "fmt"

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

func GetGameName(game Game) (name string) {
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
	name = fmt.Sprintf("Final Fantasy %s PR Mods", name)
	return
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
