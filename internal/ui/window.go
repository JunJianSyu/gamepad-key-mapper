package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	appPkg "gamepad-key-mapper/internal/app"
)

// MainWindow 主窗口
type MainWindow struct {
	app         fyne.App
	window      fyne.Window
	appCtrl     *appPkg.App
	statusLabel *widget.Label
	startBtn    *widget.Button
	stopBtn     *widget.Button
	mappingList *MappingList
	tray        *Tray
}

// Run 运行应用
func Run(appCtrl *appPkg.App) {
	a := app.New()
	a.SetIcon(resourceIcon128Png) // 设置应用图标
	
	w := a.NewWindow("游戏手柄按键映射工具")
	w.SetIcon(resourceIcon128Png) // 设置窗口图标
	w.Resize(fyne.NewSize(500, 400))

	// 加载配置
	appCtrl.LoadConfig()

	mw := &MainWindow{
		app:     a,
		window:  w,
		appCtrl: appCtrl,
	}

	mw.setup()
	mw.setupCallbacks()

	// 设置系统托盘
	mw.tray = NewTray(a, w, appCtrl)
	mw.tray.SetupWindowClose()

	w.ShowAndRun()
}

// setup 设置界面
func (mw *MainWindow) setup() {
	// 状态栏
	mw.statusLabel = widget.NewLabel("状态: 已停止")
	mw.statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 控制按钮
	mw.startBtn = widget.NewButtonWithIcon("启动", theme.MediaPlayIcon(), mw.onStart)
	mw.stopBtn = widget.NewButtonWithIcon("停止", theme.MediaStopIcon(), mw.onStop)
	mw.stopBtn.Disable()

	controlBar := container.NewHBox(
		mw.statusLabel,
		widget.NewSeparator(),
		mw.startBtn,
		mw.stopBtn,
	)

	// 映射列表
	mw.mappingList = NewMappingList(mw.appCtrl, mw.window)

	// 添加按钮
	addBtn := widget.NewButtonWithIcon("添加映射", theme.ContentAddIcon(), mw.onAddMapping)

	// 布局
	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("按键映射列表"),
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			container.NewHBox(addBtn),
		),
		nil, nil,
		mw.mappingList.Container(),
	)

	mainLayout := container.NewBorder(
		container.NewVBox(controlBar, widget.NewSeparator()),
		nil, nil, nil,
		content,
	)

	mw.window.SetContent(mainLayout)
}

// setupCallbacks 设置回调
func (mw *MainWindow) setupCallbacks() {
	mw.appCtrl.SetOnStateChange(func(state appPkg.State) {
		mw.updateStatus(state)
	})

	mw.appCtrl.SetOnRulesChange(func() {
		mw.mappingList.Refresh()
	})

	mw.appCtrl.SetOnError(func(err error) {
		dialog.ShowError(err, mw.window)
	})
}

// updateStatus 更新状态显示
func (mw *MainWindow) updateStatus(state appPkg.State) {
	switch state {
	case appPkg.StateRunning:
		mw.statusLabel.SetText("状态: 运行中")
		mw.startBtn.Disable()
		mw.stopBtn.Enable()
	case appPkg.StateStopped:
		mw.statusLabel.SetText("状态: 已停止")
		mw.startBtn.Enable()
		mw.stopBtn.Disable()
	}
}

// onStart 启动按钮点击
func (mw *MainWindow) onStart() {
	if err := mw.appCtrl.Start(); err != nil {
		dialog.ShowError(err, mw.window)
	}
}

// onStop 停止按钮点击
func (mw *MainWindow) onStop() {
	mw.appCtrl.Stop()
}

// onAddMapping 添加映射按钮点击
func (mw *MainWindow) onAddMapping() {
	ShowMappingForm(mw.window, mw.appCtrl, nil)
}
