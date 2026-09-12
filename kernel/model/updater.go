// SiYuan - Refactor your thinking
// Copyright (c) 2020-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package model

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/88250/gulu"
	"github.com/imroc/req/v3"
	"github.com/siyuan-note/logging"
	"github.com/siyuan-note/siyuan/kernel/util"
	"golang.org/x/mod/semver"
)

func execNewVerInstallPkg(newVerInstallPkgPath string) {
	logging.LogInfof("installing the new version [%s]", newVerInstallPkgPath)
	var cmd *exec.Cmd
	if gulu.OS.IsWindows() {
		cmd = exec.Command(newVerInstallPkgPath)
	} else if gulu.OS.IsDarwin() {
		exec.Command("chmod", "+x", newVerInstallPkgPath).CombinedOutput()
		cmd = exec.Command("open", newVerInstallPkgPath)
	} else {
		logging.LogErrorf("unsupported platform for auto-installing package")
		return
	}
	gulu.CmdAttr(cmd)
	cmdErr := cmd.Run()
	if nil != cmdErr {
		logging.LogErrorf("exec install new version failed: %s", cmdErr)
		return
	}
}

func getNewVerInstallPkgPath() string {
	if skipNewVerInstallPkg() {
		return ""
	}

	downloadPkgURLs, checksum, err := getUpdatePkg()
	if err != nil {
		return ""
	}

	pkg := path.Base(downloadPkgURLs[0])
	pkgPath := filepath.Join(util.TempDir, "install", pkg)
	localChecksum, _ := sha256Hash(pkgPath)
	if checksum != localChecksum {
		return ""
	}
	return pkgPath
}

var checkDownloadInstallPkgLock = sync.Mutex{}

// TryLockCheckDownloadInstallPkg 尝试获取下载安装包互斥锁，用于手动更新入口避免与自动检查并发下载同一安装包。
func TryLockCheckDownloadInstallPkg() bool {
	return checkDownloadInstallPkgLock.TryLock()
}

// UnlockCheckDownloadInstallPkg 释放下载安装包互斥锁。
func UnlockCheckDownloadInstallPkg() {
	checkDownloadInstallPkgLock.Unlock()
}

func checkDownloadInstallPkg() {
	defer logging.Recover()

	if skipNewVerInstallPkg() {
		return
	}

	if !checkDownloadInstallPkgLock.TryLock() {
		return
	}
	defer checkDownloadInstallPkgLock.Unlock()

	downloadPkgURLs, checksum, err := getUpdatePkg()
	if err != nil {
		return
	}

	existingPkgPath := getNewVerInstallPkgPath()
	if "" != existingPkgPath {
		// 存在经过 sha256Hash 检查的安装包
		util.PushUpdateMsg("update-pkg-ready", Conf.Language(62), 15*1000)
		return
	}

	util.PushUpdateMsg("update-pkg-downloading", Conf.Language(103), 1000*7)
	success := false
	for _, downloadPkgURL := range downloadPkgURLs {
		err = DownloadInstallPkg(downloadPkgURL, checksum)
		if err == nil {
			success = true
			break
		}
	}
	if success {
		util.PushUpdateMsg("update-pkg-ready", Conf.Language(62), 15*1000)
	} else {
		util.PushUpdateMsg("update-pkg-downloading", Conf.Language(104), 7000)
	}
}

func getUpdatePkg() (downloadPkgURLs []string, checksum string, err error) {
	defer logging.Recover()

	// [FORK-MOD] 使用 GitHub Release API 检查更新，替代原 b3log/liuyun 源
	release, fetchErr := fetchLatestRelease(context.TODO())
	if fetchErr != nil {
		err = fetchErr
		logging.LogErrorf("fetch latest release failed: %s", err)
		return
	}
	if release.Draft || release.Prerelease {
		err = errors.New("no stable release")
		return
	}

	ver, parseErr := parseTagVersion(release.TagName)
	if parseErr != nil {
		err = parseErr
		return
	}
	if isVersionUpToDate(ver) {
		err = fmt.Errorf("version is up to date")
		return
	}

	asset, assetErr := selectReleaseAsset(release.Assets, ver, runtime.GOOS, runtime.GOARCH)
	if assetErr != nil {
		err = assetErr
		return
	}
	checksum, err = parseSha256Digest(asset.Digest)
	if err != nil {
		return
	}
	downloadPkgURLs = []string{asset.BrowserDownloadURL}
	return
}

