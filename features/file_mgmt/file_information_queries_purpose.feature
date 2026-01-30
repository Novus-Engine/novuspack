@domain:file_mgmt @m2 @REQ-FILEMGMT-090 @spec(api_file_mgmt_queries.md#12-purpose)
Feature: File information queries provide file existence and property access

  @REQ-FILEMGMT-090 @happy
  Scenario: File information queries provide file existence and property access
    Given a package with file entries
    When file information queries are used
    Then file existence and property access are provided as specified
    And the purpose matches the file information queries specification
    And callers can query file existence and properties
