@domain:file_mgmt @m2 @REQ-FILEMGMT-163 @spec(api_file_management.md#543-returns)
Feature: AddFilePath Returns Updated FileEntry

  @REQ-FILEMGMT-163 @happy
  Scenario: AddFilePath returns updated fileentry
    Given an open NovusPack package
    And a valid context
    And an existing file entry
    When AddFilePath is called with context, entry, and path entry
    Then updated FileEntry is returned
    And FileEntry contains additional path
    And path entry is added to file entry metadata

  @REQ-FILEMGMT-163 @happy
  Scenario: AddFilePath adds additional path to existing file entry
    Given an open NovusPack package
    And a valid context
    And an existing file entry with one path
    And a path entry with path and metadata
    When AddFilePath is called
    Then additional path is added to file entry
    And PathCount is incremented
    And path entry is stored in file entry metadata

  @REQ-FILEMGMT-163 @happy
  Scenario: AddFilePath enables multiple paths pointing to same content
    Given an open NovusPack package
    And a valid context
    And an existing file entry
    And multiple path entries pointing to same content
    When AddFilePath is called for each path
    Then multiple paths are added to file entry
    And all paths point to identical content
    And path aliasing is supported

  @REQ-FILEMGMT-163 @happy
  Scenario: AddFilePath supports per-path metadata
    Given an open NovusPack package
    And a valid context
    And an existing file entry
    And a path entry with metadata (permissions, timestamps)
    When AddFilePath is called
    Then path entry includes per-path metadata
    And each path can have different permissions
    And each path can have different timestamps
    And per-path metadata is preserved

  @REQ-FILEMGMT-163 @happy
  Scenario: AddFilePath increments MetadataVersion when path is added
    Given an open NovusPack package
    And a valid context
    And an existing file entry with current MetadataVersion
    When AddFilePath is called
    Then MetadataVersion is incremented
    And version change indicates metadata modification
    And path addition is tracked in version

  @REQ-FILEMGMT-163 @error
  Scenario: AddFilePath returns error when package is not open
    Given a NovusPack package that is not open
    And a valid context
    And an existing file entry
    When AddFilePath is called
    Then package not open error is returned
    And error indicates package must be open

  @REQ-FILEMGMT-163 @error
  Scenario: AddFilePath returns error for invalid path
    Given an open NovusPack package
    And a valid context
    And an existing file entry
    And an invalid path entry
    When AddFilePath is called
    Then structured invalid path error is returned
    And error indicates path validation failure
