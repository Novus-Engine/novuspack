# Package Metadata API Requirements

## Comment Management

- REQ-META-001: Comments are persisted and retrievable. [api_metadata.md#1-comment-management](../tech_specs/api_metadata.md#1-comment-management)
- REQ-META-005: PackageComment methods (Size, WriteTo, ReadFrom, Validate) provide comment operations. [api_metadata.md#11-packagecomment-methods](../tech_specs/api_metadata.md#11-packagecomment-methods)
- REQ-META-006: Comment security validation prevents injection and validates encoding. [api_metadata.md#12-comment-security-validation](../tech_specs/api_metadata.md#12-comment-security-validation)
- REQ-META-010: Signature comment security validation prevents security issues. [api_metadata.md#13-signature-comment-security](../tech_specs/api_metadata.md#13-signature-comment-security)
- REQ-META-056: WriteTo parameters define comment writing interface. [api_metadata.md#111-writeto-parameters](../tech_specs/api_metadata.md#111-writeto-parameters)
- REQ-META-057: ReadFrom parameters define comment reading interface. [api_metadata.md#112-readfrom-parameters](../tech_specs/api_metadata.md#112-readfrom-parameters)
- REQ-META-058: Error conditions define comment operation errors. [api_metadata.md#113-error-conditions](../tech_specs/api_metadata.md#113-error-conditions)
- REQ-META-059: Example usage demonstrates comment operations [type: documentation-only]. [api_metadata.md#114-example-usage](../tech_specs/api_metadata.md#114-example-usage)

## AppID and VendorID Management

- REQ-META-002: AppID and VendorID constraints enforced [type: constraint]. [api_metadata.md#2-appid-management](../tech_specs/api_metadata.md#2-appid-management)
- REQ-META-004: AppID and VendorID management supports set, get, clear, and has operations. [api_metadata.md#2-appid-management](../tech_specs/api_metadata.md#2-appid-management)
- REQ-META-060: VendorID management provides vendor identifier operations. [api_metadata.md#3-vendorid-management](../tech_specs/api_metadata.md#3-vendorid-management)
- REQ-META-061: Combined management provides combined AppID/VendorID operations. [api_metadata.md#4-combined-management](../tech_specs/api_metadata.md#4-combined-management)

## Directory Metadata

- REQ-META-003: Directory metadata follows structure and validation. [api_metadata.md#8-directory-metadata-system](../tech_specs/api_metadata.md#8-directory-metadata-system)
- REQ-META-088: Directory structures define directory metadata structures. [api_metadata.md#81-directory-structures](../tech_specs/api_metadata.md#81-directory-structures)
- REQ-META-089: Directory management methods provide directory operations. [api_metadata.md#82-directory-management-methods](../tech_specs/api_metadata.md#82-directory-management-methods)
- REQ-META-097: Directory association system provides file-directory relationships [type: architectural]. [api_metadata.md#84-directory-association-system](../tech_specs/api_metadata.md#84-directory-association-system)
- REQ-META-098: Association properties define directory association attributes. [api_metadata.md#841-association-properties](../tech_specs/api_metadata.md#841-association-properties)
- REQ-META-099: FileEntry directory properties provide file entry directory information. [api_metadata.md#841-fileentry-directory-properties](../tech_specs/api_metadata.md#841-fileentry-directory-properties)
- REQ-META-100: DirectoryEntry filesystem properties provide directory filesystem information. [api_metadata.md#842-directoryentry-filesystem-properties](../tech_specs/api_metadata.md#842-directoryentry-filesystem-properties)
- REQ-META-101: Association management provides directory association operations. [api_metadata.md#843-association-management](../tech_specs/api_metadata.md#843-association-management)
- REQ-META-102: File-directory association provides file-directory relationships. [api_metadata.md#85-file-directory-association](../tech_specs/api_metadata.md#85-file-directory-association)

## Package Information

- REQ-META-007: GetPackageInfo returns comprehensive package information. [api_metadata.md#7-package-information-structures](../tech_specs/api_metadata.md#7-package-information-structures)
- REQ-META-008: GetSecurityStatus returns current security status. [api_metadata.md#7-package-information-structures](../tech_specs/api_metadata.md#7-package-information-structures)
- REQ-META-009: RefreshPackageInfo refreshes package information cache. [api_metadata.md#7-package-information-structures](../tech_specs/api_metadata.md#7-package-information-structures)
- REQ-META-084: PackageInfo structure provides package information. [api_metadata.md#71-packageinfo-structure](../tech_specs/api_metadata.md#71-packageinfo-structure)
- REQ-META-085: SignatureInfo structure provides signature information. [api_metadata.md#72-signatureinfo-structure](../tech_specs/api_metadata.md#72-signatureinfo-structure)
- REQ-META-086: SecurityStatus structure provides security status information. [api_metadata.md#73-securitystatus-structure](../tech_specs/api_metadata.md#73-securitystatus-structure)
- REQ-META-087: Package information methods provide package information access. [api_metadata.md#74-package-information-methods](../tech_specs/api_metadata.md#74-package-information-methods)

## Special Metadata File Types

- REQ-META-062: Special metadata file types define special file classifications [type: architectural]. [api_metadata.md#5-special-metadata-file-types](../tech_specs/api_metadata.md#5-special-metadata-file-types)
- REQ-META-063: Package metadata file type 65000 defines metadata file type. [api_metadata.md#51-package-metadata-file-type-65000](../tech_specs/api_metadata.md#51-package-metadata-file-type-65000)
- REQ-META-064: Package manifest file type 65001 defines manifest file type. [api_metadata.md#52-package-manifest-file-type-65001](../tech_specs/api_metadata.md#52-package-manifest-file-type-65001)
- REQ-META-065: Package index file type 65002 defines index file type. [api_metadata.md#53-package-index-file-type-65002](../tech_specs/api_metadata.md#53-package-index-file-type-65002)
- REQ-META-066: Package signature file type 65003 defines signature file type. [api_metadata.md#54-package-signature-file-type-65003](../tech_specs/api_metadata.md#54-package-signature-file-type-65003)
- REQ-META-067: Special file management provides special file operations. [api_metadata.md#55-special-file-management](../tech_specs/api_metadata.md#55-special-file-management)
- REQ-META-068: Special file data structures define special file formats. [api_metadata.md#551-special-file-data-structures](../tech_specs/api_metadata.md#551-special-file-data-structures)
- REQ-META-090: Special metadata file management provides special file operations. [api_metadata.md#83-special-metadata-file-management](../tech_specs/api_metadata.md#83-special-metadata-file-management)
- REQ-META-091: Special file requirements define special file specifications [type: constraint]. [api_metadata.md#831-special-file-requirements](../tech_specs/api_metadata.md#831-special-file-requirements)
- REQ-META-092: File type requirements define special file type rules [type: constraint]. [api_metadata.md#8311-file-type-requirements](../tech_specs/api_metadata.md#8311-file-type-requirements)
- REQ-META-093: Special file types define special file classifications. [api_metadata.md#8312-special-file-types](../tech_specs/api_metadata.md#8312-special-file-types)
- REQ-META-094: Package header flags define special file flags. [api_metadata.md#8313-package-header-flags](../tech_specs/api_metadata.md#8313-package-header-flags)
- REQ-META-095: FileEntry requirements define special file entry rules [type: constraint]. [api_metadata.md#8314-fileentry-requirements](../tech_specs/api_metadata.md#8314-fileentry-requirements)
- REQ-META-096: Implementation details provide special file implementation information [type: architectural]. [api_metadata.md#832-implementation-details](../tech_specs/api_metadata.md#832-implementation-details)

## Metadata-Only Packages

- REQ-META-069: Metadata-only packages support packages without file data. [api_metadata.md#6-metadata-only-packages](../tech_specs/api_metadata.md#6-metadata-only-packages)
- REQ-META-070: Metadata-only package definition defines metadata-only package structure [type: architectural]. [api_metadata.md#61-metadata-only-package-definition](../tech_specs/api_metadata.md#61-metadata-only-package-definition)
- REQ-META-071: Valid use cases define metadata-only package use cases [type: documentation-only]. [api_metadata.md#62-valid-use-cases](../tech_specs/api_metadata.md#62-valid-use-cases)
- REQ-META-072: Package catalogs and registries use metadata-only packages [type: documentation-only]. [api_metadata.md#621-package-catalogs-and-registries](../tech_specs/api_metadata.md#621-package-catalogs-and-registries)
- REQ-META-073: Configuration and schema packages use metadata-only packages [type: documentation-only]. [api_metadata.md#622-configuration-and-schema-packages](../tech_specs/api_metadata.md#622-configuration-and-schema-packages)
- REQ-META-074: Package management operations use metadata-only packages [type: documentation-only]. [api_metadata.md#623-package-management-operations](../tech_specs/api_metadata.md#623-package-management-operations)
- REQ-META-075: Development and build tools use metadata-only packages [type: documentation-only]. [api_metadata.md#624-development-and-build-tools](../tech_specs/api_metadata.md#624-development-and-build-tools)
- REQ-META-076: Security considerations define metadata-only package security [type: documentation-only]. [api_metadata.md#63-security-considerations](../tech_specs/api_metadata.md#63-security-considerations)
- REQ-META-077: Signature validation validates metadata-only package signatures. [api_metadata.md#631-signature-validation](../tech_specs/api_metadata.md#631-signature-validation)
- REQ-META-078: Trust and verification provide metadata-only package trust mechanisms [type: documentation-only]. [api_metadata.md#632-trust-and-verification](../tech_specs/api_metadata.md#632-trust-and-verification)
- REQ-META-079: Package integrity ensures metadata-only package integrity. [api_metadata.md#633-package-integrity](../tech_specs/api_metadata.md#633-package-integrity)
- REQ-META-080: Attack vectors define metadata-only package threats [type: documentation-only]. [api_metadata.md#634-attack-vectors](../tech_specs/api_metadata.md#634-attack-vectors)
- REQ-META-081: Metadata-only package API provides metadata-only operations. [api_metadata.md#64-metadata-only-package-api](../tech_specs/api_metadata.md#64-metadata-only-package-api)
- REQ-META-082: Metadata-only package validation validates metadata-only packages. [api_metadata.md#641-metadata-only-package-validation](../tech_specs/api_metadata.md#641-metadata-only-package-validation)
- REQ-META-083: Enhanced security requirements define metadata-only package security [type: documentation-only]. [api_metadata.md#642-enhanced-security-requirements](../tech_specs/api_metadata.md#642-enhanced-security-requirements)

## Tag System

- REQ-META-015: Tag storage format defines tag encoding structure [type: architectural]. [metadata.md#11-tag-storage-format](../tech_specs/metadata.md#11-tag-storage-format)
- REQ-META-016: Tag structure defines tag entry format. [metadata.md#111-tag-structure](../tech_specs/metadata.md#111-tag-structure)
- REQ-META-017: Tag value types define supported value types. [metadata.md#12-tag-value-types](../tech_specs/metadata.md#12-tag-value-types)
- REQ-META-018: Basic types provide string, integer, float, and boolean support. [metadata.md#121-basic-types](../tech_specs/metadata.md#121-basic-types)
- REQ-META-019: Structured data provides JSON, YAML, and string list support. [metadata.md#122-structured-data](../tech_specs/metadata.md#122-structured-data)
- REQ-META-020: Identifiers provide UUID, hash, and version support. [metadata.md#123-identifiers](../tech_specs/metadata.md#123-identifiers)
- REQ-META-021: Time provides timestamp support. [metadata.md#124-time](../tech_specs/metadata.md#124-time)
- REQ-META-022: Network/Communication provides URL and email support. [metadata.md#125-networkcommunication](../tech_specs/metadata.md#125-networkcommunication)
- REQ-META-023: File system provides path and MIME type support. [metadata.md#126-file-system](../tech_specs/metadata.md#126-file-system)
- REQ-META-024: Localization provides language code support. [metadata.md#127-localization](../tech_specs/metadata.md#127-localization)
- REQ-META-025: NovusPack special files provides metadata file reference support. [metadata.md#128-novuspack-special-files](../tech_specs/metadata.md#128-novuspack-special-files)
- REQ-META-026: Reserved value types are reserved for future use. [metadata.md#129-reserved](../tech_specs/metadata.md#129-reserved)
- REQ-META-027: Directory metadata system provides directory-based tag inheritance [type: architectural]. [metadata.md#13-directory-metadata-system](../tech_specs/metadata.md#13-directory-metadata-system)
- REQ-META-028: Directory metadata file defines directory metadata storage. [metadata.md#131-directory-metadata-file](../tech_specs/metadata.md#131-directory-metadata-file)
- REQ-META-029: Directory entry structure defines directory entry format. [metadata.md#132-directory-entry-structure](../tech_specs/metadata.md#132-directory-entry-structure)
- REQ-META-030: Tag inheritance rules define tag inheritance behavior. [metadata.md#133-tag-inheritance-rules](../tech_specs/metadata.md#133-tag-inheritance-rules)
- REQ-META-031: Inheritance examples demonstrate inheritance patterns [type: documentation-only]. [metadata.md#134-inheritance-examples](../tech_specs/metadata.md#134-inheritance-examples)
- REQ-META-032: Example 1 basic directory inheritance demonstrates basic inheritance [type: documentation-only]. [metadata.md#1341-example-1-basic-directory-inheritance](../tech_specs/metadata.md#1341-example-1-basic-directory-inheritance)
- REQ-META-033: Example 2 priority-based override demonstrates priority rules [type: documentation-only]. [metadata.md#1342-example-2-priority-based-override](../tech_specs/metadata.md#1342-example-2-priority-based-override)
- REQ-META-034: Example 3 inheritance disabled demonstrates disabled inheritance [type: documentation-only]. [metadata.md#1343-example-3-inheritance-disabled](../tech_specs/metadata.md#1343-example-3-inheritance-disabled)
- REQ-META-035: Tag validation provides tag validation rules [type: constraint]. [metadata.md#14-tag-validation](../tech_specs/metadata.md#14-tag-validation)
- REQ-META-036: Per-file tags usage examples demonstrate tag usage [type: documentation-only]. [metadata.md#15-per-file-tags-usage-examples](../tech_specs/metadata.md#15-per-file-tags-usage-examples)
- REQ-META-037: Texture file tagging demonstrates texture metadata [type: documentation-only]. [metadata.md#151-texture-file-tagging](../tech_specs/metadata.md#151-texture-file-tagging)

## Metadata File Format

- REQ-META-038: Metadata file requirements define metadata file specifications [type: architectural]. [metadata.md#21-metadata-file-requirements](../tech_specs/metadata.md#21-metadata-file-requirements)
- REQ-META-039: YAML schema structure defines metadata schema format [type: architectural]. [metadata.md#22-yaml-schema-structure](../tech_specs/metadata.md#22-yaml-schema-structure)
- REQ-META-040: Package metadata schema v1.0 defines schema version. [metadata.md#221-package-metadata-schema-v10](../tech_specs/metadata.md#221-package-metadata-schema-v10)
- REQ-META-041: Package information example demonstrates package metadata [type: documentation-only]. [metadata.md#222-package-information-example](../tech_specs/metadata.md#222-package-information-example)
- REQ-META-042: Metadata file API provides metadata file operations. [metadata.md#23-metadata-file-api](../tech_specs/metadata.md#23-metadata-file-api)
- REQ-META-043: Package metadata example demonstrates package metadata structure [type: documentation-only]. [metadata.md#24-package-metadata-example](../tech_specs/metadata.md#24-package-metadata-example)

## Metadata Examples

- REQ-META-044: Asset metadata provides asset tagging support. [metadata.md#asset-metadata](../tech_specs/metadata.md#asset-metadata)
- REQ-META-045: Asset metadata example demonstrates asset tagging [type: documentation-only]. [metadata.md#asset-metadata-example](../tech_specs/metadata.md#asset-metadata-example)
- REQ-META-046: Audio file tagging demonstrates audio metadata [type: documentation-only]. [metadata.md#audio-file-tagging](../tech_specs/metadata.md#audio-file-tagging)
- REQ-META-047: Custom metadata provides custom tag support. [metadata.md#custom-metadata](../tech_specs/metadata.md#custom-metadata)
- REQ-META-048: Custom metadata example demonstrates custom tags [type: documentation-only]. [metadata.md#custom-metadata-example](../tech_specs/metadata.md#custom-metadata-example)
- REQ-META-049: Directory tagging demonstrates directory metadata [type: documentation-only]. [metadata.md#directory-tagging](../tech_specs/metadata.md#directory-tagging)
- REQ-META-050: File search by tags demonstrates tag-based search [type: documentation-only]. [metadata.md#file-search-by-tags](../tech_specs/metadata.md#file-search-by-tags)
- REQ-META-051: Game-specific metadata demonstrates game metadata [type: documentation-only]. [metadata.md#game-specific-metadata](../tech_specs/metadata.md#game-specific-metadata)
- REQ-META-052: Game-specific metadata example demonstrates game tags [type: documentation-only]. [metadata.md#game-specific-metadata-example](../tech_specs/metadata.md#game-specific-metadata-example)
- REQ-META-053: Package information provides package metadata. [metadata.md#package-information](../tech_specs/metadata.md#package-information)
- REQ-META-054: Security metadata provides security tagging. [metadata.md#security-metadata](../tech_specs/metadata.md#security-metadata)
- REQ-META-055: Security metadata example demonstrates security tags [type: documentation-only]. [metadata.md#security-metadata-example](../tech_specs/metadata.md#security-metadata-example)

## Context Integration

- REQ-META-011: All metadata methods accept context.Context and respect cancellation/timeout [type: constraint]. [api_core.md#02-context-integration](../tech_specs/api_core.md#02-context-integration)
- REQ-META-014: Context cancellation during metadata operations returns structured context error. [api_core.md#02-context-integration](../tech_specs/api_core.md#02-context-integration)

## Validation

- REQ-META-012: Comment parameters validated (UTF-8 encoding, length limits, no injection patterns) [type: constraint]. [api_metadata.md#12-comment-security-validation](../tech_specs/api_metadata.md#12-comment-security-validation)
- REQ-META-013: AppID/VendorID parameters validated (format, length, allowed characters) [type: constraint]. [api_metadata.md#2-appid-management](../tech_specs/api_metadata.md#2-appid-management)
