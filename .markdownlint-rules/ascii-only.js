"use strict";

const {
  iterateNonFencedLines,
  stripInlineCode,
} = require("./utils.js");

const DEFAULT_UNICODE_PATTERNS = ["**/README.md"];
const DEFAULT_EMOJI_PATTERNS = ["dev_docs/**"];
const DEFAULT_ALLOWED_EMOJI = ["✅", "❌", "📊", "⚠️"];

/**
 * Match path against a single glob pattern. Supports ** (any path) and * (segment).
 * Path is normalized to forward slashes.
 */
function matchGlob(path, pattern) {
  if (!path || !pattern) {
    return false;
  }
  const normalized = path.replace(/\\/g, "/").replace(/^\.\//, "");
  const re = globToRegExp(pattern);
  return re.test(normalized);
}

function globToRegExp(pattern) {
  const parts = [];
  let i = 0;
  while (i < pattern.length) {
    if (pattern[i] === "*" && pattern[i + 1] === "*") {
      parts.push(".*");
      i += 2;
    } else if (pattern[i] === "*") {
      parts.push("[^/]*");
      i += 1;
    } else {
      parts.push(pattern[i].replace(/[.+?^${}()|[\]\\]/g, "\\$&"));
      i += 1;
    }
  }
  const source = parts.join("");
  return new RegExp("^" + source + "$");
}

function pathMatchesAny(path, patterns) {
  if (!Array.isArray(patterns) || patterns.length === 0) {
    return false;
  }
  for (const p of patterns) {
    if (typeof p === "string" && matchGlob(path, p)) {
      return true;
    }
  }
  return false;
}

function hasNonAscii(str) {
  return /[^\x00-\x7F]/.test(str);
}

function onlyAllowedEmoji(str, allowedSet) {
  const nonAscii = str.match(/[^\x00-\x7F]/g);
  if (!nonAscii) {
    return true;
  }
  for (const ch of nonAscii) {
    if (!allowedSet.has(ch)) {
      return false;
    }
  }
  return true;
}

function getConfig(params) {
  const c = params.config || {};
  return {
    allowedPathPatternsUnicode:
      c.allowedPathPatternsUnicode ?? DEFAULT_UNICODE_PATTERNS,
    allowedPathPatternsEmoji:
      c.allowedPathPatternsEmoji ?? DEFAULT_EMOJI_PATTERNS,
    allowedEmoji: Array.isArray(c.allowedEmoji)
      ? c.allowedEmoji
      : DEFAULT_ALLOWED_EMOJI,
  };
}

module.exports = {
  names: ["ascii-only"],
  description:
    "Disallow non-ASCII except in configured paths (allowed emoji list).",
  tags: ["content"],
  function: function (params, onError) {
    const filePath = params.name || "";
    const config = getConfig(params);
    const allowUnicode = pathMatchesAny(
      filePath,
      config.allowedPathPatternsUnicode,
    );
    const allowEmojiOnly = pathMatchesAny(
      filePath,
      config.allowedPathPatternsEmoji,
    );
    const allowedEmojiSet = new Set(config.allowedEmoji);

    for (const { lineNumber, line } of iterateNonFencedLines(params.lines)) {
      const scan = stripInlineCode(line);
      if (!hasNonAscii(scan)) {
        continue;
      }

      if (allowUnicode) {
        continue;
      }

      if (allowEmojiOnly && onlyAllowedEmoji(scan, allowedEmojiSet)) {
        continue;
      }

      if (allowEmojiOnly) {
        const list = config.allowedEmoji.join(", ");
        onError({
          lineNumber,
          detail: `Non-ASCII only allowed here: ${list}. Use ASCII or remove.`,
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
