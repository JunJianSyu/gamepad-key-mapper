package gamepad

import (
	"context"
	"sync"
	"time"
)

// ButtonEvent 按键事件
type ButtonEvent struct {
	Button   Button
	Pressed  bool // true=按下, false=释放
	PlayerID int  // 手柄ID (0-3)
}

// Listener 手柄事件监听器
type Listener struct {
	pollInterval time.Duration
	controllerID int
	eventChan    chan ButtonEvent
	prevState    uint16 // 上一次的按键状态
	prevLT       bool   // 上一次左扳机状态
	prevRT       bool   // 上一次右扳机状态

	// 摇杆方向状态
	prevLeftStickUp    bool
	prevLeftStickDown  bool
	prevLeftStickLeft  bool
	prevLeftStickRight bool
	prevRightStickUp    bool
	prevRightStickDown  bool
	prevRightStickLeft  bool
	prevRightStickRight bool

	running bool
	mu      sync.Mutex
	cancel  context.CancelFunc
}

// NewListener 创建新的监听器
func NewListener(controllerID int) *Listener {
	return &Listener{
		pollInterval: 10 * time.Millisecond, // 100Hz 轮询
		controllerID: controllerID,
		eventChan:    make(chan ButtonEvent, 64),
	}
}

// Events 返回事件通道
func (l *Listener) Events() <-chan ButtonEvent {
	return l.eventChan
}

// Start 开始监听
func (l *Listener) Start() error {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return nil
	}

	// 确保XInput已加载
	if err := LoadXInput(); err != nil {
		l.mu.Unlock()
		return err
	}

	// 重新创建事件通道（因为可能已被关闭）
	l.eventChan = make(chan ButtonEvent, 64)
	
	// 重置状态
	l.prevState = 0
	l.prevLT = false
	l.prevRT = false
	l.prevLeftStickUp = false
	l.prevLeftStickDown = false
	l.prevLeftStickLeft = false
	l.prevLeftStickRight = false
	l.prevRightStickUp = false
	l.prevRightStickDown = false
	l.prevRightStickLeft = false
	l.prevRightStickRight = false

	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel
	l.running = true
	l.mu.Unlock()

	go l.pollLoop(ctx)
	return nil
}

// Stop 停止监听
func (l *Listener) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.running {
		return
	}

	if l.cancel != nil {
		l.cancel()
	}
	l.running = false
}

// IsRunning 检查是否正在运行
func (l *Listener) IsRunning() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.running
}

// pollLoop 轮询循环
func (l *Listener) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(l.pollInterval)
	defer ticker.Stop()
	defer close(l.eventChan) // 停止时关闭通道

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.poll()
		}
	}
}

// poll 执行一次轮询
func (l *Listener) poll() {
	state, err := GetState(l.controllerID)
	if err != nil {
		// 手柄可能未连接，静默忽略
		return
	}

	// 检测普通按键变化
	l.pollButtons(state)

	// 检测扳机变化
	l.pollTriggers(state)

	// 检测摇杆方向变化
	l.pollSticks(state)
}

// pollButtons 检测普通按键变化
func (l *Listener) pollButtons(state *XInputState) {
	currentButtons := state.Gamepad.Buttons
	changed := currentButtons ^ l.prevState

	if changed != 0 {
		// 标准按键列表（不包括扳机和摇杆方向）
		standardBtns := []Button{
			ButtonA, ButtonB, ButtonX, ButtonY,
			ButtonLB, ButtonRB,
			ButtonStart, ButtonBack, ButtonXbox, ButtonShare,
			ButtonLeftThumb, ButtonRightThumb,
			ButtonDPadUp, ButtonDPadDown, ButtonDPadLeft, ButtonDPadRight,
		}

		for _, btn := range standardBtns {
			btnMask := uint16(btn)
			if changed&btnMask != 0 {
				pressed := currentButtons&btnMask != 0
				l.sendEvent(btn, pressed)
			}
		}
	}
	l.prevState = currentButtons
}

// pollTriggers 检测扳机变化
func (l *Listener) pollTriggers(state *XInputState) {
	// 左扳机
	currentLT := state.Gamepad.LeftTrigger > TriggerThreshold
	if currentLT != l.prevLT {
		l.sendEvent(ButtonLT, currentLT)
		l.prevLT = currentLT
	}

	// 右扳机
	currentRT := state.Gamepad.RightTrigger > TriggerThreshold
	if currentRT != l.prevRT {
		l.sendEvent(ButtonRT, currentRT)
		l.prevRT = currentRT
	}
}

// pollSticks 检测摇杆方向变化
func (l *Listener) pollSticks(state *XInputState) {
	// 左摇杆
	leftUp := state.Gamepad.ThumbLY > StickThreshold
	leftDown := state.Gamepad.ThumbLY < -StickThreshold
	leftLeft := state.Gamepad.ThumbLX < -StickThreshold
	leftRight := state.Gamepad.ThumbLX > StickThreshold

	if leftUp != l.prevLeftStickUp {
		l.sendEvent(ButtonLeftStickUp, leftUp)
		l.prevLeftStickUp = leftUp
	}
	if leftDown != l.prevLeftStickDown {
		l.sendEvent(ButtonLeftStickDown, leftDown)
		l.prevLeftStickDown = leftDown
	}
	if leftLeft != l.prevLeftStickLeft {
		l.sendEvent(ButtonLeftStickLeft, leftLeft)
		l.prevLeftStickLeft = leftLeft
	}
	if leftRight != l.prevLeftStickRight {
		l.sendEvent(ButtonLeftStickRight, leftRight)
		l.prevLeftStickRight = leftRight
	}

	// 右摇杆
	rightUp := state.Gamepad.ThumbRY > StickThreshold
	rightDown := state.Gamepad.ThumbRY < -StickThreshold
	rightLeft := state.Gamepad.ThumbRX < -StickThreshold
	rightRight := state.Gamepad.ThumbRX > StickThreshold

	if rightUp != l.prevRightStickUp {
		l.sendEvent(ButtonRightStickUp, rightUp)
		l.prevRightStickUp = rightUp
	}
	if rightDown != l.prevRightStickDown {
		l.sendEvent(ButtonRightStickDown, rightDown)
		l.prevRightStickDown = rightDown
	}
	if rightLeft != l.prevRightStickLeft {
		l.sendEvent(ButtonRightStickLeft, rightLeft)
		l.prevRightStickLeft = rightLeft
	}
	if rightRight != l.prevRightStickRight {
		l.sendEvent(ButtonRightStickRight, rightRight)
		l.prevRightStickRight = rightRight
	}
}

// sendEvent 发送事件到通道
func (l *Listener) sendEvent(button Button, pressed bool) {
	event := ButtonEvent{
		Button:   button,
		Pressed:  pressed,
		PlayerID: l.controllerID,
	}

	// 非阻塞发送
	select {
	case l.eventChan <- event:
	default:
		// 通道满了，丢弃事件
	}
}
