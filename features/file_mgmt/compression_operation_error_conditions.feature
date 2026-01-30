@domain:file_mgmt @m2 @REQ-FILEMGMT-193 @spec(api_file_mgmt_compression.md#15-fileentrycompress-error-conditions) @spec(api_file_mgmt_file_entry.md#fileentry-compression-error-conditions)
Feature: Compression operation error conditions define compression errors

  @REQ-FILEMGMT-193 @happy
  Scenario: Compression operation returns structured errors on failure
    Given a FileEntry and a compression operation that may fail
    When a compression error occurs
    Then error conditions are defined as specified
    And returned errors are structured
    And the behavior matches the compression error conditions specification
