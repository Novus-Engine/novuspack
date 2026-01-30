@domain:file_mgmt @m2 @REQ-FILEMGMT-134 @spec(api_file_mgmt_best_practices.md#1333-use-appropriate-context-timeouts)
Feature: File Management: Use Appropriate Context Timeouts

  @REQ-FILEMGMT-134 @happy
  Scenario: Appropriate context timeouts prevent hanging operations
    Given an open NovusPack package
    And file operations to perform
    When context timeouts are set appropriately
    Then operations complete within timeout period
    And hanging operations are prevented
    And timeout errors are handled gracefully

  @REQ-FILEMGMT-134 @happy
  Scenario: Context timeout values are selected based on operation type
    Given an open NovusPack package
    And different types of file operations
    When context timeouts are configured
    Then quick operations use short timeouts
    And long-running operations use longer timeouts
    And timeout values match operation duration

  @REQ-FILEMGMT-134 @happy
  Scenario: Context timeouts are set for all file management operations
    Given an open NovusPack package
    When file management operations are performed
    Then context timeouts are set for all operations
    And timeout handling is consistent
    And operations respect timeout values
