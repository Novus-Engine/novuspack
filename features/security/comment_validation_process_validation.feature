@domain:security @m2 @REQ-SEC-059 @spec(security.md#821-comment-validation-process)
Feature: Comment Validation Process

  @REQ-SEC-059 @happy
  Scenario: Comment validation process performs encoding check
    Given an open NovusPack package
    And a valid context
    And comment data requiring validation
    When comment validation process is performed
    Then encoding check verifies valid UTF-8 encoding
    And invalid UTF-8 encoding is detected
    And encoding validation is first step in process

  @REQ-SEC-059 @happy
  Scenario: Comment validation process performs length validation
    Given an open NovusPack package
    And a valid context
    And comment data of various sizes
    When comment validation process is performed
    Then length validation checks against maximum length limits
    And package comments checked against 65535 byte limit
    And signature comments checked against 4095 byte limit
    And oversized comments are rejected

  @REQ-SEC-059 @happy
  Scenario: Comment validation process performs character filtering
    Given an open NovusPack package
    And a valid context
    And comment data with characters
    When comment validation process is performed
    Then character filtering removes or escapes dangerous characters
    And control characters are filtered
    And null bytes are prohibited

  @REQ-SEC-059 @happy
  Scenario: Comment validation process performs content scanning
    Given an open NovusPack package
    And a valid context
    And comment data with potential issues
    When comment validation process is performed
    Then content scanning scans for potential injection patterns
    And script injection patterns are detected
    And command injection patterns are detected
    And SQL injection patterns are detected

  @REQ-SEC-059 @happy
  Scenario: Comment validation process performs sanitization
    Given an open NovusPack package
    And a valid context
    And comment data requiring sanitization
    When comment validation process is performed
    Then sanitization applies appropriate method based on content type
    And HTML escaping is applied for HTML content
    And URL encoding is applied for URL content
    And character removal is applied for dangerous content

  @REQ-SEC-059 @happy
  Scenario: Comment validation process performs final validation
    Given an open NovusPack package
    And a valid context
    And comment data after all validation steps
    When comment validation process completes
    Then final validation check occurs before storage
    And all validation steps must pass
    And validated comment is safe for storage

  @REQ-SEC-059 @error
  Scenario: Comment validation process rejects invalid comments
    Given an open NovusPack package
    And a valid context
    And comment data with validation failures
    When comment validation process detects issues
    Then validation error is returned
    And error indicates specific validation failure
    And error follows structured error format
