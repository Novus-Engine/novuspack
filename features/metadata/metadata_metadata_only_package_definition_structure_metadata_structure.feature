@domain:metadata @m2 @REQ-META-070 @spec(api_metadata.md#61-metadata-only-package-definition)
Feature: Metadata: Metadata-Only Package Definition (Structure)

  @REQ-META-070 @happy
  Scenario: Metadata-only package definition defines package structure
    Given a NovusPack package
    When metadata-only package definition is examined
    Then FileCount must be 0
    And IsMetadataOnly flag (Bit 7) must be set
    And TotalSize must be 0
    And package may or may not contain special metadata files

  @REQ-META-070 @happy
  Scenario: Metadata-only package has no regular content files
    Given a NovusPack package
    And a metadata-only package
    When package is examined
    Then FileCount is 0
    And no regular content files exist
    And package may contain special metadata files or be empty

  @REQ-META-070 @happy
  Scenario: Metadata-only package may contain special metadata files
    Given a NovusPack package
    And a metadata-only package with special metadata files
    When package is examined
    Then HasSpecialMetadataFiles is true
    And package contains at least one special metadata file
    And special metadata files are present

  @REQ-META-070 @happy
  Scenario: Metadata-only package has no uncompressed regular content
    Given a NovusPack package
    And a metadata-only package
    When package is examined
    Then TotalSize is 0
    And no uncompressed regular content data exists
    And FilesCompressedSize excludes special metadata files

  @REQ-META-070 @error
  Scenario: Metadata-only package definition validates structure
    Given a NovusPack package
    When invalid metadata-only package structure is provided
    Then validation detects structure violations
    And appropriate errors are returned
