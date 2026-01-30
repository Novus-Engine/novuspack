@domain:core @m2 @REQ-CORE-188 @spec(api_core.md#1113-context-usage)
Feature: Context usage defines when context.Context is required for PackageReader methods

  @REQ-CORE-188 @happy
  Scenario: Context usage is defined for PackageReader methods
    Given a PackageReader and methods that accept context
    When context is required for a method
    Then the specification defines when context.Context is required
    And the first parameter is context where required
    And the behavior matches the context usage specification
