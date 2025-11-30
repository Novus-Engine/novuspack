@domain:file_mgmt @m2 @REQ-FILEMGMT-012 @spec(api_file_management.md#61-update-file)
Feature: Update file operations

  @happy
  Scenario: UpdateFile updates existing file data
    Given an open writable NovusPack package with existing file
    When UpdateFile is called with new data
    Then file data is updated
    And file version is incremented
    And file entry is updated
    And file index is updated

  @happy
  Scenario: UpdateFile preserves file metadata by default
    Given an open writable NovusPack package with existing file
    When UpdateFile is called
    Then file paths are preserved
    And file tags are preserved
    And file type is preserved
    And file metadata is preserved

  @happy
  Scenario: UpdateFile supports compression and encryption updates
    Given an open writable NovusPack package with existing file
    When UpdateFile is called with new compression or encryption
    Then compression settings are updated
    And encryption settings are updated
    And file is re-processed accordingly

  @error
  Scenario: UpdateFile fails if file does not exist
    Given an open writable NovusPack package
    When UpdateFile is called with non-existent path
    Then structured validation error is returned
    And error indicates file not found

  @error
  Scenario: UpdateFile fails if package is read-only
    Given a read-only open NovusPack package
    When UpdateFile is called
    Then structured validation error is returned
    And error indicates read-only package
