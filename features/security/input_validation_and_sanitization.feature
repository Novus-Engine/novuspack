@domain:security @m2 @REQ-SEC-058 @spec(security.md#82-input-validation-and-sanitization)
Feature: Input Validation and Sanitization

  @REQ-SEC-058 @happy
  Scenario: Input validation validates UTF-8 encoding
    Given an open NovusPack package
    And a valid context
    When comment data with valid UTF-8 encoding is provided
    Then UTF-8 encoding is validated
    And valid UTF-8 data is accepted
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-058 @happy
  Scenario: Input validation enforces length limits
    Given an open NovusPack package
    And a valid context
    When package comment within length limit is provided
    Then package comment length is validated
    And comment within 65535 bytes is accepted
    When signature comment within length limit is provided
    Then signature comment length is validated
    And comment within 4095 bytes is accepted

  @REQ-SEC-058 @happy
  Scenario: Input validation filters dangerous characters
    Given an open NovusPack package
    And a valid context
    When comment data with control characters is provided
    Then control characters (0x00-0x1F, 0x7F-0x9F) are filtered
    And null bytes (0x00) are prohibited
    And dangerous characters are removed or escaped

  @REQ-SEC-058 @error
  Scenario: Input validation rejects invalid UTF-8 encoding
    Given an open NovusPack package
    And a valid context
    And comment data with invalid UTF-8 encoding
    When comment validation is performed
    Then encoding validation error is returned
    And error indicates invalid UTF-8 encoding
    And error follows structured error format

  @REQ-SEC-058 @error
  Scenario: Input validation rejects oversized comments
    Given an open NovusPack package
    And a valid context
    And package comment exceeding 65535 bytes
    When comment validation is performed
    Then length validation error is returned
    And error indicates comment exceeds maximum length
    And error follows structured error format

  @REQ-SEC-058 @happy
  Scenario: Input validation detects dangerous injection patterns
    Given an open NovusPack package
    And a valid context
    When comment data with script injection patterns is provided
    Then script injection patterns (<script>, javascript:, vbscript:) are detected
    And patterns are sanitized or rejected
    When comment data with command injection patterns is provided
    Then command injection patterns (;, |, &, $, `, \) are detected
    And patterns are sanitized or rejected
