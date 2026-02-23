#!/bin/bash
# ============================================================
# SiYuan Fork - Apply Patches Script
# 应用所有补丁到源代码
# ============================================================

set -e

PATCHES_DIR=".patches"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  SiYuan Fork - Apply Patches${NC}"
echo -e "${BLUE}================================================${NC}"

# 检查补丁目录
if [ ! -d "$PATCHES_DIR" ]; then
    echo -e "${RED}Error: Patches directory not found: $PATCHES_DIR${NC}"
    exit 1
fi

# 统计
TOTAL=0
SUCCESS=0
FAILED=0

# 应用每个补丁
for patch in "$PATCHES_DIR"/*.patch; do
    if [ -f "$patch" ]; then
        TOTAL=$((TOTAL + 1))
        patch_name=$(basename "$patch")
        
        echo -e "${BLUE}Applying: $patch_name${NC}"
        
        if git apply --check "$patch" 2>/dev/null; then
            git apply "$patch"
            echo -e "${GREEN}✓ Success: $patch_name${NC}"
            SUCCESS=$((SUCCESS + 1))
        else
            echo -e "${RED}✗ Failed: $patch_name${NC}"
            echo -e "${YELLOW}  This patch may have already been applied or conflicts exist.${NC}"
            FAILED=$((FAILED + 1))
        fi
    fi
done

# 输出统计
echo -e ""
echo -e "${BLUE}================================================${NC}"
echo -e "  Summary"
echo -e "${BLUE}================================================${NC}"
echo -e "  Total patches:   $TOTAL"
echo -e "  ${GREEN}Successful:      $SUCCESS${NC}"
echo -e "  ${RED}Failed:          $FAILED${NC}"
echo -e "${BLUE}================================================${NC}"

if [ $FAILED -gt 0 ]; then
    echo -e "${YELLOW}Some patches failed. Please check conflicts manually.${NC}"
    exit 1
fi

echo -e "${GREEN}All patches applied successfully!${NC}"
