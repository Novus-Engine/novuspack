@domain:security @m2 @v2 @REQ-SEC-070 @spec(security.md#911-signature-testing)
Feature: Signature Testing

  @REQ-SEC-070 @happy
  Scenario: Signature testing tests multiple signature creation
    Given an open NovusPack package
    And signature testing configuration
    When multiple signature creation testing is performed
    Then creating packages with multiple signatures is tested
    And incremental signing is validated
    And signature chain integrity is tested

  @REQ-SEC-070 @happy
  Scenario: Signature testing validates all signature types
    Given an open NovusPack package
    And signature testing configuration
    When signature validation testing is performed
    Then validation of all signature types is tested
    And ML-DSA signature validation is tested
    And SLH-DSA signature validation is tested
    And PGP signature validation is tested
    And X.509 signature validation is tested

  @REQ-SEC-070 @happy
  Scenario: Signature testing tests invalid signature handling
    Given an open NovusPack package
    And signature testing configuration
    When invalid signature handling testing is performed
    Then handling of invalid signatures is tested
    And handling of corrupted signatures is tested
    And error handling is validated

  @REQ-SEC-070 @happy
  Scenario: Signature testing tests performance with large numbers of signatures
    Given an open NovusPack package
    And signature testing configuration
    When performance testing is performed
    Then signature performance with large numbers of signatures is tested
    And performance scales appropriately
    And performance meets requirements

  @REQ-SEC-070 @happy
  Scenario: Signature testing provides comprehensive validation
    Given an open NovusPack package
    And signature testing configuration
    When comprehensive signature testing is performed
    Then all signature creation scenarios are tested
    And all signature validation scenarios are tested
    And all error conditions are tested
    And signature functionality is fully validated
