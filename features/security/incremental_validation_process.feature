@domain:security @m2 @v2 @REQ-SEC-076 @spec(api_security.md#12-incremental-validation-process)
Feature: Incremental Validation Process

  @REQ-SEC-076 @happy
  Scenario: Incremental validation process validates each signature against complete package state
    Given an open NovusPack package
    And package with multiple signatures
    When incremental validation is performed
    Then each signature validates complete package state up to its position
    And signature validation covers all content before signature
    And incremental validation ensures integrity

  @REQ-SEC-076 @happy
  Scenario: Incremental validation process ensures chain validation
    Given an open NovusPack package
    And package with signature chain
    When incremental validation is performed
    Then chain validation ensures no signatures were removed
    And chain validation ensures no signatures were modified
    And signature chain integrity is maintained

  @REQ-SEC-076 @happy
  Scenario: Incremental validation process provides detailed validation status
    Given an open NovusPack package
    And package with multiple signatures
    When incremental validation is performed
    Then results provide detailed validation status for each signature
    And individual signature results are available
    And validation status is comprehensive

  @REQ-SEC-076 @happy
  Scenario: Incremental validation process validates signatures in order
    Given an open NovusPack package
    And package with multiple signatures
    When incremental validation is performed
    Then signatures are validated in order
    And validation order matches signature creation order
    And order validation ensures chain integrity

  @REQ-SEC-076 @happy
  Scenario: Incremental validation process detects missing signatures
    Given an open NovusPack package
    And package with signature chain
    And missing signature in chain
    When incremental validation is performed
    Then missing signature is detected
    And validation error indicates missing signature
    And chain validation fails

  @REQ-SEC-076 @happy
  Scenario: Incremental validation process detects modified signatures
    Given an open NovusPack package
    And package with signature chain
    And modified signature in chain
    When incremental validation is performed
    Then modified signature is detected
    And validation error indicates signature modification
    And chain validation fails

  @REQ-SEC-011 @error
  Scenario: Incremental validation process respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When incremental validation operation is called
    Then structured context error is returned
    And error type is context cancellation
