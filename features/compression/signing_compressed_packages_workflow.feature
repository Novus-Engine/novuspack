@domain:compression @m2 @REQ-COMPR-024 @spec(api_package_compression.md#1011-supported-operation)
Feature: Signing compressed packages workflow

  @REQ-COMPR-024 @happy
  Scenario: Signing compressed packages follows the supported operation workflow
    Given a compressed package
    And a signing workflow configured for the package
    When the package is signed
    Then signing operates on the compressed package as a supported operation
    And signing produces a signature block for the package
    And signature data can be validated by signature-aware consumers
    And the signed package becomes immutable under signed package rules
    And the workflow follows the documented supported signing operation

