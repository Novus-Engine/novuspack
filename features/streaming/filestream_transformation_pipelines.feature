@domain:streaming @m2 @REQ-STREAM-066 @spec(api_file_mgmt_transform_pipelines.md#1-multi-stage-transformation-pipelines)
Feature: FileStream supports streaming through transformation pipelines for memory-efficient processing

  @REQ-STREAM-066 @happy
  Scenario: FileStream supports streaming through transformation pipelines
    Given a large file requiring sequential transformations
    When FileStream is used with transformation pipelines
    Then streaming supports memory-efficient processing of large files
    And sequential transformations are applied as specified
    And the behavior matches the multi-stage transformation pipelines specification
    And memory usage is bounded
