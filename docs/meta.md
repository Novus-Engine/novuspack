# NovusPack Documentation Meta

Date: 2025-10-30

## Purpose

Define the scope, audience, governance, and entry points for the NovusPack documentation set.

## Audience

- Engineers implementing NovusPack APIs.
- Security reviewers assessing cryptography and validation.
- Project maintainers curating specifications.

## Canonical Entry Point

- Technical specifications canonical entry: `docs/tech_specs/_main.md`

## Governance And Source Of Truth

- Each domain must have a single source of truth file.
- Related overviews or indexes must link to the source of truth without restating content.
- Changes to a domain spec must update the source of truth and cross references.

## Versioning And Change Management

- Specification versioning policy is defined in `docs/specs_versioning.md`.
- Record changes in `docs/CHANGELOG.md` (or per domain changelogs if needed).

## Standards

- Follow `.github/copilot_instructions.md` for markdown standards and document structure.
- Headings must be unique and properly numbered.
- Keep code in specs minimal: function signatures, constants, type parameters, and small usage examples.

## Related Documents

- `docs/glossary.md` for shared terminology.
- `docs/specs_versioning.md` for versioning policy.
- `docs/tech_specs/_main.md` for the full technical specifications index.
