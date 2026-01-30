@domain:security @m2 @REQ-SEC-022 @spec(security.md#321-ml-kem-crystals-kyber)
Feature: ML-KEM provides quantum-safe key exchange

  @REQ-SEC-022 @happy
  Scenario: ML-KEM provides quantum-safe key exchange
    Given encryption context requiring key exchange
    When ML-KEM is used for key exchange
    Then quantum-safe key exchange is provided as specified
    And the behavior matches the ML-KEM specification
    And key generation and encapsulation are supported
    And security levels are configurable