// DownloadInstallPkg 下载安装包并校验 checksum。
func DownloadInstallPkg(pkgURL, checksum string) (err error) {
	if "" == pkgURL || "" == checksum {
		return
	}

	pkg := path.Base(pkgURL)
	savePath := filepath.Join(util.TempDir, "install", pkg)
	if gulu.File.IsExist(savePath) {
		localChecksum, _ := sha256Hash(savePath)
		if localChecksum == checksum {
			return
		}
	}

	err = os.MkdirAll(filepath.Join(util.TempDir, "install"), 0755)
	if err != nil {
		logging.LogErrorf("create temp install dir failed: %s", err)
		return
	}

	logging.LogInfof("downloading install package [%s]", pkgURL)
	client := req.C().SetTLSHandshakeTimeout(7 * time.Second).SetTimeout(10 * time.Minute).DisableInsecureSkipVerify()
	callback := func(info req.DownloadInfo) {
		progress := fmt.Sprintf("%.2f%%", float64(info.DownloadedSize)/float64(info.Response.ContentLength)*100.0)
		// logging.LogDebugf("downloading install package [%s %s]", pkgURL, progress)
		util.PushStatusBar(fmt.Sprintf(Conf.Language(133), progress))
	}
	_, err = client.R().SetOutputFile(savePath).SetDownloadCallbackWithInterval(callback, 1*time.Second).Get(pkgURL)
	if err != nil {
		logging.LogErrorf("download install package [%s] failed: %s", pkgURL, err)
		return
	}

	localChecksum, _ := sha256Hash(savePath)
	if checksum != localChecksum {
		logging.LogErrorf("verify checksum failed, download install package [%s] checksum [%s] not equal to downloaded [%s] checksum [%s]", pkgURL, checksum, savePath, localChecksum)
		return
	}
	logging.LogInfof("downloaded install package [%s] to [%s]", pkgURL, savePath)
	util.PushStatusBar(Conf.Language(62))
	return
}

func sha256Hash(filename string) (ret string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha256.New()
	reader := bufio.NewReader(file)
	buf := make([]byte, 1024*1024*4)
	for {
		switch n, readErr := reader.Read(buf); readErr {
		case nil:
			hash.Write(buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", hash.Sum(nil)), nil
		default:
			return "", err
		}
	}
}

type Announcement struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Region int    `json:"region"`
}

func getAnnouncements() (ret []*Announcement) {
	result, err := util.GetRhyResult(context.TODO(), false)
	if err != nil {
		logging.LogErrorf("get announcement failed: %s", err)
		return
	}

	if nil == result["announcement"] {
		return
	}

	announcements := result["announcement"].([]interface{})
	for _, announcement := range announcements {
		ann := announcement.(map[string]interface{})
		ret = append(ret, &Announcement{
			Id:     ann["id"].(string),
			Title:  ann["title"].(string),
			URL:    ann["url"].(string),
			Region: int(ann["region"].(float64)),
		})
	}
	return
}

