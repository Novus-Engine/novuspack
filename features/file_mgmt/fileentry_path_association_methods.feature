@skip @domain:file_mgmt @m2 @spec(api_file_mgmt_file_entry.md#5-path-management)
Feature: FileEntry Path Association Methods

# This feature captures high-level expectations for FileEntry path management and association behavior.
# Detailed runnable scenarios live in the dedicated file_mgmt and metadata feature files.

  @documentation
  Scenario: FileEntry primary path is returned in display format
    Given a FileEntry with one or more stored paths
    When GetPrimaryPath is called
    Then the returned path has no leading slash
    And the returned path uses forward slashes as separators

  @documentation
  Scenario: FileEntry can associate paths with PathMetadataEntry records
    Given a FileEntry that includes one or more paths
    And a PathMetadataEntry with a matching stored path
    When the FileEntry associates with the PathMetadataEntry
    Then subsequent metadata lookups for that path return the associated PathMetadataEntry
