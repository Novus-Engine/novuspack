@domain:core @m2 @REQ-CORE-172 @spec(api_core.md#1034-error-transformation-functions) @spec(api_core.md#maperror-function)
Feature: Error transformation functions provide MapError for error transformation

  @REQ-CORE-172 @happy
  Scenario: MapError transforms error context
    Given an error and a transformation function
    When MapError is called
    Then the error context is transformed as specified
    And the transformation follows the MapError specification
    And the behavior matches the error transformation specification
