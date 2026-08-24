# WOL Panel — NAS 远程控制面板

基于 Go + Vue3 的轻量级 NAS 远程开关机控制面板，适用于香橙派等 ARM64 开发板。

## 功能

- **WOL 开机**：点击按钮发送 Wake-on-LAN 魔术包，远程唤醒 NAS
- **SSH 关机**：点击按钮通过 SSH 远程执行 `sudo poweroff` 安全关机
- **双层防抖**：前端按钮锁定 + 后端内存防抖，防止重复提交
- **单文件部署**：前端资源嵌入 Go 二进制，零依赖交付
- **版本管理**：内置版本号、架构、构建时间，支持 `./wol-panel version` 查看
- **国际化**：支持中英文切换

## 构建与编译

### 构建（推荐）

```bash
# Bash (Linux/macOS/Git Bash)
./bin/sh/build-linux-arm64.sh [版本号]

# PowerShell (Windows)
.\bin\ps\build-linux-arm64.ps1 [版本号]
```

不指定版本号时自动从 `version/version.go` 读取。产物统一输出到 `build/wol-panel`（Windows 为 `build/wol-panel.exe`）。

可用脚本：

| 平台 | Bash | PowerShell |
|------|------|------------|
| Linux x64 | `bin/sh/build-linux-x64.sh` | `bin/ps/build-linux-x64.ps1` |
| Linux arm64 | `bin/sh/build-linux-arm64.sh` | `bin/ps/build-linux-arm64.ps1` |
| Windows x64 | `bin/sh/build-windows-x64.sh` | `bin/ps/build-windows-x64.ps1` |
| Windows arm | `bin/sh/build-windows-arm.sh` | `bin/ps/build-windows-arm.ps1` |
| macOS Apple Silicon | `bin/sh/build-macos-apple-silicon.sh` | `bin/ps/build-macos-apple-silicon.ps1` |
| macOS Intel | `bin/sh/build-macos-intel.sh` | `bin/ps/build-macos-intel.ps1` |

### 手动交叉编译

```bash
VERSION=0.0.1
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS="-s -w \
  -X wol-panel/version.Version=${VERSION} \
  -X wol-panel/version.Arch=arm64 \
  -X wol-panel/version.BuildTime=${BUILD_TIME}"

# 先构建前端
cd frontend && npm run build && cd ..

# 再编译 Go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o build/wol-panel .
```

### 查看版本

```bash
./wol-panel version
# 输出：wol-panel 0.0.1 arm64 2026-07-04 12:00:00
```

## config.json 全字段说明

将 `config.template.json` 复制为 `config.json` 并修改：

```bash
cp config.template.json config.json
```

| 字段 | 类型 | 说明 |
|---|---|---|
| `server_port` | string | Web 服务监听端口，默认 `8080` |
| `stdout_log_level` | string | 控制台日志级别：`Off` / `Debug` / `Info` / `Warn` / `Error` |
| `file_log_level` | string | 磁盘文件日志级别：`Off` / `Debug` / `Info` / `Warn` / `Error` |
| `error_log_level` | string | 错误日志文件级别：`Off` / `Debug` / `Info` / `Warn` / `Error` |
| `enable_anti_shake` | bool | 是否开启后端内存防抖锁。`false` 则跳过防抖 |
| `nas_ip` | string | NAS 局域网 IP，用于 SSH 关机 |
| `nas_user` | string | NAS SSH 登录账号 |
| `nas_mac` | string | NAS MAC 地址，用于 WOL 唤醒 |
| `wol_broadcast` | string | WOL 魔术包广播目标，默认 `255.255.255.255:9` |

**日志级别说明**：
- `Off`：完全关闭该输出通道
- `Debug`：最详细，输出所有调试信息
- `Info`：常规信息（默认控制台级别）
- `Warn`：仅警告和错误（推荐文件级别，减少 SD 卡写入）
- `Error`：仅错误

三个通道不可共用同一级别。控制台可更详细，文件建议更高级别以保护 SD 卡寿命。

## 部署流程（完整步骤）

