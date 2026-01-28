package main

import (
	"gamepad-key-mapper/internal/app"
	"gamepad-key-mapper/internal/ui"
)

func main() {
	// 创建应用实例
	application := app.New()

	// 创建并运行UI
	ui.Run(application)
}
