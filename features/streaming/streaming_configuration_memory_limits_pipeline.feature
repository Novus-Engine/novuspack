@domain:streaming @m2 @REQ-STREAM-068 @spec(api_file_mgmt_transform_pipelines.md#14-disk-space-management)
Feature: Streaming configuration includes memory limits for pipeline streaming

  @REQ-STREAM-068 @happy
  Scenario: Streaming configuration includes memory limits for pipeline
    Given streaming configuration for multi-stage transformations
    When pipeline streaming is performed
    Then configuration includes memory limits to prevent resource exhaustion
    And disk space management is applied during multi-stage transformations
    And the behavior matches the disk-space-management specification
    And resource exhaustion is prevented
