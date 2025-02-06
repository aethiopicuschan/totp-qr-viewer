//go:build windows || freebsd || linux

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Shortcut string

const (
	ShortcutPaste Shortcut = "Paste"
)

type ShortcutWatcher struct {
	pressedKeyV       bool
	pressedKeyControl bool
	onPress           func(Shortcut)
}

func NewShortcutWatcher(onPress func(Shortcut)) *ShortcutWatcher {
	return &ShortcutWatcher{
		onPress: onPress,
	}
}

func (s *ShortcutWatcher) Update() {
	keys := inpututil.AppendPressedKeys([]ebiten.Key{})
	keyVPressed := false
	keyControlPressed := false
	for _, key := range keys {
		if key == ebiten.KeyV {
			keyVPressed = true
		} else if key == ebiten.KeyControl {
			keyControlPressed = true
		}
	}
	if keyVPressed && keyControlPressed && !s.pressedKeyV && !s.pressedKeyControl {
		if s.onPress != nil {
			s.onPress(ShortcutPaste)
		}
		s.pressedKeyV = true
		s.pressedKeyControl = true
	}
	if !keyVPressed {
		s.pressedKeyV = false
	}
	if !keyControlPressed {
		s.pressedKeyControl = false
	}
}
