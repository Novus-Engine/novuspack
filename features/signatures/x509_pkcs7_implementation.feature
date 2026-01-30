@domain:signatures @m2 @v2 @REQ-SIG-032 @spec(api_signatures.md#26-x509pkcs7-implementation)
Feature: X.509 PKCS#7 Implementation

  @REQ-SIG-032 @happy
  Scenario: X.509/PKCS#7 implementation provides certificate-based signatures
    Given a NovusPack package
    And a valid context
    When X.509/PKCS#7 implementation is used
    Then algorithm uses X.509 certificates with PKCS#7 signatures
    And multiple key types are supported (RSA, ECDSA, EdDSA)
    And full certificate chain validation is provided
    And fast verification with certificate chain validation is provided
    And X.509 certificate and private key file support is provided
    And context supports cancellation

  @REQ-SIG-032 @happy
  Scenario: X.509/PKCS#7 provides certificate chain validation
    Given a NovusPack package
    And a valid context
    When X.509/PKCS#7 validation is performed
    Then full certificate chain is validated
    And certificate chain validation ensures trust
    And enterprise PKI integration is supported

  @REQ-SIG-032 @happy
  Scenario: X.509/PKCS#7 supports certificate and key files
    Given a NovusPack package
    And a valid context
    When X.509/PKCS#7 signing is performed
    Then X.509 certificate files are supported
    And private key files are supported
    And passphrase protection is supported

  @REQ-SIG-032 @error
  Scenario: X.509/PKCS#7 implementation handles errors
    Given a NovusPack package
    When X.509/PKCS#7 operations fail
    Then appropriate errors are returned
    And errors follow structured error format
