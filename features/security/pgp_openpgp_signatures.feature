@domain:security @m2 @REQ-SEC-019 @spec(security.md#233-pgp-openpgp)
Feature: PGP OpenPGP Signatures

  @REQ-SEC-019 @happy
  Scenario: PGP provides traditional OpenPGP signatures
    Given an open NovusPack package
    And a valid context
    And PGP signature implementation
    When PGP signature is examined
    Then algorithm uses OpenPGP standard (RFC 4880)
    And algorithm provides traditional signature support
    And algorithm is compatible with existing PGP infrastructure

  @REQ-SEC-019 @happy
  Scenario: PGP supports multiple key types
    Given an open NovusPack package
    And a valid context
    And PGP signature with RSA key
    When PGP signature is created
    Then RSA keys (2048-4096 bits) are supported
    When PGP signature with DSA key is created
    Then DSA keys are supported
    When PGP signature with ECDSA key is created
    Then ECDSA keys (P-256/P-384/P-521) are supported
    When PGP signature with EdDSA key is created
    Then EdDSA keys are supported

  @REQ-SEC-019 @happy
  Scenario: PGP provides fast verification performance
    Given an open NovusPack package
    And a valid context
    And package with PGP signatures
    When PGP signature verification is performed
    Then verification performance is fast
    And verification scales with package size
    And verification maintains OpenPGP compatibility

  @REQ-SEC-019 @happy
  Scenario: PGP supports PGP keyring integration
    Given an open NovusPack package
    And a valid context
    And PGP keyring file
    When PGP signature validation uses keyring
    Then PGP keyring integration is supported
    And passphrase protection is supported
    And keyring provides key management

  @REQ-SEC-019 @happy
  Scenario: PGP provides developer and community verification
    Given an open NovusPack package
    And a valid context
    And package with PGP signatures
    When PGP signatures are used for verification
    Then developer signatures provide package authenticity
    And community verification is supported
    And multiple developer signatures are supported

  @REQ-SEC-019 @happy
  Scenario: PGP provides moderate signing speed
    Given an open NovusPack package
    And a valid context
    And PGP signature creation
    When PGP signing is performed
    Then signing speed is moderate
    And signing performance is acceptable for package creation
    And context supports cancellation
    And context supports timeout handling
