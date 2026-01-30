@skip @domain:file_mgmt @m2 @spec(api_file_mgmt_queries.md#12-purpose)
Feature: File Query Operations

# This feature captures high-level expectations for query operations and their relationship to other APIs.
# Detailed runnable scenarios live in the dedicated file_mgmt query and lookup feature files.

  @REQ-FILEMGMT-021 @documentation
  Scenario: File queries return FileEntry objects suitable for inspection and planning
    Given an open package
    When the caller queries for a file by path, ID, hash, or checksum
    Then the caller receives a FileEntry describing the file and its metadata
    And the caller does not need a separate GetFileInfo API

  @REQ-FILEMGMT-021 @documentation
  Scenario: Directory listing is provided by the Metadata API
    Given an open package
    When the caller needs a list of directories
    Then the caller uses the Metadata API directory listing methods
    And the file query API remains focused on files and file entries
