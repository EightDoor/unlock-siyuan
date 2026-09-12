# 🔓 SiYuan Unlock Edition

> 独立维护的 SiYuan 定制版本

## ✨ 定制内容

| 功能 | 说明 |
|------|------|
| **VIP 解锁** | 默认解锁所有 VIP 功能（云端同步、S3/WebDAV 等） |
| **关闭自动更新** | 默认关闭自动下载更新安装包 |
| **Docker 支持** | 自动构建多架构 Docker 镜像 (amd64/arm64) |

## 🔄 发布流程

1. 从上游 [siyuan-note/siyuan](https://github.com/siyuan-note/siyuan) 合并代码
2. 在本仓库开发、测试
3. 更新 `app/package.json` 中的 `version` 字段
4. 推送版本标签 `vX.Y.Z`（必须与 `app/package.json` 的 version 严格匹配）
5. GitHub Actions 自动构建与发布

### Tag 命名规范

```
vX.Y.Z
```

- 必须与 `app/package.json` 的 `version` 严格匹配
- 例如：`app/package.json` 中 version 为 `3.5.8`，则标签为 `v3.5.8`

### 手动触发

每个构建工作流也支持通过 `workflow_dispatch` 手动触发。

## 🐳 Docker 使用

```bash
# Docker Hub
docker pull 851708184/siyuan:latest

# GitHub Container Registry
docker pull ghcr.io/eightdoor/unlock-siyuan:latest

# 运行容器
docker run -d \
  -v /path/to/workspace:/siyuan/workspace \
  -p 6806:6806 \
  851708184/siyuan:latest \
  --workspace=/siyuan/workspace \
  --accessAuthCode=your_password
```

## 📥 下载

- [GitHub Releases](https://github.com/EightDoor/unlock-siyuan/releases)
- [Docker Hub](https://hub.docker.com/r/851708184/siyuan)

## ⚙️ GitHub Actions 配置

### 📍 配置位置

进入仓库 **Settings** → **Secrets and variables** → **Actions**

### 🔐 Secrets（必需）

| Secret | 说明 | 如何获取 |
|--------|------|----------|
| `DOCKER_USERNAME` | Docker Hub 用户名 | 你的 Docker Hub 账号用户名 |
| `DOCKERHUB_TOKEN` | Docker Hub Access Token | [创建 Token](https://hub.docker.com/settings/security) |

### 📝 Variables（可选）

| Variable | 说明 | 默认值 |
|----------|------|--------|
| `IMAGE_NAME` | Docker 镜像名称 | `siyuan` |

### 🚀 工作流说明

| 工作流 | 触发条件 | 功能 |
|--------|---------|------|
| **release-tag** | 推送 `v*` 标签 / 手动触发 | 创建 Release 并触发所有平台构建 |
| **desktop-release** | release-tag 触发 / 手动触发 | 构建桌面端（Windows/macOS/Linux） |
| **release-docker** | release-tag 触发 / 手动触发 | 构建并推送 Docker 镜像 |
| **release-android** | release-tag 触发 / 手动触发 | 构建 Android APK |
| **release-ios** | release-tag 触发 / 手动触发 | 构建 iOS IPA |
| **dev-check** | PR / push 到 main 或 dev | 编译检查（kernel + app） |

### 📋 发布新版本

```bash
# 1. 更新版本号
# 编辑 app/package.json，将 version 改为目标版本

# 2. 提交并推送
git add app/package.json
git commit -m "bump version to vX.Y.Z"
git push

# 3. 打标签并推送
git tag vX.Y.Z
git push origin vX.Y.Z
```

## 📁 目录结构

```
├── app/                    # 前端代码
├── kernel/                 # 后端 Go 代码
├── patches/                # 补丁文件（历史）
├── scripts/                # 维护脚本
├── .github/workflows/      # CI/CD 工作流
├── Dockerfile              # Docker 构建文件
└── README.md               # 本文件
```

## ⚠️ 免责声明

1. **AGPL-3.0 许可证**: 本项目遵循 AGPL-3.0 许可证
2. **仅供学习研究**: 此定制版仅供个人学习研究使用
3. **支持官方**: 如果你觉得思源笔记好用，请考虑[支持官方订阅](https://b3log.org/siyuan/en/pricing.html)

## 📜 原始项目

- 官网: https://b3log.org/siyuan/
- 源码: https://github.com/siyuan-note/siyuan
- 许可证: AGPL-3.0

---

## 📝 配置清单

```yaml
# GitHub Secrets（必需）
DOCKER_USERNAME: "your_docker_hub_username"
DOCKERHUB_TOKEN: "dckr_pat_xxxxxxxxxxxx"

# GitHub Variables（可选，有默认值）
IMAGE_NAME: "siyuan"           # 默认: siyuan
```
