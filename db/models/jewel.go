package models

import "mini-gamestate-service/db/models/jewels"

type Jewel struct {
	Red    int64 `json:"red"`
	Blue   int64 `json:"blue"`
	Green  int64 `json:"green"`
	Yellow int64 `json:"yellow"`
	Black  int64 `json:"black"`
}

func NewJewelFromMap(m map[string]int64) *Jewel {
	red, ok := m[jewels.RedField]
	if !ok {
		red = 0
	}
	blue, ok := m[jewels.BlueField]
	if !ok {
		blue = 0
	}
	green, ok := m[jewels.GreenField]
	if !ok {
		green = 0
	}
	yellow, ok := m[jewels.YellowField]
	if !ok {
		yellow = 0
	}
	black, ok := m[jewels.BlackField]
	if !ok {
		black = 0
	}

	return &Jewel{
		Red:    red,
		Blue:   blue,
		Green:  green,
		Yellow: yellow,
		Black:  black,
	}
}
