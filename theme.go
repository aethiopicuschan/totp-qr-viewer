package main

import (
	"image/color"
	"strconv"
)

type Theme struct {
	BackgroundColor color.Color
	MainColor       color.Color
	BorderColor     color.Color
}

func GetDefaultTheme() *Theme {
	return &Theme{
		BackgroundColor: hexToColor("BED0BC"),
		MainColor:       hexToColor("819E57"),
		BorderColor:     hexToColor("617641"),
	}
}

func hexToColor(hex string) color.Color {
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		panic(err)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		panic(err)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		panic(err)
	}
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}
