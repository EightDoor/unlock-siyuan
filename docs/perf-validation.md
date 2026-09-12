# 性能改动运行时验证指南

本仓库在 Electron 环境下运行，本文档描述对 `app/src/protyle/wysiwyg/transaction.ts`（编辑器渲染合并）与 `app/src/util/fetch.ts`（搜索过期响应处理）两项性能改动的运行时回归验证步骤。

## 一、编辑器渲染合并（transaction.ts）

改动概要：`promiseTransaction` 回调中，多个 `"update"` 操作的渲染从"每个操作独立全量渲染 7 次（processRender/highlightRender/avRender/blockRender + updateEmbed 内部 3 次）"合并为"循环结束后统一渲染 4 次"，并通过新增的 `updateEmbed(deferRender)` 延迟内部渲染。

预期收益：连续输入场景下，主线程长任务时间显著下降；CPU/内存占用减少。

### 验证步骤

1. 准备数据
   - 打开一份包含 5 个以上反链块（包含 `NodeBlockQueryEmbed` / 反链面板）的文档。
   - 关闭其他占用 CPU 的应用，保持环境温度稳定。

2. 获取基线（修改前）
   - 切换到修改前的提交或暂存当前改动（`git stash`）。
   - 按步骤 3 测量并记录。

3. 测量方法
   - 桌面端打开"开发者工具"（DevTools）→ Performance 面板 → 录制。
   - 连续输入约 30 个字符（触发多个 update 事务），停止录制。
   - 记录：
     - "Scripting" 总时长
     - "Rendering" 总时长
     - 主线程长任务（>50ms）数量
     - 同一帧中调用 `update`/`processRender`/`highlightRender`/`avRender`/`blockRender` 的次数

4. 应用本次改动（`git stash pop`），重复步骤 3。

5. 对比指标
   - 连续输入 30 字符，预期 Rendering 时长下降 50% 以上（保守目标）。
   - 长任务数量预期下降 30% 以上。
   - 渲染调用次数从"每操作 4 次"降为"整批 4 次"。

### 行为一致性检查

- 撤销/重做：连续输入若干字符后按 `Ctrl+Z` 多次、再 `Ctrl+Shift+Z` 多次，文档内容应完全可逆，无残留差异。
- 反链面板：触发输入后反链面板中的引用块应同步更新，无空内容或旧内容残留。
- 嵌入块：嵌入引用的块（NodeBlockQueryEmbed）应在最后一批 update 后统一刷新，无闪烁或缺失。
- 折叠/展开：折叠标题后再展开，反链面板与嵌入块渲染正确。
- 保存：触发事务后保存到磁盘，重新打开应一致。

### 风险点提示

`updateEmbed` 内部可能依赖"已渲染"的中间状态进行测量或位置计算。本次改动通过 `deferRender` 标志跳过中间渲染，理论上这些测量在最终统一渲染后仍可正常获取；若在嵌入块批处理场景发现位置或高亮异常，应回退该批处理路径，仅保留 4 次顶层渲染合并。

## 二、搜索过期响应处理（fetch.ts）

改动概要：在 `fetch.ts` 的两处 `includes(url)` 列表中加入 `/api/search/fullTextSearchAssetContent`，与 `/api/search/fullTextSearchBlock`、`/api/search/searchRefBlock` 一致，对资产搜索的快速输入场景启用 reqId 过期丢弃。

预期收益：用户在资产搜索快速输入时，旧请求的响应到达后不会覆盖新结果，列表只显示最终关键词的匹配。

### 静态验证（无 GUI）

```bash
# 确认 fetch.ts 中两处 includes 都包含该端点
grep -n "fullTextSearchAssetContent" app/src/util/fetch.ts

# 确认搜索调用方使用同一 URL
grep -n "fullTextSearchAssetContent" app/src/search/assets.ts app/src/search/util.ts
```

预期输出：fetch.ts 出现 2 次（请求发送侧 includes + 响应处理侧 includes），调用方出现 1 次。

### 运行时验证（需 Electron 桌面或浏览器入口）

1. 打开全局搜索面板（`Ctrl+P` 或顶栏搜索按钮）。
2. 切换到"资源"标签页。
3. 在搜索框中快速连续输入 `test12345`（无停顿或停顿 < 50ms）。
4. 观察结果列表：
   - 修改前：可能出现中间态结果（如 `test`/`test1`/`test12`），最终才稳定为 `test12345`。
   - 修改后：列表应直接显示 `test12345` 的结果，中间态不会出现或仅一闪而过（取决于网络）。
5. 在 DevTools Network 面板观察：旧请求（`reqId` 较小）的响应到达时，前端不应对 DOM 应用该响应的结果。

### 行为一致性检查

- 慢网络（DevTools → Network → Slow 3G）下，结果应只显示最终关键词的匹配。
- 正常网络下，与修改前无可见差异（响应顺序正常时无 reqId 过期）。
- 输入中文/特殊字符（粘贴整段文本），结果应只显示完整文本的匹配。

## 三、可选的静态回归脚本

```bash
# 运行以下命令确认 fetch.ts 的 reqId 列表覆盖完整
node scripts/check-reqid-coverage.js
```

`scripts/check-reqid-coverage.js` 读取 `app/src/util/fetch.ts`，比对两处 `includes` 列表是否一致，并扫描 `app/src/search/**` 中所有 `/api/search/*` 调用是否都在列表中。预期输出 `OK` 或列出缺失项。

## 四、回退方案

若运行时发现任一改动引入回归，可按以下顺序定位并回退：

1. 编辑器渲染合并：通过 `git diff app/src/protyle/wysiwyg/transaction.ts` 查看 `needsUpdateRender`/`needsEmbedRender` 标记和 `updateEmbed(deferRender)` 调用；回退方式为删除标记逻辑，恢复循环内逐个调用 4 次渲染。
2. 搜索 reqId 覆盖：回退 `fetch.ts` 的两处 `includes` 改动即可；该改动无破坏性，仅失去过期响应丢弃。

## 五、未覆盖的运行时检查项

由于本仓库的运行时验证必须在 Electron GUI 中进行，以下项目需在本地手动验证，不在本指南强制范围内：

- 真实数据集下的编辑器 p95 输入延迟
- 1000+ 文档时的搜索 p95 响应
- 启动到可交互的时间（受 AI 移除影响，预期下降）

建议使用 Chrome DevTools Performance / Lighthouse 记录前后对比，输出截图或报告归档。
