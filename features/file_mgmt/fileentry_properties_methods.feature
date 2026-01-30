@domain:file_mgmt @m2 @REQ-FILEMGMT-094 @spec(api_file_mgmt_queries.md#12-purpose)
Feature: FileEntry properties methods provide property access

  @REQ-FILEMGMT-094 @happy
  Scenario: FileEntry properties methods provide property access
    Given a FileEntry in a package
    When FileEntry properties methods are used
    Then property access is provided as specified
    And the methods match the FileEntry properties specification
    And callers can read and optionally set properties
