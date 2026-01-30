@domain:file_mgmt @m2 @REQ-FILEMGMT-340 @spec(api_file_mgmt_transform_pipelines.md#23-transformstage-structure)
Feature: TransformStage structure represents individual transformation stages

  @REQ-FILEMGMT-340 @happy
  Scenario: TransformStage represents individual stage
    Given a transformation pipeline with stages
    When TransformStage is used
    Then TransformStage structure represents individual stage with type, input/output sources, completion status, and error tracking
    And the behavior matches the TransformStage structure specification
    And stage type and sources are available
    And completion and error are tracked
