@domain:file_mgmt @m2 @REQ-FILEMGMT-344 @spec(api_file_mgmt_transform_pipelines.md#15-parallel-stage-execution)
Feature: Parallel stage execution is supported where safe but sequential by default

  @REQ-FILEMGMT-344 @happy
  Scenario: Parallel stage execution where safe
    Given a multi-stage transformation pipeline
    When stage execution is performed
    Then parallel stage execution is supported where safe (e.g., checksum verification alongside decompression)
    And sequential execution is default for correctness
    And the behavior matches the parallel-stage-execution specification
    And correctness is preserved when parallel is used
