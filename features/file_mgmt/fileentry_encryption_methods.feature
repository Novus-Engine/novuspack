@domain:file_mgmt @m2 @REQ-FILEMGMT-100 @spec(api_file_mgmt_file_entry.md#95-fileentry-encryption-purpose)
Feature: FileEntry encryption methods manage encryption keys and operations

  @REQ-FILEMGMT-100 @happy
  Scenario: FileEntry encryption methods manage encryption keys and operations
    Given a FileEntry that may be encrypted
    When FileEntry encryption methods are used
    Then encryption keys and operations are managed as specified
    And the behavior matches the FileEntry encryption purpose specification
    And key access follows secure key handling rules
