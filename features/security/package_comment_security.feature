@domain:security @m2 @REQ-SEC-056 @spec(security.md#812-package-comment-security)
Feature: Package Comment Security

  @REQ-SEC-056 @happy
  Scenario: Package comment security validates UTF-8 encoding
    Given an open NovusPack package
    And a valid context
    And package comment data
    When package comment security validates encoding
    Then all comment data must be valid UTF-8 encoding
    And invalid UTF-8 encoding is rejected
    And encoding validation prevents encoding-based attacks

  @REQ-SEC-056 @happy
  Scenario: Package comment security enforces length limits
    Given an open NovusPack package
    And a valid context
    And package comment data
    When package comment security enforces limits
    Then maximum 65535 bytes per package comment is enforced
    And oversized comments are rejected
    And length limits prevent buffer overflow attacks

  @REQ-SEC-056 @happy
  Scenario: Package comment security filters control characters
    Given an open NovusPack package
    And a valid context
    And package comment data with control characters
    When package comment security filters characters
    Then control characters (0x00-0x1F, 0x7F-0x9F) are filtered
    And null bytes (0x00) are prohibited
    And dangerous characters are removed or escaped

  @REQ-SEC-056 @happy
  Scenario: Package comment security prevents script injection
    Given an open NovusPack package
    And a valid context
    And package comment data with script tags
    When package comment security validates content
    Then HTML/XML script tags are escaped or filtered
    And script injection attacks are prevented
    And comment content is safe

  @REQ-SEC-056 @happy
  Scenario: Package comment security prevents command injection
    Given an open NovusPack package
    And a valid context
    And package comment data with command characters
    When package comment security validates content
    Then shell command characters are escaped
    And command injection attacks are prevented
    And comment content is safe for processing

  @REQ-SEC-056 @error
  Scenario: Package comment security rejects invalid comments
    Given an open NovusPack package
    And a valid context
    And package comment data with security issues
    When package comment security validates content
    Then validation error is returned
    And error indicates specific security issue
    And error follows structured error format
