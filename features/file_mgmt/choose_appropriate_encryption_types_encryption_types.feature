@domain:file_mgmt @m2 @REQ-FILEMGMT-129 @spec(api_file_mgmt_best_practices.md#1321-choose-appropriate-encryption-types)
Feature: Choose Appropriate Encryption Types

  @REQ-FILEMGMT-129 @happy
  Scenario: Appropriate encryption types are chosen for sensitive data
    Given an open NovusPack package
    And a valid context
    And sensitive data to encrypt
    When encryption type is selected
    Then strong encryption is chosen for sensitive data
    And appropriate algorithm is selected
    And security requirements are met

  @REQ-FILEMGMT-129 @happy
  Scenario: Encryption type selection considers security requirements
    Given an open NovusPack package
    And a valid context
    And files with different security needs
    When encryption types are selected
    Then high-security encryption is used for critical data
    And standard encryption may be used for less sensitive data
    And encryption type matches security requirements

  @REQ-FILEMGMT-129 @happy
  Scenario: Encryption type selection considers performance requirements
    Given an open NovusPack package
    And a valid context
    And performance requirements
    When encryption types are selected
    Then encryption performance is considered
    And encryption overhead is acceptable
    And performance requirements are balanced with security
