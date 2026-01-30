@domain:file_mgmt @m2 @REQ-FILEMGMT-035 @spec(api_file_mgmt_compression.md#1-fileentrycompress-method)
Feature: CompressFile compresses existing file in package

  @REQ-FILEMGMT-035 @happy
  Scenario: CompressFile compresses an existing file in the package
    Given a package with an uncompressed file entry
    When CompressFile is called for the file
    Then the file content is compressed
    And the FileEntry is updated with compression metadata
    And the behavior matches the CompressFile specification
