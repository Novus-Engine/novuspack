@domain:security @m2 @REQ-SEC-053 @spec(security.md#8-comment-security-and-injection-prevention)
Feature: Comment security and injection prevention protect against injection attacks

  @REQ-SEC-053 @happy
  Scenario: Comment security protects against injection
    Given package or signature comments
    When comment content is validated or sanitized
    Then comment security and injection prevention protect as specified
    And the behavior matches the comment security specification
    And dangerous patterns are detected and rejected
    And sanitization is applied before persistence
