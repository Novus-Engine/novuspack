@domain:file_mgmt @m2 @REQ-FILEMGMT-017 @spec(api_file_management.md#66-add-file-hash)
Feature: Add file hash operations

  @happy
  Scenario: AddFileHash adds hash to file entry
    Given an open writable NovusPack package with existing file
    When AddFileHash is called with hash type, purpose, and data
    Then hash is added to file entry
    And HashCount is incremented
    And hash is stored in hash data section
    And hash is accessible

  @happy
  Scenario: AddFileHash supports multiple hash types and purposes
    Given an open writable NovusPack package with existing file
    When multiple hashes are added with different types and purposes
    Then all hashes are stored
    And hash types are preserved
    And hash purposes are preserved
    And hashes are accessible individually

  @error
  Scenario: AddFileHash fails if file does not exist
    Given an open writable NovusPack package
    When AddFileHash is called with non-existent file
    Then structured validation error is returned

  @error
  Scenario: AddFileHash validates hash data
    Given invalid hash data
    When AddFileHash is called
    Then structured validation error is returned
    And error indicates invalid hash
