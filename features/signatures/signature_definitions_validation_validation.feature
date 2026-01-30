@skip @domain:signatures @m2 @spec(api_signatures.md#122-signature-validation-process)
Feature: Signature Definitions

# This feature captures signature validation process expectations from the signatures API specification.
# Detailed runnable scenarios live in the dedicated signatures feature files.

  @REQ-SIG-023 @documentation
  Scenario: Signature validation is incremental and includes signature metadata
    Given a signed package with one or more signatures appended
    When signatures are validated
    Then each signature validates all content up to its creation point
    And each signature also validates its own metadata header and signature comment

  @REQ-SIG-010 @documentation
  Scenario: ValidateAllSignatures validates signatures in order
    Given a signed package with multiple signatures
    When ValidateAllSignatures is executed
    Then signatures are read sequentially from the header SignatureOffset
    And the validation result preserves signature ordering