func CheckUpdate(showMsg bool) {
	if !showMsg {
		return
	}

	if Conf.System.IsMicrosoftStore {
		return
	}

	// [FORK-MOD] 使用 GitHub Release API 检查更新
	release, err := fetchLatestRelease(context.TODO())
	if err != nil {
		logging.LogErrorf("check update failed: %s", err)
		util.PushUpdateMsg("update-notify", Conf.Language(10), 3000)
		return
	}
	if release.Draft || release.Prerelease {
		return
	}

	ver, parseErr := parseTagVersion(release.TagName)
	if parseErr != nil {
		logging.LogErrorf("check update invalid tag: %s", parseErr)
		return
	}
	if isVersionUpToDate(ver) {
		util.PushUpdateMsg("update-notify", Conf.Language(10), 3000)
		return
	}

	link := "<a href=\"" + release.HTMLURL + "\">" + release.TagName + "</a>"
	util.PushUpdateMsg("update-notify", fmt.Sprintf(Conf.Language(9), link), 15000)

	if showMsg {
		// [FORK-MOD] 通知前端弹出更新确认对话框
		util.BroadcastByType("main", "update-confirm", 0, "", map[string]interface{}{
			"version": release.TagName,
			"url":     release.HTMLURL,
		})
	}

	go func() {
		defer logging.Recover()
		checkDownloadInstallPkg()
	}()
}

func isVersionUpToDate(releaseVer string) bool {
	return semver.Compare("v"+releaseVer, "v"+util.Ver) <= 0
}

// skipInstallPkgPlatformCached 缓存平台相关判断，-1 未初始化，0 表示不跳过，1 表示跳过
var skipInstallPkgPlatformCached = -1

func skipNewVerInstallPkg() bool {
	if skipInstallPkgPlatformCached == -1 {
		skipInstallPkgPlatformCached = 0
		if !gulu.OS.IsWindows() && !gulu.OS.IsDarwin() {
			skipInstallPkgPlatformCached = 1
		} else if util.ISMicrosoftStore || util.ContainerStd != util.Container {
			skipInstallPkgPlatformCached = 1
		} else if gulu.OS.IsWindows() {
			plat := strings.ToLower(Conf.System.OSPlatform)
			// Windows 7, 8 and Server 2012 are no longer supported https://github.com/siyuan-note/siyuan/issues/7347
			if strings.Contains(plat, " 7 ") || strings.Contains(plat, " 8 ") || strings.Contains(plat, "2012") {
				skipInstallPkgPlatformCached = 1
			}
		}
	}

	if skipInstallPkgPlatformCached == 1 || !Conf.System.DownloadInstallPkg {
		return true
	}
	return false
}

// [FORK-MOD] 以下为 GitHub Release 更新相关函数，替代原 b3log/liuyun 更新源

// githubReleaseRepo 是 fork 更新使用的 GitHub Release 仓库
const githubReleaseRepo = "EightDoor/unlock-siyuan"

// githubReleaseAPI 是 GitHub latest release API 端点
const githubReleaseAPI = "https://api.github.com/repos/" + githubReleaseRepo + "/releases/latest"

type releasePayload struct {
	TagName    string         `json:"tag_name"`
	HTMLURL    string         `json:"html_url"`
	Draft      bool           `json:"draft"`
	Prerelease bool           `json:"prerelease"`
	Assets     []releaseAsset `json:"assets"`
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"`
}

func parseTagVersion(tag string) (string, error) {
	if "" == tag {
		return "", errors.New("empty tag")
	}
	v := strings.TrimPrefix(tag, "v")
	full := "v" + v
	if !semver.IsValid(full) || "" != semver.Prerelease(full) || "" != semver.Build(full) {
		return "", fmt.Errorf("invalid semver tag %q", tag)
	}
	return v, nil
}

func parseSha256Digest(digest string) (string, error) {
	const prefix = "sha256:"
	if !strings.HasPrefix(digest, prefix) {
		return "", fmt.Errorf("unsupported digest %q", digest)
	}
	hex := strings.TrimPrefix(digest, prefix)
	if len(hex) != 64 {
		return "", fmt.Errorf("invalid sha256 length for digest %q", digest)
	}
	return strings.ToLower(hex), nil
}

