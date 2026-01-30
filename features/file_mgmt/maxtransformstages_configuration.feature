@domain:file_mgmt @REQ-FILEMGMT-347 @REQ-PIPELINE-018 @REQ-PIPELINE-027 @spec(api_file_mgmt_transform_pipelines.md#19-configuration-options) @spec(api_file_mgmt_addition.md#211-multi-stage-transformation-pipelines)
Feature: MaxTransformStages Configuration

  @REQ-PIPELINE-027 @REQ-FILEMGMT-347 @happy
  Scenario: MaxTransformStages default is 10
    Given AddFileOptions without MaxTransformStages specified
    When file with pipeline is added
    Then MaxTransformStages defaults to 10
    And pipeline with up to 10 stages is allowed
    And default covers typical 3-stage operations with headroom

  @REQ-PIPELINE-027 @happy
  Scenario: MaxTransformStages configurable via AddFileOptions
    Given AddFileOptions with MaxTransformStages set to 15
    When file requiring pipeline is added
    Then pipeline depth limit is 15 stages
    And configuration overrides default
    And allows custom limits for advanced use cases

  @REQ-PIPELINE-018 @REQ-FILEMGMT-347 @happy
  Scenario: Pipeline within MaxTransformStages succeeds
    Given MaxTransformStages set to 10
    And file requiring 3 transformation stages
    When AddFile is called
    Then pipeline is created with 3 stages
    And all stages execute successfully
    And file is added to package

  @REQ-PIPELINE-019 @REQ-FILEMGMT-347 @error
  Scenario: Pipeline exceeding MaxTransformStages fails
    Given MaxTransformStages set to 2
    And file requiring 3 transformation stages
    When AddFile is called
    Then validation error is returned
    And error type is ErrTypeValidation
    And error indicates pipeline exceeds maximum stages
    And error includes actual count (3) and limit (2)

  @REQ-PIPELINE-018 @happy
  Scenario: MaxTransformStages prevents resource exhaustion
    Given MaxTransformStages configuration
    When complex transformation sequence is requested
    Then limit prevents unbounded pipeline depth
    And protects against memory leaks
    And prevents runaway resource usage
