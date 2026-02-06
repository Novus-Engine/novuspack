"use strict";

/**
 * Parse placement options for one anchor regex. No idPattern here; placement is attached to the pattern.
 */
function parsePlacement(placement, patternIndex) {
  if (!placement || typeof placement !== "object") {
    return null;
  }
  let headingMatch = null;
  if (typeof placement.headingMatch === "string" && placement.headingMatch) {
    try {
      headingMatch = new RegExp(placement.headingMatch);
    } catch {
      headingMatch = null;
    }
  }
  let lineMatch = null;
  if (typeof placement.lineMatch === "string" && placement.lineMatch) {
    try {
      lineMatch = new RegExp(placement.lineMatch);
    } catch {
      lineMatch = null;
    }
  }
  const standaloneLine = placement.standaloneLine === true;
  const requireAfter = Array.isArray(placement.requireAfter)
    ? placement.requireAfter.filter((x) =>
        ["blank", "fencedBlock", "list"].includes(x),
      )
    : [];
  const anchorImmediatelyAfterHeading =
    placement.anchorImmediatelyAfterHeading === true;
  const maxPerSection =
    typeof placement.maxPerSection === "number" &&
    placement.maxPerSection >= 1
      ? placement.maxPerSection
      : null;
  return {
    patternIndex,
    headingMatch,
    lineMatch,
    standaloneLine,
    requireAfter,
    anchorImmediatelyAfterHeading,
    maxPerSection,
  };
}

function getConfig(params) {
  const c = params.config || {};
  const raw = Array.isArray(c.allowedIdPatterns) ? c.allowedIdPatterns : [];
  const allowedEntries = [];

  for (let i = 0; i < raw.length; i++) {
    const item = raw[i];
    let patternStr = null;
    let placementRaw = null;

    if (typeof item === "string" && item.length > 0) {
      patternStr = item;
    } else if (item && typeof item === "object" && typeof item.pattern === "string" && item.pattern.length > 0) {
      patternStr = item.pattern;
      placementRaw = item.placement;
    }

    if (!patternStr) {
      continue;
    }

    let pattern;
    try {
      pattern = new RegExp(patternStr);
    } catch {
      continue;
    }

    const placement = parsePlacement(placementRaw, i);
    allowedEntries.push({ pattern, placement });
  }

  const strictPlacement = c.strictPlacement !== false;
  return { allowedEntries, strictPlacement };
}

