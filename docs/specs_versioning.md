# NovusPack Specifications Versioning

Date: 2025-10-30

## Goals

Provide a predictable process for proposing, reviewing, tagging, and changing specifications.

## Version Labels

- Draft: Unstable and under active editing.
- Review: Open for broader feedback.
- Accepted: Approved for implementation.
- Deprecated: Superseded; kept for historical reference.

## Change Types

- Breaking: Alters behavior or structure requiring consumer changes.
- Additive: Adds new sections without breaking existing behavior.
- Clarification: Improves wording without changing behavior.

## Process

1. Propose changes with a short summary and rationale.
2. Update the single source of truth file for the domain.
3. Update cross references and affected indexes.
4. Record the change in `docs/CHANGELOG.md`.
5. Update the spec state label in the document header.

## Tagging

- Use repository tags to mark major accepted milestones when appropriate.
- Include the spec state in the document header for clarity.

## Automation

- Enforce markdown and link checks in CI.
- Validate heading uniqueness and numbering.
- Optionally validate TOCs.
