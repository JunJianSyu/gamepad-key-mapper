package gamepad

// Button 表示游戏手柄按键
type Button uint32

// Xbox手柄按键常量（与XInput API对应）
const (
	// 主按键
	ButtonA Button = 0x1000
	ButtonB Button = 0x2000
	ButtonX Button = 0x4000
	ButtonY Button = 0x8000

	// 肩键 (Bumpers)
	ButtonLB Button = 0x0100 // 左肩键 (Left Bumper)
	ButtonRB Button = 0x0200 // 右肩键 (Right Bumper)

	// 功能键
	ButtonStart  Button = 0x0010 // 开始/菜单键 (Menu)
	ButtonBack   Button = 0x0020 // 返回/视图键 (View)
	ButtonXbox   Button = 0x0400 // Xbox/Guide按钮 (需要特殊处理)
	ButtonShare  Button = 0x0800 // 分享按钮 (Xbox Series X|S)

	// 摇杆按下
	ButtonLeftThumb  Button = 0x0040 // 左摇杆按下 (LS)
	ButtonRightThumb Button = 0x0080 // 右摇杆按下 (RS)

	// 方向键 (D-Pad)
	ButtonDPadUp    Button = 0x0001
	ButtonDPadDown  Button = 0x0002
	ButtonDPadLeft  Button = 0x0004
	ButtonDPadRight Button = 0x0008
)

// 扳机键使用特殊值（因为扳机是模拟量，这里定义阈值触发）
const (
	ButtonLT Button = 0x10000 // 左扳机 (Left Trigger)
	ButtonRT Button = 0x20000 // 右扳机 (Right Trigger)
)

// Xbox精英版手柄背部拨片 (Paddles)
const (
	ButtonPaddle1 Button = 0x40000  // 背部拨片1 (P1/上左)
	ButtonPaddle2 Button = 0x80000  // 背部拨片2 (P2/上右)
	ButtonPaddle3 Button = 0x100000 // 背部拨片3 (P3/下左)
	ButtonPaddle4 Button = 0x200000 // 背部拨片4 (P4/下右)
)

// 摇杆方向映射（虚拟按键，用于将摇杆方向映射为按键）
const (
	ButtonLeftStickUp    Button = 0x400000  // 左摇杆上
	ButtonLeftStickDown  Button = 0x800000  // 左摇杆下
	ButtonLeftStickLeft  Button = 0x1000000 // 左摇杆左
	ButtonLeftStickRight Button = 0x2000000 // 左摇杆右

	ButtonRightStickUp    Button = 0x4000000  // 右摇杆上
	ButtonRightStickDown  Button = 0x8000000  // 右摇杆下
	ButtonRightStickLeft  Button = 0x10000000 // 右摇杆左
	ButtonRightStickRight Button = 0x20000000 // 右摇杆右
)

// TriggerThreshold 扳机触发阈值（0-255，超过此值视为按下）
const TriggerThreshold byte = 128

// StickThreshold 摇杆触发阈值（-32768到32767，超过此值视为推动）
const StickThreshold int16 = 16384

