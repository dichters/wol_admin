# WOL Admin — NAS 远程控制面板

基于 Go + Vue3 的轻量级 NAS 远程开关机控制面板，适用于香橙派等 ARM64 开发板。

## 功能

- **WOL 开机**：点击按钮发送 Wake-on-LAN 魔术包，远程唤醒 NAS
- **SSH 关机**：点击按钮通过 SSH 远程执行 `sudo poweroff` 安全关机
- **双层防抖**：前端按钮锁定 + 后端 Redis/内存防抖，防止重复提交
- **单文件部署**：前端资源嵌入 Go 二进制，零依赖交付

## 交叉编译

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/wol_admin main.go
```

## Armbian Redis 安装配置步骤

```bash
# 安装 Redis
sudo apt update
sudo apt install redis-server -y

# 启动并设置开机自启
sudo systemctl enable redis-server
sudo systemctl start redis-server

# 验证运行
redis-cli ping
# 应返回 PONG
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
| `enable_anti_shake` | bool | 是否开启后端 Redis 防抖锁。`false` 则跳过 Redis |
| `redis.ip` | string | Redis 地址（仅 enable_anti_shake=true 时生效） |
| `redis.port` | string | Redis 端口 |
| `redis.password` | string | Redis 密码，空字符串表示无密码 |
| `nas_ip` | string | NAS 局域网 IP，用于 SSH 关机 |
| `nas_user` | string | NAS SSH 登录账号 |
| `nas_mac` | string | NAS MAC 地址，用于 WOL 唤醒 |

**日志级别说明**：
- `Off`：完全关闭该输出通道
- `Debug`：最详细，输出所有调试信息
- `Info`：常规信息（默认控制台级别）
- `Warn`：仅警告和错误（推荐文件级别，减少 SD 卡写入）
- `Error`：仅错误

三个通道不可共用同一级别。控制台可更详细，文件建议更高级别以保护 SD 卡寿命。

## 部署流程（完整步骤）

> 以下所有操作均以**普通用户**身份在 Armbian 开发板上执行，无需 root 权限。
> 部署路径为 `~/.local/share/wol_admin/`，无需修改系统目录权限。

### 1. 上传文件

将编译好的二进制和配置文件传到开发板：

```bash
# 在开发机上执行（替换 <板子IP> 和 <用户名>）
scp build/wol_admin config.json <用户名>@<板子IP>:~/wol_admin_tmp/
```

### 2. 部署到用户目录

```bash
# SSH 登录开发板后执行
mkdir -p ~/.local/share/wol_admin
mv ~/wol_admin_tmp/wol_admin ~/.local/share/wol_admin/
mv ~/wol_admin_tmp/config.json ~/.local/share/wol_admin/
chmod +x ~/.local/share/wol_admin/wol_admin
rmdir ~/wol_admin_tmp
```

### 3. 安装用户级 systemd 服务

```bash
# 创建用户级服务目录
mkdir -p ~/.config/systemd/user/

# 复制 service 文件
cp wol_admin.service ~/.config/systemd/user/

# 重载并启用
systemctl --user daemon-reload
systemctl --user enable wol_admin

# 确保登出后服务仍运行（重要！）
loginctl enable-linger $(whoami)
```

### 4. 启动与验证

```bash
# 启动服务
systemctl --user start wol_admin

# 查看状态
systemctl --user status wol_admin

# 查看实时日志
journalctl --user -u wol_admin -f
```

浏览器访问 `http://<开发板IP>:8080`，看到控制面板即部署成功。

### 日常操作

```bash
systemctl --user start wol_admin      # 启动
systemctl --user stop wol_admin       # 停止
systemctl --user restart wol_admin    # 重启
systemctl --user status wol_admin     # 状态
journalctl --user -u wol_admin -f     # 实时日志
```

### 更新版本

```bash
# 1. 上传新二进制到开发板
scp build/wol_admin <用户名>@<板子IP>:~/

# 2. SSH 登录后替换并重启
cp ~/wol_admin ~/.local/share/wol_admin/wol_admin
chmod +x ~/.local/share/wol_admin/wol_admin
systemctl --user restart wol_admin
```

### TF 卡存储（可选）

如果 TF 卡挂载在 `/home`，`~/.local/share/` 天然就在 TF 卡上，无需额外操作。

如果需要把数据放在 TF 卡的其他位置，可以用软链接：

```bash
# 例：TF 卡挂载在 /mnt/tf
mkdir -p /mnt/tf/wol_admin
ln -sf /mnt/tf/wol_admin ~/.local/share/wol_admin
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

## WOL 依赖

```bash
# 在开发板上安装 wakeonlan 工具
sudo apt install wakeonlan -y
```

## 项目结构

```
wol_admin/
├── main.go                # 程序入口：配置加载、日志初始化、HTTP 服务
├── config/config.go       # 配置读取独立包
├── logger/logger.go       # 三渠道日志封装
├── antishake/antishake.go # Redis/内存防抖锁
├── nas/nas.go             # NAS 操作（WOL、SSH关机）
├── handler/handler.go     # HTTP 接口处理器
├── frontend/              # Vue3 前端源码
│   ├── src/
│   │   ├── App.vue
│   │   ├── main.ts
│   │   ├── api/index.ts       # API 请求封装
│   │   └── utils/debounce.ts  # 通用防抖工具
│   └── ...
├── config.template.json   # 配置模板
├── wol_admin.service      # systemd 用户级服务配置
└── README.md
```
