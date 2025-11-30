@domain:security @m2 @REQ-SEC-023 @spec(security.md#322-aes-256-gcm)
Feature: AES-256-GCM Encryption

  @REQ-SEC-023 @happy
  Scenario: AES-256-GCM uses Advanced Encryption Standard algorithm
    Given an open NovusPack package
    And AES-256-GCM encryption configuration
    When encryption is performed
    Then Advanced Encryption Standard with Galois/Counter Mode is used
    And algorithm provides strong encryption
    And GCM mode provides authentication

  @REQ-SEC-023 @happy
  Scenario: AES-256-GCM uses 256-bit keys
    Given an open NovusPack package
    And AES-256-GCM encryption configuration
    When encryption keys are generated
    Then 256-bit keys are used
    And key size provides strong security
    And keys are securely generated

  @REQ-SEC-023 @happy
  Scenario: AES-256-GCM provides high-speed encryption
    Given an open NovusPack package
    And large file for encryption
    And AES-256-GCM encryption configuration
    When encryption is performed
    Then high-speed encryption is achieved
    And performance is optimized for large files
    And encryption completes efficiently

  @REQ-SEC-023 @happy
  Scenario: AES-256-GCM provides industry standard compatibility
    Given an open NovusPack package
    And AES-256-GCM encryption configuration
    When encryption is performed
    Then industry standard encryption is used
    And maximum compatibility is provided
    And encrypted data is compatible with standard tools

  @REQ-SEC-023 @happy
  Scenario: AES-256-GCM is available as per-file encryption option
    Given an open NovusPack package
    And file entry
    When encryption type is selected
    Then AES-256-GCM can be selected per file
    And per-file encryption selection works
    And encryption type is preserved

  @REQ-SEC-027 @happy
  Scenario: AES-256-GCM supports backward compatibility
    Given an open NovusPack package
    And existing AES-encrypted package
    When package is opened
    Then AES-encrypted packages continue to work
    And backward compatibility is maintained
    And existing packages are supported

  @REQ-SEC-008 @error
  Scenario: AES-256-GCM validation fails with invalid key size
    Given an open NovusPack package
    And invalid encryption key size
    When AES-256-GCM encryption is attempted
    Then structured validation error is returned
    And error indicates invalid key size
