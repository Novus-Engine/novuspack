# Canonical Requirements Domains

This document is the **single source of truth** for NovusPack requirements domains.
Spec IDs (`NP.<Domain>.<Path>`) and requirement IDs (`REQ-<DOMAIN>-NNN`) MUST use only domains listed here.

## Domain tag convention

Domain tags use a **fixed-length pattern**: exactly **6 uppercase letters** (A–Z), no underscores.
All canonical domain tags in this document follow this pattern.

## Domain list

| Domain tag        | Requirements file                                                        | Scope                                          |
| ----------------- | ------------------------------------------------------------------------ | ---------------------------------------------- |
| [APIBAS](#apibas) | [basic_ops.md](../requirements/basic_ops.md)                             | Basic operations API (create, open, lifecycle) |
| [COMPRS](#comprs) | [compression.md](../requirements/compression.md)                         | Package compression API                        |
| [COREPK](#corepk) | [core.md](../requirements/core.md)                                       | Core API (read/write contracts, paths)         |
| [CRYPTO](#crypto) | [security_encryption.md](../requirements/security_encryption.md)         | Security and encryption                        |
| [DEDUPL](#dedupl) | [dedup.md](../requirements/dedup.md)                                     | Multi-layer deduplication                      |
| [FILEFM](#filefm) | [file_format.md](../requirements/file_format.md)                         | File format (header, layout)                   |
| [FILEMG](#filemg) | [file_mgmt.md](../requirements/file_mgmt.md)                             | File management API (add, remove, entries)     |
| [FILETY](#filety) | [file_types.md](../requirements/file_types.md)                           | File types system                              |
| [GENERI](#generi) | [generics.md](../requirements/generics.md)                               | Generic types and patterns                     |
| [METADT](#metadt) | [metadata.md](../requirements/metadata.md)                               | Package metadata API                           |
| [METSYS](#metsys) | [metadata_system.md](../requirements/metadata_system.md)                 | Metadata system                                |
| [PIPELN](#pipeln) | [transformation_pipeline.md](../requirements/transformation_pipeline.md) | Multi-stage transformation pipeline            |
| [SECURE](#secure) | [security.md](../requirements/security.md)                               | Security validation API                        |
| [SIGNAT](#signat) | [signatures.md](../requirements/signatures.md)                           | Digital signature API                          |
| [STREAM](#stream) | [streaming.md](../requirements/streaming.md)                             | Streaming and buffer management                |
| [TESTNG](#testng) | [testing.md](../requirements/testing.md)                                 | Testing                                        |
| [VALIDT](#validt) | [validation.md](../requirements/validation.md)                           | File validation                                |
| [WRITEP](#writep) | [writing.md](../requirements/writing.md)                                 | Package writing API                            |

## Domain sections

Each section below describes what the domain covers and links to the requirements file.

### APIBAS

Basic operations API: package create, open, lifecycle, structure, and constants.
See [basic_ops.md](../requirements/basic_ops.md).

### COMPRS

Package compression API: compression and decompression behavior and configuration.
See [compression.md](../requirements/compression.md).

### COREPK

Core API: read/write contracts, path normalization and validation, package access.
See [core.md](../requirements/core.md).

### CRYPTO

Security and encryption: encryption, decryption, and related security behavior.
See [security_encryption.md](../requirements/security_encryption.md).

### DEDUPL

Multi-layer deduplication: content and block deduplication.
See [dedup.md](../requirements/dedup.md).

### FILEFM

File format: package header, layout, magic, version, and on-disk structure.
See [file_format.md](../requirements/file_format.md).

### FILEMG

File management API: adding, removing, and managing file entries in the package.
See [file_mgmt.md](../requirements/file_mgmt.md).

### FILETY

File types system: type detection and classification.
See [file_types.md](../requirements/file_types.md).

### GENERI

Generic types and patterns: generic API patterns and type constraints.
See [generics.md](../requirements/generics.md).

### METADT

Package metadata API: metadata read/write and structure.
See [metadata.md](../requirements/metadata.md).

### METSYS

Metadata system: system-wide metadata behavior and configuration.
See [metadata_system.md](../requirements/metadata_system.md).

### PIPELN

Multi-stage transformation pipeline: pipeline configuration and stages.
See [transformation_pipeline.md](../requirements/transformation_pipeline.md).

### SECURE

Security validation API: validation and security checks.
See [security.md](../requirements/security.md).

### SIGNAT

Digital signature API: signing and verification.
See [signatures.md](../requirements/signatures.md).

### STREAM

Streaming and buffer management: streaming I/O and buffers.
See [streaming.md](../requirements/streaming.md).

### TESTNG

Testing: test support and testability requirements.
See [testing.md](../requirements/testing.md).

### VALIDT

File validation: validation rules and behavior.
See [validation.md](../requirements/validation.md).

### WRITEP

Package writing API: writing and persisting package content.
See [writing.md](../requirements/writing.md).

## Legacy domain tags (pre–6-letter)

Requirements and Spec IDs in the repo still use the tags below.
When migrating reqs and specs to 6-letter tags, use this mapping.
Do not use legacy tags for new content.

| Legacy tag  | Canonical (6-letter) |
| ----------- | -------------------- |
| `API_BASIC` | `APIBAS`             |
| `COMPR`     | `COMPRS`             |
| `CORE`      | `COREPK`             |
| `DEDUP`     | `DEDUPL`             |
| `FILEFMT`   | `FILEFM`             |
| `FILEMGMT`  | `FILEMG`             |
| `FILETYPES` | `FILETY`             |
| `GEN`       | `GENERI`             |
| `META`      | `METADT`             |
| `METASYS`   | `METSYS`             |
| `PIPELINE`  | `PIPELN`             |
| `SEC`       | `SECURE`             |
| `SIG`       | `SIGNAT`             |
| `TEST`      | `TESTNG`             |
| `VALID`     | `VALIDT`             |
| `WRITE`     | `WRITEP`             |

`CRYPTO` and `STREAM` were already 6 letters; unchanged.

## Disallowed or legacy domain tags

Do not introduce additional lookalike tags.
Use the canonical 6-letter tag from the [domain list](#domain-list).

| Do not use   | Use instead |
| ------------ | ----------- |
| `FILEFORMAT` | `FILEFM`    |
| `FILETYPE`   | `FILETY`    |
| `GENERIC`    | `GENERI`    |
| `SECURITY`   | `SECURE`    |
| `SEC_ENC`    | `CRYPTO`    |

## Adding or changing domains

1. Update this document (add row to the table, add a domain section, or add a disallowed alias).
2. Add or rename the corresponding requirements file under `docs/requirements/` if adding a new domain.
   New domain tags MUST be 6 uppercase letters.
3. Update any references (e.g.
   [Spec ID pattern](stds_prop1.md), [requirements README](../requirements/README.md)) to point to this document; do not duplicate the full domain list elsewhere.
