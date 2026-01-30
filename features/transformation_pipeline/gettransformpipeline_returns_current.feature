@domain:transformation_pipeline @m2 @REQ-PIPELINE-007 @spec(api_file_mgmt_file_entry.md#146-multi-stage-transformation-pipeline)
Feature: GetTransformPipeline returns current transformation pipeline (nil if no pipeline active)

  @REQ-PIPELINE-007 @happy
  Scenario: GetTransformPipeline returns current pipeline or nil
    Given a FileEntry that may or may not have an active transformation pipeline
    When GetTransformPipeline is invoked
    Then the current pipeline is returned when a pipeline is active
    And nil is returned when no pipeline is active
    And the behavior matches the multi-stage transformation pipeline specification
    And callers can safely check for nil before using the pipeline
