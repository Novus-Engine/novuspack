@domain:metadata @m2 @REQ-META-088 @spec(api_metadata.md#81-path-metadata-structures)
Feature: Metadata Path Structures

  @REQ-META-088 @happy
  Scenario: PathMetadataEntry structure contains path entry
    Given an open NovusPack package
    And a PathMetadataEntry
    When PathMetadataEntry structure is examined
    Then Path field contains PathEntry
    And path entry is accessible
    And path entry uses generics.PathEntry type

  @REQ-META-088 @happy
  Scenario: PathMetadataEntry structure contains properties
    Given an open NovusPack package
    And a PathMetadataEntry
    When PathMetadataEntry structure is examined
    Then Properties field contains path-specific tags
    And tags are stored as array
    And tags provide path metadata

  @REQ-META-088 @happy
  Scenario: PathMetadataEntry structure contains inheritance settings
    Given an open NovusPack package
    And a PathMetadataEntry
    When PathMetadataEntry structure is examined
    Then Inheritance field contains PathInheritance
    And Enabled property controls inheritance
    And Priority property determines inheritance priority

  @REQ-META-088 @happy
  Scenario: PathMetadataEntry structure contains metadata
    Given an open NovusPack package
    And a PathMetadataEntry
    When PathMetadataEntry structure is examined
    Then Metadata field contains PathMetadata
    And Created field contains ISO8601 timestamp
    And Modified field contains ISO8601 timestamp
    And Description field contains human-readable description

  @REQ-META-088 @happy
  Scenario: PathMetadataEntry structure contains filesystem properties
    Given an open NovusPack package
    And a PathMetadataEntry
    When PathMetadataEntry structure is examined
    Then FileSystem field contains PathFileSystem
    And filesystem properties are available
    And properties support Unix/Linux and Windows

  @REQ-META-088 @happy
  Scenario: PathMetadataEntry provides parent path pointer
    Given an open NovusPack package
    And a PathMetadataEntry with parent
    When PathMetadataEntry ParentPath is examined
    Then ParentPath pointer references parent path
    And pointer is nil for root path
    And path hierarchy is accessible

  @REQ-META-088 @happy
  Scenario: PathInheritance controls tag inheritance behavior
    Given an open NovusPack package
    And a PathMetadataEntry with inheritance settings
    When PathInheritance is examined
    Then Enabled property controls whether inheritance is provided
    And Priority property determines inheritance priority
    And higher priority overrides lower priority

  @REQ-META-088 @happy
  Scenario: PathMetadata contains path creation and modification times
    Given an open NovusPack package
    And a PathMetadataEntry with metadata
    When PathMetadata is examined
    Then Created field contains creation timestamp
    And Modified field contains modification timestamp
    And timestamps use ISO8601 format

  @REQ-META-011 @error
  Scenario: PathMetadataEntry validation fails with invalid path
    Given an open NovusPack package
    And a PathMetadataEntry with invalid path
    When PathMetadataEntry is validated
    Then structured validation error is returned
    And error indicates invalid path format

  @REQ-META-011 @error
  Scenario: PathMetadataEntry validation fails with invalid inheritance settings
    Given an open NovusPack package
    And a PathMetadataEntry with invalid inheritance
    When PathMetadataEntry is validated
    Then structured validation error is returned
    And error indicates invalid inheritance settings
