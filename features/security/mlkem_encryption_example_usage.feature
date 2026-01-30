@domain:security @m2 @REQ-SEC-106 @spec(api_security.md#525-example-usage)
Feature: ML-KEM encryption example usage demonstrates encryption

  @REQ-SEC-106 @happy
  Scenario: ML-KEM encryption example usage demonstrates encryption
    Given ML-KEM encryption API
    When example usage or documentation is followed
    Then example usage demonstrates encryption as specified
    And the behavior matches the ML-KEM example usage specification
    And encryption and decryption workflows are demonstrated
    And key handling and context are shown correctly
