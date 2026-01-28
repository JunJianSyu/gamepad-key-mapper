package gamepad

import "errors"

// XInputState XInput状态结构
type XInputState struct {
	PacketNumber uint32
	Gamepad      XInputGamepad
}

// XInputGamepad 手柄状态结构
type XInputGamepad struct {
	Buttons      uint16
	LeftTrigger  uint8
	RightTrigger uint8
	ThumbLX      int16
	ThumbLY      int16
	ThumbRX      int16
	ThumbRY      int16
}

// ErrXInputNotLoaded XInput未加载错误
var ErrXInputNotLoaded = errors.New("xinput not loaded")

// ErrControllerNotConnected 手柄未连接错误
var ErrControllerNotConnected = errors.New("controller not connected")

// IsButtonPressed 检查按键是否被按下
func (s *XInputState) IsButtonPressed(button Button) bool {
	// 处理扳机键
	if button == ButtonLT {
		return s.Gamepad.LeftTrigger > TriggerThreshold
	}
	if button == ButtonRT {
		return s.Gamepad.RightTrigger > TriggerThreshold
	}

	// 处理普通按键
	return s.Gamepad.Buttons&uint16(button) != 0
}

// GetPressedButtons 获取所有按下的按键
func (s *XInputState) GetPressedButtons() []Button {
	var pressed []Button

	for _, btn := range AllButtons() {
		if s.IsButtonPressed(btn) {
			pressed = append(pressed, btn)
		}
	}

	return pressed
}
