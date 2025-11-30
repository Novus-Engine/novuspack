@domain:metadata @m2 @REQ-META-084 @spec(api_metadata.md#71-packageinfo-structure)
Feature: PackageInfo Structure

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
    Then FileCount contains number of files
    And FilesUncompressedSize contains total uncompressed size
    And FilesCompressedSize contains total compressed size

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
    And SecurityLevel contains overall security level
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
