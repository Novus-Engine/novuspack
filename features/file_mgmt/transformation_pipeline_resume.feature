@domain:file_mgmt @REQ-PIPELINE-032 @REQ-PIPELINE-033 @spec(api_file_mgmt_transform_pipelines.md#18-error-recovery-and-resume)
Feature: Transformation Pipeline Resume

  @REQ-PIPELINE-032 @REQ-FILEMGMT-351 @happy
  Scenario: ResumeTransformation resumes from last completed stage
    Given a FileEntry with transformation pipeline
    And pipeline was interrupted after stage 2 completed
    When ResumeTransformation is called
    Then last completed stage (stage 2) is identified
    And stage 2 output is verified to exist
    And execution continues from stage 3
    And ProcessingState is updated at each stage

  @REQ-PIPELINE-033 @happy
  Scenario: Resume checks last completed stage
    Given a pipeline interrupted at various stages
    When ResumeTransformation is called
    Then system checks CurrentStage index
    And identifies last successfully completed stage
    And determines next stage to execute

  @REQ-PIPELINE-033 @happy
  Scenario: Resume verifies output exists before continuing
    Given a pipeline with stage 2 completed
    And stage 2 output temporary file exists
    When ResumeTransformation is called
    Then system verifies stage 2 output file exists
    And system uses existing output for stage 3
    And avoids re-executing stage 2

  @REQ-PIPELINE-033 @happy
  Scenario: Resume retries if output missing
    Given a pipeline with stage 2 marked completed
    But stage 2 output temporary file was deleted
    When ResumeTransformation is called
    Then system detects missing output
    And system retries stage 2
    And new temporary file is created
    And pipeline continues from stage 3

  @REQ-PIPELINE-033 @happy
  Scenario: Resume continues from next stage if output present
    Given a pipeline with stage 1 completed
    And stage 1 output exists
    When ResumeTransformation is called
    Then system skips stage 1
    And system executes stage 2
    And stage 2 reads from stage 1 output

  @REQ-PIPELINE-033 @happy
  Scenario: Resume updates ProcessingState at each transition
    Given a pipeline resuming from stage 2
    When stage 3 completes
    Then ProcessingState is updated to reflect stage 3 output
    And ProcessingState accurately represents current data state
    And ProcessingState matches actual transformations applied

  @REQ-PIPELINE-032 @error
  Scenario: Resume returns error if resume fails
    Given a pipeline with corrupted stage output
    When ResumeTransformation is called
    Then resume attempt fails
    And structured error is returned
    And error indicates resume failure reason
