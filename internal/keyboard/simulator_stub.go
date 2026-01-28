//go:build !windows

package keyboard

import "sync"

// Simulator 键盘模拟器（非Windows平台存根）
type Simulator struct {
	mu          sync.Mutex
	pressedKeys map[KeyCode]bool
	pressedMods Modifiers
}

// NewSimulator 创建键盘模拟器
func NewSimulator() (*Simulator, error) {
	return &Simulator{
		pressedKeys: make(map[KeyCode]bool),
	}, nil
}

// PressKeys 按下多个键并保持（存根）
func (s *Simulator) PressKeys(keys []KeyCode, mods Modifiers) error {
	return nil
}

// ReleaseKeys 释放多个键（存根）
func (s *Simulator) ReleaseKeys(keys []KeyCode, mods Modifiers) error {
	return nil
}

// ReleaseAllKeys 释放所有按住的键（存根）
func (s *Simulator) ReleaseAllKeys() error {
	return nil
}

// SimulateKey 模拟单个按键按下和释放（存根）
func (s *Simulator) SimulateKey(key KeyCode) error {
	return nil
}

// SimulateCombo 模拟组合键（存根）
func (s *Simulator) SimulateCombo(keys []KeyCode, mods Modifiers) error {
	return nil
}

// KeyDown 按下按键（存根）
func (s *Simulator) KeyDown(key KeyCode) error {
	return nil
}

// KeyUp 释放按键（存根）
func (s *Simulator) KeyUp(key KeyCode) error {
	return nil
}
