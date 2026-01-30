@domain:security @m2 @REQ-SEC-016 @spec(security.md#22-signature-types-and-algorithms)
Feature: Signature types and algorithms define supported signature algorithms

  @REQ-SEC-016 @happy
  Scenario: Signature types and algorithms define supported algorithms
    Given security validation or signing context
    When signature types and algorithms are queried
    Then supported signature algorithms are defined as specified
    And the behavior matches the signature types and algorithms specification
    And algorithm selection is consistent with security requirements
    And quantum-safe and traditional algorithms are documented
