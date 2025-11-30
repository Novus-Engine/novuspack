@domain:core @m2 @REQ-CORE-028 @spec(api_core.md#115-migration-from-sentinel-errors)
Feature: Migration from Sentinel Errors

  @REQ-CORE-028 @happy
  Scenario: Structured error system works alongside sentinel errors
    Given code using sentinel errors
    When structured error system is introduced
    Then structured error system works alongside existing sentinel errors
    And existing code using sentinel errors continues to work
    And coexistence is maintained

  @REQ-CORE-028 @happy
  Scenario: New code can take advantage of structured errors
    Given new code being developed
    When structured errors are available
    Then new code can take advantage of structured errors
    And better error handling and debugging is enabled
    And migration path is provided

  @REQ-CORE-028 @happy
  Scenario: Sentinel errors can be wrapped with structured information
    Given sentinel errors from existing code
    When migration to structured errors occurs
    Then sentinel errors can be wrapped with structured information
    And WrapError converts sentinel errors to structured errors
    And gradual migration is supported

  @REQ-CORE-028 @happy
  Scenario: Migration provides better error categorization
    Given error handling migration
    When structured errors are used
    Then better error categorization is provided
    And errors are grouped by type for easier handling
    And error types enable appropriate responses

  @REQ-CORE-028 @happy
  Scenario: Migration provides rich error context for debugging
    Given error handling migration
    When structured errors are used
    Then rich error context is provided for debugging
    And additional context fields are available
    And debugging information is enhanced

  @REQ-CORE-028 @happy
  Scenario: Migration maintains backward compatibility
    Given codebase migration to structured errors
    When migration is performed
    Then backward compatibility is maintained
    And existing sentinel error code continues to work
    And gradual adoption is supported
