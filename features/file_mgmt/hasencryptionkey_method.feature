@domain:file_mgmt @m2 @REQ-FILEMGMT-096 @spec(api_file_mgmt_file_entry.md#73-hasencryptionkey-returns)
Feature: HasEncryptionKey Method

  @REQ-FILEMGMT-096 @happy
  Scenario: Hasencryptionkey returns encryption key presence status
    Given a file entry
    When HasEncryptionKey is called
    Then boolean encryption key presence status is returned
    And status indicates if file has encryption key set

  @REQ-FILEMGMT-096 @happy
  Scenario: HasEncryptionKey returns true when encryption key is set
    Given a file entry with encryption key set
    When HasEncryptionKey is called
    Then true is returned
    And encryption key presence is correctly indicated

  @REQ-FILEMGMT-096 @happy
  Scenario: HasEncryptionKey returns false when encryption key is not set
    Given a file entry without encryption key
    When HasEncryptionKey is called
    Then false is returned
    And absence of encryption key is correctly indicated

  @REQ-FILEMGMT-096 @happy
  Scenario: HasEncryptionKey checks key presence from file entry metadata
    Given a file entry
    When HasEncryptionKey is called
    Then encryption key presence is checked from file entry metadata
    And status is determined from metadata
    And result reflects actual key presence state

  @REQ-FILEMGMT-096 @happy
  Scenario: HasEncryptionKey provides access to file entry properties
    Given a file entry
    When HasEncryptionKey is called
    Then file entry property is accessed
    And encryption key property is returned
    And access is efficient

  @REQ-FILEMGMT-096 @happy
  Scenario: HasEncryptionKey distinguishes between encryption type and key presence
    Given a file entry with encryption type set but no key
    When HasEncryptionKey is called
    Then false is returned when no key is present
    And encryption type and key presence are distinguished
