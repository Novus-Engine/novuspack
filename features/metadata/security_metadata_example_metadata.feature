@domain:metadata @m2 @REQ-META-055 @spec(metadata.md#security-metadata-example)
Feature: Security Metadata Example

  @REQ-META-055 @happy
  Scenario: Security metadata example demonstrates security fields
    Given a NovusPack package
    When security metadata example is examined
    Then example demonstrates encryption_level field
    And example demonstrates signature_type field
    And example demonstrates security_scan field
    And example demonstrates trusted_source field

  @REQ-META-055 @happy
  Scenario: Security metadata example shows encryption level
    Given a NovusPack package
    And security metadata example
    When encryption_level example is examined
    Then field shows "ML-KEM Level 3"
    And example demonstrates quantum-safe encryption level
    And level indicates encryption security level

  @REQ-META-055 @happy
  Scenario: Security metadata example shows signature type
    Given a NovusPack package
    And security metadata example
    When signature_type example is examined
    Then field shows "ML-DSA Level 3"
    And example demonstrates quantum-safe signature type
    And type indicates signature security level

  @REQ-META-055 @happy
  Scenario: Security metadata example shows security flags
    Given a NovusPack package
    And security metadata example
    When security flags example is examined
    Then security_scan field shows true
    And trusted_source field shows true
    And flags indicate security validation status

  @REQ-META-055 @error
  Scenario: Security metadata example validates field formats
    Given a NovusPack package
    When invalid security metadata is provided
    Then field validation detects format violations
    And appropriate errors are returned
