@domain:metadata @m2 @REQ-META-097 @spec(api_metadata.md#84-directory-association-system)
Feature: Directory Association System

  @REQ-META-097 @happy
  Scenario: Directory association system links FileEntry to DirectoryEntry
    Given a NovusPack package
    And a valid context
    When directory association system is used
    Then FileEntry objects are linked to DirectoryEntry metadata
    And tag inheritance is enabled
    And filesystem property management is enabled
    And context supports cancellation

  @REQ-META-097 @happy
  Scenario: AssociateFileWithDirectory links files to directories
    Given a NovusPack package
    And a valid context
    And a file path
    And a directory path
    When AssociateFileWithDirectory is called
    Then file entry is retrieved
    And directory entry is retrieved
    And association is set on file entry
    And parent directory is resolved if available
    And inherited tags are resolved and cached
    And context supports cancellation

  @REQ-META-097 @happy
  Scenario: UpdateFileDirectoryAssociations rebuilds all associations
    Given a NovusPack package
    And a valid context
    When UpdateFileDirectoryAssociations is called
    Then all files are retrieved
    And all directories are retrieved
    And directory path map is built
    And each file is associated with its directory
    And context supports cancellation

  @REQ-META-097 @error
  Scenario: Directory association system handles errors
    Given a NovusPack package
    When invalid file or directory paths are provided
    Then appropriate errors are returned
    And errors follow structured error format
