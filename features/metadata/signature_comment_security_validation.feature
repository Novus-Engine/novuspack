@domain:metadata @security @m2 @v2 @REQ-META-010 @spec(api_metadata.md#13-signature-comment-security)
Feature: Signature Comment Security Validation

  @happy
  Scenario: ValidateSignatureComment validates signature comment for security
    Given a signature comment with potential security issues
    When ValidateSignatureComment is called
    Then signature comment is validated for security
    And validation results indicate issues found
    And security threats are detected

  @happy
  Scenario: SanitizeSignatureComment sanitizes signature comment content
    Given a signature comment with potential injection
    When SanitizeSignatureComment is called
    Then signature comment is sanitized
    And injection patterns are removed
    And safe comment content is returned

  @happy
  Scenario: CheckSignatureCommentLength validates signature comment length
    Given a signature comment
    When CheckSignatureCommentLength is called
    Then signature comment length is validated
    And length limits are enforced
    And excessive length is detected

  @happy
  Scenario: AuditSignatureComment audits signature comment for security logging
    Given a signature comment
    When AuditSignatureComment is called
    Then signature comment is audited
    And security events are logged
    And audit trail is created

  @error
  Scenario: Invalid signature comment is rejected
    Given a signature comment with security issues
    When ValidateSignatureComment is called
    Then structured validation error is returned
    And error indicates security issue

  @error
  Scenario: Excessive signature comment length is rejected
    Given a signature comment exceeding length limit
    When CheckSignatureCommentLength is called
    Then structured validation error is returned
    And error indicates length limit exceeded