> 以下所有操作均以**普通用户**身份在 Armbian 开发板上执行，无需 root 权限。
> 部署路径为 `~/.local/share/wol-panel/`，无需修改系统目录权限。

### 1. 上传文件

将编译好的二进制和配置文件传到开发板：

```bash
# 在开发机上执行（替换 <板子IP> 和 <用户名>）
scp build/wol-panel config.json <用户名>@<板子IP>:~/wol-panel_tmp/
```

### 2. 部署到用户目录

```bash
# SSH 登录开发板后执行
mkdir -p ~/.local/share/wol-panel
mv ~/wol-panel_tmp/wol-panel ~/.local/share/wol-panel/
mv ~/wol-panel_tmp/config.json ~/.local/share/wol-panel/
chmod +x ~/.local/share/wol-panel/wol-panel
rmdir ~/wol-panel_tmp
```

### 3. 安装用户级 systemd 服务

```bash
# 创建用户级服务目录
mkdir -p ~/.config/systemd/user/

# 复制 service 文件
cp wol-panel.service ~/.config/systemd/user/

# 重载并启用
systemctl --user daemon-reload
systemctl --user enable wol-panel

# 确保登出后服务仍运行（重要！）
loginctl enable-linger $(whoami)
```

### 4. 启动与验证

```bash
# 启动服务
systemctl --user start wol-panel

# 查看状态
systemctl --user status wol-panel

# 查看实时日志
journalctl --user -u wol-panel -f
```

浏览器访问 `http://<开发板IP>:8080/wol/`，看到控制面板即部署成功。

### 日常操作

```bash
systemctl --user start wol-panel      # 启动
systemctl --user stop wol-panel       # 停止
systemctl --user restart wol-panel    # 重启
systemctl --user status wol-panel     # 状态
journalctl --user -u wol-panel -f     # 实时日志
```

### 更新版本

```bash
# 1. 上传新二进制到开发板
scp build/wol-panel <用户名>@<板子IP>:~/

# 2. SSH 登录后替换并重启
cp ~/wol-panel ~/.local/share/wol-panel/wol-panel
chmod +x ~/.local/share/wol-panel/wol-panel
systemctl --user restart wol-panel
```

### TF 卡存储（可选）

如果 TF 卡挂载在 `/home`，`~/.local/share/` 天然就在 TF 卡上，无需额外操作。

如果需要把数据放在 TF 卡的其他位置，可以用软链接：

```bash
# 例：TF 卡挂载在 /mnt/tf
mkdir -p /mnt/tf/wol-panel
ln -sf /mnt/tf/wol-panel ~/.local/share/wol-panel
```

## SSH 免密配置

关机功能需要 SSH 免密登录，请提前配置：

```bash
# 在开发板上生成密钥（如尚未生成）
ssh-keygen -t ed25519

# 将公钥复制到 NAS
ssh-copy-id <nas_user>@<nas_ip>

# 验证免密登录
ssh <nas_user>@<nas_ip> "echo ok"
```

## 项目结构

```
wol-panel/
├── main.go                # 程序入口：配置加载、日志初始化、HTTP 服务
├── config/config.go       # 配置读取独立包
├── logger/logger.go       # 三渠道日志封装
├── antishake/antishake.go # 内存防抖锁（TTL 缓存）
├── nas/nas.go             # NAS 操作（WOL、SSH关机）
├── handler/handler.go     # HTTP 接口处理器
├── version/version.go     # 版本信息（ldflags 注入）
├── frontend/              # Vue3 前端源码
│   ├── src/
│   │   ├── App.vue
│   │   ├── main.ts
│   │   ├── api/index.ts       # API 请求封装
│   │   ├── i18n/              # 国际化（中英文）
│   │   ├── router/            # 路由配置
│   │   ├── views/             # 页面组件
│   │   └── utils/debounce.ts  # 通用防抖工具
│   └── ...
├── bin/                   # 构建脚本
│   ├── sh/                # Bash 版（6 个平台各一个脚本）
│   └── ps/                # PowerShell 版（6 个平台各一个脚本）
├── config.template.json   # 配置模板
├── wol-panel.service      # systemd 用户级服务配置
└── README.md
```