// String 返回按键名称
func (b Button) String() string {
	switch b {
	// 主按键
	case ButtonA:
		return "A"
	case ButtonB:
		return "B"
	case ButtonX:
		return "X"
	case ButtonY:
		return "Y"

	// 肩键
	case ButtonLB:
		return "LB"
	case ButtonRB:
		return "RB"

	// 扳机
	case ButtonLT:
		return "LT"
	case ButtonRT:
		return "RT"

	// 功能键
	case ButtonStart:
		return "Menu"
	case ButtonBack:
		return "View"
	case ButtonXbox:
		return "Xbox"
	case ButtonShare:
		return "Share"

	// 摇杆按下
	case ButtonLeftThumb:
		return "LS (按下)"
	case ButtonRightThumb:
		return "RS (按下)"

	// 方向键
	case ButtonDPadUp:
		return "D-Pad ↑"
	case ButtonDPadDown:
		return "D-Pad ↓"
	case ButtonDPadLeft:
		return "D-Pad ←"
	case ButtonDPadRight:
		return "D-Pad →"

	// 背部拨片
	case ButtonPaddle1:
		return "P1 (左上拨片)"
	case ButtonPaddle2:
		return "P2 (右上拨片)"
	case ButtonPaddle3:
		return "P3 (左下拨片)"
	case ButtonPaddle4:
		return "P4 (右下拨片)"

	// 左摇杆方向
	case ButtonLeftStickUp:
		return "左摇杆 ↑"
	case ButtonLeftStickDown:
		return "左摇杆 ↓"
	case ButtonLeftStickLeft:
		return "左摇杆 ←"
	case ButtonLeftStickRight:
		return "左摇杆 →"

	// 右摇杆方向
	case ButtonRightStickUp:
		return "右摇杆 ↑"
	case ButtonRightStickDown:
		return "右摇杆 ↓"
	case ButtonRightStickLeft:
		return "右摇杆 ←"
	case ButtonRightStickRight:
		return "右摇杆 →"

	default:
		return "Unknown"
	}
}

// AllButtons 返回所有可用的手柄按键（按类别分组）
func AllButtons() []Button {
	return []Button{
		// 主按键
		ButtonA,
		ButtonB,
		ButtonX,
		ButtonY,

		// 肩键和扳机
		ButtonLB,
		ButtonRB,
		ButtonLT,
		ButtonRT,

		// 功能键
		ButtonStart,
		ButtonBack,
		ButtonShare,

		// 摇杆按下
		ButtonLeftThumb,
		ButtonRightThumb,

		// 方向键
		ButtonDPadUp,
		ButtonDPadDown,
		ButtonDPadLeft,
		ButtonDPadRight,

		// 精英版背部拨片
		ButtonPaddle1,
		ButtonPaddle2,
		ButtonPaddle3,
		ButtonPaddle4,

		// 左摇杆方向
		ButtonLeftStickUp,
		ButtonLeftStickDown,
		ButtonLeftStickLeft,
		ButtonLeftStickRight,

		// 右摇杆方向
		ButtonRightStickUp,
		ButtonRightStickDown,
		ButtonRightStickLeft,
		ButtonRightStickRight,
	}
}

// StandardButtons 返回标准Xbox手柄按键（不含精英版特有按键）
func StandardButtons() []Button {
	return []Button{
		ButtonA, ButtonB, ButtonX, ButtonY,
		ButtonLB, ButtonRB, ButtonLT, ButtonRT,
		ButtonStart, ButtonBack,
		ButtonLeftThumb, ButtonRightThumb,
		ButtonDPadUp, ButtonDPadDown, ButtonDPadLeft, ButtonDPadRight,
	}
}

// EliteButtons 返回精英版特有按键
func EliteButtons() []Button {
	return []Button{
		ButtonPaddle1, ButtonPaddle2, ButtonPaddle3, ButtonPaddle4,
	}
}

// StickButtons 返回摇杆方向映射按键
func StickButtons() []Button {
	return []Button{
		ButtonLeftStickUp, ButtonLeftStickDown, ButtonLeftStickLeft, ButtonLeftStickRight,
		ButtonRightStickUp, ButtonRightStickDown, ButtonRightStickLeft, ButtonRightStickRight,
	}
}

// IsStickButton 检查是否为摇杆方向按键
func (b Button) IsStickButton() bool {
	return b >= ButtonLeftStickUp && b <= ButtonRightStickRight
}

// IsPaddleButton 检查是否为背部拨片按键
func (b Button) IsPaddleButton() bool {
	return b >= ButtonPaddle1 && b <= ButtonPaddle4
}

// IsTriggerButton 检查是否为扳机按键
func (b Button) IsTriggerButton() bool {
	return b == ButtonLT || b == ButtonRT
}
