@domain:file_mgmt @REQ-PIPELINE-034 @REQ-PIPELINE-035 @spec(api_file_mgmt_transform_pipelines.md#18-error-recovery-and-resume)
Feature: Transformation Pipeline Error Recovery

  @REQ-PIPELINE-034 @happy
  Scenario: Stage failures store error in TransformStage.Error
    Given a transformation pipeline
    When stage 2 fails with error
    Then error is stored in TransformStage.Error for stage 2
    And Completed flag remains false for stage 2
    And pipeline overall Completed remains false
    And subsequent stages are not executed

  @REQ-PIPELINE-034 @happy
  Scenario: Partial temporary files cleaned up automatically on failure
    Given a pipeline with stage 1 completed
    And stage 2 fails midway
    When cleanup is triggered
    Then stage 1 temporary file is removed
    And partial stage 2 temporary file is removed
    And all temporary files are cleaned up
    And cleanup does not fail if files already removed

  @REQ-PIPELINE-034 @happy
  Scenario: Pipeline Completed remains false until all stages succeed
    Given a pipeline with 3 stages
    When stages 1 and 2 complete successfully
    And stage 3 has not yet executed
    Then TransformPipeline.Completed is false
    When stage 3 completes successfully
    Then TransformPipeline.Completed is true

  @REQ-PIPELINE-035 @happy
  Scenario: System supports retry of failed stage
    Given a pipeline where stage 2 failed
    When stage 2 is retried
    Then stage 2 executes again
    And stage reads from stage 1 output
    And stage can succeed on retry
    And pipeline can complete after retry

  @REQ-PIPELINE-035 @happy
  Scenario: System supports retry of entire pipeline
    Given a pipeline that failed at stage 2
    When entire pipeline is retried from beginning
    Then all stages execute from stage 1
    And previous temporary files are cleaned up first
    And new temporary files are created
    And pipeline can complete successfully on retry
