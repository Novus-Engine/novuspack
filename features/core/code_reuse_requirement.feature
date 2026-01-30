@domain:core @m2 @REQ-CORE-189 @spec(api_core.md#1115-code-reuse-requirement)
Feature: Code reuse requirement defines shared helper functions for operations with underlying functionality

  @REQ-CORE-189 @happy
  Scenario: Shared helpers are used for operations with underlying functionality
    Given package operations that share underlying logic
    When the operations are implemented
    Then shared helper functions are used where specified
    And code reuse follows the code reuse requirement
    And the behavior matches the code reuse specification
