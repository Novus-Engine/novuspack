# Package Metadata API Requirements

## Comment Management

- REQ-META-001: Comments are persisted and retrievable. [api_metadata.md#1-comment-management](../tech_specs/api_metadata.md#1-comment-management)
- REQ-META-147: PackageComment structure defines comment length, UTF-8 content, and reserved bytes for persisted comment data [type: architectural]. [api_metadata.md#12-packagecomment-structure](../tech_specs/api_metadata.md#12-packagecomment-structure)
- REQ-META-005: PackageComment methods (Size, WriteTo, ReadFrom, Validate) provide comment operations. [api_metadata.md#13-packagecomment-methods](../tech_specs/api_metadata.md#13-packagecomment-methods) (Exception: also [api_metadata.md#134-packagecommentvalidate-method](../tech_specs/api_metadata.md#134-packagecommentvalidate-method) for coverage.)
- REQ-META-006: Comment security validation prevents injection and validates encoding. [api_metadata.md#14-comment-security-validation](../tech_specs/api_metadata.md#14-comment-security-validation)
- REQ-META-010: Signature comment security validation prevents security issues. [api_metadata.md#15-signature-comment-security](../tech_specs/api_metadata.md#15-signature-comment-security)
- REQ-META-056: WriteTo parameters define comment writing interface. [api_metadata.md#1341-writeto-parameters](../tech_specs/api_metadata.md#1341-writeto-parameters)
- REQ-META-057: ReadFrom parameters define comment reading interface. [api_metadata.md#1342-readfrom-parameters](../tech_specs/api_metadata.md#1342-readfrom-parameters)
- REQ-META-058: Error conditions define comment operation errors. [api_metadata.md#1344-error-conditions](../tech_specs/api_metadata.md#1344-error-conditions)
- REQ-META-059: Example usage demonstrates comment operations [type: documentation-only]. [api_metadata.md#1345-example-usage](../tech_specs/api_metadata.md#1345-example-usage)
- REQ-META-113: NewPackageComment creates a new PackageComment with proper initialization. [api_metadata.md#135-newpackagecomment-function](../tech_specs/api_metadata.md#135-newpackagecomment-function)

## AppID and VendorID Management

- REQ-META-002: AppID and VendorID constraints enforced [type: constraint]. [api_metadata.md#2-appid-management](../tech_specs/api_metadata.md#2-appid-management)
- REQ-META-004: AppID and VendorID management supports set, get, clear, and has operations. [api_metadata.md#2-appid-management](../tech_specs/api_metadata.md#2-appid-management)
- REQ-META-060: VendorID management provides vendor identifier operations. [api_metadata.md#3-vendorid-management](../tech_specs/api_metadata.md#3-vendorid-management)
- REQ-META-061: Combined management provides combined AppID/VendorID operations. [api_metadata.md#4-combined-management](../tech_specs/api_metadata.md#4-combined-management)

## Path Metadata

- REQ-META-003: Path metadata follows structure and validation. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-117: All path metadata paths MUST have leading slash [type: constraint]. [api_metadata.md#81-pathmetadata-structures](../tech_specs/api_metadata.md#81-pathmetadata-structures)
- REQ-META-118: PathMetadataEntry paths stored with leading slash for package root reference [type: constraint]. [metadata.md#132-pathmetadataentry-structure](../tech_specs/metadata.md#132-pathmetadataentry-structure)
- REQ-META-119: Path metadata display operations MUST strip leading slash for user output [type: constraint]. [metadata.md#133-tag-inheritance-rules](../tech_specs/metadata.md#133-tag-inheritance-rules)
- REQ-META-121: AddDirectoryMetadata creates directory metadata entries without adding file content [type: constraint]. [api_metadata.md#82-pathmetadata-management-methods](../tech_specs/api_metadata.md#82-pathmetadata-management-methods)
- REQ-META-088: Path structures define path metadata structures. [api_metadata.md#81-pathmetadata-structures](../tech_specs/api_metadata.md#81-pathmetadata-structures)
- REQ-META-089: Path management methods provide path operations. [api_metadata.md#82-pathmetadata-management-methods](../tech_specs/api_metadata.md#82-pathmetadata-management-methods)
- REQ-META-097: Path association system provides file-path relationships [type: architectural]. [api_metadata.md#84-path-association-system](../tech_specs/api_metadata.md#84-path-association-system)
- REQ-META-098: Association properties define path association attributes. [api_metadata.md#841-fileentry-path-properties](../tech_specs/api_metadata.md#841-fileentry-path-properties)
- REQ-META-099: FileEntry path properties provide file entry path information. [api_metadata.md#841-fileentry-path-properties](../tech_specs/api_metadata.md#841-fileentry-path-properties)
- REQ-META-100: PathMetadataEntry filesystem properties provide path filesystem information. [api_metadata.md#842-pathmetadataentry-filesystem-properties](../tech_specs/api_metadata.md#842-pathmetadataentry-filesystem-properties)
- REQ-META-124: PathFileSystem.IsExecutable field tracks file execute permissions. [api_metadata.md#842-pathmetadataentry-filesystem-properties](../tech_specs/api_metadata.md#842-pathmetadataentry-filesystem-properties)
- REQ-META-101: Association management provides path association operations. [api_metadata.md#843-association-management](../tech_specs/api_metadata.md#843-association-management)
- REQ-META-102: File-path association provides file-path relationships. [api_metadata.md#85-file-path-association](../tech_specs/api_metadata.md#85-file-path-association)
- REQ-META-103: Path tag operations use standalone functions due to Go generic method limitations [type: architectural]. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-104: GetPathMetaTags returns all tags as typed tags for a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-105: GetPathMetaTagsByType returns all tags of a specific type for a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-106: AddPathMetaTags adds multiple new tags with type safety to a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-107: SetPathMetaTags updates existing tags from a slice of typed tags for a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-108: GetPathMetaTag retrieves a type-safe tag by key from a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-109: AddPathMetaTag adds a new tag with type safety to a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-110: SetPathMetaTag updates an existing tag with type safety for a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-111: RemovePathMetaTag removes a tag by key from a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-112: HasPathMetaTag checks if a tag with the specified key exists on a PathMetadataEntry. [api_metadata.md#8-pathmetadata-system](../tech_specs/api_metadata.md#8-pathmetadata-system)
- REQ-META-152: PathMetadataEntry tag management defines tag storage and typed tag operations for path metadata entries. [api_metadata.md#817-pathmetadataentry-tag-management](../tech_specs/api_metadata.md#817-pathmetadataentry-tag-management)
- REQ-META-125: PathMetadataEntry stores DestPath and DestPathWin fields for persistent destination extraction overrides. [api_metadata.md#81-pathmetadata-structures](../tech_specs/api_metadata.md#81-pathmetadata-structures)
- REQ-META-126: DestPath and DestPathWin support relative paths (resolved from default extraction directory) and absolute paths. [api_metadata.md#81-pathmetadata-structures](../tech_specs/api_metadata.md#81-pathmetadata-structures)
- REQ-META-127: DestPathWin is used for Windows-specific destination paths; if only DestPath is absolute on Windows, root is treated as C:\\. [api_metadata.md#81-pathmetadata-structures](../tech_specs/api_metadata.md#81-pathmetadata-structures)
- REQ-META-128: SetDestPath sets destination extraction overrides for a stored path, creating PathMetadataEntry if missing. [api_metadata.md#8216-packagesetdestpath-method](../tech_specs/api_metadata.md#8216-packagesetdestpath-method)
- REQ-META-129: SetDestPath accepts DestPathOverride struct with optional DestPath and DestPathWin pointers. [api_metadata.md#8111-destpathoverride-structure](../tech_specs/api_metadata.md#8111-destpathoverride-structure)
- REQ-META-130: SetDestPathTyped is a generic helper that accepts string or map[string]string input and converts to DestPathOverride. [api_metadata.md#8216-packagesetdestpath-method](../tech_specs/api_metadata.md#8216-packagesetdestpath-method)
- REQ-META-131: SetDestPath normalizes storedPath by prefixing leading slash if missing before matching. [api_metadata.md#8216-packagesetdestpath-method](../tech_specs/api_metadata.md#8216-packagesetdestpath-method)
- REQ-META-132: SetDestPath parses string input to determine Windows-only paths (drive letter or UNC) and stores in appropriate field. [api_metadata.md#8216-packagesetdestpath-method](../tech_specs/api_metadata.md#8216-packagesetdestpath-method)

- REQ-META-153: Directory metadata convenience methods provide helpers for directory path metadata operations. [api_metadata.md#8218-directory-metadata-convenience-methods](../tech_specs/api_metadata.md#8218-directory-metadata-convenience-methods)
- REQ-META-154: UpdateDirectoryMetadata updates directory path metadata without modifying files. [api_metadata.md#82111-packageupdatedirectorymetadata-method](../tech_specs/api_metadata.md#82111-packageupdatedirectorymetadata-method)

### Path Information Queries and Associations

- REQ-META-133: GetDirectoryCount returns total directory count and error. [api_metadata.md#8251-getdirectorycount-returns](../tech_specs/api_metadata.md#8251-getdirectorycount-returns)
- REQ-META-134: GetPathHierarchy returns parent => children mapping and error. [api_metadata.md#8261-getpathhierarchy-returns](../tech_specs/api_metadata.md#8261-getpathhierarchy-returns)
- REQ-META-135: AssociateFileWithPath links a FileEntry to its PathMetadataEntry and sets ParentPath for hierarchy traversal without failing on missing parents. [api_metadata.md#827-packageassociatefilewithpath-method](../tech_specs/api_metadata.md#827-packageassociatefilewithpath-method)
- REQ-META-136: UpdateFilePathAssociations rebuilds file-path associations for all files and establishes ParentPath links for path hierarchy traversal without failing on missing parents. [api_metadata.md#828-packageupdatefilepathassociations-method](../tech_specs/api_metadata.md#828-packageupdatefilepathassociations-method)

- REQ-META-137: UpdateSpecialMetadataFlags updates header flags and PackageInfo fields based on presence of special metadata files, per-file tags, and extended attributes. [api_metadata.md#8352-packageupdatespecialmetadataflags-method](../tech_specs/api_metadata.md#8352-packageupdatespecialmetadataflags-method)
- REQ-META-138: Data flow architecture defines header flags and PackageInfo source-of-truth across open, in-memory operations, and write paths [type: architectural]. [api_metadata.md#8353-data-flow-architecture](../tech_specs/api_metadata.md#8353-data-flow-architecture)

- REQ-META-155: Path information query methods provide path information query operations (GetPathInfo, ListPaths, ListDirectories, counts, and hierarchy). [api_metadata.md#path-information-query-methods](../tech_specs/api_metadata.md#path-information-query-methods)
- REQ-META-156: Path association methods provide file-path association operations (AssociateFileWithPath and UpdateFilePathAssociations). [api_metadata.md#path-association-methods](../tech_specs/api_metadata.md#path-association-methods)

- REQ-META-157: GetPathInfo returns PathInfo and error. [api_metadata.md#8222-getpathinfo-returns](../tech_specs/api_metadata.md#8222-getpathinfo-returns)
- REQ-META-158: ListPaths returns slice of PathInfo and error. [api_metadata.md#8231-listpaths-returns](../tech_specs/api_metadata.md#8231-listpaths-returns)
- REQ-META-159: ListDirectories returns slice of PathInfo and error. [api_metadata.md#8241-listdirectories-returns](../tech_specs/api_metadata.md#8241-listdirectories-returns)

### Symlink Validation and Conversion Rules

- REQ-META-146: ValidateSymlinkPaths validation details define that source and target paths are validated for package boundaries. [api_metadata.md#85111-validatesymlinkpaths-method-validation-details](../tech_specs/api_metadata.md#85111-validatesymlinkpaths-method-validation-details)
- REQ-META-139: Symlink validation rules require package-relative paths, package-root boundary enforcement, target existence, and defined structured error types. [api_metadata.md#85112-validation-rules](../tech_specs/api_metadata.md#85112-validation-rules)
- REQ-META-140: TargetExists checks whether a path exists as a FileEntry or directory PathMetadataEntry. [api_metadata.md#85113-targetexists-method-validation-details](../tech_specs/api_metadata.md#85113-targetexists-method-validation-details)
- REQ-META-141: Symlink validation use cases document target checks and validation workflows [type: documentation-only]. [api_metadata.md#85114-use-cases](../tech_specs/api_metadata.md#85114-use-cases)
- REQ-META-142: ValidatePathWithinPackageRoot validates and normalizes a path and prevents package root boundary escapes. [api_metadata.md#85115-validatepathwithinpackageroot-method-validation-details](../tech_specs/api_metadata.md#85115-validatepathwithinpackageroot-method-validation-details)
- REQ-META-143: ValidatePathWithinPackageRoot validation rules define required format, normalization, and structured error types. [api_metadata.md#85116-validatepathwithinpackageroot-validation-rules](../tech_specs/api_metadata.md#85116-validatepathwithinpackageroot-validation-rules)
- REQ-META-144: Symlink creation conversion process selects a primary path, creates symlink entries, creates a symlink PathMetadataEntry, and updates the FileEntry to a single primary path. [api_metadata.md#85121-conversion-process](../tech_specs/api_metadata.md#85121-conversion-process)
- REQ-META-145: Symlink target existence validation requires targets to exist before symlink creation and returns ErrTypeNotFound when missing. [api_metadata.md#85123-target-existence-validation](../tech_specs/api_metadata.md#85123-target-existence-validation)

## Package Information

- REQ-META-007: GetPackageInfo returns comprehensive package information. [api_metadata.md#7-package-information-structures](../tech_specs/api_metadata.md#7-package-information-structures)
- REQ-META-008: GetSecurityStatus returns current security status. [api_metadata.md#7-package-information-structures](../tech_specs/api_metadata.md#7-package-information-structures)
- REQ-META-009: RefreshPackageInfo refreshes package information cache. [api_metadata.md#7-package-information-structures](../tech_specs/api_metadata.md#7-package-information-structures)
- REQ-META-084: PackageInfo structure provides package information. [api_metadata.md#71-packageinfo-structure](../tech_specs/api_metadata.md#71-packageinfo-structure)
- REQ-META-149: PackageInfo scope and exclusions define which files and sizes are counted (excluding special metadata files). [api_metadata.md#711-packageinfo-scope-and-exclusions](../tech_specs/api_metadata.md#711-packageinfo-scope-and-exclusions)
- REQ-META-150: PackageInfo is the source of truth for in-memory package state and is used to update header flags for write operations [type: architectural]. [api_metadata.md#713-packageinfo-as-source-of-truth](../tech_specs/api_metadata.md#713-packageinfo-as-source-of-truth)
- REQ-META-085: SignatureInfo structure provides signature information. [api_metadata.md#72-signatureinfo-structure](../tech_specs/api_metadata.md#72-signatureinfo-structure)
- REQ-META-086: SecurityStatus structure provides security status information. [api_metadata.md#73-securitystatus-structure](../tech_specs/api_metadata.md#73-securitystatus-structure)
- REQ-META-087: Package information methods provide package information access. [api_metadata.md#74-package-information-methods](../tech_specs/api_metadata.md#74-package-information-methods)
- REQ-META-114: NewPackageInfo creates a new PackageInfo with default values. [api_metadata.md#712-newpackageinfo-function](../tech_specs/api_metadata.md#712-newpackageinfo-function)
- REQ-META-115: PackageInfo.FromHeader synchronizes PackageInfo fields from a PackageHeader. [api_metadata.md#714-packageinfofromheader-method](../tech_specs/api_metadata.md#714-packageinfofromheader-method)
- REQ-META-151: PackageHeader structure defines package header fields and flags for PackageInfo synchronization [type: architectural]. [api_metadata.md#715-packageheader-structure](../tech_specs/api_metadata.md#715-packageheader-structure)
- REQ-META-116: PackageHeader.ToHeader synchronizes PackageHeader fields from a PackageInfo. [api_metadata.md#716-packageheadertoheader-method](../tech_specs/api_metadata.md#716-packageheadertoheader-method)

## Special Metadata File Types

- REQ-META-062: Special metadata file types define special file classifications [type: architectural]. [api_metadata.md#5-special-metadata-file-types](../tech_specs/api_metadata.md#5-special-metadata-file-types)
- REQ-META-063: Package metadata file type 65000 defines metadata file type. [api_metadata.md#51-package-metadata-file-type-65000](../tech_specs/api_metadata.md#51-package-metadata-file-type-65000) (Exception: also [api_metadata.md#515-packagehasmetadatafile-method](../tech_specs/api_metadata.md#515-packagehasmetadatafile-method) for coverage.)
- REQ-META-064: Package manifest file type 65001 defines manifest file type. [api_metadata.md#52-package-manifest-file-type-65001](../tech_specs/api_metadata.md#52-package-manifest-file-type-65001) (Exception: also [api_metadata.md#525-packagehasmanifestfile-method](../tech_specs/api_metadata.md#525-packagehasmanifestfile-method) for coverage.)
- REQ-META-065: Package index file type 65002 defines index file type. [api_metadata.md#53-package-index-file-type-65002](../tech_specs/api_metadata.md#53-package-index-file-type-65002) (Exception: also [api_metadata.md#535-packagehasindexfile-method](../tech_specs/api_metadata.md#535-packagehasindexfile-method) for coverage.)
- REQ-META-066: Package signature file type 65003 defines signature file type. [api_metadata.md#54-package-signature-file-type-65003](../tech_specs/api_metadata.md#54-package-signature-file-type-65003) (Exception: also [api_metadata.md#545-packagehassignaturefile-method](../tech_specs/api_metadata.md#545-packagehassignaturefile-method) for coverage.)
- REQ-META-067: Special file management provides special file operations. [api_metadata.md#55-special-file-management](../tech_specs/api_metadata.md#55-special-file-management)
- REQ-META-068: Special file data structures define special file formats. [api_metadata.md#555-special-file-data-structures](../tech_specs/api_metadata.md#555-special-file-data-structures)
- REQ-META-090: Special metadata file management provides special file operations. [api_metadata.md#83-special-metadata-file-management](../tech_specs/api_metadata.md#83-special-metadata-file-management)
- REQ-META-091: Special file requirements define special file specifications [type: constraint]. [api_metadata.md#831-special-file-type-requirements](../tech_specs/api_metadata.md#831-special-file-type-requirements)
- REQ-META-092: File type requirements define special file type rules [type: constraint]. [api_metadata.md#831-special-file-type-requirements](../tech_specs/api_metadata.md#831-special-file-type-requirements)
- REQ-META-093: Special file types define special file classifications. [api_metadata.md#832-special-file-types](../tech_specs/api_metadata.md#832-special-file-types)
- REQ-META-094: Package header flags define special file flags. [api_metadata.md#833-packageheader-flags](../tech_specs/api_metadata.md#833-packageheader-flags)
- REQ-META-095: FileEntry requirements define special file entry rules [type: constraint]. [api_metadata.md#834-fileentry-requirements](../tech_specs/api_metadata.md#834-fileentry-requirements)
- REQ-META-096: Implementation details provide special file implementation information [type: architectural]. [api_metadata.md#835-implementation-details](../tech_specs/api_metadata.md#835-implementation-details)

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
- REQ-META-082: Metadata-only package validation validates metadata-only packages. [api_metadata.md#645-packagevalidatemetadataonlypackage-method](../tech_specs/api_metadata.md#645-packagevalidatemetadataonlypackage-method)
- REQ-META-148: Write operations automatically detect metadata-only status, set header Bit 7 when FileCount = 0, and synchronize PackageInfo.IsMetadataOnly before writing [type: constraint]. [api_metadata.md#646-write-operation-requirements](../tech_specs/api_metadata.md#646-write-operation-requirements)
- REQ-META-083: Enhanced security requirements define metadata-only package security [type: documentation-only]. [api_metadata.md#647-security-considerations-for-metadata-only-packages](../tech_specs/api_metadata.md#647-security-considerations-for-metadata-only-packages)

## Tag System

- REQ-META-122: Per-file tags system specification defines tag-based metadata architecture for file entries [type: architectural]. [metadata.md#1-per-file-tags-system-specification](../tech_specs/metadata.md#1-per-file-tags-system-specification)
- REQ-META-015: Tag storage format defines tag encoding structure [type: architectural]. [metadata.md#11-tag-storage-format](../tech_specs/metadata.md#11-tag-storage-format)
- REQ-META-016: Tag structure defines tag entry format. [metadata.md#111-tag-structure](../tech_specs/metadata.md#111-tag-structure)
- REQ-META-017: Tag value types define supported value types. [metadata.md#12-tag-value-types](../tech_specs/metadata.md#12-tag-value-types)
- REQ-META-018: Basic types provide string, integer, float, and boolean support. [metadata.md#121-basic-types](../tech_specs/metadata.md#121-basic-types)
- REQ-META-019: Structured data provides JSON, YAML, and string list support. [metadata.md#122-structured-data](../tech_specs/metadata.md#122-structured-data)
- REQ-META-020: Identifiers provide UUID, hash, and version support. [metadata.md#123-identifier-value-types](../tech_specs/metadata.md#123-identifier-value-types)
- REQ-META-021: Time provides timestamp support. [metadata.md#124-time-value-types](../tech_specs/metadata.md#124-time-value-types)
- REQ-META-022: Network/Communication provides URL and email support. [metadata.md#125-network-and-communication-value-types](../tech_specs/metadata.md#125-network-and-communication-value-types)
- REQ-META-023: File system provides path and MIME type support. [metadata.md#126-file-system](../tech_specs/metadata.md#126-file-system)
- REQ-META-024: Localization provides language code support. [metadata.md#127-localization-value-types](../tech_specs/metadata.md#127-localization-value-types)
- REQ-META-025: NovusPack special files provides metadata file reference support. [metadata.md#128-novuspack-special-files](../tech_specs/metadata.md#128-novuspack-special-files)
- REQ-META-026: Reserved value types are reserved for future use. [metadata.md#129-reserved-value-types](../tech_specs/metadata.md#129-reserved-value-types)
- REQ-META-027: Path metadata system provides path-based tag inheritance [type: architectural]. [metadata.md#13-pathmetadata-system](../tech_specs/metadata.md#13-pathmetadata-system)
- REQ-META-028: Path metadata file defines path metadata storage. [metadata.md#131-pathmetadata-file](../tech_specs/metadata.md#131-pathmetadata-file)
- REQ-META-029: Path entry structure defines path entry format. [metadata.md#132-pathmetadataentry-structure](../tech_specs/metadata.md#132-pathmetadataentry-structure)
- REQ-META-030: Tag inheritance rules define tag inheritance behavior. [metadata.md#133-tag-inheritance-rules](../tech_specs/metadata.md#133-tag-inheritance-rules)
- REQ-META-031: Inheritance examples demonstrate inheritance patterns [type: documentation-only]. [metadata.md#134-inheritance-examples](../tech_specs/metadata.md#134-inheritance-examples)
- REQ-META-032: Example 1 basic path inheritance demonstrates basic inheritance [type: documentation-only]. [metadata.md#1341-example-1-basic-pathinheritance](../tech_specs/metadata.md#1341-example-1-basic-pathinheritance)
- REQ-META-033: Example 2 priority-based override demonstrates priority rules [type: documentation-only]. [metadata.md#1342-example-2-priority-based-override](../tech_specs/metadata.md#1342-example-2-priority-based-override)
- REQ-META-034: Example 3 inheritance disabled demonstrates disabled inheritance [type: documentation-only]. [metadata.md#1343-example-3-inheritance-disabled](../tech_specs/metadata.md#1343-example-3-inheritance-disabled)
- REQ-META-035: Tag validation provides tag validation rules [type: constraint]. [metadata.md#14-tag-validation](../tech_specs/metadata.md#14-tag-validation)
- REQ-META-036: Per-file tags usage examples demonstrate tag usage [type: documentation-only]. [metadata.md#15-per-file-tags-usage-examples](../tech_specs/metadata.md#15-per-file-tags-usage-examples)
- REQ-META-037: Texture file tagging demonstrates texture metadata [type: documentation-only]. [metadata.md#151-texture-file-tagging](../tech_specs/metadata.md#151-texture-file-tagging)

## Metadata File Format

- REQ-META-038: Metadata file requirements define metadata file specifications [type: architectural]. [metadata.md#21-metadata-file-requirements](../tech_specs/metadata.md#21-metadata-file-requirements)
- REQ-META-039: YAML schema structure defines metadata schema format [type: architectural]. [metadata.md#22-yaml-schema-structure](../tech_specs/metadata.md#22-yaml-schema-structure)
- REQ-META-040: Package metadata schema v1.0 defines schema version. [metadata.md#221-package-metadata-schema-v10](../tech_specs/metadata.md#221-package-metadata-schema-v10)
- REQ-META-041: Package information example demonstrates package metadata [type: documentation-only]. [metadata.md#241-package-information-example](../tech_specs/metadata.md#241-package-information-example)
- REQ-META-123: Package metadata file specification defines package-level metadata storage format [type: architectural]. [metadata.md#2-package-metadata-file-specification](../tech_specs/metadata.md#2-package-metadata-file-specification)
- REQ-META-042: Metadata file API provides metadata file operations. [metadata.md#23-metadata-file-api](../tech_specs/metadata.md#23-metadata-file-api)
- REQ-META-043: Package metadata example demonstrates package metadata structure [type: documentation-only]. [metadata.md#24-package-metadata-example](../tech_specs/metadata.md#24-package-metadata-example)

## Metadata Examples

- REQ-META-044: Asset metadata provides asset tagging support. [metadata.md#asset-metadata](../tech_specs/metadata.md#asset-metadata)
- REQ-META-045: Asset metadata example demonstrates asset tagging [type: documentation-only]. [metadata.md#asset-metadata-example](../tech_specs/metadata.md#asset-metadata-example)
- REQ-META-046: Audio file tagging demonstrates audio metadata [type: documentation-only]. [metadata.md#audio-file-tagging](../tech_specs/metadata.md#audio-file-tagging)
- REQ-META-047: Custom metadata provides custom tag support. [metadata.md#custom-metadata](../tech_specs/metadata.md#custom-metadata)
- REQ-META-048: Custom metadata example demonstrates custom tags [type: documentation-only]. [metadata.md#custom-metadata-example](../tech_specs/metadata.md#custom-metadata-example)
- REQ-META-049: Path tagging demonstrates path metadata [type: documentation-only]. [metadata.md#path-tagging](../tech_specs/metadata.md#path-tagging)
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

- REQ-META-012: Comment parameters validated (UTF-8 encoding, length limits, no injection patterns) [type: constraint]. [api_metadata.md#14-comment-security-validation](../tech_specs/api_metadata.md#14-comment-security-validation)
- REQ-META-013: AppID/VendorID parameters validated (format, length, allowed characters) [type: constraint]. [api_metadata.md#2-appid-management](../tech_specs/api_metadata.md#2-appid-management)
