@domain:metadata @m2 @REQ-META-096 @spec(api_metadata.md#832-implementation-details)
Feature: Metadata Implementation Details

  @REQ-META-096 @happy
  Scenario: Implementation details provide special file implementation information
    Given a NovusPack package
    When special file implementation is examined
    Then SaveDirectoryMetadataFile creates and saves directory metadata file
    And implementation details describe file creation process
    And implementation details describe flag updates
    And implementation details describe metadata marshaling

  @REQ-META-096 @happy
  Scenario: SaveDirectoryMetadataFile implementation process
    Given a NovusPack package
    And a valid context
    And directory metadata entries
    When SaveDirectoryMetadataFile is called
    Then current directory metadata is retrieved
    And metadata is marshaled to YAML
    And special metadata file entry is created
    And appropriate tags are set on file entry
    And package flags are updated
    And context supports cancellation

  @REQ-META-096 @happy
  Scenario: UpdateSpecialMetadataFlags implementation process
    Given a NovusPack package
    And a valid context
    When UpdateSpecialMetadataFlags is called
    Then special metadata files are checked
    And per-file tags are checked
    And package header flags are updated
    And context supports cancellation

  @REQ-META-096 @error
  Scenario: Implementation details handle errors during file operations
    Given a NovusPack package
    When implementation encounters errors
    Then errors are handled gracefully
    And appropriate errors are returned
    And errors follow structured error format
