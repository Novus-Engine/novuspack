@domain:file_mgmt @m2 @REQ-FILEMGMT-303 @REQ-FILEMGMT-304 @spec(api_file_mgmt_file_entry.md#5-path-management) @spec(api_generics.md#1335-path-conversion-on-extractionviewing)
Feature: Path Display Format

  Background:
    Given an open NovusPack package
    And a valid context

  @REQ-FILEMGMT-303 @happy
  Scenario: ListFiles strips leading slash from displayed paths
    Given files stored with leading slash in package
    When ListFiles is called
    Then returned paths do not have leading slash
    And paths are in display format
    And storage format is not exposed to user

  @REQ-FILEMGMT-303 @happy
  Scenario: GetPrimaryPath returns path without leading slash
    Given a FileEntry with stored path "/assets/textures/button.png"
    When GetPrimaryPath is called
    Then returned path is "assets/textures/button.png"
    And leading slash is stripped for display
    And path is user-friendly format

  @REQ-FILEMGMT-304 @happy
  Scenario: Error messages strip leading slash from paths
    Given an operation that fails with path reference
    When error message is generated
    Then error message shows path without leading slash
    And path is user-friendly in error
    And internal storage format is not exposed

  @REQ-FILEMGMT-303 @happy
  Scenario: File extraction uses paths without leading slash
    Given files stored with leading slash
    When files are extracted
    Then extraction paths do not have leading slash
    And extracted paths are relative to extraction directory
    And filesystem operations use display format

  @REQ-CORE-060 @happy
  Scenario: Platform-specific path display strips leading slash
    Given a stored path "/assets/textures/ui/button.png"
    When path is converted for display on Unix
    Then displayed path is "assets/textures/ui/button.png"
    And leading slash is not shown to user

  @REQ-CORE-060 @happy
  Scenario: Windows path display strips leading slash and converts separators
    Given a stored path "/assets/textures/ui/button.png"
    When path is converted for display on Windows
    Then displayed path is "assets\textures\ui\button.png"
    And leading slash is not shown to user
    And separators are converted to backslashes

  @REQ-GEN-029 @happy
  Scenario: API methods returning paths strip leading slash
    Given various API methods that return paths
    When paths are returned from API
    Then all returned paths have leading slash stripped
    And paths are in user-facing format
    And internal storage format is abstracted

  @REQ-GEN-028 @happy
  Scenario: Path normalization adds leading slash for storage
    Given user input path "assets/textures/button.png"
    When path is normalized for storage
    Then normalized path is "/assets/textures/button.png"
    And leading slash is added for internal storage
    But display methods will strip it for user output
