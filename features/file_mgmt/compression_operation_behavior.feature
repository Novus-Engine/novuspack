@domain:file_mgmt @m2 @REQ-FILEMGMT-192 @spec(api_file_mgmt_compression.md#14-fileentrycompress-behavior)
Feature: Compression operation behavior defines compression process

  @REQ-FILEMGMT-192 @happy
  Scenario: Compression operation behavior defines compression process
    Given a FileEntry and a compression operation
    When the compression operation is performed
    Then the compression process follows the defined behavior
    And the behavior matches the FileEntry.Compress behavior specification
    And compression metadata is updated as specified
