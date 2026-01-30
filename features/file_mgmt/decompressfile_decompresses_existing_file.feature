@domain:file_mgmt @m2 @REQ-FILEMGMT-036 @spec(api_file_mgmt_compression.md#1-fileentrycompress-method)
Feature: DecompressFile decompresses existing file in package

  @REQ-FILEMGMT-036 @happy
  Scenario: DecompressFile decompresses an existing file in the package
    Given a package with a compressed file entry
    When DecompressFile is called for the file
    Then the file content is decompressed
    And the FileEntry is updated to reflect decompressed state
    And the behavior matches the DecompressFile specification
