@domain:security @m2 @REQ-SEC-012 @spec(security.md#11-security-layers)
Feature: Security layers define multi-layered protection architecture

  @REQ-SEC-012 @happy
  Scenario: Security layers define multi-layered protection
    Given a package or encryption context
    When security layers are applied
    Then multi-layered protection architecture is defined as specified
    And layers are applied consistently
    And the behavior matches the security layers specification
    And protection covers encryption, integrity, and access control
