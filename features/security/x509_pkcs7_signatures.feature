@domain:security @m2 @REQ-SEC-020 @spec(security.md#234-x509pkcs7)
Feature: X.509 PKCS#7 Signatures

  @REQ-SEC-020 @happy
  Scenario: X.509/PKCS#7 uses certificate-based signatures
    Given an open NovusPack package
    And X.509/PKCS#7 signature
    When signature is examined
    Then signature uses X.509 certificates with PKCS#7 signatures
    And certificate-based signatures are provided
    And signature type is 0x04

  @REQ-SEC-020 @happy
  Scenario: X.509/PKCS#7 supports enterprise PKI integration
    Given an open NovusPack package
    And X.509/PKCS#7 signature
    When signature is examined
    Then enterprise PKI integration is supported
    And certificate chains are supported
    And PKI infrastructure is compatible

  @REQ-SEC-020 @happy
  Scenario: X.509/PKCS#7 uses X.509 certificate chains
    Given an open NovusPack package
    And X.509/PKCS#7 signature
    When signature is examined
    Then X.509 certificate chains are used
    And key management uses X.509 certificates
    And certificate chain validation is supported

  @REQ-SEC-020 @happy
  Scenario: X.509/PKCS#7 supports corporate signing
    Given an open NovusPack package
    And X.509/PKCS#7 signature
    When signature use cases are examined
    Then corporate signing is supported
    And enterprise signing workflows are enabled
    And certificate-based signing is provided

  @REQ-SEC-020 @happy
  Scenario: X.509/PKCS#7 supports code signing certificates
    Given an open NovusPack package
    And X.509/PKCS#7 signature
    When signature use cases are examined
    Then code signing certificates are supported
    And code signing workflows are enabled
    And certificate-based code signing is provided

  @REQ-SEC-020 @happy
  Scenario: X.509/PKCS#7 provides fast verification
    Given an open NovusPack package
    And X.509/PKCS#7 signature
    When signature verification is performed
    Then fast verification is achieved
    And certificate chain validation is performed
    And verification performance is optimized

  @REQ-SEC-011 @error
  Scenario: X.509/PKCS#7 signature validation fails with invalid certificate
    Given an open NovusPack package
    And X.509/PKCS#7 signature with invalid certificate
    When signature validation is performed
    Then structured validation error is returned
    And error indicates invalid certificate
