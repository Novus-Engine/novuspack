@domain:security @m2 @REQ-SEC-060 @spec(security.md#822-dangerous-pattern-detection)
Feature: Dangerous Pattern Detection

  @REQ-SEC-060 @happy
  Scenario: Dangerous pattern detection identifies script injection patterns
    Given an open NovusPack package
    And content with script injection patterns
    When dangerous pattern detection is performed
    Then script injection patterns are detected
    And patterns like "<script>", "javascript:", "vbscript:" are identified
    And security threats are identified

  @REQ-SEC-060 @happy
  Scenario: Dangerous pattern detection identifies command injection patterns
    Given an open NovusPack package
    And content with command injection patterns
    When dangerous pattern detection is performed
    Then command injection patterns are detected
    And shell metacharacters like ";", "|", "&", "$", "`", "\" are identified
    And security threats are identified

  @REQ-SEC-060 @happy
  Scenario: Dangerous pattern detection identifies SQL injection patterns
    Given an open NovusPack package
    And content with SQL injection patterns
    When dangerous pattern detection is performed
    Then SQL injection patterns are detected
    And SQL keywords and special characters are identified
    And security threats are identified

  @REQ-SEC-060 @happy
  Scenario: Dangerous pattern detection identifies path traversal patterns
    Given an open NovusPack package
    And content with path traversal patterns
    When dangerous pattern detection is performed
    Then path traversal patterns are detected
    And patterns like "../", "..\\", absolute paths are identified
    And security threats are identified

  @REQ-SEC-060 @happy
  Scenario: Dangerous pattern detection identifies Unicode attack patterns
    Given an open NovusPack package
    And content with Unicode attack patterns
    When dangerous pattern detection is performed
    Then Unicode attack patterns are detected
    And Unicode normalization attacks and homograph attacks are identified
    And security threats are identified

  @REQ-SEC-060 @happy
  Scenario: Dangerous pattern detection identifies control characters
    Given an open NovusPack package
    And content with control characters
    When dangerous pattern detection is performed
    Then control characters are detected
    And all control characters and escape sequences are identified
    And security threats are identified

  @REQ-SEC-060 @error
  Scenario: Dangerous pattern detection prevents injection attacks
    Given an open NovusPack package
    And content with dangerous patterns
    When content is processed
    Then dangerous patterns are prevented
    And injection attacks are blocked
    And structured validation error is returned
