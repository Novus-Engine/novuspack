@domain:security @m2 @REQ-SEC-096 @spec(api_security.md#521-purpose)
Feature: ML-KEM key generation purpose defines key creation functionality

  @REQ-SEC-096 @happy
  Scenario: ML-KEM key generation purpose defines functionality
    Given ML-KEM key generation operations
    When key generation is performed
    Then the purpose defines key creation functionality as specified
    And the behavior matches the ML-KEM key generation purpose specification
    And key pairs are created with correct security levels
    And context cancellation is supported
