---
alwaysApply: true
---

# AI Instructions

## General Rules

- Check the actual date before writing the date.
- See [../ai_files/](../ai_files/) for AI assisted coding instructions.
- See [../docs/tech_specs/](../docs/tech_specs/) for technical specifications.
- For markdown files, always abide by the [markdown standards](#markdown-standards) below.
- Avoid writing scripts; use existing tools in the repository.
- Avoid using commands which require approval.
- Commands that do not require approval:
  - awk
  - basename
  - cat
  - cd
  - date
  - echo
  - find
  - git diff
  - git log
  - git show
  - git status
  - go build
  - go test
  - go tool cover
  - grep
  - head
  - ls
  - make bdd
  - pwd
  - realpath
  - sha256sum
  - sort
  - tail
  - test
  - uniq
  - wc
- Unapproved commands:
  - sed
  - xargs
  - tee

## Markdown Standards

- Avoid pseudo-headings or heading-like lines.
  - Use proper Markdown headings instead of `^\*\*.*:\*\*$`, `^\*\*.*\*\*:$`, `^[0-9]+\. \*\*.*\*\*$` or similar.
- Use "=>" instead of "→".
- Avoid using non-ASCII characters, with the following exceptions:
  - The following may be used in AI work tracking documents: ✅, ❌, 📊, ⚠️
- Bulleted lists that are indented below numbered lists should be indented with 4 spaces, not 3.
- Always include a blank line after any heading lines.
- Always include a blank line before and after a list.
- Put one sentence on a line, except in tables.
- Code blocks that are part of a list should be indented by four spaces.
- All headings within a document must be unique
- All headings must have proper numbering.
- Ensure there is no content duplication; use references instead.
- When referencing a file or path, make it a link.

## Tech Specs Docs

- Use prose to describe the technical implementation.
- Tech specs docs should have minimal code aside from:
  - Function signatures
  - Constants definitions
  - Generic type definitions (type parameters, constraints)
  - Basic usage examples
- There must be only one source of truth for each specification.
  Any related specifications should refer to the singular source of truth and link to it instead of re-stating it.
- All error handling should be using the latest structured errors approach.
  References to legacy sentinel errors should be removed and replaced with structured errors.
