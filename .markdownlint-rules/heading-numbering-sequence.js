"use strict";

const {
  extractHeadings,
  parseHeadingNumberPrefix,
} = require("./utils.js");

/**
 * Build parent index for each heading (0-based index of parent in sorted list).
 * Siblings are same-level headings under the same parent; parent is previous heading with level one less.
 */
function buildParentIndex(headings) {
  const sorted = headings.slice().sort((a, b) => a.lineNumber - b.lineNumber);
  const parentIndex = [];
  const stack = [];

  for (let i = 0; i < sorted.length; i++) {
    const h = sorted[i];
    while (stack.length > 0 && stack[stack.length - 1].level >= h.level) {
      stack.pop();
    }
    const parent = stack.length > 0 ? stack[stack.length - 1].index : null;
    parentIndex[i] = parent;
    stack.push({ level: h.level, index: i });
  }

  return { sorted, parentIndex };
}

/**
 * Get expected number for a heading at index i: parent prefix + sibling sequence.
 * H2 may be 0-based (0, 1, 2...) or 1-based (1, 2, 3...); deeper levels are 1-based.
 */
function getExpectedNumber(sorted, parentIndex, i, firstH2Numbering) {
  const h = sorted[i];
  const { numbering } = parseHeadingNumberPrefix(h.rawText);
  if (numbering == null) {
    return null;
  }

  const parent = parentIndex[i] != null ? sorted[parentIndex[i]] : null;
  const parentNum =
    parent != null ? parseHeadingNumberPrefix(parent.rawText).numbering : null;

  const siblings = [];
  for (let j = 0; j < sorted.length; j++) {
    if (parentIndex[j] !== parentIndex[i]) {
      continue;
    }
    if (sorted[j].level !== h.level) {
      continue;
    }
    siblings.push(sorted[j]);
  }
  siblings.sort((a, b) => a.lineNumber - b.lineNumber);

  const myIdx = siblings.findIndex((s) => s.lineNumber === h.lineNumber);
  if (myIdx < 0) {
    return null;
  }

  const isH2 = h.level === 2;
  const startAtZero = isH2 && firstH2Numbering === "0";
  const nextNum = startAtZero ? myIdx : myIdx + 1;
  const prefix = parentNum ? parentNum + "." : "";
  return prefix + String(nextNum);
}

module.exports = {
  names: ["heading-numbering-sequence"],
  description:
    "Numbering must be sequential within a parent and match parent prefix; H2 punctuation consistent.",
  tags: ["headings"],
  function: function (params, onError) {
    const headings = extractHeadings(params.lines);
    const withNumbering = headings
      .map((h) => ({
        ...h,
        parsed: parseHeadingNumberPrefix(h.rawText),
      }))
      .filter((h) => h.parsed.numbering != null);

    if (withNumbering.length === 0) {
      return;
    }

    const firstH2 = headings.find((h) => h.level === 2);
    const firstH2Num =
      firstH2 != null ? parseHeadingNumberPrefix(firstH2.rawText) : null;
    const docUsesNumbering =
      firstH2Num != null && firstH2Num.numbering != null;
    if (!docUsesNumbering) {
      return;
    }

    const { sorted, parentIndex } = buildParentIndex(headings);

    const firstH2ByLine = sorted
      .filter((x) => x.level === 2)
      .sort((a, b) => a.lineNumber - b.lineNumber)[0];
    const firstH2Numbering =
      firstH2ByLine != null
        ? parseHeadingNumberPrefix(firstH2ByLine.rawText).numbering
        : null;

    const firstH2Dot =
      firstH2ByLine != null
        ? parseHeadingNumberPrefix(firstH2ByLine.rawText).hasH2Dot
        : false;

    for (let i = 0; i < sorted.length; i++) {
      const h = sorted[i];
      const { numbering, hasH2Dot } = parseHeadingNumberPrefix(h.rawText);
      if (numbering == null) {
        continue;
      }

      if (h.level === 2 && hasH2Dot !== firstH2Dot) {
        onError({
          lineNumber: h.lineNumber,
          detail: `H2 period inconsistency: use ${firstH2Dot ? "period" : "no period"} after number to match first H2.`,
          context: params.lines[h.lineNumber - 1],
        });
      }

      const expected = getExpectedNumber(
        sorted,
        parentIndex,
        i,
        firstH2Numbering
      );
      if (expected != null && numbering !== expected) {
        onError({
          lineNumber: h.lineNumber,
          detail: `Non-sequential numbering: got '${numbering}', expected '${expected}'.`,
          context: params.lines[h.lineNumber - 1],
        });
      }
    }
  },
};
