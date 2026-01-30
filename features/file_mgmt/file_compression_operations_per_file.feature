@domain:file_mgmt @m2 @REQ-FILEMGMT-018 @spec(api_file_mgmt_compression.md#1-fileentrycompress-method)
Feature: File compression operations support per-file compression

  @REQ-FILEMGMT-018 @happy
  Scenario: File compression operations support per-file compression
    Given a package with file entries
    When file compression operations are used
    Then per-file compression is supported
    And the FileEntry.Compress method is available for compression
    And the behavior matches the file compression operations specification
