package main

import (
	"crypto/rand"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.design/x/clipboard"
)

const (
	screenWidth      = 640
	screenHeight     = 396 + 16
	labelSecret      = "Secret"
	labelIssuer      = "Issuer"
	labelAccountName = "AccountName"
	labelRandom      = "Random"
)

type Game struct {
	shortcutWatcher *ShortcutWatcher
	theme           *Theme
	qrcode          *QRCode
	textFields      []*TextField
	buttons         []*Button
}

func (g *Game) Update() (err error) {
	g.shortcutWatcher.Update()
	if g.textFields == nil {
		g.textFields = append(g.textFields, NewTextField(image.Rect(16, 16, screenWidth-100, 16+textFieldHeight), WithLabel(labelSecret)))
		g.textFields = append(g.textFields, NewTextField(image.Rect(16, 16+textFieldHeight+16, screenWidth-100, 16+textFieldHeight*2+16), WithLabel(labelIssuer)))
		g.textFields = append(g.textFields, NewTextField(image.Rect(16, 16+textFieldHeight*2+16*2, screenWidth-100, 16+textFieldHeight*3+16*2), WithLabel(labelAccountName)))
	}
	if g.buttons == nil {
		g.buttons = append(g.buttons, NewButton(image.Rect(screenWidth-100+16, 16, screenWidth-16, 46), labelRandom, g.theme.MainColor, g.theme.BorderColor, func() {
			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "Issuer",
				AccountName: "AccountName",
				Rand:        rand.Reader,
			})
			if err != nil {
				log.Fatal(err)
			}
			for _, tf := range g.textFields {
				if tf.label == labelSecret {
					tf.field.SetTextAndSelection(string(key.Secret()), 0, len(key.Secret()))
				}
			}
		}))
	}

	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(ids) > 0 {
		var x, y int
		if len(ids) > 0 {
			x, y = ebiten.TouchPosition(ids[0])
		} else {
			x, y = ebiten.CursorPosition()
		}
		for _, tf := range g.textFields {
			if tf.Contains(x, y) {
				tf.Focus()
				tf.SetSelectionStartByCursorPosition(x, y)
			} else {
				tf.Blur()
			}
		}
	}

	for _, tf := range g.textFields {
		if err := tf.Update(); err != nil {
			return err
		}
	}
	for _, bt := range g.buttons {
		if err := bt.Update(); err != nil {
			return err
		}
	}

	x, y := ebiten.CursorPosition()
	var inTextField bool
	var inButton bool
	for _, tf := range g.textFields {
		if tf.Contains(x, y) {
			inTextField = true
			break
		}
	}
	for _, bt := range g.buttons {
		if bt.Contains(x, y) {
			inButton = true
			break
		}
	}
	if inTextField {
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	} else if inButton {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	var secret, issuer, accountName string
	for _, textField := range g.textFields {
		switch textField.label {
		case labelSecret:
			secret = textField.field.Text()
		case labelIssuer:
			issuer = textField.field.Text()
		case labelAccountName:
			accountName = textField.field.Text()
		}
	}
	if secret != "" && issuer != "" && accountName != "" {
		uri := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, accountName, secret, issuer)
		key, err := otp.NewKeyFromURL(uri)
		if err != nil {
			return err
		}
		g.qrcode = NewQRCode(image.Rect(16, 140, 16+256, 140+256), key)
	} else if g.qrcode != nil {
		g.qrcode = nil
	}

	return
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.theme.BackgroundColor)
	for _, tf := range g.textFields {
		tf.Draw(screen)
	}
	for _, bt := range g.buttons {
		bt.Draw(screen)
	}
	if g.qrcode != nil {
		g.qrcode.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("TOTP QR Viewer")

	err := clipboard.Init()
	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		theme: GetDefaultTheme(),
	}
	g.shortcutWatcher = NewShortcutWatcher(func(s Shortcut) {
		switch s {
		case ShortcutPaste:
			for _, tf := range g.textFields {
				if tf.field.IsFocused() {
					b := clipboard.Read(clipboard.FmtText)
					tf.Paste(string(b))
				}
			}
		}
	})

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
