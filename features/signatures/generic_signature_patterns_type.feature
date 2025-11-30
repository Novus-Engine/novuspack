@domain:signatures @m2 @REQ-SIG-046 @REQ-SIG-048 @REQ-SIG-049 @spec(api_signatures.md#4-generic-signature-patterns)
Feature: Generic Signature Patterns

  @REQ-SIG-046 @happy
  Scenario: Generic signature patterns provide type-safe signature support
    Given a NovusPack package
    When generic signature patterns are used
    Then SignatureStrategy interface provides type-safe signing
    And SignatureConfig provides type-safe configuration
    And SignatureValidator provides type-safe validation
    And patterns extend generic configuration patterns from api_generics.md
    And patterns enable type-safe signature operations

  @REQ-SIG-046 @happy
  Scenario: Generic signature patterns support multiple signature types
    Given a NovusPack package
    When generic signature patterns are used
    Then patterns support ML-DSA signature type
    And patterns support SLH-DSA signature type
    And patterns support PGP signature type
    And patterns support X.509 signature type
    And patterns enable unified signature operations

  @REQ-SIG-048 @happy
  Scenario: Generic signature configuration provides type-safe configuration
    Given a NovusPack package
    When generic signature configuration is used
    Then SignatureConfig structure provides type-safe configuration
    And configuration extends Config for type-safe settings
    And SignatureConfigBuilder provides fluent configuration building
    And configuration supports optional fields with Option types

  @REQ-SIG-048 @happy
  Scenario: Generic signature configuration enables flexible configuration
    Given a NovusPack package
    When signature configuration is built
    Then SignatureType can be configured
    And KeySize can be configured
    And UseTimestamp can be configured
    And IncludeMetadata can be configured
    And CompressionLevel can be configured

  @REQ-SIG-049 @happy
  Scenario: Generic signature validation provides type-safe validation
    Given a NovusPack package
    When generic signature validation is used
    Then SignatureValidator structure provides type-safe validation
    And validator extends Validator for type-safe validation
    And validation rules can be added with AddSignatureRule
    And ValidateSignatureData validates signature data
    And ValidateSignatureKey validates signature key
    And ValidateSignatureFormat validates signature format

  @REQ-SIG-049 @happy
  Scenario: Generic signature validation enables comprehensive validation
    Given a NovusPack package
    When signature validation is performed
    Then validation rules are applied in sequence
    And validation provides type-safe error information
    And validation supports custom validation logic
    And validation integrates with generic validation system
