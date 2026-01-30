@domain:metadata @m2 @v2 @REQ-META-077 @spec(api_metadata.md#631-signature-validation)
Feature: Signature Validation for Metadata

  @REQ-META-077 @happy
  Scenario: Signature validation validates metadata-only package signatures
    Given a NovusPack package
    And a metadata-only package
    When signature validation is performed
    Then metadata integrity is validated
    And signatures must validate all metadata files
    And empty content handling is special validated
    And signature scope is clearly defined

  @REQ-META-077 @happy
  Scenario: Signature validation checks metadata integrity
    Given a NovusPack package
    And a metadata-only package with signatures
    When signature validation is performed
    Then all metadata files are validated against signatures
    And metadata integrity is verified
    And validation succeeds if signatures are valid

  @REQ-META-077 @happy
  Scenario: Signature validation handles empty content
    Given a NovusPack package
    And a metadata-only package with no content files
    When signature validation is performed
    Then special validation for packages with no content is performed
    And signature scope clearly defines what gets signed
    And empty content handling is correct

  @REQ-META-077 @error
  Scenario: Signature validation detects invalid signatures
    Given a NovusPack package
    And a metadata-only package with invalid signatures
    When signature validation is performed
    Then validation detects invalid signatures
    And appropriate errors are returned
    And metadata integrity violations are reported
