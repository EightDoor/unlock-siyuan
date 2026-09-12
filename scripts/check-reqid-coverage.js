#!/usr/bin/env node
// check-reqid-coverage.js
// 静态检查：
//   1. app/src/util/fetch.ts 中两处 [.includes(url) 数组内容完全一致；
//   2. app/src/search/** 中引用的所有 /api/search/* 端点都已包含在 fetch.ts 的列表里。
//
// 不检查 fetch.ts 列表中是否有 graph / recent 等非搜索端点，这些由其他模块使用。
//
// 用法：node scripts/check-reqid-coverage.js

const fs = require("fs");
const path = require("path");

const ROOT = path.resolve(__dirname, "..");
const FETCH = path.join(ROOT, "app/src/util/fetch.ts");
const SEARCH_DIR = path.join(ROOT, "app/src/search");

function readText(p) {
  return fs.readFileSync(p, "utf8");
}

function listTsFiles(dir) {
  const out = [];
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      out.push(...listTsFiles(p));
    } else if (entry.isFile() && /\.(ts|tsx)$/.test(entry.name)) {
      out.push(p);
    }
  }
  return out;
}

// 字符串感知的数组扫描：找到每个 [ ... ] 且后续紧跟 .includes(
function extractIncludesArrays(src) {
  const arrays = [];
  let i = 0;
  while (i < src.length) {
    const start = src.indexOf("[", i);
    if (start === -1) break;
    let depth = 1;
    let j = start + 1;
    while (j < src.length && depth > 0) {
      const c = src[j];
      if (c === '"' || c === "'" || c === "`") {
        const quote = c;
        j++;
        while (j < src.length && src[j] !== quote) {
          if (src[j] === "\\") j += 2;
          else j++;
        }
        j++;
        continue;
      }
      if (c === "[") depth++;
      else if (c === "]") depth--;
      j++;
    }
    if (depth !== 0) break;
    const after = src.slice(j, j + 20);
    if (/^\s*\.includes\s*\(/.test(after)) {
      const arrayText = src.slice(start, j);
      const items = [];
      for (const sm of arrayText.matchAll(/['"`]([^'"`]+)['"`]/g)) {
        items.push(sm[1]);
      }
      arrays.push(items);
    }
    i = j + 1;
  }
  return arrays;
}

function extractSearchUrls(src) {
  const urls = new Set();
  for (const m of src.matchAll(/['"`](\/api\/search\/[A-Za-z0-9_]+)['"`]/g)) {
    urls.add(m[1]);
  }
  return urls;
}

function pickReqIdLists(lists) {
  // reqId 列表的特征：同时含 /api/search/、/api/graph/ 或 /api/block/ 中的至少一个
  return lists.filter((l) =>
    l.some(
      (u) =>
        u.startsWith("/api/search/") ||
        u.startsWith("/api/graph/") ||
        u.startsWith("/api/block/")
    )
  );
}

function main() {
  if (!fs.existsSync(FETCH)) {
    console.error("找不到", FETCH);
    process.exit(1);
  }
  const fetchSrc = readText(FETCH);
  const allLists = extractIncludesArrays(fetchSrc);
  const reqIdLists = pickReqIdLists(allLists);

  if (reqIdLists.length < 2) {
    console.error(
      "fetch.ts 中未找到至少两处 reqId includes 列表，实际找到",
      reqIdLists.length
    );
    process.exit(1);
  }
  if (reqIdLists.length > 2) {
    console.error(
      "fetch.ts 中存在超过两处形似的 reqId 列表，请人工核对：",
      reqIdLists.map((l) => l.length)
    );
    process.exit(1);
  }

  const sorted = reqIdLists.map((l) => [...l].sort());
  if (sorted[0].length !== sorted[1].length ||
      !sorted[0].every((u, i) => u === sorted[1][i])) {
    console.error("两处 reqId includes 列表内容不一致：");
    console.error("  [0]", reqIdLists[0]);
    console.error("  [1]", reqIdLists[1]);
    process.exit(1);
  }

  const covered = new Set(reqIdLists[0]);

  if (!fs.existsSync(SEARCH_DIR)) {
    console.error("找不到", SEARCH_DIR);
    process.exit(1);
  }
  const searchFiles = listTsFiles(SEARCH_DIR);
  const used = new Set();
  for (const f of searchFiles) {
    for (const u of extractSearchUrls(readText(f))) {
      used.add(u);
    }
  }
  const missing = [...used].filter((u) => !covered.has(u));
  if (missing.length) {
    console.error("以下搜索端点未被 fetch.ts 的 reqId 列表覆盖：");
    for (const u of missing) console.error("  -", u);
    process.exit(1);
  }

  console.log(
    "OK",
    reqIdLists[0].length,
    "URLs covered consistently;",
    used.size,
    "/api/search/* endpoints all present."
  );
}

main();
