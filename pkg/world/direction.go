package world

import "errors"

const (
	North = iota
	East
	South
	West
)

type Direction = int

func parseDirection(s string) (_ Direction, ok bool) {
	switch s {
	case "north":
		return North, true
	case "east":
		return East, true
	case "south":
		return South, true
	case "west":
		return West, true
	default:
		return -1, false
	}
}

func directionString(d Direction) (string, error) {
	switch d {
	case North:
		return "north", nil
	case East:
		return "east", nil
	case South:
		return "south", nil
	case West:
		return "west", nil
	}
	return "", errors.New("invalid direction")
}
