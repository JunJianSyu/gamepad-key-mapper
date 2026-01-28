package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	"gamepad-key-mapper/internal/app"
)

// Tray 系统托盘管理
type Tray struct {
	app     fyne.App
	window  fyne.Window
	appCtrl *app.App
	menu    *fyne.Menu
	
	startItem *fyne.MenuItem
	stopItem  *fyne.MenuItem
}

// NewTray 创建系统托盘
func NewTray(a fyne.App, w fyne.Window, appCtrl *app.App) *Tray {
	t := &Tray{
		app:     a,
		window:  w,
		appCtrl: appCtrl,
	}
	
	t.setup()
	return t
}

// setup 设置托盘
func (t *Tray) setup() {
	desk, ok := t.app.(desktop.App)
	if !ok {
		// 不支持桌面特性
		return
	}

	// 创建菜单项
	t.startItem = fyne.NewMenuItem("启动", t.onStart)
	t.stopItem = fyne.NewMenuItem("停止", t.onStop)
	t.stopItem.Disabled = true

	separator := fyne.NewMenuItemSeparator()
	showItem := fyne.NewMenuItem("显示窗口", t.onShow)
	quitItem := fyne.NewMenuItem("退出", t.onQuit)

	// 创建菜单
	t.menu = fyne.NewMenu("GamepadKeyMapper",
		t.startItem,
		t.stopItem,
		separator,
		showItem,
		separator,
		quitItem,
	)

	// 设置托盘菜单
	desk.SetSystemTrayMenu(t.menu)

	// 设置状态变更回调以同步菜单状态
	t.appCtrl.SetOnStateChange(func(state app.State) {
		t.updateMenuState(state)
	})
}

// updateMenuState 更新菜单状态
func (t *Tray) updateMenuState(state app.State) {
	switch state {
	case app.StateRunning:
		t.startItem.Disabled = true
		t.stopItem.Disabled = false
	case app.StateStopped:
		t.startItem.Disabled = false
		t.stopItem.Disabled = true
	}
	t.menu.Refresh()
}

// onStart 启动映射
func (t *Tray) onStart() {
	t.appCtrl.Start()
}

// onStop 停止映射
func (t *Tray) onStop() {
	t.appCtrl.Stop()
}

// onShow 显示窗口
func (t *Tray) onShow() {
	t.window.Show()
}

// onQuit 退出程序
func (t *Tray) onQuit() {
	t.appCtrl.Stop()
	t.app.Quit()
}

// SetupWindowClose 设置窗口关闭行为（最小化到托盘）
func (t *Tray) SetupWindowClose() {
	t.window.SetCloseIntercept(func() {
		t.window.Hide()
	})
}
