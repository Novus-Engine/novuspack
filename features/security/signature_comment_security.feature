@domain:security @m2 @v2 @REQ-SEC-057 @spec(security.md#813-signature-comment-security)
Feature: Signature Comment Security

  @REQ-SEC-057 @happy
  Scenario: Signature comment security validates UTF-8 encoding
    Given an open NovusPack package
    And a valid context
    And signature comment data
    When signature comment security validates encoding
    Then all signature comment data must be valid UTF-8 encoding
    And invalid UTF-8 encoding is rejected
    And encoding validation prevents encoding-based attacks

  @REQ-SEC-057 @happy
  Scenario: Signature comment security enforces length limits
    Given an open NovusPack package
    And a valid context
    And signature comment data
    When signature comment security enforces limits
    Then maximum 4095 bytes per signature comment is enforced
    And oversized comments are rejected
    And length limits prevent buffer overflow attacks

  @REQ-SEC-057 @happy
  Scenario: Signature comment security filters dangerous characters
    Given an open NovusPack package
    And a valid context
    And signature comment data with control characters
    When signature comment security filters characters
    Then same filtering as package comments is applied
    And additional restrictions are applied
    And dangerous characters are removed or escaped

  @REQ-SEC-057 @happy
  Scenario: Signature comment security prevents executable content
    Given an open NovusPack package
    And a valid context
    And signature comment data with potential executable content
    When signature comment security validates content
    Then no executable code, scripts, or commands are allowed
    And executable content is rejected
    And comment content is safe

  @REQ-SEC-057 @happy
  Scenario: Signature comment security includes comments in signature validation
    Given an open NovusPack package
    And a valid context
    And package with signature comments
    When signature comment security is applied
    Then comments are included in signature validation to prevent tampering
    And signature validation ensures comment integrity
    And comment tampering is detected

  @REQ-SEC-057 @happy
  Scenario: Signature comment security provides audit trail
    Given an open NovusPack package
    And a valid context
    And package with signature comments
    When signature comment security is applied
    Then all signature comments are logged for security auditing
    And audit trail tracks comment modifications
    And security events are recorded

  @REQ-SEC-057 @error
  Scenario: Signature comment security rejects invalid comments
    Given an open NovusPack package
    And a valid context
    And signature comment data with security issues
    When signature comment security validates content
    Then validation error is returned
    And error indicates specific security issue
    And error follows structured error format
