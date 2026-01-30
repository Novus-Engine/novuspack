@domain:metadata @m2 @REQ-META-153 @REQ-META-154 @spec(api_metadata.md#8218-directory-metadata-convenience-methods) @spec(api_metadata.md#82111-packageupdatedirectorymetadata-method)
Feature: Directory Metadata Convenience Methods

  @REQ-META-153 @happy
  Scenario: Directory metadata convenience methods provide helpers for directory metadata operations
    Given an open NovusPack package
    And a valid context
    When using directory metadata convenience methods
    Then directory metadata convenience helpers are available
    And helpers operate on directory paths

  @REQ-META-154 @happy
  Scenario: UpdateDirectoryMetadata updates directory metadata without modifying files
    Given an open writable NovusPack package
    And a directory path "/docs/"
    And a valid context
    When UpdateDirectoryMetadata is called for "/docs/"
    Then directory path metadata is updated
    And file content is not modified

