@domain:metadata @m2 @REQ-META-102 @spec(api_metadata.md#85-file-directory-association)
Feature: File-Directory Relationship Management

  @REQ-META-102 @happy
  Scenario: File-directory association provides file-directory relationships
    Given a NovusPack package
    When file-directory association is used
    Then files can be associated with directory metadata
    And directory associations enable tag inheritance
    And file-directory relationships are managed

  @REQ-META-102 @happy
  Scenario: AssociateFileWithDirectory links files to directories
    Given a NovusPack package
    And a valid context
    And a file path
    And a directory path
    When AssociateFileWithDirectory is called
    Then file is linked to directory metadata
    And association enables tag inheritance
    And context supports cancellation

  @REQ-META-102 @happy
  Scenario: DisassociateFileFromDirectory removes associations
    Given a NovusPack package
    And a valid context
    And a file with directory association
    When DisassociateFileFromDirectory is called
    Then file-directory association is removed
    And tag inheritance is no longer active
    And context supports cancellation

  @REQ-META-102 @happy
  Scenario: GetFileDirectoryAssociations retrieves all associations
    Given a NovusPack package
    And a valid context
    And files with directory associations
    When GetFileDirectoryAssociations is called
    Then all file-directory associations are returned
    And associations map file paths to DirectoryEntry
    And context supports cancellation

  @REQ-META-102 @error
  Scenario: File-directory association handles invalid paths
    Given a NovusPack package
    When invalid file or directory path is provided
    Then appropriate error is returned
    And error indicates path problem
    And error follows structured error format
