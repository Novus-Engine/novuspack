@domain:security @m2 @REQ-SEC-025 @spec(security.md#34-encryption-implementation-details)
Feature: Encryption Implementation Details

  @REQ-SEC-025 @happy
  Scenario: Encryption implementation details provide ML-KEM key management
    Given an open NovusPack package
    And a valid context
    And encryption implementation with ML-KEM
    When encryption implementation details are examined
    Then ML-KEM key generation at security levels 1-5 is supported
    And ML-KEM encryption and decryption operations are provided
    And ML-KEM key access for public keys and security levels is supported
    And ML-KEM key structure for secure storage is defined

  @REQ-SEC-025 @happy
  Scenario: Encryption implementation details provide per-file encryption operations
    Given an open NovusPack package
    And a valid context
    And per-file encryption system
    When encryption implementation details are examined
    Then AddFileWithEncryption adds files with specific encryption types
    And encryption types support None, AES-256-GCM, and ML-KEM
    And GetFileEncryptionType retrieves encryption type for files
    And GetEncryptedFiles lists all encrypted files

  @REQ-SEC-025 @happy
  Scenario: Encryption implementation details provide dual encryption strategy
    Given an open NovusPack package
    And a valid context
    And dual encryption implementation
    When encryption implementation details are examined
    Then ML-KEM is default encryption method for new packages
    And ML-KEM provides full quantum resistance
    And AES-256-GCM is maintained for compatibility
    And per-file selection allows ML-KEM, AES, or none
    And backward compatibility with existing packages is maintained

  @REQ-SEC-025 @happy
  Scenario: Encryption implementation details provide implementation considerations
    Given an open NovusPack package
    And a valid context
    And encryption implementation
    When encryption implementation details are examined
    Then CIRCL library is used for quantum-safe algorithms
    And standard library (crypto/aes, crypto/cipher) is used for AES
    And secure storage for both key types is provided
    And both encryption methods are optimized for large archives
    And backward compatibility is maintained

  @REQ-SEC-025 @happy
  Scenario: Encryption implementation details support hybrid approach
    Given an open NovusPack package
    And a valid context
    And hybrid encryption approach
    When encryption implementation details are examined
    Then ML-KEM is used for key exchange
    And AES-256-GCM is used for data encryption
    And hybrid approach combines quantum-safe and traditional encryption
    And performance is optimized for package archives

  @REQ-SEC-025 @error
  Scenario: Encryption implementation details handle invalid encryption operations
    Given an open NovusPack package
    And a valid context
    And invalid encryption operation
    When encryption implementation details are used incorrectly
    Then structured error is returned
    And error indicates specific encryption failure
    And error follows structured error format
