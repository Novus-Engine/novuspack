@domain:security @m2 @REQ-SEC-088 @spec(api_security.md#42-generic-encryption-configuration)
Feature: Generic Encryption Configuration

  @REQ-SEC-088 @happy
  Scenario: EncryptionConfig provides type-safe encryption configuration
    Given an open NovusPack package
    And encryption requirements
    When EncryptionConfig is created
    Then type-safe encryption configuration is provided
    And configuration extends Config base type
    And configuration is generic over data type

  @REQ-SEC-088 @happy
  Scenario: EncryptionConfig supports encryption type option
    Given an open NovusPack package
    And encryption configuration
    When EncryptionConfig is configured
    Then EncryptionType option is available
    And encryption algorithm type can be set
    And encryption type option provides type-safe configuration

  @REQ-SEC-088 @happy
  Scenario: EncryptionConfig supports key size option
    Given an open NovusPack package
    And encryption configuration
    When EncryptionConfig is configured
    Then KeySize option is available
    And key size in bits can be set
    And key size option provides type-safe configuration

  @REQ-SEC-088 @happy
  Scenario: EncryptionConfig supports random IV option
    Given an open NovusPack package
    And encryption configuration
    When EncryptionConfig is configured
    Then UseRandomIV option is available
    And random initialization vector setting can be configured
    And random IV option provides type-safe configuration

  @REQ-SEC-088 @happy
  Scenario: EncryptionConfig supports authentication tag option
    Given an open NovusPack package
    And encryption configuration
    When EncryptionConfig is configured
    Then AuthenticationTag option is available
    And authentication tag inclusion can be configured
    And authentication tag option provides type-safe configuration

  @REQ-SEC-088 @happy
  Scenario: EncryptionConfig supports compression level option
    Given an open NovusPack package
    And encryption configuration
    When EncryptionConfig is configured
    Then CompressionLevel option is available
    And compression level for encrypted data can be set
    And compression level option provides type-safe configuration

  @REQ-SEC-088 @happy
  Scenario: EncryptionConfigBuilder provides fluent configuration building
    Given an open NovusPack package
    And encryption configuration requirements
    When EncryptionConfigBuilder is used
    Then NewEncryptionConfigBuilder creates new builder
    And WithEncryptionType sets encryption type
    And WithKeySize sets key size
    And WithRandomIV sets random IV setting
    And WithAuthenticationTag sets authentication tag setting
    And Build creates EncryptionConfig instance

  @REQ-SEC-088 @happy
  Scenario: EncryptionConfigBuilder supports fluent method chaining
    Given an open NovusPack package
    And encryption configuration requirements
    When EncryptionConfigBuilder methods are chained
    Then fluent method chaining is supported
    And configuration can be built in single expression
    And builder pattern enables type-safe configuration

  @REQ-SEC-011 @error
  Scenario: EncryptionConfig operations respect context cancellation
    Given an open NovusPack package
    And a cancelled context
    When encryption configuration operation is called
    Then structured context error is returned
    And error type is context cancellation
