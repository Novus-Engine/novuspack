@domain:security @m2 @REQ-SEC-108 @spec(api_security.md#532-purpose)
Feature: ML-KEM key information purpose defines key metadata functionality

  @REQ-SEC-108 @happy
  Scenario: ML-KEM key information purpose defines functionality
    Given ML-KEM key information operations
    When key metadata is accessed
    Then the purpose defines key metadata functionality as specified
    And the behavior matches the ML-KEM key information purpose specification
    And key metadata is retrievable without exposing key material
    And key validity and expiration are queryable
