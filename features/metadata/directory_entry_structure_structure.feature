@domain:metadata @m2 @REQ-META-029 @spec(metadata.md#132-directory-entry-structure)
Feature: Directory Entry Structure

  @REQ-META-029 @happy
  Scenario: Directory entry structure contains path field
    Given an open NovusPack package
    And a directory entry
    When directory entry structure is examined
    Then path field contains directory path
    And path must end with "/"
    And path entry is accessible

  @REQ-META-029 @happy
  Scenario: Directory entry structure contains properties field
    Given an open NovusPack package
    And a directory entry
    When directory entry structure is examined
    Then properties field contains directory-specific tags
    And tags are stored as array
    And tags provide directory metadata
    And properties support various tag value types

  @REQ-META-029 @happy
  Scenario: Directory entry structure contains inheritance field
    Given an open NovusPack package
    And a directory entry
    When directory entry structure is examined
    Then inheritance field contains inheritance settings
    And enabled property controls inheritance
    And priority property determines inheritance priority
    And inheritance settings are accessible

  @REQ-META-029 @happy
  Scenario: Directory entry structure contains metadata field
    Given an open NovusPack package
    And a directory entry
    When directory entry structure is examined
    Then metadata field contains directory metadata
    And created field contains ISO8601 creation timestamp
    And modified field contains ISO8601 modification timestamp
    And description field contains human-readable description

  @REQ-META-029 @happy
  Scenario: Directory entry example demonstrates structure format
    Given an open NovusPack package
    And directory metadata file with example entry
    When directory entry is examined
    Then path is "/assets/"
    And properties contain category="texture" and compression="lossless"
    And inheritance is enabled with priority 1
    And metadata contains creation and modification times
    And example demonstrates complete structure

  @REQ-META-029 @happy
  Scenario: Directory entry structure supports YAML format
    Given an open NovusPack package
    And directory metadata file
    When directory entry is parsed
    Then entry follows YAML format
    And YAML structure is valid
    And entry is parseable from YAML

  @REQ-META-011 @error
  Scenario: Directory entry validation fails with invalid path
    Given an open NovusPack package
    And directory entry with invalid path
    When directory entry is validated
    Then structured validation error is returned
    And error indicates invalid path format

  @REQ-META-011 @error
  Scenario: Directory entry validation fails with invalid properties
    Given an open NovusPack package
    And directory entry with invalid properties
    When directory entry is validated
    Then structured validation error is returned
    And error indicates invalid properties
