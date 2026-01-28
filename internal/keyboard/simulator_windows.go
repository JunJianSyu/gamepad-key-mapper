//go:build windows

package keyboard

import (
	"sync"
	"syscall"
	"unsafe"
)

var (
	user32           *syscall.DLL
	procSendInput    *syscall.Proc
	keyboardInitOnce sync.Once
	keyboardInitErr  error
)

// INPUT 结构体类型
const (
	INPUT_KEYBOARD = 1
)

// KEYEVENTF 标志
const (
	KEYEVENTF_KEYUP       = 0x0002
	KEYEVENTF_EXTENDEDKEY = 0x0001
)

// INPUT 结构体
type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
}

// KEYBDINPUT 结构体
type KEYBDINPUT struct {
	Vk        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
	_         [8]byte // padding
}

// Simulator 键盘模拟器
type Simulator struct {
	mu           sync.Mutex
	pressedKeys  map[KeyCode]bool // 当前按住的键
	pressedMods  Modifiers        // 当前按住的修饰键
}

// NewSimulator 创建键盘模拟器
func NewSimulator() (*Simulator, error) {
	keyboardInitOnce.Do(func() {
		user32, keyboardInitErr = syscall.LoadDLL("user32.dll")
		if keyboardInitErr != nil {
			return
		}
		procSendInput, keyboardInitErr = user32.FindProc("SendInput")
	})

	if keyboardInitErr != nil {
		return nil, keyboardInitErr
	}

	return &Simulator{
		pressedKeys: make(map[KeyCode]bool),
	}, nil
}

// PressKeys 按下多个键并保持（不释放）
func (s *Simulator) PressKeys(keys []KeyCode, mods Modifiers) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var inputs []INPUT

	// 按下修饰键（如果之前没有按下）
	if mods.Ctrl && !s.pressedMods.Ctrl {
		inputs = append(inputs, s.makeKeyInput(0x11, false)) // VK_CONTROL
		s.pressedMods.Ctrl = true
	}
	if mods.Alt && !s.pressedMods.Alt {
		inputs = append(inputs, s.makeKeyInput(0x12, false)) // VK_MENU
		s.pressedMods.Alt = true
	}
	if mods.Shift && !s.pressedMods.Shift {
		inputs = append(inputs, s.makeKeyInput(0x10, false)) // VK_SHIFT
		s.pressedMods.Shift = true
	}
	if mods.Win && !s.pressedMods.Win {
		inputs = append(inputs, s.makeKeyInput(0x5B, false)) // VK_LWIN
		s.pressedMods.Win = true
	}

	// 按下目标键（如果之前没有按下）
	for _, key := range keys {
		if !s.pressedKeys[key] {
			inputs = append(inputs, s.makeKeyInput(uint16(key), false))
			s.pressedKeys[key] = true
		}
	}

	if len(inputs) == 0 {
		return nil // 所有键都已经按下
	}

	return s.sendInputs(inputs)
}

// ReleaseKeys 释放多个键
func (s *Simulator) ReleaseKeys(keys []KeyCode, mods Modifiers) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var inputs []INPUT

	// 释放目标键
	for _, key := range keys {
		if s.pressedKeys[key] {
			inputs = append(inputs, s.makeKeyInput(uint16(key), true))
			delete(s.pressedKeys, key)
		}
	}

	// 释放修饰键（如果需要）
	if mods.Win && s.pressedMods.Win {
		inputs = append(inputs, s.makeKeyInput(0x5B, true))
		s.pressedMods.Win = false
	}
	if mods.Shift && s.pressedMods.Shift {
		inputs = append(inputs, s.makeKeyInput(0x10, true))
		s.pressedMods.Shift = false
	}
	if mods.Alt && s.pressedMods.Alt {
		inputs = append(inputs, s.makeKeyInput(0x12, true))
		s.pressedMods.Alt = false
	}
	if mods.Ctrl && s.pressedMods.Ctrl {
		inputs = append(inputs, s.makeKeyInput(0x11, true))
		s.pressedMods.Ctrl = false
	}

	if len(inputs) == 0 {
		return nil
	}

	return s.sendInputs(inputs)
}

