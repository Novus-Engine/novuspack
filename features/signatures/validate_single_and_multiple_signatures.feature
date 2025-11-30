@domain:signatures @security @m2 @REQ-SIG-002 @spec(api_signatures.md#27-signature-validation)
Feature: Validate single and multiple signatures

  @happy
  Scenario: Validation returns detailed status per signature
    Given a package with multiple signatures
    When I validate the signatures
    Then I should receive status for each signature indicating validity

  @happy
  Scenario: ValidateSignature validates single signature
    Given a package with signature
    When ValidateSignature is called with signature index
    Then signature validation is performed
    And validation status is returned
    And status indicates validity

  @happy
  Scenario: ValidateAllSignatures validates all signatures
    Given a package with multiple signatures
    When ValidateAllSignatures is called
    Then all signatures are validated
    And validation status for each signature is returned
    And overall validation status is provided

  @happy
  Scenario: Signature validation checks content integrity
    Given a package with signature
    When signature validation is performed
    Then all content up to signature is validated
    And signature metadata is validated
    And signature comment is validated

  @happy
  Scenario: Incremental signatures validate correctly
    Given a package with multiple incremental signatures
    When signature validation is performed
    Then first signature validates all content up to it
    And subsequent signatures validate all content including previous signatures
    And all signatures validate correctly

  @error
  Scenario: Signature validation fails for corrupted content
    Given a package with corrupted content
    When signature validation is performed
    Then validation fails
    And structured corruption error is returned
    And error indicates content corruption

  @error
  Scenario: Signature validation fails for invalid signature
    Given a package with invalid signature
    When signature validation is performed
    Then validation fails
    And structured signature error is returned
    And error indicates signature invalidity

  @error
  Scenario: Signature validation respects context cancellation
    Given a package with signatures
    And a cancelled context
    When signature validation is called
    Then structured context error is returned
