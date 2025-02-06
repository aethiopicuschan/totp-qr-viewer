//go:build darwin

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
	pressedKeyV bool
	pressedMeta bool
	onPress     func(Shortcut)
}

func NewShortcutWatcher(onPress func(Shortcut)) *ShortcutWatcher {
	return &ShortcutWatcher{
		onPress: onPress,
	}
}

func (s *ShortcutWatcher) Update() {
	keys := inpututil.AppendPressedKeys([]ebiten.Key{})
	keyVPressed := false
	metaPressed := false
	for _, key := range keys {
		if key == ebiten.KeyV {
			keyVPressed = true
		} else if key == ebiten.KeyMetaLeft || key == ebiten.KeyMetaRight {
			metaPressed = true
		}
	}
	if keyVPressed && metaPressed && !s.pressedKeyV && !s.pressedMeta {
		if s.onPress != nil {
			s.onPress(ShortcutPaste)
		}
		s.pressedKeyV = true
		s.pressedMeta = true
	}
	if !keyVPressed {
		s.pressedKeyV = false
	}
	if !metaPressed {
		s.pressedMeta = false
	}
}