// ReleaseAllKeys 释放所有按住的键
func (s *Simulator) ReleaseAllKeys() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var inputs []INPUT

	// 释放所有按住的目标键
	for key := range s.pressedKeys {
		inputs = append(inputs, s.makeKeyInput(uint16(key), true))
	}
	s.pressedKeys = make(map[KeyCode]bool)

	// 释放所有修饰键
	if s.pressedMods.Win {
		inputs = append(inputs, s.makeKeyInput(0x5B, true))
	}
	if s.pressedMods.Shift {
		inputs = append(inputs, s.makeKeyInput(0x10, true))
	}
	if s.pressedMods.Alt {
		inputs = append(inputs, s.makeKeyInput(0x12, true))
	}
	if s.pressedMods.Ctrl {
		inputs = append(inputs, s.makeKeyInput(0x11, true))
	}
	s.pressedMods = Modifiers{}

	if len(inputs) == 0 {
		return nil
	}

	return s.sendInputs(inputs)
}

// SimulateKey 模拟单个按键按下和释放（一次性触发）
func (s *Simulator) SimulateKey(key KeyCode) error {
	return s.SimulateCombo([]KeyCode{key}, Modifiers{})
}

// SimulateCombo 模拟组合键（一次性触发：按下-释放）
func (s *Simulator) SimulateCombo(keys []KeyCode, mods Modifiers) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var inputs []INPUT

	// 按下修饰键
	if mods.Ctrl {
		inputs = append(inputs, s.makeKeyInput(0x11, false))
	}
	if mods.Alt {
		inputs = append(inputs, s.makeKeyInput(0x12, false))
	}
	if mods.Shift {
		inputs = append(inputs, s.makeKeyInput(0x10, false))
	}
	if mods.Win {
		inputs = append(inputs, s.makeKeyInput(0x5B, false))
	}

	// 按下所有目标键
	for _, key := range keys {
		inputs = append(inputs, s.makeKeyInput(uint16(key), false))
	}

	// 发送按下事件
	if err := s.sendInputs(inputs); err != nil {
		return err
	}

	// 立即释放（反向顺序）
	var releaseInputs []INPUT

	// 释放目标键
	for i := len(keys) - 1; i >= 0; i-- {
		releaseInputs = append(releaseInputs, s.makeKeyInput(uint16(keys[i]), true))
	}

	// 释放修饰键
	if mods.Win {
		releaseInputs = append(releaseInputs, s.makeKeyInput(0x5B, true))
	}
	if mods.Shift {
		releaseInputs = append(releaseInputs, s.makeKeyInput(0x10, true))
	}
	if mods.Alt {
		releaseInputs = append(releaseInputs, s.makeKeyInput(0x12, true))
	}
	if mods.Ctrl {
		releaseInputs = append(releaseInputs, s.makeKeyInput(0x11, true))
	}

	return s.sendInputs(releaseInputs)
}

// makeKeyInput 创建键盘输入结构
func (s *Simulator) makeKeyInput(vk uint16, keyUp bool) INPUT {
	input := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Vk: vk,
		},
	}

	if keyUp {
		input.Ki.Flags = KEYEVENTF_KEYUP
	}

	// 扩展键（如功能键）
	if vk >= 0x70 && vk <= 0x7B { // F1-F12
		input.Ki.Flags |= KEYEVENTF_EXTENDEDKEY
	}

	return input
}

// sendInputs 发送输入事件
func (s *Simulator) sendInputs(inputs []INPUT) error {
	if len(inputs) == 0 {
		return nil
	}

	inputSize := unsafe.Sizeof(INPUT{})
	ret, _, err := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(inputSize),
	)

	if ret == 0 {
		return err
	}

	return nil
}

// KeyDown 按下按键（不释放）
func (s *Simulator) KeyDown(key KeyCode) error {
	return s.PressKeys([]KeyCode{key}, Modifiers{})
}

// KeyUp 释放按键
func (s *Simulator) KeyUp(key KeyCode) error {
	return s.ReleaseKeys([]KeyCode{key}, Modifiers{})
}
