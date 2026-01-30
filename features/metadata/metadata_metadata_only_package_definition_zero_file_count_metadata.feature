@domain:metadata @m2 @REQ-META-069 @REQ-META-081 @REQ-META-082 @spec(api_metadata.md#6-metadata-only-packages)
Feature: Metadata: Metadata-Only Package Definition (Zero File Count)

  @REQ-META-069 @happy
  Scenario: Metadata-only package has zero file count
    Given an open NovusPack package
    And a metadata-only package
    When package is examined
    Then FileCount is 0
    And no regular content files exist
    And package may contain special metadata files or be empty

  @REQ-META-069 @happy
  Scenario: Metadata-only package may have special metadata files
    Given an open NovusPack package
    And a metadata-only package with special metadata files
    When package is examined
    Then HasSpecialMetadataFiles is true
    And package contains at least one special metadata file
    And special metadata files are present

  @REQ-META-069 @happy
  Scenario: Metadata-only package has zero total size for regular files
    Given an open NovusPack package
    And a metadata-only package
    When package is examined
    Then TotalSize is 0
    And no uncompressed regular content data exists
    And FilesCompressedSize excludes special metadata files

  @REQ-META-081 @happy
  Scenario: IsMetadataOnlyPackage checks if package is metadata-only
    Given an open NovusPack package
    When IsMetadataOnlyPackage is called
    Then true is returned if package has no regular files
    And false is returned if package has regular files
    And metadata-only status is determined

  @REQ-META-081 @happy
  Scenario: AddMetadataOnlyFile adds special metadata file
    Given an open writable metadata-only package
    And special metadata file data
    When AddMetadataOnlyFile is called
    Then special metadata file is added
    And file is stored in package
    And file type is preserved

  @REQ-META-081 @happy
  Scenario: GetMetadataOnlyFiles returns all metadata files
    Given an open NovusPack package
    And metadata-only package with files
    When GetMetadataOnlyFiles is called
    Then array of SpecialFileInfo is returned
    And all metadata files are included
    And file information is complete

  @REQ-META-081 @happy
  Scenario: ValidateMetadataOnlyIntegrity validates package integrity
    Given an open NovusPack package
    And a metadata-only package
    When ValidateMetadataOnlyIntegrity is called
    Then metadata integrity is validated
    And validation result is returned
    And integrity issues are identified

  @REQ-META-082 @happy
  Scenario: ValidateMetadataOnlyPackage performs comprehensive validation
    Given an open NovusPack package
    And a metadata-only package
    When ValidateMetadataOnlyPackage is called
    Then FileCount is validated to be 0
    And IsMetadataOnly flag (Bit 7) is validated to be set
    And all special metadata files are validated (if present)
    And malicious metadata patterns are checked (if metadata files present)
    And signature scope is verified (if signatures present)
    And metadata consistency is ensured
    And package structure is validated

  @REQ-META-082 @error
  Scenario: ValidateMetadataOnlyPackage fails if package has regular files
    Given an open NovusPack package
    And package with regular content files
    When ValidateMetadataOnlyPackage is called
    Then structured validation error is returned
    And error indicates package is not metadata-only

  @REQ-META-082 @happy
  Scenario: ValidateMetadataOnlyPackage allows empty packages without special files
    Given an open NovusPack package
    And package with FileCount 0 and no special metadata files
    When ValidateMetadataOnlyPackage is called
    Then validation succeeds
    And package is recognized as valid empty/placeholder package
