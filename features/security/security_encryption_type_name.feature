@domain:security @m2 @REQ-SEC-005 @spec(api_security.md#31-encryption-type-definition)
Feature: Get encryption type name

  @happy
  Scenario: GetEncryptionTypeName returns name for all valid types
    Given all valid encryption types
    When GetEncryptionTypeName is called for each type
    Then human-readable name is returned for each
    And names are consistent
    And names are documented
