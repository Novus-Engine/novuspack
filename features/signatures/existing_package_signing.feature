@domain:signatures @m2 @v2 @REQ-SIG-033 @spec(api_signatures.md#28-existing-package-signing)
Feature: Existing Package Signing

  @REQ-SIG-033 @happy
  Scenario: Existing package signing supports signing existing packages
    Given a NovusPack package
    And a valid context
    And an existing unsigned package
    When SignPackage is called
    Then package is signed with generated signature
    Then signature is added using AddSignature internally
    And signature bit is set
    And SignatureOffset is set
    And context supports cancellation

  @REQ-SIG-033 @happy
  Scenario: SignPackage generates signature from private key
    Given a NovusPack package
    And a valid context
    And an existing package
    And a private key
    When SignPackage is called with private key
    Then signature is generated using private key
    Then signature data is added using AddSignature
    And signature validates all package content

  @REQ-SIG-033 @happy
  Scenario: SignPackageWithKeyFile loads key and signs package
    Given a NovusPack package
    And a valid context
    And an existing package
    And a key file
    When SignPackageWithKeyFile is called
    Then key is loaded from file
    Then signature is generated using loaded key
    Then signature is added to package

  @REQ-SIG-033 @error
  Scenario: Existing package signing handles errors
    Given a NovusPack package
    When signature generation fails
    Then appropriate errors are returned
    And errors follow structured error format
