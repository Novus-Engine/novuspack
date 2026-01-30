@domain:file_mgmt @m2 @REQ-FILEMGMT-339 @spec(api_file_mgmt_transform_pipelines.md#22-transformpipeline-structure)
Feature: TransformPipeline structure tracks multi-stage transformation pipelines

  @REQ-FILEMGMT-339 @happy
  Scenario: TransformPipeline tracks ordered stages and completion
    Given a FileEntry with multi-stage transformation
    When TransformPipeline is used
    Then TransformPipeline structure tracks ordered stages, current stage index, and completion status
    And the behavior matches the TransformPipeline structure specification
    And stage list and current index are available
    And completion status is tracked
