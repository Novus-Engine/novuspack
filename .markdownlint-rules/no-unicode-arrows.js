"use strict";

const {
  iterateNonFencedLines,
  stripInlineCode,
} = require("./utils.js");

const UNICODE_ARROWS = /[\u2190-\u21FF\u2900-\u297F]/u;

module.exports = {
  names: ["no-unicode-arrows"],
  description: "Disallow Unicode arrow characters in prose; use '=>' instead.",
  tags: ["content"],
  function: function (params, onError) {
    for (const { lineNumber, line } of iterateNonFencedLines(params.lines)) {
      const scan = stripInlineCode(line);
      const match = scan.match(UNICODE_ARROWS);
      if (!match) {
        continue;
      }

      onError({
        lineNumber,
        detail: "Unicode arrow characters are not allowed. Use '=>' instead.",
        context: line,
      });
    }
  },
};
