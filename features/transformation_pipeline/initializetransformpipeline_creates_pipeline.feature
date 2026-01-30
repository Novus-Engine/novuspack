@domain:transformation_pipeline @m2 @REQ-PIPELINE-006 @spec(api_file_mgmt_file_entry.md#146-multi-stage-transformation-pipeline)
Feature: InitializeTransformPipeline creates new transformation pipeline with ordered stages

  @REQ-PIPELINE-006 @happy
  Scenario: InitializeTransformPipeline creates pipeline with stages
    Given a FileEntry or package context requiring multi-stage transformation
    When InitializeTransformPipeline is invoked with ordered stages
    Then a new transformation pipeline is created
    And stages are in the specified order
    And the behavior matches the multi-stage transformation pipeline specification
    And the pipeline is ready for stage execution
