"use strict";

const DEFAULT_ALLOWED_ID_PATTERNS = [
  "^spec-[a-z0-9-]+$",
  "^ref-[a-z0-9]+-[a-z0-9-]+$",
  "^algo-[a-z0-9-]+$",
  "^algo-[a-z0-9-]+-step-[0-9]+(?:-[0-9]+)*$",
];

function getConfig(params) {
  const c = params.config || {};
  const raw = Array.isArray(c.allowedIdPatterns)
    ? c.allowedIdPatterns
    : DEFAULT_ALLOWED_ID_PATTERNS;
  const allowedIdPatterns = raw
    .filter((p) => typeof p === "string")
    .map((p) => new RegExp(p));
  const strictPlacement = c.strictPlacement !== false;
  return { allowedIdPatterns, strictPlacement };
}

module.exports = {
  names: ["allow-custom-anchors"],
  description:
    "Allow only configured <a id=\"...\"></a> anchor id patterns; optional placement rules.",
  tags: ["html", "anchors"],
  function: function (params, onError) {
    const { allowedIdPatterns, strictPlacement } = getConfig(params);

    /**
     * Strictly allow only: <a id="SOME_ID"></a>
     * - double-quotes required
     * - only attribute is id
     * - no inner content
     */
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
    let inAlgorithmSection = false;
    let algorithmHeadingLevel = null;
    let seenAlgorithmAnchorInSection = false;

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
        if (inAlgorithmSection && algorithmHeadingLevel != null && level <= algorithmHeadingLevel) {
          inAlgorithmSection = false;
          algorithmHeadingLevel = null;
          seenAlgorithmAnchorInSection = false;
        }

        if (/`[^`]+`.*\bAlgorithm\b/.test(trimmed)) {
          inAlgorithmSection = true;
          algorithmHeadingLevel = level;
          seenAlgorithmAnchorInSection = false;
        }
      }

      const scanLine = stripInlineCode(line);
      const anchorIndex = scanLine.indexOf("<a");
      if (anchorIndex === -1) {
        continue;
      }

      if (scanLine.indexOf("<a", anchorIndex + 1) !== -1) {
        onError({
          lineNumber,
          detail: "Only one <a id=\"...\"></a> anchor is allowed per line.",
          context: line,
        });
        continue;
      }

      const match = scanLine.match(anchorTagRegex);
      if (!match) {
        onError({
          lineNumber,
          detail:
            "Only <a id=\"...\"></a> anchors are allowed, with id as the only attribute.",
          context: line,
        });
        continue;
      }

      const id = match[1];

      const ok = allowedIdPatterns.some((re) => re.test(id));
      if (!ok) {
        onError({
          lineNumber,
          detail:
            "Anchor id must match one of the configured allowedIdPatterns.",
          context: line,
        });
        continue;
      }

      const anchorMatch = scanLine.match(anchorAtEndOfLineRegex);
      if (!anchorMatch) {
        onError({
          lineNumber,
          detail:
            "Anchors must appear at the end of the line (or be a standalone reference anchor line above a fenced code block).",
          context: line,
        });
        continue;
      }

      if (!strictPlacement) {
        continue;
      }

      const anchorPosOriginal = line.lastIndexOf("<a");
      const beforeOriginal = (anchorPosOriginal >= 0
        ? line.slice(0, anchorPosOriginal)
        : line).trim();

      const startsWithList =
        /^\s*(?:[-*+]\s+|\d+[.)]\s+)/.test(beforeOriginal);

      if (id.startsWith("spec-")) {
        if (!/^\s*-\s+Spec ID:\s+`NP\.[^`]+`/.test(beforeOriginal)) {
          onError({
            lineNumber,
            detail:
              "Spec anchors must be appended to the end of the '- Spec ID: `NP....`' list item line.",
            context: line,
          });
        }
        continue;
      }

      if (id.startsWith("ref-")) {
        const expected = `<a id="${id}"></a>`;
        if (trimmed !== expected) {
          onError({
            lineNumber,
            detail:
              "Reference anchors must be on their own line directly above a fenced code block.",
            context: line,
          });
          continue;
        }

        const next = params.lines[index + 1];
        const next2 = params.lines[index + 2];
        if (next == null || next.trim() !== "" || next2 == null || !next2.trim().match(/^(```+|~~~+)/)) {
          onError({
            lineNumber,
            detail:
              "Reference anchor line must be followed by a blank line and then a fenced code block.",
            context: line,
          });
        }
        continue;
      }

      if (id.startsWith("algo-") && !id.includes("-step-")) {
        const expected = `<a id="${id}"></a>`;
        if (trimmed !== expected) {
          onError({
            lineNumber,
            detail:
              "Algorithm anchors must be on their own line at the start of an Algorithm section.",
            context: line,
          });
          continue;
        }

        if (!inAlgorithmSection) {
          onError({
            lineNumber,
            detail:
              "Algorithm anchors must appear within an Algorithm section.",
            context: line,
          });
          continue;
        }

        if (seenAlgorithmAnchorInSection) {
          onError({
            lineNumber,
            detail:
              "Only one Algorithm anchor is allowed per Algorithm section.",
            context: line,
          });
          continue;
        }
        seenAlgorithmAnchorInSection = true;

        let prev = index - 1;
        while (prev >= 0 && params.lines[prev].trim() === "") {
          prev--;
        }
        if (prev < 0 || !/^\s*#{1,6}\s+.*`[^`]+`.*\bAlgorithm\b/.test(params.lines[prev].trim())) {
          onError({
            lineNumber,
            detail:
              "Algorithm anchor line must appear immediately after the Algorithm heading (allowing blank lines).",
            context: line,
          });
        }

        const next = params.lines[index + 1];
        const next2 = params.lines[index + 2];
        const startsWithListNext2 =
          next2 != null && /^\s*(?:[-*+]\s+|\d+[.)]\s+)/.test(next2);

        if (next == null || next.trim() !== "" || next2 == null || !startsWithListNext2) {
          onError({
            lineNumber,
            detail:
              "Algorithm anchor line must be followed by a blank line and then the procedure list (ordered or unordered).",
            context: line,
          });
        }
        continue;
      }

      if (id.includes("-step-")) {
        if (!beforeOriginal || !inAlgorithmSection || !startsWithList) {
          onError({
            lineNumber,
            detail:
              "Algorithm step anchors must be appended to the end of an ordered or unordered list item line within an Algorithm section.",
            context: line,
          });
        }
      }
    }
  },
};
