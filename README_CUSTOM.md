# SiYuan Fork 定制版工具包

基于 [siyuan-note/siyuan](https://github.com/siyuan-note/siyuan) 的定制化版本。

## 🎯 定制内容

| 功能 | 说明 |
|------|------|
| **VIP 解锁** | 默认解锁所有 VIP 功能（云端同步、S3/WebDAV 等） |
| **关闭自动更新下载** | 默认关闭自动下载更新安装包 |
| **Docker 支持** | 自动构建多架构 Docker 镜像 |
| **GitHub Release** | 自动发布版本 |

## 📦 快速开始

### 步骤 1: Fork 仓库

1. 访问 https://github.com/siyuan-note/siyuan
2. 点击右上角 "Fork" 按钮
3. 等待 Fork 完成

### 步骤 2: 克隆你的 Fork

```bash
# 替换 YOUR_USERNAME 为你的 GitHub 用户名
git clone https://github.com/YOUR_USERNAME/siyuan.git
cd siyuan
```

### 步骤 3: 复制工具包文件

将本工具包中的文件复制到你的仓库根目录：

```bash
# 复制补丁目录
cp -r siyuan-fork-toolkit/.patches .

# 复制脚本
cp -r siyuan-fork-toolkit/scripts .

# 复制 GitHub Actions 工作流
mkdir -p .github/workflows
cp siyuan-fork-toolkit/.github/workflows/build-release.yml .github/workflows/
```

### 步骤 4: 应用补丁

```bash
chmod +x scripts/*.sh
./scripts/apply-patches.sh
```

### 步骤 5: 提交更改

```bash
git add .
git commit -m "Apply custom patches: VIP unlock, disable auto-update"
git push origin master
```

### 步骤 6: 配置 GitHub Secrets

进入你的 Fork 仓库 → Settings → Secrets and variables → Actions

添加以下 secrets:
- `DOCKERHUB_TOKEN`: Docker Hub Access Token

### 步骤 7: 创建发布

```bash
# 创建标签并推送
git tag v3.x.x-custom
git push origin v3.x.x-custom
```

GitHub Actions 会自动：
1. 构建 Docker 镜像并推送到 Docker Hub 和 ghcr.io
2. 创建 GitHub Release

## 🔄 同步上游更新

当上游 siyuan-note/siyuan 有新版本时：

```bash
./scripts/sync-upstream.sh
```

脚本会自动：
1. 拉取上游最新代码
2. 合并到当前分支
3. 重新应用补丁

## 🐳 Docker 使用

```bash
# 从 Docker Hub 拉取
docker pull eightdoor/siyuan:latest

# 从 GitHub Container Registry 拉取
docker pull ghcr.io/YOUR_USERNAME/siyuan:latest

# 运行容器
docker run -d \
  -v /path/to/workspace:/siyuan/workspace \
  -p 6806:6806 \
  eightdoor/siyuan:latest \
  --workspace=/siyuan/workspace \
  --accessAuthCode=your_password
```

## 📁 文件结构

```
siyuan/
├── .patches/
│   ├── 001-vip-bypass.patch        # VIP 解锁补丁
│   └── 002-disable-auto-update.patch # 禁用自动更新补丁
├── .github/
│   └── workflows/
│       └── build-release.yml       # CI/CD 工作流
├── scripts/
│   ├── apply-patches.sh            # 应用补丁脚本
│   └── sync-upstream.sh            # 同步上游脚本
└── ... (原有代码)
```

## ⚠️ 注意事项

1. **AGPL-3.0 许可证**: 思源笔记使用 AGPL-3.0 许可证，你的 Fork 也必须遵守该协议
2. **仅供学习研究**: 此定制版仅供个人学习研究使用
3. **支持官方**: 如果你觉得思源笔记好用，请考虑[支持官方订阅](https://b3log.org/siyuan/en/pricing.html)
4. **更新维护**: 每次同步上游后，请检查补丁是否仍然适用

## 🔧 手动修改指南

如果你想手动修改代码，以下是关键位置：

### VIP 解锁
文件: `kernel/model/conf.go`

```go
func IsSubscriber() bool {
    return true  // 修改为始终返回 true
}

func IsPaidUser() bool {
    return true  // 修改为始终返回 true
}
```

### 关闭自动更新下载
文件: `kernel/conf/conf.go`

搜索 `DownloadInstallPkg` 或 `downloadInstallPkg`，将默认值改为 `false`

## 📝 更新日志

### v1.0.0
- 初始版本
- VIP 功能解锁
- CI/CD 工作流
- 同步上游脚本
