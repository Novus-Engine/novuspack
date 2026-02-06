# Custom markdownlint rules

This directory contains custom rules used by [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2) for this repo. Rules are registered in the root [.markdownlint-cli2.jsonc](../.markdownlint-cli2.jsonc) and configured in [.markdownlint.yml](../.markdownlint.yml).

## Overview

- **Rule modules**: Each `*.js` file under this directory (except `utils.js`) is a custom rule.
  `utils.js` is a shared helper and is not a rule.
- **Config**: Rule-specific options are set in `.markdownlint.yml` under the rule name.
  Only rules that accept options are documented with a config section below.

## Rules

### allow-custom-anchors

**File:** `allow-custom-anchors.js`

**Description:** Allow only configured `<a id="..."></a>` anchor id patterns; optional configurable placement rules (heading match, line match, etc.).

**Configuration:** In `.markdownlint.yml` under `allow-custom-anchors`:

| Option              | Type                                | Default         | Meaning                                                                                                                               |
| ------------------- | ----------------------------------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| `allowedIdPatterns` | array of strings or pattern objects | none (required) | Each entry is a regex string or `{ pattern: string, placement?: object }`. No built-in default.                                       |
| `strictPlacement`   | boolean                             | `true`          | If `true`, enforce placement when the matching pattern has a `placement` object; if `false`, only id match and anchor at end of line. |

**Per-pattern placement** (optional `placement` on an entry in `allowedIdPatterns`):

| Property                        | Type    | Meaning                                                                                                                                                                                                                                                                                                  |
| ------------------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `headingMatch`                  | string  | Optional. Regex for the heading line. Anchor must be inside a section whose heading matches (sections tracked by heading level).                                                                                                                                                                         |
| `lineMatch`                     | string  | Optional. Regex for the line content before the anchor. The line (before the anchor) must match.                                                                                                                                                                                                         |
| `standaloneLine`                | boolean | Optional. If true, anchor must be the only content on its line.                                                                                                                                                                                                                                          |
| `requireAfter`                  | array   | Optional. Sequence after anchor line: `["blank"]`, `["blank", "fencedBlock"]`, or `["blank", "list"]`.                                                                                                                                                                                                   |
| `anchorImmediatelyAfterHeading` | boolean | Optional. If true, anchor line must follow (with optional blank lines) a heading. When `headingMatch` is set, that heading must match it; otherwise the previous non-blank line may be any ATX heading (`#`–`######`). Works when the anchor shares a line with other content (e.g. end of a list item). |
| `maxPerSection`                 | number  | Optional. Max anchors of this pattern per `headingMatch` section (e.g. 1).                                                                                                                                                                                                                               |

Order of entries matters: the first pattern that matches the anchor id is used. Put more specific patterns (e.g. algo-step) before general ones (e.g. algo). Entries may be a plain regex string (no placement) or `{ pattern: "regex", placement: { ... } }`.

**Behavior:**

