# SiYuan Fork - Apply Patches Script (PowerShell)
# 用于在 Windows 环境下应用补丁

$PATCHES_DIR = ".patches"

Write-Host "================================================" -ForegroundColor Blue
Write-Host "  SiYuan Fork - Apply Patches (Windows)" -ForegroundColor Blue
Write-Host "================================================" -ForegroundColor Blue

# 检查补丁目录
if (-not (Test-Path $PATCHES_DIR)) {
    Write-Host "Error: Patches directory not found: $PATCHES_DIR" -ForegroundColor Red
    exit 1
}

# 统计
$TOTAL = 0
$SUCCESS = 0
$FAILED = 0

# 检查 git 是否可用
$gitAvailable = $null -ne (Get-Command git -ErrorAction SilentlyContinue)

if (-not $gitAvailable) {
    Write-Host "Warning: Git not found. Please apply patches manually." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Manual steps:" -ForegroundColor Cyan
    Write-Host "1. Open kernel/model/conf.go"
    Write-Host "2. Find IsSubscriber() function"
    Write-Host "3. Change it to return true"
    Write-Host "4. Find IsPaidUser() function"
    Write-Host "5. Change it to return true"
    exit 0
}

# 应用每个补丁
Get-ChildItem -Path $PATCHES_DIR -Filter "*.patch" | ForEach-Object {
    $patch = $_.FullName
    $patchName = $_.Name
    
    $TOTAL++
    
    Write-Host "Applying: $patchName" -ForegroundColor Blue
    
    # 使用 git apply
    $checkResult = git apply --check $patch 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        git apply $patch
        Write-Host "  Success: $patchName" -ForegroundColor Green
        $SUCCESS++
    } else {
        Write-Host "  Failed: $patchName" -ForegroundColor Red
        Write-Host "  This patch may have already been applied or conflicts exist." -ForegroundColor Yellow
        $FAILED++
    }
}

# 输出统计
Write-Host ""
Write-Host "================================================" -ForegroundColor Blue
Write-Host "  Summary" -ForegroundColor Blue
Write-Host "================================================" -ForegroundColor Blue
Write-Host "  Total patches:   $TOTAL"
Write-Host "  Successful:      $SUCCESS" -ForegroundColor Green
Write-Host "  Failed:          $FAILED" -ForegroundColor Red
Write-Host "================================================" -ForegroundColor Blue

if ($FAILED -gt 0) {
    Write-Host "Some patches failed. Please check conflicts manually." -ForegroundColor Yellow
    exit 1
}

Write-Host "All patches applied successfully!" -ForegroundColor Green
