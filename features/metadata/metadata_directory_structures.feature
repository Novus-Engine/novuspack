@domain:metadata @m2 @REQ-META-088 @spec(api_metadata.md#81-directory-structures)
Feature: Metadata Directory Structures

  @REQ-META-088 @happy
  Scenario: DirectoryEntry structure contains path entry
    Given an open NovusPack package
    And a DirectoryEntry
    When DirectoryEntry structure is examined
    Then Path field contains PathEntry
    And path entry must end with "/"
    And path entry is accessible

  @REQ-META-088 @happy
  Scenario: DirectoryEntry structure contains properties
    Given an open NovusPack package
    And a DirectoryEntry
    When DirectoryEntry structure is examined
    Then Properties field contains directory-specific tags
    And tags are stored as array
    And tags provide directory metadata

  @REQ-META-088 @happy
  Scenario: DirectoryEntry structure contains inheritance settings
    Given an open NovusPack package
    And a DirectoryEntry
    When DirectoryEntry structure is examined
    Then Inheritance field contains DirectoryInheritance
    And Enabled property controls inheritance
    And Priority property determines inheritance priority

  @REQ-META-088 @happy
  Scenario: DirectoryEntry structure contains metadata
    Given an open NovusPack package
    And a DirectoryEntry
    When DirectoryEntry structure is examined
    Then Metadata field contains DirectoryMetadata
    And Created field contains ISO8601 timestamp
    And Modified field contains ISO8601 timestamp
    And Description field contains human-readable description

  @REQ-META-088 @happy
  Scenario: DirectoryEntry structure contains filesystem properties
    Given an open NovusPack package
    And a DirectoryEntry
    When DirectoryEntry structure is examined
    Then FileSystem field contains DirectoryFileSystem
    And filesystem properties are available
    And properties support Unix/Linux and Windows

  @REQ-META-088 @happy
  Scenario: DirectoryEntry provides parent directory pointer
    Given an open NovusPack package
    And a DirectoryEntry with parent
    When DirectoryEntry ParentDirectory is examined
    Then ParentDirectory pointer references parent directory
    And pointer is nil for root directory
    And directory hierarchy is accessible

  @REQ-META-088 @happy
  Scenario: DirectoryInheritance controls tag inheritance behavior
    Given an open NovusPack package
    And a DirectoryEntry with inheritance settings
    When DirectoryInheritance is examined
    Then Enabled property controls whether inheritance is provided
    And Priority property determines inheritance priority
    And higher priority overrides lower priority

  @REQ-META-088 @happy
  Scenario: DirectoryMetadata contains directory creation and modification times
    Given an open NovusPack package
    And a DirectoryEntry with metadata
    When DirectoryMetadata is examined
    Then Created field contains creation timestamp
    And Modified field contains modification timestamp
    And timestamps use ISO8601 format

  @REQ-META-011 @error
  Scenario: DirectoryEntry validation fails with invalid path
    Given an open NovusPack package
    And a DirectoryEntry with invalid path
    When DirectoryEntry is validated
    Then structured validation error is returned
    And error indicates invalid path format

  @REQ-META-011 @error
  Scenario: DirectoryEntry validation fails with invalid inheritance settings
    Given an open NovusPack package
    And a DirectoryEntry with invalid inheritance
    When DirectoryEntry is validated
    Then structured validation error is returned
    And error indicates invalid inheritance settings
