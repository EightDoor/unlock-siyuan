# SiYuan Fork 定制版

基于 [siyuan-note/siyuan](https://github.com/siyuan-note/siyuan) 的独立维护定制版。所有定制行为已直接集成到源码中。

## 🎯 定制内容

| 功能 | 说明 |
|------|------|
| **VIP 解锁** | 默认解锁所有 VIP 功能（云端同步、S3/WebDAV 等） |
| **默认中文** | 界面语言默认为简体中文 |
| **关闭按钮最小化到托盘** | 关闭窗口时默认最小化到系统托盘 |
| **默认 S3 同步** | 同步提供者默认为 S3，开启冲突文档生成 |
| **关闭自动更新下载** | 默认关闭自动下载更新安装包（可手动开启） |
| **自定义更新源** | 更新检查指向 `EightDoor/unlock-siyuan` GitHub Release |
| **VIP 本地用户** | 默认本地 VIP 用户，无需网络请求 |

## 🐳 Docker 使用

```bash
# 从 Docker Hub 拉取
docker pull eightdoor/siyuan:latest

# 从 GitHub Container Registry 拉取
docker pull ghcr.io/EightDoor/unlock-siyuan:latest

# 运行容器
docker run -d \
  -v /path/to/workspace:/siyuan/workspace \
  -p 6806:6806 \
  eightdoor/siyuan:latest \
  --workspace=/siyuan/workspace \
  --accessAuthCode=your_password
```

## ⚠️ 注意事项

1. **AGPL-3.0 许可证**: 思源笔记使用 AGPL-3.0 许可证，此 Fork 也必须遵守该协议
2. **仅供学习研究**: 此定制版仅供个人学习研究使用
3. **支持官方**: 如果你觉得思源笔记好用，请考虑[支持官方订阅](https://b3log.org/siyuan/en/pricing.html)

## 🔧 关键修改位置

### VIP 解锁
- `kernel/model/conf.go` — `IsSubscriber()` 和 `IsPaidUser()` 始终返回 `true`

### 默认配置
- `kernel/conf/account.go` — 默认不显示头衔和 VIP 标识
- `kernel/conf/appearance.go` — 默认中文，关闭按钮最小化到托盘
- `kernel/conf/sync.go` — 默认 S3 同步，生成冲突文档
- `kernel/conf/system.go` — 默认关闭自动下载更新安装包

### VIP 本地用户
- `kernel/api/mock_vip_user.go` — 本地默认 VIP 用户
- `kernel/api/setting.go` — `getCloudUser` 短路返回本地用户

### GitHub Release 更新
- `kernel/model/updater.go` — GitHub Release 检查与下载
- `kernel/api/system.go` — `downloadUpdate` 接口
- `kernel/api/router.go` — `/api/system/downloadUpdate` 路由
- `app/src/util/processMessage.ts` — 前端更新确认对话框
