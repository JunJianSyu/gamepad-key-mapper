package app

import (
	"fmt"
	"sync"
	"time"

	"gamepad-key-mapper/internal/config"
	"gamepad-key-mapper/internal/gamepad"
	"gamepad-key-mapper/internal/keyboard"
	"gamepad-key-mapper/internal/mapper"
)

// State 应用状态
type State int

const (
	StateStopped State = iota
	StateRunning
)

// App 应用主控制器
type App struct {
	mapper   *mapper.Mapper
	listener *gamepad.Listener
	state    State
	mu       sync.RWMutex

	// 状态变更回调
	onStateChange func(State)
	onRulesChange func()
	onError       func(error)
}

// New 创建新的应用实例
func New() *App {
	m, _ := mapper.New()
	return &App{
		mapper:   m,
		listener: gamepad.NewListener(0), // 默认监听第一个手柄
		state:    StateStopped,
	}
}

// Start 启动映射
func (a *App) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state == StateRunning {
		return nil
	}

	// 启动手柄监听
	if err := a.listener.Start(); err != nil {
		if a.onError != nil {
			a.onError(err)
		}
		return err
	}

	// 启动事件处理协程
	go a.eventLoop()

	a.state = StateRunning
	if a.onStateChange != nil {
		a.onStateChange(StateRunning)
	}

	return nil
}

// Stop 停止映射
func (a *App) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state == StateStopped {
		return
	}

	// 先释放所有按住的键
	a.mapper.ReleaseAll()

	a.listener.Stop()
	a.state = StateStopped

	if a.onStateChange != nil {
		a.onStateChange(StateStopped)
	}
}

// IsRunning 检查是否正在运行
func (a *App) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state == StateRunning
}

// GetState 获取当前状态
func (a *App) GetState() State {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state
}

// eventLoop 事件处理循环
func (a *App) eventLoop() {
	for event := range a.listener.Events() {
		a.mapper.HandleEvent(event)
	}
}

// AddRule 添加映射规则（单个目标键）
func (a *App) AddRule(source gamepad.Button, target keyboard.KeyCode, mods keyboard.Modifiers) (*mapper.MappingRule, error) {
	return a.AddRuleMultiKeys(source, []keyboard.KeyCode{target}, mods)
}

// generateRuleID 生成唯一规则ID
func (a *App) generateRuleID() string {
	return fmt.Sprintf("rule_%d", time.Now().UnixNano())
}

// AddRuleMultiKeys 添加键盘映射规则（多个目标键）
func (a *App) AddRuleMultiKeys(source gamepad.Button, targets []keyboard.KeyCode, mods keyboard.Modifiers) (*mapper.MappingRule, error) {
	// 检查冲突
	if a.mapper.HasConflict(source, "") {
		return nil, fmt.Errorf("源按键 %s 已存在映射规则", source.String())
	}

	// 生成唯一ID
	id := a.generateRuleID()

	rule := mapper.NewRuleMultiKeys(id, source, targets, mods)
	a.mapper.AddRule(rule)

	// 自动保存配置
	a.SaveConfig()

	if a.onRulesChange != nil {
		a.onRulesChange()
	}

	return rule, nil
}

// AddRuleGamepad 添加手柄到手柄的映射规则
func (a *App) AddRuleGamepad(source gamepad.Button, targets []gamepad.Button) (*mapper.MappingRule, error) {
	// 检查冲突
	if a.mapper.HasConflict(source, "") {
		return nil, fmt.Errorf("源按键 %s 已存在映射规则", source.String())
	}

	// 检查是否映射到自己
	for _, target := range targets {
		if target == source {
			return nil, fmt.Errorf("不能将按键映射到自己")
		}
	}

	// 生成唯一ID
	id := a.generateRuleID()

	rule := mapper.NewRuleGamepad(id, source, targets)
	a.mapper.AddRule(rule)

	// 自动保存配置
	a.SaveConfig()

	if a.onRulesChange != nil {
		a.onRulesChange()
	}

	return rule, nil
}

// RemoveRule 删除映射规则
func (a *App) RemoveRule(id string) bool {
	removed := a.mapper.RemoveRule(id)
	if removed {
		// 自动保存配置
		a.SaveConfig()
		
		if a.onRulesChange != nil {
			a.onRulesChange()
		}
	}
	return removed
}

// GetRules 获取所有规则
func (a *App) GetRules() []*mapper.MappingRule {
	return a.mapper.GetRules()
}

// SetRules 设置所有规则
func (a *App) SetRules(rules []*mapper.MappingRule) {
	a.mapper.SetRules(rules)
	if a.onRulesChange != nil {
		a.onRulesChange()
	}
}

// HasConflict 检查源按键是否冲突
func (a *App) HasConflict(source gamepad.Button, excludeID string) bool {
	return a.mapper.HasConflict(source, excludeID)
}

// SetOnStateChange 设置状态变更回调
func (a *App) SetOnStateChange(callback func(State)) {
	a.onStateChange = callback
}

// SetOnRulesChange 设置规则变更回调
func (a *App) SetOnRulesChange(callback func()) {
	a.onRulesChange = callback
}

// SetOnError 设置错误回调
func (a *App) SetOnError(callback func(error)) {
	a.onError = callback
}

// StateString 返回状态的字符串表示
func (s State) String() string {
	switch s {
	case StateStopped:
		return "已停止"
	case StateRunning:
		return "运行中"
	default:
		return "未知"
	}
}

// LoadConfig 加载配置
func (a *App) LoadConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	a.mapper.SetRules(cfg.Rules)
	return nil
}

// SaveConfig 保存配置
func (a *App) SaveConfig() error {
	cfg := &config.Config{
		Rules:          a.mapper.GetRules(),
		MinimizeToTray: true,
		StartMinimized: false,
	}
	return config.Save(cfg)
}
