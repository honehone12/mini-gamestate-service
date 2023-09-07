package jewels

import "errors"

type ColorCode byte

const (
	Red ColorCode = iota
	Blue
	Green
	Yellow
	Black
)

const (
	RedField    = "red"
	BlueField   = "blue"
	GreenField  = "green"
	YellowField = "yellow"
	BlackField  = "black"
)

var (
	ErrorInvalidColorCode = errors.New("no such color code")
)

func (c ColorCode) ColorCodeToString() (string, error) {
	switch c {
	case Red:
		return RedField, nil
	case Blue:
		return BlueField, nil
	case Green:
		return GreenField, nil
	case Yellow:
		return YellowField, nil
	case Black:
		return BlackField, nil
	default:
		return "", ErrorInvalidColorCode
	}
}
