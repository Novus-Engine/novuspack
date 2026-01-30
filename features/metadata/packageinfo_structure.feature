@domain:metadata @m2 @REQ-META-084 @REQ-META-114 @REQ-META-115 @REQ-META-116 @REQ-META-149 @REQ-META-150 @REQ-META-151 @spec(api_metadata.md#71-packageinfo-structure) @spec(api_metadata.md#711-packageinfo-scope-and-exclusions) @spec(api_metadata.md#713-packageinfo-as-source-of-truth) @spec(api_metadata.md#714-packageinfofromheader-method) @spec(api_metadata.md#715-packageheader-structure) @spec(api_metadata.md#716-packageheadertoheader-method)
Feature: PackageInfo Structure

  @REQ-META-114 @happy
  Scenario: NewPackageInfo creates new PackageInfo with default values
    Given NewPackageInfo is called
    Then a PackageInfo is returned
    And all numeric fields are set to 0
    And all boolean fields are set to false
    And string fields are empty
    And slice fields are empty slices
    And time fields are zero time
    And PackageInfo is ready for use

  @REQ-META-084 @happy
  Scenario: PackageInfo structure provides comprehensive package information
    Given a NovusPack package
    When PackageInfo structure is examined
    Then structure contains basic package information
    And structure contains package identity fields
    And structure contains package comment fields
    And structure contains digital signature information
    And structure contains security information
    And structure contains timestamps
    And structure contains package feature flags
    And structure contains package compression information

  @REQ-META-084 @happy
  Scenario: PackageInfo contains basic package information
    Given a NovusPack package
    And PackageInfo structure
    When basic package information is examined
    Then FormatVersion contains package file format version
    And FileCount contains number of regular content files (excludes special metadata files types 65000-65535)
    And FilesUncompressedSize contains total uncompressed size of regular files (excludes special metadata files)
    And FilesCompressedSize contains total compressed size of regular files (excludes special metadata files)

  @REQ-META-149 @happy
  Scenario: PackageInfo scope and exclusions exclude special metadata files from counts and sizes
    Given a NovusPack package
    When evaluating PackageInfo file counts and sizes
    Then special metadata files (types 65000-65535) are excluded from FileCount
    And special metadata files are excluded from FilesUncompressedSize and FilesCompressedSize

  @REQ-META-115 @happy
  Scenario: PackageInfo.FromHeader synchronizes fields from PackageHeader
    Given a PackageHeader with populated fields and flags
    And an empty PackageInfo
    When PackageInfo.FromHeader is called with the PackageHeader
    Then PackageInfo.FormatVersion equals PackageHeader.FormatVersion
    And PackageInfo.VendorID equals PackageHeader.VendorID
    And PackageInfo.AppID equals PackageHeader.AppID
    And PackageInfo.PackageDataVersion equals PackageHeader.PackageDataVersion
    And PackageInfo.MetadataVersion equals PackageHeader.MetadataVersion
    And PackageInfo feature flags reflect PackageHeader.Flags
    And PackageInfo compression fields reflect PackageHeader.Flags compression type

  @REQ-META-116 @happy
  Scenario: PackageHeader.ToHeader synchronizes header fields from PackageInfo
    Given a PackageInfo with populated fields and flags
    And an empty PackageHeader
    When PackageHeader.ToHeader is called with the PackageInfo
    Then PackageHeader.FormatVersion equals PackageInfo.FormatVersion
    And PackageHeader.VendorID equals PackageInfo.VendorID
    And PackageHeader.AppID equals PackageInfo.AppID
    And PackageHeader.PackageDataVersion equals PackageInfo.PackageDataVersion
    And PackageHeader.MetadataVersion equals PackageInfo.MetadataVersion
    And PackageHeader.Flags reflect PackageInfo feature flags and compression type

  @REQ-META-150 @happy
  Scenario: PackageInfo is treated as source of truth for in-memory state and write serialization
    Given a NovusPack package
    When performing in-memory operations
    Then PackageInfo is used as the source of truth for package state
    When writing package header flags
    Then PackageInfo is used to update header flags before serialization

  @REQ-META-084 @happy
  Scenario: PackageInfo contains package identity and comment
    Given a NovusPack package
    And PackageInfo structure
    When package identity is examined
    Then VendorID contains vendor/platform identifier
    And AppID contains application identifier
    And HasComment indicates if comment exists
    And Comment contains actual comment content

  @REQ-META-084 @happy
  Scenario: PackageInfo contains signature and security information
    Given a NovusPack package
    And PackageInfo structure
    When signature and security information is examined
    Then HasSignatures indicates if signatures exist
    And SignatureCount contains number of signatures
    And Signatures contains detailed signature information
    And IsImmutable indicates if package is signed

  @REQ-META-084 @happy
  Scenario: PackageInfo contains timestamps and feature flags
    Given a NovusPack package
    And PackageInfo structure
    When timestamps and features are examined
    Then Created contains package creation timestamp
    And Modified contains package modification timestamp
    And HasMetadataFiles indicates if metadata files exist
    And HasEncryptedData indicates if encrypted files exist
    And HasCompressedData indicates if compressed files exist
    And IsMetadataOnly indicates if package is metadata-only