module.exports = {
  names: ["allow-custom-anchors"],
  description:
    "Allow only configured <a id=\"...\"></a> anchor id patterns; optional placement rules.",
  tags: ["html", "anchors"],
  function: function (params, onError) {
    const { allowedEntries, strictPlacement } = getConfig(params);
    const allowedPatterns = allowedEntries.map((e) => e.pattern);

    const anchorTagRegex = /<a id="([^"]+)"><\/a>/;
    const anchorAtEndOfLineRegex = /<a id="([^"]+)"><\/a>\s*$/;

    function stripInlineCode(line) {
      let out = "";
      let inCode = false;
      let fence = "";

      for (let i = 0; i < line.length; i++) {
        const ch = line[i];
        if (ch !== "`") {
          out += inCode ? " " : ch;
          continue;
        }

        let j = i;
        while (j < line.length && line[j] === "`") {
          j++;
        }
        const run = line.slice(i, j);

        if (!inCode) {
          inCode = true;
          fence = run;
        } else if (run === fence) {
          inCode = false;
          fence = "";
        }

        out += run;
        i = j - 1;
      }

      return out;
    }

    let inFence = false;
    let fenceMarker = null;

    /** Stack of { patternIndex, level } for sections we're inside (by headingMatch). */
    const sectionStack = [];
    /** For each pattern index with maxPerSection: count of anchors seen in current section. */
    const sectionAnchorCount = new Map();

    for (let index = 0; index < params.lines.length; index++) {
      const lineNumber = index + 1;
      const line = params.lines[index];
      const trimmed = line.trim();

      const fenceMatch = trimmed.match(/^(```+|~~~+)/);
      if (fenceMatch) {
        const marker = fenceMatch[1][0] === "`" ? "```" : "~~~";
        if (!inFence) {
          inFence = true;
          fenceMarker = marker;
        } else if (fenceMarker === marker) {
          inFence = false;
          fenceMarker = null;
        }
        continue;
      }

      if (inFence) {
        continue;
      }

      const headingMatch = trimmed.match(/^(#{1,6})\s+/);
      if (headingMatch) {
        const level = headingMatch[1].length;
        while (sectionStack.length > 0 && sectionStack[sectionStack.length - 1].level >= level) {
          sectionStack.pop();
        }
        for (let pi = 0; pi < allowedEntries.length; pi++) {
          const entry = allowedEntries[pi];
          const pl = entry.placement;
          if (pl && pl.headingMatch && pl.headingMatch.test(trimmed)) {
            sectionStack.push({ patternIndex: pi, level });
            sectionAnchorCount.set(pi, 0);
          }
        }
        continue;
      }

      const scanLine = stripInlineCode(line);
      const anchorIndex = scanLine.indexOf("<a");
      if (anchorIndex === -1) {
        continue;
      }

      if (scanLine.indexOf("<a", anchorIndex + 1) !== -1) {
        onError({
          lineNumber,
          detail: "[one-per-line] Only one <a id=\"...\"></a> anchor is allowed per line.",
          context: line,
        });
        continue;
      }

      const match = scanLine.match(anchorTagRegex);
      if (!match) {
        onError({
          lineNumber,
          detail:
            "[anchor-format] Only <a id=\"...\"></a> anchors are allowed, with id as the only attribute.",
          context: line,
        });
        continue;
      }

      const id = match[1];

      const allowed = allowedPatterns.some((re) => re.test(id));
      if (!allowed) {
        onError({
          lineNumber,
          detail:
            "[allowedIdPatterns] Anchor id must match one of the configured allowedIdPatterns.",
          context: line,
        });
        continue;
      }

      const anchorMatch = scanLine.match(anchorAtEndOfLineRegex);
      if (!anchorMatch) {
        onError({
          lineNumber,
          detail:
            "[end-of-line] Anchors must appear at the end of the line (or be a standalone reference anchor line above a fenced code block).",
          context: line,
        });
        continue;
      }

      const matchIndex = allowedEntries.findIndex((e) => e.pattern.test(id));
      const entry = allowedEntries[matchIndex];
      const rule = entry.placement;

      if (!strictPlacement || !rule) {
        continue;
      }

      const anchorPosOriginal = line.lastIndexOf("<a");
      const beforeOriginal = (anchorPosOriginal >= 0
        ? line.slice(0, anchorPosOriginal)
        : line).trim();

      if (rule.lineMatch && !rule.lineMatch.test(beforeOriginal)) {
        onError({
          lineNumber,
          detail:
            "[lineMatch] Anchor line must match the configured lineMatch pattern for this id.",
          context: line,
        });
        continue;
      }

      if (rule.standaloneLine && trimmed !== `<a id="${id}"></a>`) {
        onError({
          lineNumber,
          detail:
            "[standaloneLine] This anchor must be on its own line (no other content).",
          context: line,
        });
        continue;
      }

      if (rule.headingMatch) {
        const inSection = sectionStack.some(
          (s) => s.patternIndex === matchIndex,
        );
        if (!inSection) {
          onError({
            lineNumber,
            detail:
              "[headingMatch] This anchor must appear within a section whose heading matches the configured headingMatch.",
            context: line,
          });
          continue;
        }

        if (rule.maxPerSection != null) {
          const count = sectionAnchorCount.get(matchIndex) || 0;
          if (count >= rule.maxPerSection) {
            onError({
              lineNumber,
              detail: `[maxPerSection] Only ${rule.maxPerSection} anchor(s) of this type allowed per section.`,
              context: line,
            });
            continue;
          }
          sectionAnchorCount.set(matchIndex, count + 1);
        }
      }

      if (rule.anchorImmediatelyAfterHeading) {
        let prev = index - 1;
        while (prev >= 0 && params.lines[prev].trim() === "") {
          prev--;
        }
        const prevLine = prev >= 0 ? params.lines[prev].trim() : "";
        const isAnyAtxHeading = /^\s*#{1,6}\s+/.test(prevLine);
        const matchesHeading =
          rule.headingMatch
            ? rule.headingMatch.test(prevLine)
            : isAnyAtxHeading;
        if (prev < 0 || !matchesHeading) {
          onError({
            lineNumber,
            detail:
              "[anchorImmediatelyAfterHeading] This anchor must appear immediately after the section heading (blank lines allowed).",
            context: line,
          });
          continue;
        }
      }

      if (rule.requireAfter.length > 0) {
        const next = params.lines[index + 1];
        const next2 = params.lines[index + 2];
        const needBlank = rule.requireAfter[0] === "blank";
        const needFenced = rule.requireAfter.includes("fencedBlock");
        const needList = rule.requireAfter.includes("list");

        if (needBlank && (next == null || next.trim() !== "")) {
          onError({
            lineNumber,
            detail: "[requireAfter] Anchor line must be followed by a blank line.",
            context: line,
          });
          continue;
        }

        const checkLine = needBlank ? next2 : next;

        if (needFenced && (checkLine == null || !checkLine.trim().match(/^(```+|~~~+)/))) {
          onError({
            lineNumber,
            detail:
              "[requireAfter] Anchor line must be followed by a blank line and then a fenced code block.",
            context: line,
          });
          continue;
        }

        if (needList && (checkLine == null || !/^\s*(?:[-*+]\s+|\d+[.)]\s+)/.test(checkLine.trim()))) {
          onError({
            lineNumber,
            detail:
              "[requireAfter] Anchor line must be followed by a blank line and then a list (ordered or unordered).",
            context: line,
          });
          continue;
        }
      }
    }
  },
};
