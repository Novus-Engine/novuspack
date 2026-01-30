@domain:metadata @m2 @REQ-META-029 @spec(metadata.md#132-path-entry-structure)
Feature: Path Entry Structure

  @REQ-META-029 @happy
  Scenario: Path entry structure contains path field
    Given an open NovusPack package
    And a PathMetadataEntry
    When path entry structure is examined
    Then path field contains path entry
    And path entry uses generics.PathEntry type
    And path entry is accessible

  @REQ-META-029 @happy
  Scenario: Path entry structure contains properties field
    Given an open NovusPack package
    And a PathMetadataEntry
    When path entry structure is examined
    Then properties field contains path-specific tags
    And tags are stored as array
    And tags provide path metadata
    And properties support various tag value types

  @REQ-META-029 @happy
  Scenario: Path entry structure contains inheritance field
    Given an open NovusPack package
    And a PathMetadataEntry
    When path entry structure is examined
    Then inheritance field contains inheritance settings
    And enabled property controls inheritance
    And priority property determines inheritance priority
    And inheritance settings are accessible

  @REQ-META-029 @happy
  Scenario: Path entry structure contains metadata field
    Given an open NovusPack package
    And a PathMetadataEntry
    When path entry structure is examined
    Then metadata field contains path metadata
    And created field contains ISO8601 creation timestamp
    And modified field contains ISO8601 modification timestamp
    And description field contains human-readable description

  @REQ-META-029 @happy
  Scenario: Path entry example demonstrates structure format
    Given an open NovusPack package
    And path metadata file with example entry
    When path entry is examined
    Then path is "/assets/" with leading slash
    And properties contain category="texture" and compression="lossless"
    And inheritance is enabled with priority 1
    And metadata contains creation and modification times
    And example demonstrates complete structure
    And path format includes mandatory leading slash

  @REQ-META-029 @happy
  Scenario: Path entry structure supports YAML format
    Given an open NovusPack package
    And path metadata file
    When path entry is parsed
    Then entry follows YAML format
    And YAML structure is valid
    And entry is parseable from YAML

  @REQ-META-011 @error
  Scenario: Path entry validation fails with invalid path
    Given an open NovusPack package
    And PathMetadataEntry with invalid path
    When path entry is validated
    Then structured validation error is returned
    And error indicates invalid path format

  @REQ-META-011 @error
  Scenario: Path entry validation fails with invalid properties
    Given an open NovusPack package
    And PathMetadataEntry with invalid properties
    When path entry is validated
    Then structured validation error is returned
    And error indicates invalid properties
