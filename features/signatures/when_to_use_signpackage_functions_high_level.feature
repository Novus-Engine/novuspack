@domain:signatures @m2 @REQ-SIG-037 @spec(api_signatures.md#2822-when-to-use-signpackage-functions-high-level)
Feature: When to Use SignPackage Functions High-Level

  @REQ-SIG-037 @happy
  Scenario: SignPackage functions are used when private key is available
    Given a NovusPack package
    And a private key
    When SignPackage is called with private key
    Then signature is automatically generated
    And signature is automatically added to package
    And function provides convenience of automatic signature generation
    And function handles signature generation internally

  @REQ-SIG-037 @happy
  Scenario: SignPackage functions support standard signature types
    Given a NovusPack package
    And a private key
    When SignPackage functions are used
    Then ML-DSA signature type is supported
    And SLH-DSA signature type is supported
    And PGP signature type is supported
    And X.509 signature type is supported
    And automatic key management is provided

  @REQ-SIG-037 @happy
  Scenario: SignPackage functions provide automatic signature generation
    Given a NovusPack package
    And a private key
    When SignPackageWithKeyFile is called
    Then signature is generated from private key file
    And signature is automatically added to package
    And key management is handled automatically
    And function provides high-level signing interface

  @REQ-SIG-037 @happy
  Scenario: SignPackage functions support key generation
    Given a NovusPack package
    When SignPackageWithNewKey is called
    Then new signing key is generated
    And signature is generated with new key
    And signature is added to package
    And generated key is returned for future use

  @REQ-SIG-037 @error
  Scenario: SignPackage functions handle key errors
    Given a NovusPack package
    And an invalid private key
    When SignPackage is called with invalid key
    Then structured error is returned
    And error indicates invalid key or key format
    And error follows structured error format
