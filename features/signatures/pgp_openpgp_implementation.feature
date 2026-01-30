@domain:signatures @m2 @v2 @REQ-SIG-031 @spec(api_signatures.md#25-pgp-openpgp-implementation)
Feature: PGP OpenPGP Implementation

  @REQ-SIG-031 @happy
  Scenario: PGP implementation provides OpenPGP signature support
    Given a NovusPack package
    And a valid context
    When PGP implementation is used
    Then algorithm follows OpenPGP standard (RFC 4880)
    And multiple key types are supported (RSA, DSA, ECDSA, EdDSA)
    And variable key sizes are supported
    And fast verification is provided
    And PGP keyring support is provided with passphrase protection
    And context supports cancellation

  @REQ-SIG-031 @happy
  Scenario: PGP supports multiple key types and sizes
    Given a NovusPack package
    And a valid context
    When PGP key types are examined
    Then RSA keys support 2048-4096 bits
    And ECDSA keys support P-256/P-384/P-521 curves
    And DSA and EdDSA keys are supported

  @REQ-SIG-031 @happy
  Scenario: PGP provides keyring support
    Given a NovusPack package
    And a valid context
    When PGP keyring is used
    Then PGP keyring support is provided
    And passphrase protection is supported
    And key management follows OpenPGP standards

  @REQ-SIG-031 @error
  Scenario: PGP implementation handles errors
    Given a NovusPack package
    When PGP operations fail
    Then appropriate errors are returned
    And errors follow structured error format
