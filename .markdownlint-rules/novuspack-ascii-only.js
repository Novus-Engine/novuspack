"use strict";

const {
  iterateNonFencedLines,
  stripInlineCode,
} = require("./novuspack-utils.js");

const DEV_DOCS_ALLOWED = new Set(["✅", "❌", "📊", "⚠️"]);

function isUnderDevDocs(filePath) {
  if (!filePath || typeof filePath !== "string") {
    return false;
  }
  const normalized = filePath.replace(/\\/g, "/");
  return normalized.includes("dev_docs/");
}

function isReadme(filePath) {
  if (!filePath || typeof filePath !== "string") {
    return false;
  }
  const normalized = filePath.replace(/\\/g, "/");
  return normalized.endsWith("README.md");
}

function hasNonAscii(str) {
  return /[^\x00-\x7F]/.test(str);
}

function onlyAllowedDevDocsEmoji(str) {
  const nonAscii = str.match(/[^\x00-\x7F]/g);
  if (!nonAscii) {
    return true;
  }
  for (const ch of nonAscii) {
    if (!DEV_DOCS_ALLOWED.has(ch)) {
      return false;
    }
  }
  return true;
}

module.exports = {
  names: ["novuspack-ascii-only"],
  description:
    "Disallow non-ASCII except in dev_docs (limited emoji) and README.md.",
  tags: ["content"],
  function: function (params, onError) {
    const filePath = params.name || "";

    const allowUnicode = isReadme(filePath);
    const allowDevDocsEmoji = isUnderDevDocs(filePath);

    for (const { lineNumber, line } of iterateNonFencedLines(params.lines)) {
      const scan = stripInlineCode(line);
      if (!hasNonAscii(scan)) {
        continue;
      }

      if (allowUnicode) {
        continue;
      }

      if (allowDevDocsEmoji && onlyAllowedDevDocsEmoji(scan)) {
        continue;
      }

      if (allowDevDocsEmoji) {
        onError({
          lineNumber,
          detail:
            "Non-ASCII only allowed in dev_docs for ✅, ❌, 📊, ⚠️. Use ASCII or remove.",
          context: line,
        });
      } else {
        onError({
          lineNumber,
          detail: "Non-ASCII characters are not allowed. Use ASCII only.",
          context: line,
        });
      }
    }
  },
};
