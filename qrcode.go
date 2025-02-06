package main

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type QRCode struct {
	bounds image.Rectangle
	key    *otp.Key
	img    *ebiten.Image
}

func NewQRCode(bounds image.Rectangle, key *otp.Key) *QRCode {
	q := &QRCode{
		bounds: bounds,
		key:    key,
	}
	q.Update(key)
	return q
}

func (q *QRCode) Update(key *otp.Key) (err error) {
	if q.img == nil || q.key.URL() != key.URL() {
		q.key = key
		qr, err := qrcode.New(q.key.URL(), qrcode.Medium)
		if err != nil {
			return err
		}
		qrImg := qr.Image(q.bounds.Dx())
		q.img = ebiten.NewImageFromImage(qrImg)
	}
	return
}

func (q *QRCode) Draw(screen *ebiten.Image) {
	if q.img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(q.bounds.Min.X), float64(q.bounds.Min.Y))
	screen.DrawImage(q.img, op)
	code, err := totp.GenerateCode(q.key.Secret(), time.Now())
	if err != nil {
		return
	}
	top := &text.DrawOptions{}
	scale := 6
	top.GeoM.Scale(float64(scale), float64(scale))
	top.GeoM.Translate(float64(350), float64(220))
	text.Draw(screen, code, fontFace, top)
}
