@domain:file_mgmt @m2 @REQ-FILEMGMT-233 @spec(api_file_mgmt_file_entry.md#9-fileentry-encryption)
Feature: 11.2 File Entry Encryption is specified and implemented

  @REQ-FILEMGMT-233 @happy
  Scenario: File Entry Encryption is specified and implemented
    Given an open NovusPack package with encrypted file entries
    When FileEntry encryption is used or queried
    Then FileEntry encryption specification is implemented
    And encryption operations match the specification
    And the behavior matches the fileentry-encryption specification
    And encryption key and type are available as specified
