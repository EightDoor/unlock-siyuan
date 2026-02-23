# 🔓 SiYuan Unlock Edition

> 基于 [siyuan-note/siyuan](https://github.com/siyuan-note/siyuan) 的定制版本

## ✨ 定制内容

| 功能 | 说明 |
|------|------|
| **VIP 解锁** | 默认解锁所有 VIP 功能（云端同步、S3/WebDAV 等） |
| **关闭自动更新** | 默认关闭自动下载更新安装包 |
| **Docker 支持** | 自动构建多架构 Docker 镜像 (amd64/arm64) |

## 🐳 Docker 使用

```bash
# 拉取镜像 (替换 <DOCKER_USERNAME> 为你的 Docker Hub 用户名)
docker pull <DOCKER_USERNAME>/siyuan:latest

# 或从 GitHub Container Registry 拉取
docker pull ghcr.io/eightdoor/unlock-siyuan:latest

# 运行容器
docker run -d \
  -v /path/to/workspace:/siyuan/workspace \
  -p 6806:6806 \
  <DOCKER_USERNAME>/siyuan:latest \
  --workspace=/siyuan/workspace \
  --accessAuthCode=your_password
```

## 📥 下载

- [GitHub Releases](https://github.com/EightDoor/unlock-siyuan/releases)
- [Docker Hub](https://hub.docker.com/r/<DOCKER_USERNAME>/siyuan)

## 🔄 同步上游

当上游 siyuan-note/siyuan 有新版本时：

```bash
./scripts/sync-upstream.sh
```

## ⚙️ GitHub Actions 配置

### 📍 配置位置

进入仓库 **Settings** → **Secrets and variables** → **Actions**

### 🔐 Secrets（必需）

| Secret | 说明 | 如何获取 |
|--------|------|----------|
| `DOCKERHUB_TOKEN` | Docker Hub Access Token | [创建 Token](https://hub.docker.com/settings/security) |

### 📝 Variables（可选）

| Variable | 说明 | 默认值 |
|----------|------|--------|
| `DOCKER_USERNAME` | Docker Hub 用户名 | `eightdoor` |
| `IMAGE_NAME` | Docker 镜像名称 | `siyuan` |

---

### 🔧 配置步骤

#### 1. 创建 Docker Hub Token

1. 登录 [Docker Hub](https://hub.docker.com/)
2. 点击右上角头像 → **Account Settings**
3. 左侧菜单 → **Security**
4. 点击 **New Access Token**
5. 设置：
   - Name: `github-actions`
   - Permissions: `Read, Write, Delete`
6. 点击 **Generate** 并**复制 Token**（只显示一次）

#### 2. 添加 GitHub Secret

1. 进入仓库 **Settings** → **Secrets and variables** → **Actions**
2. 点击 **New repository secret**
3. 填写：
   - Name: `DOCKERHUB_TOKEN`
   - Secret: 粘贴刚才复制的 Token
4. 点击 **Add secret**

#### 3. 添加 GitHub Variable（可选）

1. 在同一页面点击 **Variables** 标签
2. 点击 **New repository variable**
3. 填写：
   - Name: `DOCKER_USERNAME`
   - Value: 你的 Docker Hub 用户名
4. 点击 **Add variable**

---

### 🚀 工作流说明

| 工作流 | 触发条件 | 功能 |
|--------|---------|------|
| **Build and Release** | Tag 推送 / 手动触发 | 构建 Docker 镜像 + 创建 Release |
| **Sync Upstream** | 每7天自动 / 手动触发 | 同步上游代码 + 应用补丁 |

### 📋 手动触发构建

1. 进入 **Actions** 页面
2. 选择 **Build and Release**
3. 点击 **Run workflow**
4. 配置选项：

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `push_docker` | 推送到 Docker Hub | ✅ true |
| `push_ghcr` | 推送到 GitHub Container Registry | ✅ true |
| `build_platforms` | 构建平台 | linux/amd64,linux/arm64 |
| `create_release` | 创建 GitHub Release | ✅ true |
| `docker_username` | Docker Hub 用户名（覆盖变量） | 空（使用变量） |

---

## 📁 目录结构

```
├── .patches/              # 补丁文件
│   ├── 001-vip-bypass.patch
│   └── 002-disable-auto-update.patch
├── scripts/               # 维护脚本
│   ├── apply-patches.sh
│   ├── apply-patches.ps1
│   └── sync-upstream.sh
├── .github/workflows/     # CI/CD
│   ├── build-release.yml
│   └── sync-upstream.yml
├── README.md              # 本文件
└── ...                    # SiYuan 源码
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
DOCKERHUB_TOKEN: "dckr_pat_xxxxxxxxxxxx"

# GitHub Variables（可选，有默认值）
DOCKER_USERNAME: "851708184"   # 默认: eightdoor
IMAGE_NAME: "siyuan"           # 默认: siyuan
```

### 你需要配置：

| 类型 | 名称 | 你的值 |
|------|------|--------|
| Secret | `DOCKERHUB_TOKEN` | 你的 Docker Hub Token |
| Variable | `DOCKER_USERNAME` | `851708184` |
