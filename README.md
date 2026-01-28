# 游戏手柄按键映射工具 (GamepadKeyMapper)

一个使用 Go 语言开发的 Windows 桌面应用程序，用于将游戏手柄（Xbox/XInput兼容）的按键映射到键盘快捷键或其他手柄按键。

## 功能特性

### 核心功能
- **键盘映射**: 将手柄按键映射到键盘按键（支持组合键如 Ctrl+Alt+F1）
- **手柄映射**: 将手柄按键映射到其他手柄按键（触发对应的映射规则）
- **多键映射**: 单个手柄按键可映射到多个目标键
- **按键保持**: 手柄按键按住时，目标键也保持按住状态

### 手柄支持
- Xbox 360 / Xbox One / Xbox Series X|S 手柄
- Xbox 精英版手柄（支持背部拨片 P1-P4）
- 所有 XInput 兼容手柄
- 摇杆方向映射（将摇杆方向当作按键使用）

### 界面功能
- 图形化界面，易于配置
- 系统托盘支持，可最小化运行
- 右键托盘菜单快速控制启动/停止
- 配置自动保存和加载

## 系统要求

- Windows 10/11 (x64)
- Xbox 手柄或 XInput 兼容手柄

## 支持的按键

### 标准按键
| 按键 | 说明 |
|------|------|
| A, B, X, Y | 主按键 |
| LB, RB | 肩键 (Bumper) |
| LT, RT | 扳机 (Trigger) |
| Menu, View | 功能键 (Start/Back) |
| Share | 分享按钮 (Xbox Series) |
| LS, RS | 摇杆按下 |
| D-Pad | 方向键 (上/下/左/右) |

### 精英版按键
| 按键 | 说明 |
|------|------|
| P1, P2, P3, P4 | 背部拨片 |

### 摇杆方向
| 按键 | 说明 |
|------|------|
| 左摇杆 ↑↓←→ | 左摇杆四方向 |
| 右摇杆 ↑↓←→ | 右摇杆四方向 |

## 使用方法

### 1. 添加键盘映射

1. 点击「添加映射」按钮
2. 选择目标类型：「键盘按键」
3. 选择源按键（手柄按键）
4. 选择目标按键（可多选键盘按键）
5. 可选：勾选修饰键（Ctrl/Alt/Shift）
6. 点击「确定」保存

### 2. 添加手柄映射

1. 点击「添加映射」按钮
2. 选择目标类型：「手柄按键」
3. 选择源按键（触发按键）
4. 选择目标按键（要触发的手柄按键，可多选）
5. 点击「确定」保存

> 手柄映射会触发目标按键对应的映射规则，实现连锁映射效果

### 3. 启动/停止映射

- 点击「启动」按钮开始监听手柄输入
- 点击「停止」按钮暂停映射
- 停止时会自动释放所有按住的键

### 4. 系统托盘

- 关闭窗口会最小化到系统托盘
- 右键托盘图标可快速启动/停止映射
- 双击托盘图标恢复窗口
- 选择「退出」完全关闭程序

## 使用示例

### 示例1: 简单快捷键
```
RB → Ctrl+C (复制)
LB → Ctrl+V (粘贴)
```

### 示例2: 多键组合
```
RT → X+Y+B (同时按下三个键的效果)
```

### 示例3: 连锁映射
```
配置:
  RT → 手柄 X+B      (RT触发X和B)
  X  → 键盘 F1       (X映射到F1)
  B  → 键盘 F2       (B映射到F2)

效果:
  按下RT → 同时触发F1和F2
```

### 示例4: 摇杆方向映射
```
左摇杆↑ → W (前进)
左摇杆↓ → S (后退)
左摇杆← → A (左移)
左摇杆→ → D (右移)
```

## 构建

### 环境要求

- Go 1.21+
- Windows 系统（用于实际运行）

### 编译命令

```bash
# 安装依赖
go mod tidy

# 编译（隐藏控制台窗口）
go build -ldflags "-H=windowsgui" -o GamepadKeyMapper.exe

# 或者普通编译（显示控制台，便于调试）
go build -o GamepadKeyMapper.exe
```

### 跨平台编译

```bash
# 从 macOS/Linux 编译 Windows 版本
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-H=windowsgui" -o GamepadKeyMapper.exe
```

## 配置文件

配置文件自动保存在用户配置目录：
- Windows: `%APPDATA%\GamepadKeyMapper\gamepad-key-mapper.json`

配置格式示例：
```json
{
  "rules": [
    {
      "id": "rule_1234567890",
      "source_key": 512,
      "target_type": 0,
      "target_keys": [112],
      "modifiers": {"ctrl": true, "alt": false, "shift": false, "win": false},
      "enabled": true
    }
  ],
  "minimize_to_tray": true
}
```

## 技术栈

- **语言**: Go 1.21+
- **GUI框架**: [Fyne](https://fyne.io/) v2.4+
- **手柄输入**: Windows XInput API
- **键盘模拟**: Windows SendInput API

## 已知问题

1. **管理员权限**: 某些游戏可能需要以管理员权限运行才能接收模拟的键盘输入
2. **安全软件**: 部分安全软件可能会拦截键盘模拟功能，需要添加白名单
3. **精英版拨片**: 标准XInput API对拨片支持有限，可能需要Xbox Accessories应用配置

## 许可证

MIT License
