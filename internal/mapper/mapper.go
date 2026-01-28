package mapper

import (
	"sync"

	"gamepad-key-mapper/internal/gamepad"
	"gamepad-key-mapper/internal/keyboard"
)

// Mapper 映射引擎
type Mapper struct {
	rules     []*MappingRule
	simulator *keyboard.Simulator
	mu        sync.RWMutex
	
	// 用于防止循环映射的处理中标记
	processing   map[gamepad.Button]bool
	processingMu sync.Mutex // 保护 processing map
}

// New 创建新的映射引擎
func New() (*Mapper, error) {
	sim, err := keyboard.NewSimulator()
	if err != nil {
		return nil, err
	}

	return &Mapper{
		rules:      make([]*MappingRule, 0),
		simulator:  sim,
		processing: make(map[gamepad.Button]bool),
	}, nil
}

// AddRule 添加映射规则
func (m *Mapper) AddRule(rule *MappingRule) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules = append(m.rules, rule)
}

// RemoveRule 删除映射规则
func (m *Mapper) RemoveRule(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, rule := range m.rules {
		if rule.ID == id {
			m.rules = append(m.rules[:i], m.rules[i+1:]...)
			return true
		}
	}
	return false
}

// GetRules 获取所有规则
func (m *Mapper) GetRules() []*MappingRule {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回副本
	rules := make([]*MappingRule, len(m.rules))
	copy(rules, m.rules)
	return rules
}

// SetRules 设置所有规则
func (m *Mapper) SetRules(rules []*MappingRule) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules = rules
}

// ClearRules 清空所有规则
func (m *Mapper) ClearRules() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules = make([]*MappingRule, 0)
}

// HasConflict 检查是否存在源按键冲突
func (m *Mapper) HasConflict(sourceKey gamepad.Button, excludeID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rule := range m.rules {
		if rule.SourceKey == sourceKey && rule.ID != excludeID {
			return true
		}
	}
	return false
}

// HandleEvent 处理手柄按键事件
func (m *Mapper) HandleEvent(event gamepad.ButtonEvent) {
	// 检查是否正在处理（防止循环）
	m.processingMu.Lock()
	if m.processing[event.Button] {
		m.processingMu.Unlock()
		return
	}
	m.processingMu.Unlock()

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rule := range m.rules {
		if rule.Enabled && rule.SourceKey == event.Button {
			if rule.TargetType == TargetKeyboard {
				// 键盘映射
				if event.Pressed {
					m.simulator.PressKeys(rule.TargetKeys, rule.Modifiers)
				} else {
					m.simulator.ReleaseKeys(rule.TargetKeys, rule.Modifiers)
				}
			} else if rule.TargetType == TargetGamepad {
				// 手柄到手柄映射：触发目标按键的映射
				m.handleGamepadMapping(rule, event.Pressed, event.PlayerID)
			}
			break // 每个源按键只匹配一个规则
		}
	}
}

// handleGamepadMapping 处理手柄到手柄的映射
func (m *Mapper) handleGamepadMapping(rule *MappingRule, pressed bool, playerID int) {
	// 标记源按键正在处理，防止循环
	m.processingMu.Lock()
	m.processing[rule.SourceKey] = true
	m.processingMu.Unlock()
	
	defer func() {
		m.processingMu.Lock()
		delete(m.processing, rule.SourceKey)
		m.processingMu.Unlock()
	}()

	// 为每个目标手柄按键触发映射
	for _, targetBtn := range rule.TargetButtons {
		// 查找目标按键的映射规则
		for _, targetRule := range m.rules {
			if targetRule.Enabled && targetRule.SourceKey == targetBtn && targetRule.ID != rule.ID {
				if targetRule.TargetType == TargetKeyboard {
					// 目标按键映射到键盘
					if pressed {
						m.simulator.PressKeys(targetRule.TargetKeys, targetRule.Modifiers)
					} else {
						m.simulator.ReleaseKeys(targetRule.TargetKeys, targetRule.Modifiers)
					}
				} else if targetRule.TargetType == TargetGamepad {
					// 目标按键也是手柄映射，递归处理（有循环保护）
					m.processingMu.Lock()
					isProcessing := m.processing[targetBtn]
					if !isProcessing {
						m.processing[targetBtn] = true
					}
					m.processingMu.Unlock()
					
					if !isProcessing {
						m.handleGamepadMapping(targetRule, pressed, playerID)
						m.processingMu.Lock()
						delete(m.processing, targetBtn)
						m.processingMu.Unlock()
					}
				}
				break
			}
		}
	}
}

// ReleaseAll 释放所有按键（用于停止映射时）
func (m *Mapper) ReleaseAll() {
	m.simulator.ReleaseAllKeys()
}

// FindRuleBySource 根据源按键查找规则
func (m *Mapper) FindRuleBySource(button gamepad.Button) *MappingRule {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rule := range m.rules {
		if rule.SourceKey == button {
			return rule
		}
	}
	return nil
}

// GetRuleByID 根据ID获取规则
func (m *Mapper) GetRuleByID(id string) *MappingRule {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rule := range m.rules {
		if rule.ID == id {
			return rule
		}
	}
	return nil
}
