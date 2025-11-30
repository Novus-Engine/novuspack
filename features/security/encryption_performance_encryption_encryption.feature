@domain:security @m2 @REQ-SEC-052 @spec(security.md#722-encryption-performance)
Feature: Encryption Performance

  @REQ-SEC-052 @happy
  Scenario: Encryption performance supports selective encryption
    Given an open NovusPack package
    And a valid context
    And package with mixed encrypted/unencrypted files
    When encryption performance is measured
    Then only sensitive assets are encrypted
    And selective encryption reduces processing overhead
    And unencrypted files accessed without decryption overhead
    And encryption performance matches security requirements

  @REQ-SEC-052 @happy
  Scenario: Encryption performance supports hardware acceleration
    Given an open NovusPack package
    And a valid context
    And hardware acceleration available
    When encryption operations are performed
    Then hardware acceleration improves encryption speed
    And hardware acceleration is used when available
    And encryption correctness is maintained with acceleration

  @REQ-SEC-052 @happy
  Scenario: Encryption performance supports streaming for large files
    Given an open NovusPack package
    And a valid context
    And large file requiring encryption
    When streaming encryption is performed
    Then large files are encrypted using streaming
    And memory usage is optimized during encryption
    And streaming maintains encryption security
    And streaming improves performance for large files

  @REQ-SEC-052 @happy
  Scenario: Encryption performance supports compression integration
    Given an open NovusPack package
    And a valid context
    And file requiring encryption and compression
    When encryption and compression are combined
    Then compression and encryption are combined efficiently
    And combined operations maintain security and compression ratios
    And performance is optimized for dual operations

  @REQ-SEC-052 @happy
  Scenario: Encryption performance optimizes for large archive packages
    Given an open NovusPack package
    And a valid context
    And large archive package with encryption
    When encryption performance is measured
    Then ML-KEM is optimized for large archive packages
    And AES-256-GCM is optimized for large archive packages
    And performance scales with package size

  @REQ-SEC-052 @happy
  Scenario: Encryption performance maintains acceptable overhead
    Given an open NovusPack package
    And a valid context
    And package with encryption
    When encryption performance is measured
    Then encryption overhead is acceptable
    And performance impact is minimized
    And encryption doesn't significantly degrade package operations
