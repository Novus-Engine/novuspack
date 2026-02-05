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
 * Level of the "numbering root" for heading at index i: the nearest ancestor that has no numbering (or 1 at doc root).
 * Segment count for a numbered heading = heading.level - numberingRootLevel.
 */
function getNumberingRootLevel(sorted, parentIndex, i) {
  const parentIdx = parentIndex[i];
  if (parentIdx == null) {
    return 1;
  }
  const parentNum = parseHeadingNumberPrefix(sorted[parentIdx].rawText).numbering;
  if (parentNum == null) {
    return sorted[parentIdx].level;
  }
  return getNumberingRootLevel(sorted, parentIndex, parentIdx);
}

/**
 * Get siblings of heading at index i (same parent, same level), sorted by line.
 */
function getSiblings(sorted, parentIndex, i) {
  const h = sorted[i];
  const siblings = [];
  for (let j = 0; j < sorted.length; j++) {
    if (parentIndex[j] !== parentIndex[i]) {
      continue;
    }
    if (sorted[j].level !== h.level) {
      continue;
    }
    siblings.push({ index: j, ...sorted[j] });
  }
  siblings.sort((a, b) => a.lineNumber - b.lineNumber);
  return siblings;
}

/**
 * Expected number for heading at index i within its section (parent prefix + sibling sequence).
 * Section-scoped: only used when at least one sibling has numbering.
 * 0-based when the first numbered sibling has numbering "0".
 */
function getExpectedNumberInSection(sorted, parentIndex, i) {
  const h = sorted[i];
  const parentIdx = parentIndex[i];
  const parent = parentIdx != null ? sorted[parentIdx] : null;
  const parentNum =
    parent != null ? parseHeadingNumberPrefix(parent.rawText).numbering : null;

  const siblings = getSiblings(sorted, parentIndex, i);
  const myIdx = siblings.findIndex((s) => s.lineNumber === h.lineNumber);
  if (myIdx < 0) {
    return null;
  }

  const firstNumbered = siblings.find((s) =>
    parseHeadingNumberPrefix(s.rawText).numbering != null
  );
  const firstNumbering =
    firstNumbered != null
      ? parseHeadingNumberPrefix(firstNumbered.rawText).numbering
      : null;
  const startAtZero = firstNumbering === "0";
  const nextNum = startAtZero ? myIdx : myIdx + 1;
  const prefix = parentNum ? parentNum + "." : "";
  return prefix + String(nextNum);
}

/**
 * Whether any sibling in the section (same parent) has numbering.
 */
function sectionUsesNumbering(sorted, parentIndex, i) {
  const siblings = getSiblings(sorted, parentIndex, i);
  return siblings.some(
    (s) => parseHeadingNumberPrefix(s.rawText).numbering != null
  );
}

/**
 * First period style (hasH2Dot) among numbered siblings in this section; null if none numbered.
 */
function getSectionPeriodStyle(sorted, parentIndex, i) {
  const siblings = getSiblings(sorted, parentIndex, i);
  const firstNumbered = siblings.find((s) =>
    parseHeadingNumberPrefix(s.rawText).numbering != null
  );
  if (firstNumbered == null) {
    return null;
  }
  return parseHeadingNumberPrefix(firstNumbered.rawText).hasH2Dot;
}

module.exports = {
  names: ["heading-numbering"],
  description:
    "Numbered headings: segment count by numbering root; numbering consistent within each section; period style consistent within section.",
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

    const { sorted, parentIndex } = buildParentIndex(headings);

    for (let i = 0; i < sorted.length; i++) {
      const h = sorted[i];
      const { numbering, hasH2Dot } = parseHeadingNumberPrefix(h.rawText);

      if (sectionUsesNumbering(sorted, parentIndex, i) && numbering == null) {
        onError({
          lineNumber: h.lineNumber,
          detail:
            "This section uses numbering; add a number prefix to match siblings.",
          context: params.lines[h.lineNumber - 1],
        });
        continue;
      }

      if (numbering == null) {
        continue;
      }

      const rootLevel = getNumberingRootLevel(sorted, parentIndex, i);
      const expectedSegmentCount = h.level - rootLevel;
      const segments = numbering.split(".");

      if (segments.length !== expectedSegmentCount) {
        onError({
          lineNumber: h.lineNumber,
          detail: `H${h.level} heading has ${segments.length} number(s), expected ${expectedSegmentCount} (level - numbering root).`,
          context: params.lines[h.lineNumber - 1],
        });
        continue;
      }

      if (!sectionUsesNumbering(sorted, parentIndex, i)) {
        continue;
      }

      const sectionPeriodStyle = getSectionPeriodStyle(sorted, parentIndex, i);
      if (sectionPeriodStyle != null && hasH2Dot !== sectionPeriodStyle) {
        onError({
          lineNumber: h.lineNumber,
          detail: `Period inconsistency in this section: use ${sectionPeriodStyle ? "period" : "no period"} after number to match sibling.`,
          context: params.lines[h.lineNumber - 1],
        });
      }

      const expected = getExpectedNumberInSection(sorted, parentIndex, i);
      if (expected != null && numbering !== expected) {
        onError({
          lineNumber: h.lineNumber,
          detail: `Non-sequential numbering in this section: got '${numbering}', expected '${expected}'.`,
          context: params.lines[h.lineNumber - 1],
        });
      }
    }
  },
};
