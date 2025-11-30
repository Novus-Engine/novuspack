@domain:metadata @m2 @REQ-META-043 @spec(metadata.md#24-package-metadata-example)
Feature: Package Metadata Example

  @REQ-META-043 @happy
  Scenario: Package metadata example contains package information
    Given an open NovusPack package
    And a metadata file
    When package metadata example is examined
    Then name field contains package name
    And version field contains package version
    And description field contains package description
    And author field contains package author
    And license field contains package license
    And created field contains ISO8601 creation timestamp
    And modified field contains ISO8601 modification timestamp

  @REQ-META-043 @happy
  Scenario: Package metadata example contains game-specific metadata
    Given an open NovusPack package
    And a metadata file
    When package metadata example is examined
    Then engine field contains game engine name
    And platform field contains array of target platforms
    And genre field contains game genre
    And rating field contains age rating
    And requirements object contains system requirements

  @REQ-META-043 @happy
  Scenario: Package metadata example contains asset metadata
    Given an open NovusPack package
    And a metadata file
    When package metadata example is examined
    Then textures field contains number of texture files
    And sounds field contains number of sound files
    And models field contains number of 3D model files
    And scripts field contains number of script files
    And total_size field contains total asset size in bytes

  @REQ-META-043 @happy
  Scenario: Package metadata example contains security metadata
    Given an open NovusPack package
    And a metadata file
    When package metadata example is examined
    Then encryption_level field contains encryption level
    And signature_type field contains signature type
    And security_scan field contains boolean scan status
    And trusted_source field contains boolean trusted status

  @REQ-META-043 @happy
  Scenario: Package metadata example contains custom metadata
    Given an open NovusPack package
    And a metadata file
    When package metadata example is examined
    Then custom object provides extensible key-value pairs
    And custom fields can store additional metadata
    And custom metadata extends package information

  @REQ-META-043 @happy
  Scenario: Package metadata example demonstrates complete metadata structure
    Given an open NovusPack package
    And metadata file with example data
    When package metadata is retrieved
    Then metadata structure matches example format
    And all metadata fields are populated
    And example demonstrates typical usage

  @REQ-META-011 @error
  Scenario: Package metadata example validation fails with invalid structure
    Given an open NovusPack package
    And metadata file with invalid structure
    When package metadata is validated
    Then structured validation error is returned
    And error indicates invalid metadata structure
