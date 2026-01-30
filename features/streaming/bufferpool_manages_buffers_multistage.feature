@domain:streaming @m2 @REQ-STREAM-067 @spec(api_file_mgmt_transform_pipelines.md#16-intermediate-stage-cleanup)
Feature: BufferPool manages buffers for multi-stage operations with appropriate sizing

  @REQ-STREAM-067 @happy
  Scenario: BufferPool manages buffers for multi-stage operations
    Given multi-stage transformation operations
    When BufferPool is used for intermediate stages
    Then buffers are managed with appropriate sizing for intermediate stages
    And intermediate stage cleanup is performed as specified
    And the behavior matches the intermediate-stage-cleanup specification
    And buffer reuse reduces allocations
