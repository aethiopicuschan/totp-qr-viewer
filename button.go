package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	bounds      image.Rectangle
	label       string
	mainColor   color.Color
	borderColor color.Color
	onClick     func()
}

func NewButton(bounds image.Rectangle, label string, mainColor color.Color, borderColor color.Color, onClick func()) *Button {
	return &Button{
		bounds:      bounds,
		label:       label,
		mainColor:   mainColor,
		borderColor: borderColor,
		onClick:     onClick,
	}
}

func (b *Button) Contains(x, y int) bool {
	return image.Pt(x, y).In(b.bounds)
}

func (b *Button) Update() (err error) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		if b.Contains(mouseX, mouseY) {
			b.onClick()
		}
	}
	return
}

func (b *Button) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(b.bounds.Min.X), float32(b.bounds.Min.Y), float32(b.bounds.Dx()), float32(b.bounds.Dy()), b.mainColor, false)
	x, y := ebiten.CursorPosition()
	if b.Contains(x, y) {
		vector.StrokeRect(screen, float32(b.bounds.Min.X), float32(b.bounds.Min.Y), float32(b.bounds.Dx()), float32(b.bounds.Dy()), 1, b.borderColor, false)
	}

	m := fontFace.Metrics()
	lineHeight := m.HLineGap + m.HAscent + m.HDescent
	textWidth := text.Advance(b.label, fontFace)
	textHeight := float32(lineHeight)

	textX := float32(b.bounds.Min.X + (b.bounds.Dx()-int(textWidth))/2)
	textY := float32(b.bounds.Min.Y + (b.bounds.Dy()-int(textHeight))/2)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(textX), float64(textY))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, b.label, fontFace, op)
}
