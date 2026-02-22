package main

type Color int

const (
	RED Color = iota
	GREEN
	YELLOW
	BLUE
	CYAN
	PURPLE
	RESET
)

func (c Color) String() string {
	switch c {
	case RED:
		return "\u001B[0;31m"
	case GREEN:
		return "\u001B[0;32m"
	case YELLOW:
		return "\u001B[0;33m"
	case BLUE:
		return "\033[0;34m"
	case CYAN:
		return "\u001B[0;36m"
	case PURPLE:
		return "\033[0;35m"
	case RESET:
		return "\u001B[0m"
	default:
		return "\u001B"
	}
}
