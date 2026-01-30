@domain:security @m2 @REQ-SEC-104 @spec(api_security.md#523-returns)
Feature: ML-KEM encryption returns define encryption results

  @REQ-SEC-104 @happy
  Scenario: ML-KEM encryption returns define results
    Given ML-KEM encryption or decryption operations
    When operations complete successfully
    Then returns define encryption results as specified
    And the behavior matches the ML-KEM encryption returns specification
    And ciphertext or plaintext is returned correctly
    And error returns indicate failure conditions
