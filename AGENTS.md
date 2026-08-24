# AGENTS.md

## 项目概述
wol-panel：面向香橙派等 ARM64 开发板（Armbian）的 NAS 远程开关机控制面板。Go 后端通过发起网络魔术包开机、SSH 执行 `sudo poweroff` 关机；Vue3 前端构建产物嵌入 Go 二进制，单文件部署。

## 技术栈
- 后端：Go 1.25（module 名 `wol-panel`），标准库 `net/http` + `log/slog`；依赖 ttlcache（防抖）、lumberjack（日志滚动）
- 前端：Vue3 + Vite 6 + vue-router + vue-i18n + TypeScript（`frontend/`）
- 部署：systemd 用户级服务（`wol-panel.service`）；WOL 由 Go 原生实现无需外部工具，仅需配置到 NAS 的 SSH 免密

## 目录结构（核心）
| 路径 | 职责 |
|---|---|
| `main.go` | 入口：加载配置、初始化日志、注册路由、`//go:embed dist/*` 托管前端 |
| `config/` | 读取运行目录下 `config.json`（模板见 `config.template.json`） |
| `logger/` | 三通道日志（控制台/文件/错误文件，级别独立配置） |
| `antishake/` | 后端防抖锁（`enable_anti_shake` 开启时用内存 TTL 缓存） |
| `nas/` | WOL（Go 原生魔术包）与 SSH 关机实现 |
| `handler/` | HTTP 接口：`POST /wol/api/wol`、`POST /wol/api/shutdown`、`GET /wol/api/version` |
| `version/` | 版本号/架构/构建时间，由构建脚本经 ldflags 注入 |
| `bin/sh`、`bin/ps` | 各平台构建脚本（sh/PowerShell 双版本）；`build-all.sh` 仅供 CI |
| `.github/workflows/release.yml` | 合并 master / 打 tag 后自动构建 6 平台并发布 |
| `.docs/` | 需求与调研文档，按 `日期-主题` 目录组织；`about.md` 项目背景、`notice.md` 注意事项。只保存在本地，不提交到代码库。 |

## 常用命令
开发环境为 **Windows，优先使用 PowerShell**，勿直接执行 bash 命令。

```powershell
# 前端（在 frontend/ 目录）
npm install
npm run dev        # 本地开发
npm run build      # 产物输出到仓库根目录 dist/

# 一键构建（dist 缺失时自动先构建前端）
.\bin\ps\build-linux-aarch64.ps1 [版本号]   # 其他平台见 bin\ps\

# 手动编译（前置：已执行 npm run build，dist/ 必须存在）
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="arm64"; go build -o build\wol-panel .

# 运行与验证
.\build\wol-panel.exe            # 读取当前目录 config.json，监听 server_port
.\build\wol-panel.exe version    # 输出版本号/架构/构建时间
go vet .\...
```

## 约定与注意事项
- 前端产物 `dist/` 已 gitignore，**改前端后必须重新 `npm run build` 再编译 Go**，否则二进制内嵌资源是旧的。
- 所有路由与静态资源带 `/wol/` 前缀；英文界面路由为 `/wol/en`（文案在 `frontend/src/i18n/`，路由在 `frontend/src/router/`）。
- 版本号、架构、构建时间只通过构建脚本的 ldflags 注入，默认版本从 `version/version.go` 读取，勿硬编码。
- `config.json` 含敏感信息，禁止提交；仓库只保留 `config.template.json`。
- 新需求先在 `.docs/<yyyymmdd-主题>/` 下建需求文档，并引用 `.docs/about.md`、`.docs/notice.md`。
