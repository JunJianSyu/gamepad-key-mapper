package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"gamepad-key-mapper/internal/app"
	"gamepad-key-mapper/internal/gamepad"
	"gamepad-key-mapper/internal/keyboard"
)

// ShowMappingForm 显示添加/编辑映射对话框
func ShowMappingForm(parent fyne.Window, appCtrl *app.App, editID *string) {
	// 源按键选择
	sourceButtons := gamepad.AllButtons()
	sourceOptions := make([]string, len(sourceButtons))
	for i, btn := range sourceButtons {
		sourceOptions[i] = btn.String()
	}
	sourceSelect := widget.NewSelect(sourceOptions, nil)
	sourceSelect.PlaceHolder = "选择手柄按键"

	// 目标类型选择
	targetTypeSelect := widget.NewSelect([]string{"键盘按键", "手柄按键"}, nil)
	targetTypeSelect.SetSelected("键盘按键")

	// ===== 键盘目标部分 =====
	targetKeys := keyboard.AllKeys()
	targetKeyOptions := make([]string, len(targetKeys))
	for i, key := range targetKeys {
		targetKeyOptions[i] = key.String()
	}

	selectedKeyTargets := make(map[int]bool)
	keyboardCheckGroup := widget.NewCheckGroup(targetKeyOptions, func(selected []string) {
		selectedKeyTargets = make(map[int]bool)
		for _, sel := range selected {
			for i, opt := range targetKeyOptions {
				if opt == sel {
					selectedKeyTargets[i] = true
					break
				}
			}
		}
	})
	keyboardScroll := container.NewVScroll(keyboardCheckGroup)
	keyboardScroll.SetMinSize(fyne.NewSize(200, 120))

	// 修饰键复选框
	ctrlCheck := widget.NewCheck("Ctrl", nil)
	altCheck := widget.NewCheck("Alt", nil)
	shiftCheck := widget.NewCheck("Shift", nil)
	modifiersBox := container.NewHBox(
		widget.NewLabel("修饰键:"),
		ctrlCheck,
		altCheck,
		shiftCheck,
	)

	keyboardContainer := container.NewVBox(
		widget.NewLabel("目标按键 (键盘) - 可多选"),
		keyboardScroll,
		modifiersBox,
	)

	// ===== 手柄目标部分 =====
	targetButtons := gamepad.AllButtons()
	targetBtnOptions := make([]string, len(targetButtons))
	for i, btn := range targetButtons {
		targetBtnOptions[i] = btn.String()
	}

	selectedBtnTargets := make(map[int]bool)
	gamepadCheckGroup := widget.NewCheckGroup(targetBtnOptions, func(selected []string) {
		selectedBtnTargets = make(map[int]bool)
		for _, sel := range selected {
			for i, opt := range targetBtnOptions {
				if opt == sel {
					selectedBtnTargets[i] = true
					break
				}
			}
		}
	})
	gamepadScroll := container.NewVScroll(gamepadCheckGroup)
	gamepadScroll.SetMinSize(fyne.NewSize(200, 120))

	gamepadContainer := container.NewVBox(
		widget.NewLabel("目标按键 (手柄) - 可多选"),
		gamepadScroll,
		widget.NewLabel("提示: 会触发目标按键的映射规则"),
	)
	gamepadContainer.Hide() // 默认隐藏

	// 目标容器（切换显示）
	targetContainer := container.NewStack(keyboardContainer, gamepadContainer)

	// 目标类型切换逻辑
	targetTypeSelect.OnChanged = func(selected string) {
		if selected == "键盘按键" {
			keyboardContainer.Show()
			gamepadContainer.Hide()
		} else {
			keyboardContainer.Hide()
			gamepadContainer.Show()
		}
		targetContainer.Refresh()
	}

	// 提示标签
	tipLabel := widget.NewLabel("按住源按键时，目标按键也会保持按住状态")
	tipLabel.Wrapping = fyne.TextWrapWord

	// 表单内容
	formContent := container.NewVBox(
		widget.NewLabel("源按键 (手柄)"),
		sourceSelect,
		widget.NewSeparator(),
		widget.NewLabel("目标类型"),
		targetTypeSelect,
		widget.NewSeparator(),
		targetContainer,
		widget.NewSeparator(),
		tipLabel,
	)

	// 创建对话框
	d := dialog.NewCustomConfirm(
		"添加按键映射",
		"确定",
		"取消",
		formContent,
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// 验证源按键
			if sourceSelect.Selected == "" {
				dialog.ShowError(errors.New("请选择源按键"), parent)
				return
			}
			sourceIdx := sourceSelect.SelectedIndex()
			if sourceIdx < 0 {
				return
			}
			sourceKey := sourceButtons[sourceIdx]

			// 检查冲突
			excludeID := ""
			if editID != nil {
				excludeID = *editID
			}
			if appCtrl.HasConflict(sourceKey, excludeID) {
				dialog.ShowError(errors.New("源按键已存在映射，请选择其他按键"), parent)
				return
			}

			if targetTypeSelect.Selected == "键盘按键" {
				// 键盘映射
				if len(selectedKeyTargets) == 0 {
					dialog.ShowError(errors.New("请至少选择一个目标按键"), parent)
					return
				}

				var targets []keyboard.KeyCode
				for idx := range selectedKeyTargets {
					targets = append(targets, targetKeys[idx])
				}

				mods := keyboard.Modifiers{
					Ctrl:  ctrlCheck.Checked,
					Alt:   altCheck.Checked,
					Shift: shiftCheck.Checked,
				}

				_, err := appCtrl.AddRuleMultiKeys(sourceKey, targets, mods)
				if err != nil {
					dialog.ShowError(err, parent)
					return
				}
			} else {
				// 手柄映射
				if len(selectedBtnTargets) == 0 {
					dialog.ShowError(errors.New("请至少选择一个目标按键"), parent)
					return
				}

				var targets []gamepad.Button
				for idx := range selectedBtnTargets {
					targets = append(targets, targetButtons[idx])
				}

				_, err := appCtrl.AddRuleGamepad(sourceKey, targets)
				if err != nil {
					dialog.ShowError(err, parent)
					return
				}
			}
		},
		parent,
	)

	d.Resize(fyne.NewSize(420, 550))
	d.Show()
}
