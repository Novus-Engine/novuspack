@domain:metadata @m2 @REQ-META-126 @spec(api_metadata.md#81-path-metadata-structures)
Feature: DestPath and DestPathWin support relative and absolute paths

  @REQ-META-126 @happy
  Scenario: DestPath and DestPathWin support relative and absolute paths
    Given DestPath or DestPathWin destination overrides
    When paths are set or resolved
    Then relative paths are resolved from default extraction directory
    And absolute paths are supported as specified
    And the behavior matches the path metadata structures specification
    And path resolution is consistent
