@domain:metadata @m2 @REQ-META-053 @spec(metadata.md#package-information)
Feature: Package information provides package metadata

  @REQ-META-053 @happy
  Scenario: Package information provides package metadata
    Given a package with package information
    When package metadata is accessed
    Then package information provides package metadata as specified
    And metadata fields are populated correctly
    And the behavior matches the package information specification
    And metadata is consistent with package state
