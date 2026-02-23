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
# 拉取镜像
docker pull eightdoor/siyuan:latest

# 或从 GitHub Container Registry 拉取
docker pull ghcr.io/eightdoor/unlock-siyuan:latest

# 运行容器
docker run -d \
  -v /path/to/workspace:/siyuan/workspace \
  -p 6806:6806 \
  eightdoor/siyuan:latest \
  --workspace=/siyuan/workspace \
  --accessAuthCode=your_password
```

## 📥 下载

- [GitHub Releases](https://github.com/EightDoor/unlock-siyuan/releases)
- [Docker Hub](https://hub.docker.com/r/eightdoor/siyuan)

## 🔄 同步上游

当上游 siyuan-note/siyuan 有新版本时：

```bash
./scripts/sync-upstream.sh
```

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
├── README.md              # 本文件
├── README_ORIGINAL.md     # 原始 README
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
