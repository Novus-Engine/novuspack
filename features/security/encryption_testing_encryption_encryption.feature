@domain:security @m2 @REQ-SEC-071 @spec(security.md#912-encryption-testing)
Feature: Encryption Testing

  @REQ-SEC-071 @happy
  Scenario: Encryption testing validates encryption and decryption
    Given an open NovusPack package
    And encryption testing configuration
    When encryption testing is performed
    Then encryption and decryption for all algorithms is tested
    And ML-KEM encryption/decryption is validated
    And AES-256-GCM encryption/decryption is validated
    And encryption correctness is verified

  @REQ-SEC-071 @happy
  Scenario: Encryption testing validates key management
    Given an open NovusPack package
    And encryption testing configuration
    When encryption testing is performed
    Then key generation is tested
    And key storage is tested
    And key retrieval is tested
    And key management correctness is verified

  @REQ-SEC-071 @happy
  Scenario: Encryption testing validates performance with various file sizes
    Given an open NovusPack package
    And encryption testing configuration
    When encryption testing is performed
    Then encryption performance with small files is tested
    And encryption performance with medium files is tested
    And encryption performance with large files is tested
    And performance meets requirements across file sizes

  @REQ-SEC-071 @happy
  Scenario: Encryption testing validates compatibility with existing packages
    Given an open NovusPack package
    And encryption testing configuration
    When encryption testing is performed
    Then compatibility with existing AES-encrypted packages is tested
    And backward compatibility is verified
    And existing packages continue to work correctly

  @REQ-SEC-071 @happy
  Scenario: Encryption testing validates error handling
    Given an open NovusPack package
    And encryption testing configuration
    When encryption testing is performed
    Then error handling for invalid keys is tested
    And error handling for corrupted data is tested
    And error handling for missing keys is tested
    And structured errors are properly returned

  @REQ-SEC-071 @happy
  Scenario: Encryption testing provides comprehensive validation
    Given an open NovusPack package
    And encryption testing configuration
    When comprehensive encryption testing is performed
    Then all encryption algorithms are tested
    And all key management operations are tested
    And all error conditions are tested
    And encryption functionality is fully validated
