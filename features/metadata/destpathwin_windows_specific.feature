@domain:metadata @m2 @REQ-META-127 @spec(api_metadata.md#81-path-metadata-structures)
Feature: DestPathWin is used for Windows-specific destination paths

  @REQ-META-127 @happy
  Scenario: DestPathWin used for Windows-specific paths
    Given destination overrides on Windows
    When DestPathWin is set or only DestPath is absolute
    Then DestPathWin is used for Windows-specific destination paths
    And if only DestPath is absolute on Windows, root is treated as C:\
    And the behavior matches the path metadata structures specification
    And platform-specific paths are handled correctly
