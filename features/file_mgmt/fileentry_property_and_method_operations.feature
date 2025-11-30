@domain:file_mgmt @m2 @REQ-FILEMGMT-022 @spec(api_file_management.md#11-fileentry-methods)
Feature: FileEntry property and method operations

  @happy
  Scenario: IsCompressed indicates compression status
    Given a FileEntry instance
    When IsCompressed is called
    Then true is returned if file is compressed
    And false is returned if file is not compressed
    And status matches CompressionType

  @happy
  Scenario: HasEncryptionKey indicates encryption key presence
    Given a FileEntry instance
    When HasEncryptionKey is called
    Then true is returned if encryption key is set
    And false is returned if no encryption key
    And status matches encryption configuration

  @happy
  Scenario: GetEncryptionType returns encryption type
    Given a FileEntry instance with encryption
    When GetEncryptionType is called
    Then encryption type is returned
    And type matches EncryptionType field

  @happy
  Scenario: IsEncrypted indicates encryption status
    Given a FileEntry instance
    When IsEncrypted is called
    Then true is returned if file is encrypted
    And false is returned if file is not encrypted
    And status matches EncryptionType

  @happy
  Scenario: ToBinaryFormat serializes file entry
    Given a FileEntry instance
    When ToBinaryFormat is called
    Then file entry is serialized to binary
    And binary format matches specification
    And serialized data can be deserialized

  @happy
  Scenario: SetEncryptionKey configures encryption key
    Given a FileEntry instance
    When SetEncryptionKey is called with key
    Then encryption key is set
    And key is stored securely
    And encryption can use key

  @happy
  Scenario: Encrypt method encrypts file data
    Given a FileEntry instance with unencrypted data
    When Encrypt is called
    Then file data is encrypted
    And EncryptionType is set
    And encrypted data is stored

  @happy
  Scenario: Decrypt method decrypts file data
    Given a FileEntry instance with encrypted data
    When Decrypt is called with correct key
    Then file data is decrypted
    And original data is restored

  @happy
  Scenario: LoadData loads file data into memory
    Given a FileEntry instance
    When LoadData is called
    Then file data is loaded into Data field
    And IsDataLoaded is set to true
    And data is accessible

  @happy
  Scenario: ProcessData processes file data
    Given a FileEntry instance with data
    When ProcessData is called
    Then data processing occurs
    And compression is applied if configured
    And encryption is applied if configured
    And ProcessingState tracks progress

  @error
  Scenario: Decrypt fails with incorrect key
    Given a FileEntry instance with encrypted data
    When Decrypt is called with incorrect key
    Then structured encryption error is returned

  @error
  Scenario: FileEntry methods respect context cancellation
    Given a FileEntry instance
    And a cancelled context
    When FileEntry method is called
    Then structured context error is returned
