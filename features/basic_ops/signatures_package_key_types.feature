@domain:basic_ops @m2 @REQ-API_BASIC-185 @spec(api_basic_operations.md#2115-signatures-package-key-types)
Feature: signatures package key types

  @REQ-API_BASIC-185 @happy
  Scenario: signatures package defines key types for signatures and signature metadata
    Given the signatures package API
    When key types are used
    Then Signature is defined as a key type
    And SignatureInfo is defined as a key type
    And key types support describing signature data and properties
    And key types are used consistently by signing and validation APIs
    And key types align with the signature file format constraints