- Only `<a id="..."></a>` is allowed (no other attributes, no inner content).
- Anchor `id` must match one of the configured patterns in `allowedIdPatterns`.
- Anchors must appear at the end of the line (or on a standalone line where required by that pattern's placement).
- When `strictPlacement` is true and the matching pattern has a `placement` object, the anchor is validated against that placement (heading match, line match, standalone, require-after, etc.).
- Error messages are prefixed with a sub-rule tag in brackets (e.g. `[lineMatch]`, `[headingMatch]`, `[requireAfter]`, `[allowedIdPatterns]`) so you can see which check failed.

### no-heading-like-lines

**File:** `no-heading-like-lines.js`

**Description:** Disallow heading-like lines that should be proper Markdown headings.

**Configuration:** None.

**Behavior:** Reports lines that look like headings but are not (e.g. `**Text:**`, `**Text**:`, `1. **Text**`, and italic variants). Prompts use of real `#` headings instead.

### ascii-only

**File:** `ascii-only.js`

**Description:** Disallow non-ASCII except in configured paths; optional replacement suggestions via `unicodeReplacements`.

**Configuration:** In `.markdownlint.yml` under `ascii-only`:

| Option                       | Type                                   | Default  | Meaning                                                                                                                                                |
| ---------------------------- | -------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `allowedPathPatternsUnicode` | list of strings                        | none     | Glob patterns for files where any non-ASCII is allowed.                                                                                                |
| `allowedPathPatternsEmoji`   | list of strings                        | none     | Glob patterns for files where only `allowedEmoji` characters are allowed.                                                                              |
| `allowedEmoji`               | list of strings                        | none     | Emoji (or other chars) allowed in paths matching `allowedPathPatternsEmoji`; each entry may be multi-codepoint (e.g. ⚠️); all code points are allowed. |
| `allowedUnicode`             | list of single-character strings       | none     | Optional. Characters allowed in all files (global allowlist).                                                                                          |
| `unicodeReplacements`        | object or array of [char, replacement] | built-in | Map of single Unicode character to suggested ASCII replacement in error messages. When omitted, rule uses built-in defaults (arrows, quotes, ≤≥×).     |

Glob matching supports `**` (any path) and `*` (within a segment). Paths are normalized (forward slashes, leading `./` removed). Relative patterns (no leading `/` or `*`) match both path-prefix (e.g. `dev_docs/foo.md`) and mid-path (e.g. absolute paths containing `dev_docs/`).

**Behavior:**

- No built-in path or emoji defaults; configure `allowedPathPatternsUnicode`, `allowedPathPatternsEmoji`, and `allowedEmoji` as needed.
- If the file path matches `allowedPathPatternsUnicode`, any non-ASCII is allowed in that file.
- If the file path matches `allowedPathPatternsEmoji`, only characters in `allowedEmoji` (and Unicode variation selectors U+FE00–U+FE0F) are allowed; other non-ASCII is reported with message: "Only the listed emoji (...) are allowed in this path. Replace or remove other non-ASCII characters."
- Characters in `allowedUnicode` (when configured) are allowed in all files.
- Non-ASCII is detected by code-point iteration (surrogate pairs treated as one character) and compared after NFC normalization.
- When reporting a disallowed non-ASCII line, any character present in `unicodeReplacements` is mentioned in the error with its suggested replacement.
- Inline code (backticks) is stripped before scanning.

### no-duplicate-headings-normalized

**File:** `no-duplicate-headings-normalized.js`

**Description:** Disallow duplicate heading titles after stripping numbering and normalizing.

**Configuration:** None.

**Behavior:** Extracts all headings, strips numeric prefixes (e.g. `1.2.3`), normalizes the title (case/whitespace), and reports any heading whose normalized title appears more than once in the document. The first occurrence is the reference; duplicates are reported with the line number of the first.

### heading-numbering

**File:** `heading-numbering.js`

**Description:** Enforces structure and consistency of numbered headings: segment count by numbering root; numbering sequential within each section; period style consistent within section.

**Configuration:** None.

**Behavior:**

1. **Segment count by numbering root:** For each heading with a numeric prefix (e.g. `### 1.2 Title`), the number of segments (split on `.`) must equal heading level minus the numbering root level. The numbering root is the nearest ancestor heading that has no numbering (or document root, level 1). Example: H2 under doc root → 1 segment; H3 under unnumbered `## Section` → 1 segment; H4 under `### 1. First` → 1 segment. Headings without a numeric prefix are ignored.
2. **Section-scoped consistency:** For each section (siblings under the same parent), if any sibling has numbering then all siblings at that level must be numbered sequentially (e.g. 1., 2., 3.) and use consistent period style (all `## 1. Title` or all `## 1 Title`). Unnumbered siblings in a numbered section are reported.

## Shared helper

**utils.js** is not a rule. It provides utilities used by several rules (e.g. `extractHeadings`, `iterateNonFencedLines`, `stripInlineCode`, `parseHeadingNumberPrefix`, `normalizedTitleForDuplicate`). Do not list it in `customRules` in `.markdownlint-cli2.jsonc`.
