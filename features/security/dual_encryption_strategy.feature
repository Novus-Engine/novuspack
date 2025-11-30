@domain:security @m2 @REQ-SEC-027 @spec(security.md#343-dual-encryption-strategy)
Feature: Dual Encryption Strategy

  @REQ-SEC-027 @happy
  Scenario: Dual encryption strategy uses ML-KEM as default
    Given an open NovusPack package
    And new package creation
    When encryption strategy is examined
    Then ML-KEM is the default encryption method
    And default provides quantum resistance
    And default is optimized for file archives

  @REQ-SEC-027 @happy
  Scenario: Dual encryption strategy maintains AES-256-GCM support
    Given an open NovusPack package
    And encryption configuration
    When encryption strategy is examined
    Then AES-256-GCM is maintained for compatibility
    And AES support provides user preference option
    And traditional encryption remains available

  @REQ-SEC-027 @happy
  Scenario: Dual encryption strategy allows per-file selection
    Given an open NovusPack package
    And file entries
    When encryption type is selected
    Then users can choose encryption type per file
    And ML-KEM can be selected per file
    And AES-256-GCM can be selected per file
    And no encryption can be selected per file

  @REQ-SEC-027 @happy
  Scenario: Dual encryption strategy maintains backward compatibility
    Given an open NovusPack package
    And existing AES-encrypted package
    When package is opened
    Then existing AES-encrypted packages continue to work
    And backward compatibility is maintained
    And existing packages are supported

  @REQ-SEC-027 @happy
  Scenario: Dual encryption strategy supports hybrid approach
    Given an open NovusPack package
    And hybrid encryption configuration
    When encryption is performed
    Then ML-KEM can be used for key exchange
    And AES-256-GCM can be used for data encryption
    And hybrid approach provides flexibility

  @REQ-SEC-027 @happy
  Scenario: Dual encryption strategy optimizes both methods
    Given an open NovusPack package
    And encryption configuration
    When encryption performance is examined
    Then ML-KEM is optimized for large archive packages
    And AES-256-GCM is optimized for large archive packages
    And both methods provide good performance

  @REQ-SEC-027 @error
  Scenario: Dual encryption strategy validation fails with invalid combination
    Given an open NovusPack package
    And invalid encryption combination
    When encryption strategy is validated
    Then structured validation error is returned
    And error indicates invalid encryption combination
