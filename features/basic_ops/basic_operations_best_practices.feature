@domain:basic_ops @m2 @REQ-API_BASIC-077 @spec(api_basic_operations.md#9-best-practices)
Feature: Basic Operations Best Practices

  @REQ-API_BASIC-077 @happy
  Scenario: Best practices define recommended usage patterns
    Given NovusPack package operations
    When best practices are followed
    Then package lifecycle is managed properly
    And error handling is appropriate
    And resource management is correct

  @REQ-API_BASIC-077 @happy
  Scenario: Best practices include package lifecycle management
    Given package creation and usage
    When lifecycle best practices are applied
    Then packages are created properly
    And packages are opened correctly
    And packages are closed with defer
    And state is checked before operations

  @REQ-API_BASIC-077 @happy
  Scenario: Best practices include error handling patterns
    Given package operations that may fail
    When error handling best practices are applied
    Then errors are always checked
    And structured errors are used
    And context cancellation is handled
    And error types guide handling strategy

  @REQ-API_BASIC-077 @happy
  Scenario: Best practices include resource management
    Given package operations with resources
    When resource management best practices are applied
    Then context is used for resource management
    And cleanup errors are handled gracefully
    And defer ensures cleanup
    And resource leaks are prevented
