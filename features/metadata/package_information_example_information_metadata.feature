@domain:metadata @m2 @REQ-META-041 @spec(metadata.md#222-package-information-example)
Feature: Package Information Example

  @REQ-META-041 @happy
  Scenario: Package information example demonstrates package metadata structure
    Given a NovusPack package
    When package information example is examined
    Then example demonstrates package information fields
    And example shows name, version, description, author, license
    And example shows created and modified timestamps

  @REQ-META-041 @happy
  Scenario: Package information example shows basic package fields
    Given a NovusPack package
    And package information example
    When example is examined
    Then name field shows "MyAwesomeGame"
    And version field shows "1.2.0"
    And description field shows package description
    And author field shows "GameStudio Inc."
    And license field shows "Commercial"

  @REQ-META-041 @happy
  Scenario: Package information example shows timestamps
    Given a NovusPack package
    And package information example
    When timestamps are examined
    Then created field shows ISO8601 timestamp like "2024-01-15T10:30:00Z"
    And modified field shows ISO8601 timestamp like "2024-01-20T14:45:00Z"
    And timestamps demonstrate temporal tracking

  @REQ-META-041 @error
  Scenario: Package information example validates field formats
    Given a NovusPack package
    When invalid package information is provided
    Then field validation detects format violations
    And appropriate errors are returned
