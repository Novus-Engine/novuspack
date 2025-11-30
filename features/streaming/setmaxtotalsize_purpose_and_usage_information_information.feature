@domain:streaming @m2 @REQ-STREAM-047 @spec(api_streaming.md#251-purpose)
Feature: SetMaxTotalSize Purpose and Usage

  @REQ-STREAM-047 @happy
  Scenario: BufferPool information purpose defines buffer information access
    Given a NovusPack package
    And a BufferPool
    When buffer pool information purpose is examined
    Then purpose provides additional buffer pool management capabilities
    And purpose enables buffer pool monitoring
    And purpose extends buffer pool functionality beyond basic operations

  @REQ-STREAM-047 @happy
  Scenario: Additional methods provide buffer pool management
    Given a NovusPack package
    And a BufferPool
    When additional buffer pool methods are used
    Then methods provide buffer pool management beyond Get and Put
    And methods enable monitoring and configuration
    And management capabilities support operational needs

  @REQ-STREAM-047 @happy
  Scenario: Buffer pool monitoring capabilities
    Given a NovusPack package
    And a BufferPool
    When monitoring capabilities are used
    Then TotalSize enables memory usage monitoring
    And monitoring supports proactive memory management
    And capabilities enable operational insights

  @REQ-STREAM-047 @happy
  Scenario: Dynamic configuration capabilities
    Given a NovusPack package
    And a BufferPool
    When dynamic configuration is needed
    Then SetMaxTotalSize enables runtime limit adjustment
    And configuration changes take effect immediately
    And capabilities support adaptive memory management

  @REQ-STREAM-047 @happy
  Scenario: Purpose extends core buffer pool operations
    Given a NovusPack package
    And a BufferPool with core operations
    When extended operations are needed
    Then additional methods complement core Get and Put operations
    And extensions provide enhanced functionality
    And purpose defines complete buffer pool interface
