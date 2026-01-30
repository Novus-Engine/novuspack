@domain:signatures @m2 @v2 @REQ-SIG-064 @REQ-SIG-065 @REQ-SIG-066 @REQ-SIG-067 @REQ-SIG-069 @REQ-SIG-070 @REQ-SIG-071 @REQ-SIG-072 @spec(api_signatures.md#4151-getkey-behavior) @spec(api_signatures.md#4152-setkey-behavior) @spec(api_signatures.md#4153-error-conditions) @spec(api_signatures.md#4154-operation-requirements) @spec(api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret)
Feature: SigningKey GetKey and SetKey behavior

  @REQ-SIG-069 @happy
  Scenario: GetKey returns a copy and should be used within runtime/secret.Do
    Given a SigningKey with a key set
    When GetKey is called
    Then a copy of the key value is returned
    And slice key types are deep-copied
    And the returned key copy should be used within runtime/secret.Do

  @REQ-SIG-070 @happy
  Scenario: SetKey overwrites existing key and validates type and expiration
    Given a SigningKey with an existing key set
    When SetKey is called with a new key value
    Then the stored key value is overwritten
    And the stored key value is copied (deep copy for slices)
    And CreatedAt is not automatically updated
    And SetKey validates KeyType compatibility and expiration

  @REQ-SIG-071 @error
  Scenario: GetKey returns ErrTypeSignature PackageError for missing invalid expired or mismatched key
    Given a SigningKey that is not set or is invalid or expired or mismatched
    When GetKey is called
    Then a *PackageError is returned
    And error type is ErrTypeSignature
    And error message indicates the specific failure reason

  @REQ-SIG-072 @error
  Scenario: Operations using SigningKey re-validate before use and return ErrTypeSignature on failure
    Given a SigningKey used for signing operations
    When an operation is executed using the key
    Then key presence is checked before use
    And key validity and expiry are checked before use
    And *PackageError with ErrTypeSignature is returned when validation fails

  @REQ-SIG-064 @REQ-SIG-065 @REQ-SIG-066 @REQ-SIG-067 @happy
  Scenario: SigningKey private key material is handled within runtime/secret.Do
    Given SigningKey instances that contain private key material
    When key material is retrieved set or used for signing operations
    Then runtime/secret.Do is used for key handling operations
    And key material is not retained outside the secret context

