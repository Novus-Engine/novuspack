@domain:metadata @m2 @REQ-META-096 @REQ-META-137 @REQ-META-138 @spec(api_metadata.md#832-implementation-details) @spec(api_metadata.md#packageupdatespecialmetadataflags-method) @spec(api_metadata.md#data-flow-architecture)
Feature: Metadata Implementation Details

  @REQ-META-096 @happy
  Scenario: Implementation details provide special file implementation information
    Given a NovusPack package
    When special file implementation is examined
    Then SavePathMetadataFile creates and saves path metadata file
    And implementation details describe file creation process
    And implementation details describe flag updates
    And implementation details describe metadata marshaling

  @REQ-META-096 @happy
  Scenario: SavePathMetadataFile implementation process
    Given a NovusPack package
    And a valid context
    And path metadata entries
    When SavePathMetadataFile is called
    Then current path metadata is retrieved
    And metadata is marshaled to YAML
    And special metadata file entry is created
    And appropriate tags are set on file entry
    And package flags are updated
    And context supports cancellation

  @REQ-META-096 @REQ-META-137 @happy
  Scenario: UpdateSpecialMetadataFlags implementation process
    Given a NovusPack package
    And a valid context
    When UpdateSpecialMetadataFlags is called
    Then special metadata files are checked
    And per-file tags are checked
    And extended attributes are checked
    And package header flags are updated
    And PackageInfo reflects the updated header flags
    And context supports cancellation

  @REQ-META-138 @happy
  Scenario: Header flags and PackageInfo use defined sources of truth across lifecycle
    Given a NovusPack package
    When opening a package from disk
    Then header flags are read from disk and populate PackageInfo
    When performing in-memory operations
    Then PackageInfo is the source of truth for package state
    When writing a package to disk
    Then PackageInfo is used to update header flags before serialization

  @REQ-META-096 @error
  Scenario: Implementation details handle errors during file operations
    Given a NovusPack package
    When implementation encounters errors
    Then errors are handled gracefully
    And appropriate errors are returned
    And errors follow structured error format
