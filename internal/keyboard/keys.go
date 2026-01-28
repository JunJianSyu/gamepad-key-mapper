package keyboard

// KeyCode 表示键盘按键码
type KeyCode int

// 常用键盘按键码（Windows Virtual Key Codes）
const (
	// 功能键
	KeyF1  KeyCode = 0x70
	KeyF2  KeyCode = 0x71
	KeyF3  KeyCode = 0x72
	KeyF4  KeyCode = 0x73
	KeyF5  KeyCode = 0x74
	KeyF6  KeyCode = 0x75
	KeyF7  KeyCode = 0x76
	KeyF8  KeyCode = 0x77
	KeyF9  KeyCode = 0x78
	KeyF10 KeyCode = 0x79
	KeyF11 KeyCode = 0x7A
	KeyF12 KeyCode = 0x7B

	// 字母键
	KeyA KeyCode = 0x41
	KeyB KeyCode = 0x42
	KeyC KeyCode = 0x43
	KeyD KeyCode = 0x44
	KeyE KeyCode = 0x45
	KeyF KeyCode = 0x46
	KeyG KeyCode = 0x47
	KeyH KeyCode = 0x48
	KeyI KeyCode = 0x49
	KeyJ KeyCode = 0x4A
	KeyK KeyCode = 0x4B
	KeyL KeyCode = 0x4C
	KeyM KeyCode = 0x4D
	KeyN KeyCode = 0x4E
	KeyO KeyCode = 0x4F
	KeyP KeyCode = 0x50
	KeyQ KeyCode = 0x51
	KeyR KeyCode = 0x52
	KeyS KeyCode = 0x53
	KeyT KeyCode = 0x54
	KeyU KeyCode = 0x55
	KeyV KeyCode = 0x56
	KeyW KeyCode = 0x57
	KeyX KeyCode = 0x58
	KeyY KeyCode = 0x59
	KeyZ KeyCode = 0x5A

	// 数字键
	Key0 KeyCode = 0x30
	Key1 KeyCode = 0x31
	Key2 KeyCode = 0x32
	Key3 KeyCode = 0x33
	Key4 KeyCode = 0x34
	Key5 KeyCode = 0x35
	Key6 KeyCode = 0x36
	Key7 KeyCode = 0x37
	Key8 KeyCode = 0x38
	Key9 KeyCode = 0x39

	// 特殊键
	KeySpace     KeyCode = 0x20
	KeyEnter     KeyCode = 0x0D
	KeyTab       KeyCode = 0x09
	KeyEscape    KeyCode = 0x1B
	KeyBackspace KeyCode = 0x08
	KeyDelete    KeyCode = 0x2E
	KeyInsert    KeyCode = 0x2D
	KeyHome      KeyCode = 0x24
	KeyEnd       KeyCode = 0x23
	KeyPageUp    KeyCode = 0x21
	KeyPageDown  KeyCode = 0x22

	// 方向键
	KeyUp    KeyCode = 0x26
	KeyDown  KeyCode = 0x28
	KeyLeft  KeyCode = 0x25
	KeyRight KeyCode = 0x27

	// 小键盘
	KeyNumpad0 KeyCode = 0x60
	KeyNumpad1 KeyCode = 0x61
	KeyNumpad2 KeyCode = 0x62
	KeyNumpad3 KeyCode = 0x63
	KeyNumpad4 KeyCode = 0x64
	KeyNumpad5 KeyCode = 0x65
	KeyNumpad6 KeyCode = 0x66
	KeyNumpad7 KeyCode = 0x67
	KeyNumpad8 KeyCode = 0x68
	KeyNumpad9 KeyCode = 0x69
)

// String 返回按键名称
func (k KeyCode) String() string {
	switch k {
	case KeyF1:
		return "F1"
	case KeyF2:
		return "F2"
	case KeyF3:
		return "F3"
	case KeyF4:
		return "F4"
	case KeyF5:
		return "F5"
	case KeyF6:
		return "F6"
	case KeyF7:
		return "F7"
	case KeyF8:
		return "F8"
	case KeyF9:
		return "F9"
	case KeyF10:
		return "F10"
	case KeyF11:
		return "F11"
	case KeyF12:
		return "F12"
	case KeyA:
		return "A"
	case KeyB:
		return "B"
	case KeyC:
		return "C"
	case KeyD:
		return "D"
	case KeyE:
		return "E"
	case KeyF:
		return "F"
	case KeyG:
		return "G"
	case KeyH:
		return "H"
	case KeyI:
		return "I"
	case KeyJ:
		return "J"
	case KeyK:
		return "K"
	case KeyL:
		return "L"
	case KeyM:
		return "M"
	case KeyN:
		return "N"
	case KeyO:
		return "O"
	case KeyP:
		return "P"
	case KeyQ:
		return "Q"
	case KeyR:
		return "R"
	case KeyS:
		return "S"
	case KeyT:
		return "T"
	case KeyU:
		return "U"
	case KeyV:
		return "V"
	case KeyW:
		return "W"
	case KeyX:
		return "X"
	case KeyY:
		return "Y"
	case KeyZ:
		return "Z"
	case Key0:
		return "0"
	case Key1:
		return "1"
	case Key2:
		return "2"
	case Key3:
		return "3"
	case Key4:
		return "4"
	case Key5:
		return "5"
	case Key6:
		return "6"
	case Key7:
		return "7"
	case Key8:
		return "8"
	case Key9:
		return "9"
	case KeySpace:
		return "Space"
	case KeyEnter:
		return "Enter"
	case KeyTab:
		return "Tab"
	case KeyEscape:
		return "Escape"
	case KeyBackspace:
		return "Backspace"
	case KeyDelete:
		return "Delete"
	case KeyInsert:
		return "Insert"
	case KeyHome:
		return "Home"
	case KeyEnd:
		return "End"
	case KeyPageUp:
		return "PageUp"
	case KeyPageDown:
		return "PageDown"
	case KeyUp:
		return "Up"
	case KeyDown:
		return "Down"
	case KeyLeft:
		return "Left"
	case KeyRight:
		return "Right"
	default:
		return "Unknown"
	}
}

// Modifiers 表示修饰键组合
type Modifiers struct {
	Ctrl  bool `json:"ctrl"`
	Alt   bool `json:"alt"`
	Shift bool `json:"shift"`
	Win   bool `json:"win"`
}

// AllKeys 返回所有可选的目标键
func AllKeys() []KeyCode {
	return []KeyCode{
		KeyF1, KeyF2, KeyF3, KeyF4, KeyF5, KeyF6,
		KeyF7, KeyF8, KeyF9, KeyF10, KeyF11, KeyF12,
		KeyA, KeyB, KeyC, KeyD, KeyE, KeyF, KeyG, KeyH,
		KeyI, KeyJ, KeyK, KeyL, KeyM, KeyN, KeyO, KeyP,
		KeyQ, KeyR, KeyS, KeyT, KeyU, KeyV, KeyW, KeyX, KeyY, KeyZ,
		Key0, Key1, Key2, Key3, Key4, Key5, Key6, Key7, Key8, Key9,
		KeySpace, KeyEnter, KeyTab, KeyEscape,
		KeyUp, KeyDown, KeyLeft, KeyRight,
	}
}
