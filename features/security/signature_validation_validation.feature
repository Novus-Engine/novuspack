@domain:security @m2 @v2 @REQ-SEC-049 @spec(security.md#712-signature-validation)
Feature: Signature Validation

  @REQ-SEC-049 @happy
  Scenario: Signature validation validates multiple signatures
    Given an open NovusPack package
    And a valid context
    And package with multiple signatures
    When signature validation is performed
    Then all signatures are validated, not just one
    And multiple signature validation provides comprehensive verification
    And validation results indicate validity of each signature

  @REQ-SEC-049 @happy
  Scenario: Signature validation provides trust verification
    Given an open NovusPack package
    And a valid context
    And package with signatures
    When signature validation is performed
    Then trust verification is implemented
    And trust chain validation is performed
    And trusted signatures are identified

  @REQ-SEC-049 @happy
  Scenario: Signature validation provides timestamp verification
    Given an open NovusPack package
    And a valid context
    And package with timestamped signatures
    When signature validation is performed
    Then timestamp verification is performed
    And signature timestamps are validated
    And timestamp validity is confirmed

  @REQ-SEC-049 @happy
  Scenario: Signature validation provides revocation checking
    Given an open NovusPack package
    And a valid context
    And package with signatures
    When signature validation is performed
    Then revocation checking is implemented
    And revoked certificates are detected
    And revoked keys are detected
    And revocation status is reported

  @REQ-SEC-049 @happy
  Scenario: Signature validation provides comprehensive verification
    Given an open NovusPack package
    And a valid context
    And package with signatures
    When signature validation is performed
    Then multiple signatures are validated
    And trust verification is performed
    And timestamp verification is performed
    And revocation checking is performed
    And comprehensive verification ensures signature authenticity
