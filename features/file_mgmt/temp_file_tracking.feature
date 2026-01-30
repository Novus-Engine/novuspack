@domain:file_mgmt @REQ-PIPELINE-021 @REQ-PIPELINE-020 @spec(api_file_mgmt_transform_pipelines.md#16-intermediate-stage-cleanup)
Feature: Temporary File Tracking

  @REQ-PIPELINE-020 @happy
  Scenario: Temporary files created in system temp directory
    Given a transformation stage requiring temp file
    When temporary file is created
    Then file is created in system temp directory
    And file has unique name for isolation
    And file permissions follow OS temp directory defaults

  @REQ-PIPELINE-021 @happy
  Scenario: All temporary files tracked in TransformPipeline
    Given a pipeline with 4 stages
    And each stage creates temporary output file
    When pipeline executes
    Then all 4 temporary files are tracked in pipeline
    And each stage stores OutputSource with temp file path
    And TransformPipeline maintains complete cleanup list

  @REQ-PIPELINE-021 @happy
  Scenario: Tracking prevents orphaned temporary files
    Given a multi-stage transformation
    When any stage creates temporary file
    Then file path is immediately stored in TransformStage
    And file can be cleaned up even if operation fails
    And no temporary files are orphaned or lost

  @REQ-PIPELINE-021 @happy
  Scenario: Temporary file tracking survives interruptions
    Given a pipeline interrupted mid-execution
    And some stages have created temporary files
    When system recovers
    Then temporary file paths are available in pipeline
    And cleanup can proceed based on stored paths
    And all tracked files can be removed
