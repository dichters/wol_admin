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

## 服务启停命令（用户级 systemd 服务）

```bash
# 复制二进制和配置到部署目录
mkdir -p /opt/wol_admin
cp wol_admin /opt/wol_admin/
cp config.json /opt/wol_admin/

# 安装用户级 systemd 服务
mkdir -p ~/.config/systemd/user/
cp wol_admin.service ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable wol_admin

# 确保登出后服务仍运行
loginctl enable-linger $(whoami)

# 启动 / 停止 / 重启
systemctl --user start wol_admin
systemctl --user stop wol_admin
systemctl --user restart wol_admin

# 查看状态和日志
systemctl --user status wol_admin
journalctl --user -u wol_admin -f
```

## 部署流程

1. 在开发机上安装 Go 1.25.8+、Node.js 18+
2. 构建前端：`cd frontend && npm install && npm run build`
3. 编译后端：`CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o wol_admin main.go`
4. 将 `wol_admin`、`config.json` 传至 Armbian 开发板
5. 按上述命令安装用户级 systemd 服务
6. 浏览器访问 `http://<开发板IP>:8080`

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
├── main.go              # 程序入口：配置加载、日志初始化、HTTP 服务
├── config/config.go     # 配置读取独立包
├── logger/logger.go     # 三渠道日志封装
├── antishake/antishake.go # Redis/内存防抖锁
├── nas/nas.go           # NAS 操作（WOL、SSH关机）
├── handler/handler.go   # HTTP 接口处理器
├── frontend/            # Vue3 前端源码
│   ├── src/
│   │   ├── App.vue
│   │   ├── main.ts
│   │   ├── api/index.ts       # API 请求封装
│   │   └── utils/debounce.ts  # 通用防抖工具
│   └── ...
├── config.template.json # 配置模板
├── wol_admin.service    # systemd 用户级服务配置
└── README.md
```