func platformAssetNames(version, goos, goarch string) []string {
	switch goos {
	case "windows":
		if goarch == "arm64" {
			return []string{"siyuan-v" + version + "-win-arm64.exe"}
		}
		return []string{"siyuan-v" + version + "-win.exe"}
	case "darwin":
		if goarch == "arm64" {
			return []string{"siyuan-v" + version + "-mac-arm64.dmg"}
		}
		return []string{"siyuan-v" + version + "-mac.dmg"}
	case "linux":
		if goarch == "arm64" {
			return []string{"siyuan-v" + version + "-linux-arm64.tar.gz"}
		}
		return []string{
			"siyuan-v" + version + "-linux.tar.gz",
			"siyuan-v" + version + "-linux.AppImage",
		}
	}
	return nil
}

func selectReleaseAsset(assets []releaseAsset, version, goos, goarch string) (releaseAsset, error) {
	names := platformAssetNames(version, goos, goarch)
	if len(names) == 0 {
		return releaseAsset{}, fmt.Errorf("unsupported platform %s/%s", goos, goarch)
	}
	index := make(map[string]releaseAsset, len(assets))
	for _, a := range assets {
		index[a.Name] = a
	}
	for _, n := range names {
		if a, ok := index[n]; ok {
			return a, nil
		}
	}
	return releaseAsset{}, fmt.Errorf("no asset matches platform %s/%s in version %s", goos, goarch, version)
}

func fetchLatestRelease(ctx context.Context) (releasePayload, error) {
	var payload releasePayload
	client := req.C().SetTLSHandshakeTimeout(7 * time.Second).SetTimeout(15 * time.Second)
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/vnd.github+json").
		SetSuccessResult(&payload).
		Get(githubReleaseAPI)
	if err != nil {
		return payload, err
	}
	if !resp.IsSuccessState() {
		return payload, fmt.Errorf("github release api status %d", resp.StatusCode)
	}
	if 0 == len(payload.TagName) {
		if err := json.Unmarshal(resp.Bytes(), &payload); err != nil {
			return payload, fmt.Errorf("decode release payload: %w", err)
		}
	}
	return payload, nil
}

// DownloadUpdatePkg 由前端确认更新后调用，忽略 DownloadInstallPkg 开关，成功后推送 update-pkg-ready。
func DownloadUpdatePkg() (downloadPkgURLs []string, checksum string, err error) {
	defer logging.Recover()

	if skipUpdatePkgPlatform() {
		err = errors.New("current platform does not support auto install")
		return
	}

	release, fetchErr := fetchLatestRelease(context.TODO())
	if fetchErr != nil {
		err = fetchErr
		logging.LogErrorf("download update: fetch latest release failed: %s", err)
		return
	}
	if release.Draft || release.Prerelease {
		err = errors.New("no stable release")
		return
	}

	ver, parseErr := parseTagVersion(release.TagName)
	if parseErr != nil {
		err = parseErr
		return
	}
	if isVersionUpToDate(ver) {
		err = errors.New("version is up to date")
		return
	}

	asset, assetErr := selectReleaseAsset(release.Assets, ver, runtime.GOOS, runtime.GOARCH)
	if assetErr != nil {
		err = assetErr
		return
	}
	checksum, err = parseSha256Digest(asset.Digest)
	if err != nil {
		return
	}
	downloadPkgURLs = []string{asset.BrowserDownloadURL}
	return
}

// skipUpdatePkgPlatform 仅判断平台是否支持自动安装，不读取 DownloadInstallPkg 设置。
func skipUpdatePkgPlatform() bool {
	if !gulu.OS.IsWindows() && !gulu.OS.IsDarwin() {
		return true
	}
	if util.ISMicrosoftStore || util.ContainerStd != util.Container {
		return true
	}
	if gulu.OS.IsWindows() {
		plat := strings.ToLower(Conf.System.OSPlatform)
		if strings.Contains(plat, " 7 ") || strings.Contains(plat, " 8 ") || strings.Contains(plat, "2012") {
			return true
		}
	}
	return false
}
