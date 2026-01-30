@domain:file_mgmt @m2 @REQ-FILEMGMT-124 @spec(api_file_mgmt_best_practices.md#13-best-practices)
Feature: File Management Best Practices

  @REQ-FILEMGMT-124 @happy
  Scenario: Best practices cover file path management
    Given an open NovusPack package
    And a valid context
    When file path operations are performed
    Then consistent path formats are used
    And paths are validated before use
    And path management best practices are followed

  @REQ-FILEMGMT-124 @happy
  Scenario: Best practices cover encryption management
    Given an open NovusPack package
    And a valid context
    When encryption operations are performed
    Then appropriate encryption types are chosen
    And secure key management is used
    And encryption best practices are followed

  @REQ-FILEMGMT-124 @happy
  Scenario: Best practices cover performance considerations
    Given an open NovusPack package
    And a valid context
    When file operations are performed
    Then patterns are used for bulk operations
    And streaming is used for large files
    And appropriate context timeouts are set
