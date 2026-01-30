@domain:signatures @m2 @v2 @REQ-SIG-068 @spec(api_signatures.md#2851-secure-signing-operations-with-runtimesecret)
Feature: Secure signing operations with runtime/secret

  @REQ-SIG-068 @happy
  Scenario: Signing operations with key material execute within runtime/secret.Do
    Given a signing operation that uses private key material
    When signature generation is performed
    Then key access and signature computation occur within runtime/secret.Do
    And key material and intermediate values are not retained outside the secret context

  @REQ-SIG-068 @happy
  Scenario: Signing operations load keys within runtime/secret.Do
    Given a signing operation that loads private keys from files
    When key material is loaded for signing
    Then key loading occurs within runtime/secret.Do
    And key material does not persist in memory outside the secret context

