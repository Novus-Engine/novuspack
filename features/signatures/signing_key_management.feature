@domain:signatures @m2 @REQ-SIG-039 @spec(api_signatures.md#29-signing-key-management)
Feature: Signing Key Management

  @REQ-SIG-039 @happy
  Scenario: Signing key management provides key handling operations
    Given a NovusPack package
    And a valid context
    When signing key management is used
    Then SigningKey structure provides key information
    And GenerateSigningKey generates new signing keys
    And SaveSigningKey saves keys to files
    And LoadSigningKey loads keys from files
    And context supports cancellation

  @REQ-SIG-039 @happy
  Scenario: GenerateSigningKey generates new signing keys
    Given a NovusPack package
    And a valid context
    And a signature type
    And a security level
    When GenerateSigningKey is called
    Then new SigningKey is generated
    And key contains PrivateKey and PublicKey
    And key contains Type and Level
    And key is ready for use

  @REQ-SIG-039 @happy
  Scenario: SigningKey provides key management
    Given a NovusPack package
    And a valid context
    And a SigningKey
    When key management operations are performed
    Then SaveSigningKey saves key to file
    And LoadSigningKey loads key from file
    And key can be validated and checked for expiration

  @REQ-SIG-039 @error
  Scenario: Signing key management handles errors
    Given a NovusPack package
    When key operations fail
    Then appropriate errors are returned
    And errors follow structured error format
