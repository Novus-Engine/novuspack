@domain:file_mgmt @REQ-FILEMGMT-458 @spec(api_file_mgmt_extraction.md#16-extractpath-path-handling) @spec(api_file_mgmt_extraction.md#1-extractpath-package-method)
Feature: ExtractPath Path Handling

  @REQ-FILEMGMT-458 @happy
  Scenario: ExtractPath treats storedPath as a stored package path
    Given an open package
    When ExtractPath is called with storedPath that does not begin with "/"
    Then implementation prefixes "/" before lookup and matching

  @REQ-FILEMGMT-458 @happy
  Scenario: ExtractPath converts separators based on target platform
    Given an open package
    When extracting to Windows targets
    Then forward slashes are converted to backslashes for filesystem operations
    When extracting to Unix-like targets
    Then forward slashes are used

