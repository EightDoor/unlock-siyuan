# Custom Update Source Patch (Optional)
# 自定义更新源补丁（可选）
#
# 如果你想让更新检查指向你自己的仓库，请修改以下文件：
# 
# 文件: kernel/model/updater.go
# 
# 修改下载 URL（约第 151-165 行）：
#
# 原始代码：
#   b3logURL := "https://release.b3log.org/siyuan/" + pkg
#   liuyunURL := "https://release.liuyun.io/siyuan/" + pkg
#   githubURL := "https://github.com/siyuan-note/siyuan/releases/download/v" + ver + "/" + pkg
#
# 修改为（替换 YOUR_USERNAME 为你的 GitHub 用户名）：
#   githubURL := "https://github.com/YOUR_USERNAME/siyuan/releases/download/v" + ver + "/" + pkg
#   // 移除 b3log 和 liuyun URL，只使用你自己的 GitHub
#
# ---------------------------------------------------
#
# 文件: kernel/util/rhy.go
#
# 版本检查 API（约第 64 行）：
# 
# 原始代码：
#   resp, err := request.SetContext(ctx).SetSuccessResult(&cachedRhyResult).Get(GetCloudServer() + "/apis/siyuan/version?ver=" + Ver)
#
# 修改为（使用你自己的版本 API）：
#   // 方案1: 完全禁用版本检查，直接返回空结果
#   // 方案2: 使用你自己的版本检查服务器
#
# ---------------------------------------------------
#
# 文件: kernel/util/cloud.go
#
# 云服务器地址（第 67-84 行）：
#
# 如果你想完全自主控制，可以修改这些常量指向你自己的服务器

# 此补丁仅为说明文档，实际修改需要根据你的需求手动调整
