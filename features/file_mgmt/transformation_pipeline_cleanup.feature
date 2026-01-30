@domain:file_mgmt @REQ-PIPELINE-021 @REQ-PIPELINE-024 @spec(api_file_mgmt_transform_pipelines.md#16-intermediate-stage-cleanup) @spec(api_file_mgmt_file_entry.md#45-multi-stage-transformation-pipeline)
Feature: Transformation Pipeline Cleanup

  @REQ-PIPELINE-021 @happy
  Scenario: All intermediate temporary files tracked in pipeline
    Given a TransformPipeline with 3 stages
    When each stage creates temporary file
    Then all temporary file paths stored in pipeline stages
    And TransformPipeline tracks complete cleanup list
    And no temporary files are lost or forgotten

  @REQ-PIPELINE-022 @happy
  Scenario: Success cleanup removes all intermediate files
    Given a completed transformation pipeline
    And all stages succeeded
    When cleanup is triggered
    Then all intermediate temporary files are removed
    And only final output is retained
    And CurrentSource is restored to OriginalSource

  @REQ-PIPELINE-023 @happy
  Scenario: Failure cleanup removes all created temporary files
    Given a transformation pipeline that failed at stage 2
    And stage 1 completed (temp file created)
    And stage 2 partially executed (partial temp file)
    When cleanup is triggered
    Then stage 1 temporary file is removed
    And stage 2 partial temporary file is removed
    And no temporary files remain
    And cleanup handles missing files gracefully

  @REQ-PIPELINE-024 @happy
  Scenario: CleanupTransformPipeline provides manual cleanup
    Given a FileEntry with TransformPipeline
    And pipeline has created temporary files
    When CleanupTransformPipeline is called
    Then all temporary files in pipeline are removed
    And cleanup returns error only if removal fails
    And missing files are handled gracefully (no error)

  @REQ-PIPELINE-024 @happy
  Scenario: Cleanup handles nil pipeline gracefully
    Given a FileEntry without TransformPipeline
    When CleanupTransformPipeline is called
    Then no error is returned
    And operation completes successfully

  @REQ-PIPELINE-023 @happy
  Scenario: Cleanup removes files even if stages incomplete
    Given a pipeline with 3 stages
    And only stage 1 completed
    When cleanup is triggered
    Then stage 1 temporary file is removed
    And no attempt to remove non-existent stage 2 or 3 files
    And cleanup succeeds without errors
