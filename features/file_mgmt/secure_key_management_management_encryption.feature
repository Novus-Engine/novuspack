@domain:file_mgmt @m2 @REQ-FILEMGMT-130 @spec(api_file_mgmt_best_practices.md#1322-secure-key-management)
Feature: Secure Key Management

  @REQ-FILEMGMT-130 @happy
  Scenario: Secure key management protects encryption keys
    Given an open NovusPack package
    And a valid context
    And encryption keys are used
    When key management operations are performed
    Then encryption keys are protected
    And key storage is secure
    And key access is controlled

  @REQ-FILEMGMT-130 @happy
  Scenario: Secure key management clears sensitive data
    Given an open NovusPack package
    And a valid context
    And encryption keys are used
    When keys are no longer needed
    Then sensitive key data is cleared from memory
    And key cleanup is performed
    And security is maintained

  @REQ-FILEMGMT-130 @happy
  Scenario: Secure key management follows security best practices
    Given an open NovusPack package
    And a valid context
    When key management is performed
    Then security best practices are followed
    And key lifecycle is properly managed
    And key security is maintained
