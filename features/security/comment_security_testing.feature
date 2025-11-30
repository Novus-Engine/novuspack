@domain:security @m2 @REQ-SEC-066 @spec(security.md#841-comment-security-testing)
Feature: Comment Security Testing

  @REQ-SEC-066 @happy
  Scenario: Comment security testing tests injection resistance
    Given an open NovusPack package
    And comment security testing configuration
    When injection testing is performed
    Then resistance to script injection attacks is tested
    And resistance to command injection attacks is tested
    And resistance to SQL injection attacks is tested
    And resistance to path traversal attacks is tested
    And all injection attack types are validated

  @REQ-SEC-066 @happy
  Scenario: Comment security testing tests encoding resistance
    Given an open NovusPack package
    And comment security testing configuration
    When encoding testing is performed
    Then various character encodings are tested
    And Unicode attacks are tested
    And Unicode normalization attacks are tested
    And homograph attacks are tested
    And encoding-based vulnerabilities are detected

  @REQ-SEC-066 @happy
  Scenario: Comment security testing tests length limits
    Given an open NovusPack package
    And comment security testing configuration
    When length testing is performed
    Then maximum comment lengths are tested
    And oversized comment lengths are tested
    And buffer overflow prevention is validated
    And length limit enforcement is verified

  @REQ-SEC-066 @happy
  Scenario: Comment security testing tests pattern detection
    Given an open NovusPack package
    And comment security testing configuration
    When pattern testing is performed
    Then known malicious patterns are tested
    And dangerous sequences are tested
    And pattern detection is validated
    And malicious content is properly identified

  @REQ-SEC-066 @happy
  Scenario: Comment security testing validates sanitization
    Given an open NovusPack package
    And comment security testing configuration
    When sanitization testing is performed
    Then proper sanitization of dangerous content is verified
    And sanitization methods are validated
    And sanitized content is safe for storage
    And sanitization prevents injection attacks

  @REQ-SEC-066 @happy
  Scenario: Comment security testing provides comprehensive coverage
    Given an open NovusPack package
    And comment security testing configuration
    When comprehensive testing is performed
    Then all injection types are tested
    And all encoding attacks are tested
    And all length scenarios are tested
    And all malicious patterns are tested
    And sanitization is thoroughly validated
