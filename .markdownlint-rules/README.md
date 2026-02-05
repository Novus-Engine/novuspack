# Custom markdownlint rules

This directory contains custom rules used by [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2) for this repo. Rules are registered in the root [.markdownlint-cli2.jsonc](../.markdownlint-cli2.jsonc) and configured in [.markdownlint.yml](../.markdownlint.yml).

## Overview

- **Rule modules**: Each `*.js` file under this directory (except `utils.js`) is a custom rule. `utils.js` is a shared helper and is not a rule.
- **Config**: Rule-specific options are set in `.markdownlint.yml` under the rule name. Only rules that accept options are documented with a config section below.

## Rules

### allow-custom-anchors

**File:** `allow-custom-anchors.js`

**Description:** Allow only configured `<a id="..."></a>` anchor id patterns; optional placement rules (spec/ref/algo).

**Configuration:** In `.markdownlint.yml` under `allow-custom-anchors`:

| Option               | Type    | Default | Meaning |
| -------------------- | ------- | ------- | ------- |
| `allowedIdPatterns`  | list of regex strings | spec/ref/algo patterns | Regex strings for allowed anchor `id` values. |
| `strictPlacement`    | boolean | `true`  | If `true`, enforce spec/ref/algo placement rules; if `false`, only validate id match and anchor at end of line. |

**Behavior:**

- Only `<a id="..."></a>` is allowed (no other attributes, no inner content).
- Anchor `id` must match one of the configured `allowedIdPatterns` (each string is compiled to a RegExp).
- Anchors must appear at the end of the line (or standalone for ref anchors).
- When `strictPlacement` is true: spec anchors on Spec ID line; ref anchors alone above a fenced block; algo anchors in Algorithm sections with placement rules.

### no-heading-like-lines

**File:** `no-heading-like-lines.js`

**Description:** Disallow heading-like lines that should be proper Markdown headings.

**Configuration:** None.

**Behavior:** Reports lines that look like headings but are not (e.g. `**Text:**`, `**Text**:`, `1. **Text**`, and italic variants). Prompts use of real `#` headings instead.

### ascii-only

**File:** `ascii-only.js`

**Description:** Disallow non-ASCII except in configured paths; optional replacement suggestions via a pairing config.

**Configuration:** In `.markdownlint.yml` under `ascii-only`:

| Option                         | Type                             | Default | Meaning                                                                   |
| ------------------------------ | -------------------------------- | ------- | ------------------------------------------------------------------------- |
| `allowedPathPatternsUnicode`   | list of strings                  | none    | Glob patterns for files where any non-ASCII is allowed.                   |
| `allowedPathPatternsEmoji`     | list of strings                  | none    | Glob patterns for files where only `allowedEmoji` characters are allowed. |
| `allowedEmoji`                 | list of single-character strings | none    | Characters allowed in paths matching `allowedPathPatternsEmoji`.          |
| `allowedUnicode`               | list of single-character strings | none    | Optional. Characters allowed in all files (global allowlist).             |
| `unicodeReplacements`          | object or array of [char, replacement] | built-in | Map of single Unicode character to suggested ASCII replacement in error messages. When omitted, rule uses built-in defaults (arrows, quotes, ≤≥×). |

Glob matching supports `**` (any path) and `*` (within a segment). Paths are normalized (forward slashes, leading `./` removed) before matching.

**Behavior:**

- No built-in path or emoji defaults; configure `allowedPathPatternsUnicode`, `allowedPathPatternsEmoji`, and `allowedEmoji` as needed.
- If the file path matches `allowedPathPatternsUnicode`, any non-ASCII is allowed in that file.
- If the file path matches `allowedPathPatternsEmoji`, only characters in `allowedEmoji` are allowed; other non-ASCII is reported.
- Characters in `allowedUnicode` (when configured) are allowed in all files.
- When reporting a disallowed non-ASCII line, any character present in `unicodeReplacements` is mentioned in the error with its suggested replacement.
- Inline code (backticks) is stripped before scanning.

### no-duplicate-headings-normalized

**File:** `no-duplicate-headings-normalized.js`

**Description:** Disallow duplicate heading titles after stripping numbering and normalizing.

**Configuration:** None.

**Behavior:** Extracts all headings, strips numeric prefixes (e.g. `1.2.3`), normalizes the title (case/whitespace), and reports any heading whose normalized title appears more than once in the document. The first occurrence is the reference; duplicates are reported with the line number of the first.

### heading-numbering

**File:** `heading-numbering.js`

**Description:** Enforces structure and consistency of numbered headings: segment count matches level; numbering is sequential within a parent and matches parent prefix; H2 punctuation is consistent.

**Configuration:** None.

**Behavior:**

1. **Segment count:** For each heading with a numeric prefix (e.g. `### 1.2 Title`), the number of segments (split on `.`) must equal heading level minus one. H2 → 1 segment; H3 → 2 segments; H4 → 3 segments. Headings without a numeric prefix are ignored.
2. **Sequence (when the doc uses numbering):** If at least one H2 has a number, the rule also checks: sibling headings are numbered sequentially (e.g. 8.2.1, 8.2.2); child number extends parent prefix (e.g. under `## 8.2` use `### 8.2.1`, `### 8.2.2`); all H2 headings use the same style for a period after the number (all `## 1. Title` or all `## 1 Title`).

## Shared helper

**utils.js** is not a rule. It provides utilities used by several rules (e.g. `extractHeadings`, `iterateNonFencedLines`, `stripInlineCode`, `parseHeadingNumberPrefix`, `normalizedTitleForDuplicate`). Do not list it in `customRules` in `.markdownlint-cli2.jsonc`.
