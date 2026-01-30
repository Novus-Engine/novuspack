@domain:compression @m2 @REQ-COMPR-166 @spec(api_file_mgmt_transform_pipelines.md#211-multi-stage-transformation-pipelines)
Feature: Compression integrates as a transformation pipeline stage

  @REQ-COMPR-166 @happy
  Scenario: Compression and decompression operate as pipeline stages with intermediate files
    Given a multi-stage transformation pipeline for package operations
    When a compress stage is included in the pipeline
    Then the compress stage writes to a temporary file as an intermediate artifact
    When a decompress stage is included in the pipeline
    Then the decompress stage reads from a temporary file as an intermediate artifact
    And stages can be composed with other transformations as documented
    And intermediate artifacts are managed by the pipeline machinery

