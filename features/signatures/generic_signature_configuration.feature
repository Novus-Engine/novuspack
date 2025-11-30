@domain:signatures @m2 @REQ-SIG-048 @spec(api_signatures.md#42-generic-signature-configuration)
Feature: Generic Signature Configuration

  @REQ-SIG-048 @happy
  Scenario: SignatureConfig provides type-safe signature configuration
    Given a NovusPack package
    When SignatureConfig structure is examined
    Then structure extends Config for type-safe configuration
    And structure contains SignatureType option with signature algorithm type
    And structure contains KeySize option with key size in bits
    And structure contains UseTimestamp option for timestamp inclusion
    And structure contains IncludeMetadata option for metadata inclusion
    And structure contains CompressionLevel option for compression level

  @REQ-SIG-048 @happy
  Scenario: SignatureConfigBuilder provides fluent configuration building
    Given a NovusPack package
    When SignatureConfigBuilder is used
    Then NewSignatureConfigBuilder creates new builder instance
    And WithSignatureType sets signature type configuration
    And WithKeySize sets key size configuration
    And WithTimestamp sets timestamp configuration
    And WithMetadata sets metadata inclusion configuration
    And Build returns configured SignatureConfig instance

  @REQ-SIG-048 @happy
  Scenario: Generic signature configuration enables type-safe operations
    Given a NovusPack package
    When signature configuration is used
    Then configuration provides type-safe signature settings
    And configuration enables fluent builder pattern
    And configuration supports optional fields with Option types
    And configuration enables flexible signature configuration

  @REQ-SIG-048 @error
  Scenario: Generic signature configuration handles invalid configuration
    Given a NovusPack package
    When invalid configuration is provided
    Then validation error is returned
    And error indicates invalid configuration parameter
    And error follows structured error format
