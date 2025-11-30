@domain:metadata @m2 @REQ-META-011 @REQ-META-014 @REQ-META-098 @REQ-META-101 @spec(api_metadata.md#841-association-properties,api_metadata.md#843-association-management)
Feature: File-Directory Association Management

  @REQ-META-098 @happy
  Scenario: FileEntry directory properties link to directory metadata
    Given an open NovusPack package
    And a file entry with directory association
    When FileEntry directory properties are examined
    Then DirectoryEntry pointer is available
    And ParentDirectory pointer is available
    And InheritedTags are cached

  @REQ-META-098 @happy
  Scenario: DirectoryEntry filesystem properties provide filesystem information
    Given an open NovusPack package
    And a DirectoryEntry with filesystem properties
    When DirectoryEntry filesystem properties are examined
    Then Mode property is available for Unix/Linux permissions
    And UID and GID properties are available
    And ACL entries are available
    And WindowsAttrs property is available for Windows
    And ExtendedAttrs map is available
    And Flags property is available

  @REQ-META-101 @happy
  Scenario: AssociateFileWithDirectory links file to directory metadata
    Given an open writable NovusPack package
    And a file path
    And a directory path with metadata
    And a valid context
    When AssociateFileWithDirectory is called
    Then file is linked to directory metadata
    And file inherits directory tags
    And association is persisted

  @REQ-META-101 @happy
  Scenario: DisassociateFileFromDirectory removes file-directory link
    Given an open writable NovusPack package
    And a file with directory association
    And a valid context
    When DisassociateFileFromDirectory is called
    Then file-directory link is removed
    And file no longer inherits directory tags
    And association is cleared

  @REQ-META-101 @happy
  Scenario: UpdateFileDirectoryAssociations updates all file associations
    Given an open writable NovusPack package
    And multiple files with directory paths
    And a valid context
    When UpdateFileDirectoryAssociations is called
    Then all file associations are updated
    And files are linked to correct directories
    And associations match file paths

  @REQ-META-101 @happy
  Scenario: GetFileDirectoryAssociations retrieves all file associations
    Given an open NovusPack package
    And files with directory associations
    And a valid context
    When GetFileDirectoryAssociations is called
    Then map of file paths to DirectoryEntry is returned
    And all file associations are included
    And DirectoryEntry pointers are valid

  @REQ-META-101 @happy
  Scenario: Association management maintains inheritance relationships
    Given an open writable NovusPack package
    And files with directory associations
    And directory hierarchy with tags
    When associations are maintained
    Then files inherit tags from associated directories
    And inheritance relationships are preserved
    And tag inheritance works correctly

  @REQ-META-011 @REQ-META-014 @error
  Scenario: Association management operations respect context cancellation
    Given an open writable NovusPack package
    And a cancelled context
    When association management operation is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-META-011 @error
  Scenario: AssociateFileWithDirectory fails if file does not exist
    Given an open writable NovusPack package
    And a non-existent file path
    And a directory path
    And a valid context
    When AssociateFileWithDirectory is called
    Then structured validation error is returned
    And error indicates file not found

  @REQ-META-011 @error
  Scenario: AssociateFileWithDirectory fails if directory does not exist
    Given an open writable NovusPack package
    And a file path
    And a non-existent directory path
    And a valid context
    When AssociateFileWithDirectory is called
    Then structured validation error is returned
    And error indicates directory not found
