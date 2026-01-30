@skip @domain:file_mgmt @m2 @spec(api_file_mgmt_addition.md#29-usage-notes)
Feature: File Management Usage Examples

# This feature captures usage-oriented scenarios derived from the file management specs.
# Detailed runnable scenarios live in the dedicated file_mgmt feature files.

  @documentation
  Scenario: Typical workflow for creating a package and adding files
    Given a new package is created in memory and configured with a target path
    When the caller adds one or more files using the file management API
    Then the package contains new FileEntry objects for those files
    And the caller writes the package to disk using an appropriate write method

  @documentation
  Scenario: Typical workflow for inspecting a package before extraction
    Given an existing package is open
    When the caller queries for files by path, tag, type, or content hash
    Then the caller receives FileEntry objects with metadata needed to plan extraction
    And directory listings are retrieved from the Metadata API
