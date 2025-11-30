@domain:core @m2 @REQ-CORE-046 @spec(api_core.md#error-type-categories)
Feature: Error Type Categories

  @REQ-CORE-046 @happy
  Scenario: ErrTypeValidation categorizes input validation errors
    Given an error operation
    And validation error occurs
    When error is categorized
    Then ErrTypeValidation category is used
    And category includes invalid parameters and format errors
    And category groups validation-related errors

  @REQ-CORE-046 @happy
  Scenario: ErrTypeIO categorizes file I/O errors
    Given an error operation
    And I/O error occurs
    When error is categorized
    Then ErrTypeIO category is used
    And category includes permission errors and disk space issues
    And category groups I/O-related errors

  @REQ-CORE-046 @happy
  Scenario: ErrTypeSecurity categorizes security-related errors
    Given an error operation
    And security error occurs
    When error is categorized
    Then ErrTypeSecurity category is used
    And category includes access denied and authentication failures
    And category groups security-related errors

  @REQ-CORE-046 @happy
  Scenario: ErrTypeCorruption categorizes data corruption errors
    Given an error operation
    And corruption error occurs
    When error is categorized
    Then ErrTypeCorruption category is used
    And category includes checksum failures and integrity violations
    And category groups corruption-related errors

  @REQ-CORE-046 @happy
  Scenario: Error type categories include all major error types
    Given error type categories
    When categories are examined
    Then ErrTypeValidation is defined
    And ErrTypeIO is defined
    And ErrTypeSecurity is defined
    And ErrTypeCorruption is defined
    And ErrTypeUnsupported is defined
    And ErrTypeContext is defined
    And ErrTypeEncryption is defined
    And ErrTypeCompression is defined
    And ErrTypeSignature is defined
