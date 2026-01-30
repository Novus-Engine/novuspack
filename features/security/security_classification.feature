@domain:security @m2 @REQ-SEC-038 @spec(security.md#511-security-classification)
Feature: Security classification defines file security levels

  @REQ-SEC-038 @happy
  Scenario: Security classification defines file security levels
    Given a file or package with security metadata
    When security classification is applied or queried
    Then file security levels are defined as specified
    And the behavior matches the security classification specification
    And classification is consistent with access control
    And security levels are persisted and retrievable
