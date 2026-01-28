package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"gamepad-key-mapper/internal/app"
)

// MappingList 映射列表组件
type MappingList struct {
	appCtrl   *app.App
	parent    fyne.Window
	container *fyne.Container
	list      *widget.List
}

// NewMappingList 创建映射列表
func NewMappingList(appCtrl *app.App, parent fyne.Window) *MappingList {
	ml := &MappingList{
		appCtrl: appCtrl,
		parent:  parent,
	}

	ml.list = widget.NewList(
		func() int {
			return len(ml.appCtrl.GetRules())
		},
		func() fyne.CanvasObject {
			return ml.createListItem()
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			ml.updateListItem(id, item)
		},
	)

	ml.container = container.NewStack(ml.list)

	return ml
}

// Container 返回容器
func (ml *MappingList) Container() *fyne.Container {
	return ml.container
}

// Refresh 刷新列表
func (ml *MappingList) Refresh() {
	ml.list.Refresh()
}

// createListItem 创建列表项模板
func (ml *MappingList) createListItem() fyne.CanvasObject {
	label := widget.NewLabel("映射规则")
	deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
	deleteBtn.Importance = widget.LowImportance

	return container.NewBorder(nil, nil, nil, deleteBtn, label)
}

// updateListItem 更新列表项内容
func (ml *MappingList) updateListItem(id widget.ListItemID, item fyne.CanvasObject) {
	rules := ml.appCtrl.GetRules()
	if id >= len(rules) {
		return
	}

	rule := rules[id]
	border := item.(*fyne.Container)

	// 更新标签
	label := border.Objects[0].(*widget.Label)
	label.SetText(rule.String())

	// 更新删除按钮
	deleteBtn := border.Objects[1].(*widget.Button)
	ruleID := rule.ID // 捕获当前规则ID
	deleteBtn.OnTapped = func() {
		ml.confirmDelete(ruleID, rule.String())
	}
}

// confirmDelete 确认删除对话框
func (ml *MappingList) confirmDelete(ruleID string, ruleDesc string) {
	dialog.ShowConfirm(
		"确认删除",
		"确定要删除映射规则 \""+ruleDesc+"\" 吗？",
		func(confirmed bool) {
			if confirmed {
				ml.appCtrl.RemoveRule(ruleID)
			}
		},
		ml.parent,
	)
}
