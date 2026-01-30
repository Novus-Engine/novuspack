@domain:file_mgmt @m2 @REQ-FILEMGMT-131 @spec(api_file_mgmt_best_practices.md#13-performance-considerations)
Feature: File Management Performance Considerations

  @REQ-FILEMGMT-131 @happy
  Scenario: Performance considerations define performance optimization patterns
    Given an open NovusPack package
    When performance considerations are applied
    Then performance optimization patterns are defined
    And best practices guide performance optimization

  @REQ-FILEMGMT-131 @happy
  Scenario: Use patterns for bulk operations improves performance
    Given an open package
    And multiple files to add
    When AddFilePattern is used for bulk operations
    Then performance is improved compared to individual AddFile calls
    And bulk operations are more efficient
    And resource usage is optimized

  @REQ-FILEMGMT-131 @happy
  Scenario: Handle large files with streaming for better performance
    Given an open package
    And very large file to process
    When streaming API is used
    Then large files are handled efficiently
    And memory usage is optimized
    And performance benefits from streaming approach

  @REQ-FILEMGMT-131 @happy
  Scenario: Use appropriate context timeouts prevents indefinite blocking
    Given an open package
    And operation that may take long time
    When context.WithTimeout is used
    Then indefinite blocking is prevented
    And timeout handling is appropriate
    And performance characteristics are predictable

  @REQ-FILEMGMT-131 @happy
  Scenario: Performance considerations guide resource usage
    Given an open package
    When performance best practices are followed
    Then resource usage is optimized
    And memory usage is efficient
    And CPU usage is appropriate

  @REQ-FILEMGMT-131 @error
  Scenario: Performance considerations handle timeout scenarios
    Given an open package
    And operation exceeds timeout
    When operation is performed
    Then context timeout error is returned
    And timeout is handled gracefully
    And error follows structured error format
