"use strict";

const {
  iterateNonFencedLines,
  stripInlineCode,
} = require("./utils.js");

const DEFAULT_UNICODE_REPLACEMENTS = {
  "\u2192": "->",
  "\u2190": "<-",
  "\u2194": "<=>",
  "\u21D2": "=>",
  "\u21D0": "<=",
  "\u21D4": "<=>",
  "\u2264": "<=",
  "\u2265": ">=",
  "\u00D7": "*",
  "\u2033": "\"",
  "\u2032": "'",
  "\u201C": "\"",
  "\u201D": "\"",
  "\u2019": "'",
  "\u2018": "'",
};

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

function buildReplacementsMap(unicodeReplacements) {
  const map = new Map();
  if (!unicodeReplacements || typeof unicodeReplacements !== "object") {
    return map;
  }
  if (Array.isArray(unicodeReplacements)) {
    for (const entry of unicodeReplacements) {
      if (Array.isArray(entry) && entry.length >= 2 && typeof entry[0] === "string" && entry[0].length === 1) {
        map.set(entry[0], String(entry[1]));
      }
    }
    return map;
  }
  for (const [ch, replacement] of Object.entries(unicodeReplacements)) {
    if (typeof ch === "string" && ch.length === 1 && replacement != null) {
      map.set(ch, String(replacement));
    }
  }
  return map;
}

function toCharSet(arr) {
  const set = new Set();
  if (!Array.isArray(arr)) {
    return set;
  }
  for (const item of arr) {
    if (typeof item === "string" && item.length === 1) {
      set.add(item);
    }
  }
  return set;
}

function getConfig(params) {
  const c = params.config || {};
  return {
    allowedPathPatternsUnicode: Array.isArray(c.allowedPathPatternsUnicode)
      ? c.allowedPathPatternsUnicode
      : [],
    allowedPathPatternsEmoji: Array.isArray(c.allowedPathPatternsEmoji)
      ? c.allowedPathPatternsEmoji
      : [],
    allowedEmoji: Array.isArray(c.allowedEmoji) ? c.allowedEmoji : [],
    allowedUnicode: toCharSet(c.allowedUnicode),
    unicodeReplacements: buildReplacementsMap(
      c.unicodeReplacements ?? DEFAULT_UNICODE_REPLACEMENTS,
    ),
  };
}

module.exports = {
  names: ["ascii-only"],
  description:
    "Disallow non-ASCII except in configured paths; optional replacement suggestions via unicodeReplacements.",
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
    const allowedUnicodeSet = config.allowedUnicode;

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

      const disallowedChars = [];
      const nonAsciiInLine = scan.match(/[^\x00-\x7F]/g) || [];
      for (const ch of nonAsciiInLine) {
        if (allowedUnicodeSet.has(ch)) {
          continue;
        }
        if (allowEmojiOnly && allowedEmojiSet.has(ch)) {
          continue;
        }
        if (!disallowedChars.includes(ch)) {
          disallowedChars.push(ch);
        }
      }

      if (disallowedChars.length === 0) {
        continue;
      }

      const suggestions = [];
      for (const ch of disallowedChars) {
        const replacement = config.unicodeReplacements.get(ch);
        if (replacement !== undefined) {
          suggestions.push(`'${ch}' with '${replacement}'`);
        }
      }

      let detail;
      if (suggestions.length > 0) {
        detail = `Replace ${suggestions.join("; ")}. Non-ASCII not allowed here.`;
      } else if (allowEmojiOnly) {
        const list = config.allowedEmoji.join(", ");
        detail = `Non-ASCII only allowed here: ${list}. Use ASCII or remove.`;
      } else {
        detail = "Non-ASCII characters are not allowed. Use ASCII only.";
      }

      onError({
        lineNumber,
        detail,
        context: line,
      });
    }
  },
};
