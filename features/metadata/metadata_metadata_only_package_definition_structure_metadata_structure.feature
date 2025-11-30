@domain:metadata @m2 @REQ-META-070 @spec(api_metadata.md#61-metadata-only-package-definition)
Feature: Metadata: Metadata-Only Package Definition (Structure)

  @REQ-META-070 @happy
  Scenario: Metadata-only package definition defines package structure
    Given a NovusPack package
    When metadata-only package definition is examined
    Then FileCount must be 0
    And HasSpecialMetadataFiles must be true
    And TotalSize must be 0
    And CompressedSize must be greater than 0

  @REQ-META-070 @happy
  Scenario: Metadata-only package has no regular content files
    Given a NovusPack package
    And a metadata-only package
    When package is examined
    Then FileCount is 0
    And no regular content files exist
    And package contains only special metadata files

  @REQ-META-070 @happy
  Scenario: Metadata-only package contains special metadata files
    Given a NovusPack package
    And a metadata-only package
    When package is examined
    Then HasSpecialMetadataFiles is true
    And package contains at least one special metadata file
    And special metadata files are present

  @REQ-META-070 @happy
  Scenario: Metadata-only package has no uncompressed content
    Given a NovusPack package
    And a metadata-only package
    When package is examined
    Then TotalSize is 0
    And no uncompressed content data exists
    And CompressedSize is greater than 0

  @REQ-META-070 @error
  Scenario: Metadata-only package definition validates structure
    Given a NovusPack package
    When invalid metadata-only package structure is provided
    Then validation detects structure violations
    And appropriate errors are returned
