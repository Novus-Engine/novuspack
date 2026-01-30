@domain:core @m2 @REQ-CORE-026 @spec(api_core.md#1044-sentinel-error-compatibility)
Feature: Sentinel Error Compatibility

  @REQ-CORE-026 @happy
  Scenario: Sentinel errors are still supported and can be wrapped
    Given an operation
    And a sentinel error
    When sentinel error is used
    Then sentinel errors are still supported
    And sentinel errors can be wrapped with structured information
    And existing code using sentinel errors continues to work

  @REQ-CORE-026 @happy
  Scenario: Sentinel errors can be converted to structured errors
    Given a sentinel error
    When sentinel error is wrapped
    Then sentinel error can be converted to structured error
    And conversion preserves error information
    And structured error wraps sentinel error as cause

  @REQ-CORE-026 @happy
  Scenario: Sentinel error compatibility enables migration
    Given existing code using sentinel errors
    When migration to structured errors is needed
    Then sentinel errors continue to work
    And code can migrate gradually
    And compatibility is maintained during migration

  @REQ-CORE-026 @happy
  Scenario: Is method supports sentinel error matching
    Given a PackageError
    And a sentinel error target
    When Is method is called
    Then Is method checks underlying cause for sentinel error
    And sentinel error matching works with error unwrapping
    And compatibility with errors.Is is maintained
