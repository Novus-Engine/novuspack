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

**Description:** Disallow non-ASCII except in configured paths, where either full Unicode or a fixed set of emoji is allowed.

**Configuration:** In `.markdownlint.yml` under `ascii-only`:

| Option                         | Type                             | Default                    | Meaning                                                                   |
| ------------------------------ | -------------------------------- | -------------------------- | ------------------------------------------------------------------------- |
| `allowedPathPatternsUnicode`   | list of strings                  | `["**/README.md"]`         | Glob patterns for files where any non-ASCII is allowed.                   |
| `allowedPathPatternsEmoji`     | list of strings                  | `["dev_docs/**"]`          | Glob patterns for files where only `allowedEmoji` characters are allowed. |
| `allowedEmoji`                 | list of single-character strings | `["✅", "❌", "📊", "⚠️"]` | Characters allowed in paths matching `allowedPathPatternsEmoji`.          |

Glob matching supports `**` (any path) and `*` (within a segment). Paths are normalized (forward slashes, leading `./` removed) before matching.

**Behavior:**

- If the file path matches **only** `allowedPathPatternsUnicode`, any non-ASCII is allowed.
- If the file path matches **only** `allowedPathPatternsEmoji`, only characters in `allowedEmoji` are allowed; other non-ASCII is reported.
- If the path matches neither, no non-ASCII is allowed.
- Inline code (backticks) is stripped before scanning so code snippets are still checked for non-ASCII outside of allowed paths/emoji.

### no-duplicate-headings-normalized

**File:** `no-duplicate-headings-normalized.js`

**Description:** Disallow duplicate heading titles after stripping numbering and normalizing.

**Configuration:** None.

**Behavior:** Extracts all headings, strips numeric prefixes (e.g. `1.2.3`), normalizes the title (case/whitespace), and reports any heading whose normalized title appears more than once in the document. The first occurrence is the reference; duplicates are reported with the line number of the first.

### heading-numbering-segment-count

**File:** `heading-numbering-segment-count.js`

**Description:** If a heading has a numbering prefix, the number of segments must equal (heading level - 1).

**Configuration:** None.

**Behavior:** For each heading that has a numeric prefix (e.g. `### 1.2 Title`), the prefix is split by `.`; the number of segments must equal the heading level minus one. Examples: H2 may have `1` or `2` (one segment); H3 may have `1.1`, `2.3` (two segments); H4 may have `1.1.1` (three segments). Headings without a numeric prefix are ignored.

### heading-numbering-sequence

**File:** `heading-numbering-sequence.js`

**Description:** Numbering must be sequential within a parent and match the parent prefix; H2 punctuation (period after number) must be consistent.

**Configuration:** None.

**Behavior:** For documents that use numbered headings (determined by the first H2), the rule checks: (1) sibling headings are numbered sequentially (e.g. 8.2.1, 8.2.2, 8.2.3); (2) child prefix extends parent (e.g. under `8.2` children are `8.2.1`, `8.2.2`); (3) all H2 headings use the same style for a period after the number (e.g. all `1.` or all `1`). Non-sequential or mismatched numbers are reported.

### no-unicode-arrows

**File:** `no-unicode-arrows.js`

**Description:** Disallow Unicode arrow characters in prose; use `=>` instead.

**Configuration:** None.

**Behavior:** Scans non-fenced, non-inline-code content for Unicode arrow characters (e.g. `→`, `↔`, `⇒`). Reports them and suggests using ASCII `=>` instead.

## Shared helper

**utils.js** is not a rule. It provides utilities used by several rules (e.g. `extractHeadings`, `iterateNonFencedLines`, `stripInlineCode`, `parseHeadingNumberPrefix`, `normalizedTitleForDuplicate`). Do not list it in `customRules` in `.markdownlint-cli2.jsonc`.
