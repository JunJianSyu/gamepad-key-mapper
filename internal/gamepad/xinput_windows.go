//go:build windows

package gamepad

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"
)

var (
	xinputDLL       *syscall.DLL
	procGetState    *syscall.Proc
	xinputLoaded    bool
	xinputLoadOnce  sync.Once
	xinputLoadError error
)

// LoadXInput 加载XInput DLL
func LoadXInput() error {
	xinputLoadOnce.Do(func() {
		// 尝试加载不同版本的XInput DLL
		dllNames := []string{
			"xinput1_4.dll",   // Windows 8.1+
			"xinput1_3.dll",   // Windows 7
			"xinput9_1_0.dll", // Windows Vista
		}

		for _, name := range dllNames {
			dll, err := syscall.LoadDLL(name)
			if err == nil {
				xinputDLL = dll
				procGetState, err = dll.FindProc("XInputGetState")
				if err == nil {
					xinputLoaded = true
					return
				}
			}
		}
		xinputLoadError = errors.New("failed to load xinput dll")
	})

	return xinputLoadError
}

// IsLoaded 检查XInput是否已加载
func IsLoaded() bool {
	return xinputLoaded
}

// GetState 获取指定手柄的状态
// controllerID: 0-3 (最多支持4个手柄)
func GetState(controllerID int) (*XInputState, error) {
	if !xinputLoaded {
		return nil, ErrXInputNotLoaded
	}

	if controllerID < 0 || controllerID > 3 {
		return nil, errors.New("invalid controller id (must be 0-3)")
	}

	var state XInputState
	ret, _, _ := procGetState.Call(
		uintptr(controllerID),
		uintptr(unsafe.Pointer(&state)),
	)

	// ERROR_DEVICE_NOT_CONNECTED = 1167
	if ret == 1167 {
		return nil, ErrControllerNotConnected
	}

	if ret != 0 {
		return nil, errors.New("xinput get state failed")
	}

	return &state, nil
}
