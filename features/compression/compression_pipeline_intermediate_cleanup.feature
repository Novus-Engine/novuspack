@domain:compression @m2 @REQ-COMPR-168 @spec(api_file_mgmt_transform_pipelines.md#2116-intermediate-stage-cleanup)
Feature: Compression pipeline intermediate cleanup

  @REQ-COMPR-168 @happy
  Scenario: Intermediate compressed files are cleaned up on success or failure
    Given a transformation pipeline that creates intermediate compressed artifacts
    When the pipeline completes successfully
    Then intermediate stage files are cleaned up automatically
    Given a pipeline that fails partway through execution
    When failure occurs
    Then intermediate stage files are still cleaned up automatically
    And cleanup behavior follows the documented intermediate stage cleanup rules

