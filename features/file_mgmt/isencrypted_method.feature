@domain:file_mgmt @m2 @REQ-FILEMGMT-097 @spec(api_file_management.md#1115-isencrypted-returns)
Feature: IsEncrypted Method

  @REQ-FILEMGMT-097 @happy
  Scenario: Isencrypted returns encryption status
    Given a file entry
    When IsEncrypted is called
    Then boolean encryption status is returned
    And status indicates if file is encrypted

  @REQ-FILEMGMT-097 @happy
  Scenario: IsEncrypted returns true for encrypted files
    Given a file entry with encryption enabled
    When IsEncrypted is called
    Then true is returned
    And encryption status is correctly indicated

  @REQ-FILEMGMT-097 @happy
  Scenario: IsEncrypted returns false for unencrypted files
    Given a file entry without encryption
    When IsEncrypted is called
    Then false is returned
    And unencrypted status is correctly indicated

  @REQ-FILEMGMT-097 @happy
  Scenario: IsEncrypted is equivalent to GetEncryptionType != EncryptionNone
    Given a file entry
    When IsEncrypted is called
    Then result matches GetEncryptionType() != EncryptionNone
    And both methods return consistent encryption status

  @REQ-FILEMGMT-097 @happy
  Scenario: IsEncrypted checks encryption type from file entry metadata
    Given a file entry
    When IsEncrypted is called
    Then encryption type is checked from file entry metadata
    And status is determined from metadata
    And result reflects actual encryption state

  @REQ-FILEMGMT-097 @happy
  Scenario: IsEncrypted provides access to file entry properties
    Given a file entry
    When IsEncrypted is called
    Then file entry property is accessed
    And encryption property is returned
    And access is efficient
