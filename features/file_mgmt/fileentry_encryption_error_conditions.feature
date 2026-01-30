@domain:file_mgmt @m2 @REQ-FILEMGMT-103 @spec(api_file_mgmt_file_entry.md#98-fileentry-encryption-error-conditions)
Feature: FileEntry encryption error conditions handle encryption failures

  @REQ-FILEMGMT-103 @happy
  Scenario: FileEntry encryption error conditions handle encryption failures
    Given a FileEntry and an encryption operation that may fail
    When an encryption failure occurs
    Then error conditions are handled as specified
    And returned errors are structured
    And the behavior matches the FileEntry encryption error conditions specification
