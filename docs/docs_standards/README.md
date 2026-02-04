# Documentation Standards

This directory defines **documentation standards** for the NovusPack project.

It describes how documentation (specs, requirements, READMEs, and other Markdown) should be structured, formatted, and maintained so that docs stay consistent and easy to navigate across the repo.

## Purpose

- **Single source of truth** for doc conventions (structure, style, linking, tooling).
- **Onboarding** for contributors and tooling (e.g. lint, validation scripts) so everyone follows the same rules.
- **Alignment** with existing automation (e.g. `make docs-check`) and future checks.

## Where the Standards Live

The standards themselves are **not** written in this README. They are kept in separate documents in this directory and linked from here. That keeps this file as an index and keeps each standard focused and easy to update.

When a standard is added, it will be listed below with a short description and link.

- [Markdown conventions](markdown_conventions.md) – Repository-wide Markdown formatting rules (including heading numbering and uniqueness).
- [Requirements domains](requirements_domains.md) – Canonical list of requirement/spec domains (e.g. for Spec IDs and REQ IDs); single source of truth for domain tags and requirements file mapping.
- ~~[Writing and validation](writing_and_validation.md) – Markdown and tech-spec conventions, documentation validation pipeline (`make docs-check`), and Make target reference.~~ <!-- TODO: Remove this line later. -->

## Relation to Other Docs

- **Technical specs**: [`../tech_specs/`](../tech_specs/) – follow these standards where applicable.
- **Requirements**: [`../requirements/`](../requirements/) – same for requirement docs.
- **Validation**: The [Documentation Checks](https://github.com/novus-engine/novuspack/actions/workflows/docs-check.yml) workflow and `make docs-check` enforce some of these conventions; see the scripts in `scripts/` and the workflow for details.
