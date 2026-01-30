@domain:file_mgmt @REQ-FILEMGMT-438 @REQ-FILEMGMT-023 @spec(api_file_mgmt_addition.md#addfilewithencryption-implementation) @spec(api_file_mgmt_addition.md#23-packageaddfilewithencryption-method)
Feature: AddFileWithEncryption Implementation

  AddFileWithEncryption merges the provided options with the encryption key,
  then calls AddFile. The key parameter takes precedence over options.EncryptionKey.
  If options is nil, a new AddFileOptions with only the key set is created.

  @REQ-FILEMGMT-438 @happy
  Scenario: AddFileWithEncryption merges options with key then calls AddFile
    Given an open writable package
    And a valid encryption key
    And a filesystem file path
    And AddFileOptions with compression set
    When AddFileWithEncryption is called with path, key, and options
    Then options are merged with EncryptionKey set to provided key
    And AddFile is invoked with merged options
    And file is added with encryption enabled

  @REQ-FILEMGMT-438 @happy
  Scenario: Key parameter takes precedence over options.EncryptionKey
    Given an open writable package
    And two valid encryption keys "keyA" and "keyB"
    And AddFileOptions with EncryptionKey set to "keyA"
    When AddFileWithEncryption is called with path, "keyB", and options
    Then merged options use "keyB" as EncryptionKey
    And file is encrypted with "keyB"

  @REQ-FILEMGMT-438 @happy
  Scenario: Nil options creates new AddFileOptions with only key set
    Given an open writable package
    And a valid encryption key
    And a filesystem file path
    When AddFileWithEncryption is called with path, key, and nil options
    Then new AddFileOptions is created with EncryptionKey set to key
    And AddFile is invoked with that options
    And file is added with encryption enabled
