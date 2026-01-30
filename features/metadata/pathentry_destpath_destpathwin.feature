@domain:metadata @m2 @REQ-META-125 @spec(api_metadata.md#81-path-metadata-structures)
Feature: PathMetadataEntry stores DestPath and DestPathWin for extraction overrides

  @REQ-META-125 @happy
  Scenario: PathMetadataEntry stores DestPath and DestPathWin
    Given a PathMetadataEntry with destination overrides
    When extraction destination is configured
    Then DestPath and DestPathWin fields store persistent destination overrides
    And the behavior matches the path metadata structures specification
    And overrides are applied during extraction
    And relative and absolute paths are supported
