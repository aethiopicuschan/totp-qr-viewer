package main

import (
	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const textFieldHeight = 24

var fontFace = text.NewGoXFace(bitmapfont.FaceEA)

func textFieldPadding() (int, int) {
	m := fontFace.Metrics()
	return 4, (textFieldHeight - int(m.HLineGap+m.HAscent+m.HDescent)) / 2
}
