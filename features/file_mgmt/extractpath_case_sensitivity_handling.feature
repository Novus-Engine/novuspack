@domain:file_mgmt @REQ-FILEMGMT-459 @spec(api_file_mgmt_extraction.md#161-extractpath-case-sensitivity-handling) @spec(api_file_mgmt_extraction.md#1-extractpath-package-method)
Feature: ExtractPath Case Sensitivity Handling

  @REQ-FILEMGMT-459 @error
  Scenario: ExtractPath fails on case conflicts on case-insensitive filesystems
    Given an open package
    And the package contains two paths that differ only by case
    When ExtractPath is called targeting a case-insensitive filesystem
    Then extraction fails with an error
    And no files are written to the destination directory
    And the error indicates case conflict or ambiguous path resolution

