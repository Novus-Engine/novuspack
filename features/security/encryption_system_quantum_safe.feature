@domain:security @m2 @REQ-SEC-021 @spec(security.md#3-encryption-system)
Feature: Encryption system provides quantum-safe encryption

  @REQ-SEC-021 @happy
  Scenario: Encryption system provides quantum-safe encryption
    Given a package or file requiring encryption
    When encryption system is used
    Then quantum-safe encryption is provided as specified
    And the behavior matches the encryption system specification
    And ML-KEM and traditional encryption are supported
    And algorithm selection follows security guidelines
