//go:build !windows

package gamepad

import "errors"

// LoadXInput 加载XInput DLL（非Windows平台存根）
func LoadXInput() error {
	return errors.New("xinput is only supported on Windows")
}

// IsLoaded 检查XInput是否已加载
func IsLoaded() bool {
	return false
}

// GetState 获取指定手柄的状态（非Windows平台存根）
func GetState(controllerID int) (*XInputState, error) {
	return nil, errors.New("xinput is only supported on Windows")
}
