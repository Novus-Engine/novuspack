"use strict";

const {
  extractHeadings,
  parseHeadingNumberPrefix,
} = require("./novuspack-utils.js");

module.exports = {
  names: ["novuspack-heading-numbering-segment-count"],
  description:
    "If a heading has a numbering prefix, segment count must equal (heading level - 1).",
  tags: ["headings"],
  function: function (params, onError) {
    const headings = extractHeadings(params.lines);

    for (const h of headings) {
      const { numbering } = parseHeadingNumberPrefix(h.rawText);
      if (numbering == null) {
        continue;
      }

      const segments = numbering.split(".");
      const expectedCount = h.level - 1;
      if (segments.length !== expectedCount) {
        onError({
          lineNumber: h.lineNumber,
          detail: `H${h.level} heading has ${segments.length} number(s), expected ${expectedCount}.`,
          context: params.lines[h.lineNumber - 1],
        });
      }
    }
  },
};
