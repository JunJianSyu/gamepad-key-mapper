package mapper

import (
	"strings"

	"gamepad-key-mapper/internal/gamepad"
	"gamepad-key-mapper/internal/keyboard"
)

// TargetType 目标类型
type TargetType int

const (
	TargetKeyboard TargetType = iota // 目标是键盘按键
	TargetGamepad                    // 目标是手柄按键（内部转发）
)

// MappingRule 定义一条从手柄按键到目标的映射规则
type MappingRule struct {
	ID         string             `json:"id"`          // 唯一标识
	Name       string             `json:"name"`        // 规则名称（可选）
	SourceKey  gamepad.Button     `json:"source_key"`  // 源按键（手柄）
	TargetType TargetType         `json:"target_type"` // 目标类型
	
	// 键盘目标（当 TargetType == TargetKeyboard）
	TargetKeys []keyboard.KeyCode `json:"target_keys"` // 目标按键（键盘，支持多键）
	Modifiers  keyboard.Modifiers `json:"modifiers"`   // 修饰键
	
	// 手柄目标（当 TargetType == TargetGamepad）
	TargetButtons []gamepad.Button `json:"target_buttons"` // 目标按键（手柄，支持多键）
	
	Enabled bool `json:"enabled"` // 是否启用
}

// NewRule 创建一个键盘映射规则（单个目标键）
func NewRule(id string, source gamepad.Button, target keyboard.KeyCode, mods keyboard.Modifiers) *MappingRule {
	return &MappingRule{
		ID:         id,
		SourceKey:  source,
		TargetType: TargetKeyboard,
		TargetKeys: []keyboard.KeyCode{target},
		Modifiers:  mods,
		Enabled:    true,
	}
}

// NewRuleMultiKeys 创建一个键盘映射规则（多个目标键）
func NewRuleMultiKeys(id string, source gamepad.Button, targets []keyboard.KeyCode, mods keyboard.Modifiers) *MappingRule {
	return &MappingRule{
		ID:         id,
		SourceKey:  source,
		TargetType: TargetKeyboard,
		TargetKeys: targets,
		Modifiers:  mods,
		Enabled:    true,
	}
}

// NewRuleGamepad 创建一个手柄映射规则（手柄按键到手柄按键）
func NewRuleGamepad(id string, source gamepad.Button, targets []gamepad.Button) *MappingRule {
	return &MappingRule{
		ID:            id,
		SourceKey:     source,
		TargetType:    TargetGamepad,
		TargetButtons: targets,
		Enabled:       true,
	}
}

// String 返回规则的可读描述
func (r *MappingRule) String() string {
	sourceStr := r.SourceKey.String()
	
	if r.TargetType == TargetGamepad {
		// 手柄到手柄映射
		var btnNames []string
		for _, btn := range r.TargetButtons {
			btnNames = append(btnNames, btn.String())
		}
		return sourceStr + " → 🎮 " + strings.Join(btnNames, "+")
	}
	
	// 键盘映射
	modStr := ""
	if r.Modifiers.Ctrl {
		modStr += "Ctrl+"
	}
	if r.Modifiers.Alt {
		modStr += "Alt+"
	}
	if r.Modifiers.Shift {
		modStr += "Shift+"
	}
	if r.Modifiers.Win {
		modStr += "Win+"
	}

	var keyNames []string
	for _, key := range r.TargetKeys {
		keyNames = append(keyNames, key.String())
	}
	targetStr := strings.Join(keyNames, "+")

	return sourceStr + " → ⌨️ " + modStr + targetStr
}

// IsKeyboardMapping 检查是否为键盘映射
func (r *MappingRule) IsKeyboardMapping() bool {
	return r.TargetType == TargetKeyboard
}

// IsGamepadMapping 检查是否为手柄映射
func (r *MappingRule) IsGamepadMapping() bool {
	return r.TargetType == TargetGamepad
}

// GetFirstTargetKey 获取第一个目标键（兼容旧代码）
func (r *MappingRule) GetFirstTargetKey() keyboard.KeyCode {
	if len(r.TargetKeys) > 0 {
		return r.TargetKeys[0]
	}
	return 0
}

// HasMultipleTargets 检查是否有多个目标
func (r *MappingRule) HasMultipleTargets() bool {
	if r.TargetType == TargetKeyboard {
		return len(r.TargetKeys) > 1
	}
	return len(r.TargetButtons) > 1
}
