#!/bin/bash
# ============================================================
# SiYuan Fork - Sync Upstream Script
# 同步上游仓库更新并自动应用补丁
# ============================================================

set -e

# 配置
UPSTREAM_REPO="https://github.com/siyuan-note/siyuan.git"
UPSTREAM_BRANCH="master"
PATCHES_DIR=".patches"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  SiYuan Fork - Sync Upstream${NC}"
echo -e "${BLUE}================================================${NC}"

# 检查是否在 git 仓库中
if [ ! -d ".git" ]; then
    echo -e "${RED}Error: Not a git repository. Please run this script from the root of your fork.${NC}"
    exit 1
fi

# 检查 upstream remote 是否存在
if ! git remote | grep -q "^upstream$"; then
    echo -e "${YELLOW}Adding upstream remote...${NC}"
    git remote add upstream "$UPSTREAM_REPO"
fi

# 获取上游更新
echo -e "${GREEN}📥 Fetching upstream changes...${NC}"
git fetch upstream "$UPSTREAM_BRANCH"

# 检查是否有更新
UPSTREAM_HASH=$(git rev-parse upstream/$UPSTREAM_BRANCH)
LOCAL_HASH=$(git rev-parse HEAD)

if [ "$UPSTREAM_HASH" = "$LOCAL_HASH" ]; then
    echo -e "${YELLOW}Already up to date with upstream.${NC}"
    exit 0
fi

# 显示更新信息
echo -e "${GREEN}📊 Upstream has new commits:${NC}"
git log --oneline HEAD..upstream/$UPSTREAM_BRANCH | head -10

# 询问是否继续
read -p "$(echo -e ${YELLOW}Continue with merge? [y/N]: ${NC})" -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}Aborted.${NC}"
    exit 1
fi

# 合并上游
echo -e "${GREEN}🔄 Merging upstream/$UPSTREAM_BRANCH...${NC}"
git merge upstream/$UPSTREAM_BRANCH --no-edit

# 应用补丁
echo -e "${GREEN}🩹 Applying patches...${NC}"
if [ -d "$PATCHES_DIR" ] && [ "$(ls -A $PATCHES_DIR/*.patch 2>/dev/null)" ]; then
    for patch in "$PATCHES_DIR"/*.patch; do
        echo -e "${BLUE}Applying: $(basename $patch)${NC}"
        if git apply --check "$patch" 2>/dev/null; then
            git apply "$patch"
            echo -e "${GREEN}✓ Applied: $(basename $patch)${NC}"
        else
            echo -e "${RED}✗ Conflict in: $(basename $patch)${NC}"
            echo -e "${YELLOW}Please resolve conflicts manually.${NC}"
        fi
    done
else
    echo -e "${YELLOW}No patches found in $PATCHES_DIR${NC}"
fi

echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}✅ Sync complete!${NC}"
echo -e "${GREEN}================================================${NC}"
echo -e "${YELLOW}Please review changes and commit if everything looks good.${NC}"
