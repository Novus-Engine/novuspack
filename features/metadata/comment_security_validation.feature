@domain:metadata @m2 @REQ-META-006 @spec(api_metadata.md#12-comment-security-validation)
Feature: Comment Security Validation

  @happy
  Scenario: ValidateComment validates comment content for security
    Given a package comment with potential security issues
    When ValidateComment is called
    Then comment is validated for security issues
    And validation results indicate issues found
    And malicious patterns are detected

  @happy
  Scenario: SanitizeComment sanitizes comment content
    Given a package comment with potential injection
    When SanitizeComment is called
    Then comment is sanitized
    And injection patterns are removed
    And safe comment content is returned

  @happy
  Scenario: ValidateCommentEncoding validates UTF-8 encoding
    Given a package comment
    When ValidateCommentEncoding is called
    Then comment encoding is validated
    And valid UTF-8 encoding is confirmed
    And invalid encoding is detected

  @happy
  Scenario: CheckCommentLength validates comment length
    Given a package comment
    When CheckCommentLength is called
    Then comment length is validated
    And length limits are enforced
    And excessive length is detected

  @happy
  Scenario: DetectInjectionPatterns detects malicious patterns
    Given a package comment with potential injection
    When DetectInjectionPatterns is called
    Then malicious patterns are detected
    And detected patterns are reported
    And safe content passes detection

  @error
  Scenario: Invalid comment encoding is rejected
    Given a package comment with invalid encoding
    When ValidateCommentEncoding is called
    Then structured validation error is returned
    And error indicates encoding issue

  @error
  Scenario: Excessive comment length is rejected
    Given a package comment exceeding length limit
    When CheckCommentLength is called
    Then structured validation error is returned
    And error indicates length limit exceeded
