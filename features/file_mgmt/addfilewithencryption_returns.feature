@domain:file_mgmt @REQ-FILEMGMT-440 @REQ-FILEMGMT-023 @spec(api_file_mgmt_addition.md#addfilewithencryption-returns) @spec(api_file_mgmt_addition.md#23-packageaddfilewithencryption-method)
Feature: AddFileWithEncryption Returns

  AddFileWithEncryption returns the created FileEntry with encryption enabled,
  or an error that occurred during file addition or encryption.

  @REQ-FILEMGMT-440 @happy
  Scenario: AddFileWithEncryption returns created FileEntry on success
    Given an open writable package
    And a valid encryption key
    And a filesystem file path
    When AddFileWithEncryption is called with path, key, and options
    Then created FileEntry is returned
    And FileEntry has encryption enabled
    And error is nil

  @REQ-FILEMGMT-440 @error
  Scenario: AddFileWithEncryption returns error on failure
    Given an open writable package
    When AddFileWithEncryption fails due to path validation, I/O, or encryption
    Then error is returned
    And FileEntry is nil
    And error describes the failure
