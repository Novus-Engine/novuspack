@domain:metadata @m2 @REQ-META-027 @spec(metadata.md#13-path-metadata-system)
Feature: Path metadata system provides path-based tag inheritance

  @REQ-META-027 @happy
  Scenario: Path metadata system provides path-based tag inheritance
    Given a path metadata system for package paths
    When path metadata is queried or updated
    Then path-based tag inheritance is provided as specified
    And inheritance rules are applied consistently
    And the behavior matches the path metadata system specification
    And path hierarchy is respected for tag resolution
