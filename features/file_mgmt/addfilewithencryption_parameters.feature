@domain:file_mgmt @REQ-FILEMGMT-439 @REQ-FILEMGMT-023 @spec(api_file_mgmt_addition.md#addfilewithencryption-parameters) @spec(api_file_mgmt_addition.md#23-packageaddfilewithencryption-method)
Feature: AddFileWithEncryption Parameters

  AddFileWithEncryption accepts context, path, key, and optional options.
  Path follows filesystem input path rules. Key must be valid and not expired.

  @REQ-FILEMGMT-439 @happy
  Scenario: AddFileWithEncryption accepts context path key and options
    Given an open writable package
    And a valid encryption key
    And a filesystem file path
    When AddFileWithEncryption is called with context, path, key, and nil options
    Then context is used for cancellation and timeout
    And path is used as filesystem-style input path
    And key is used for encryption
    And options may be nil for defaults

  @REQ-FILEMGMT-439 @happy
  Scenario: Path follows filesystem input path derivation rules
    Given an open writable package
    And a valid encryption key
    When AddFileWithEncryption is called with relative or absolute path
    Then stored package path is derived per AddFile path rules
    And path is validated per spec

  @REQ-FILEMGMT-439 @error
  Scenario: Invalid or expired key causes validation error
    Given an open writable package
    When AddFileWithEncryption is called with invalid or expired key
    Then operation fails with ErrTypeEncryption or validation error
    And error indicates key is invalid or expired
