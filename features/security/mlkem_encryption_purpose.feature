@domain:security @m2 @REQ-SEC-102 @spec(api_security.md#521-purpose)
Feature: ML-KEM encryption purpose defines encryption functionality

  @REQ-SEC-102 @happy
  Scenario: ML-KEM encryption purpose defines functionality
    Given ML-KEM encryption operations
    When encryption or decryption is performed
    Then the purpose defines encryption functionality as specified
    And the behavior matches the ML-KEM encryption purpose specification
    And encapsulation and decapsulation are supported
    And key material is handled securely
