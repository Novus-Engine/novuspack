@domain:file_mgmt @REQ-PIPELINE-029 @REQ-PIPELINE-030 @REQ-PIPELINE-031 @spec(api_file_mgmt_file_entry.md#15-processingstate-type) @spec(api_file_mgmt_addition.md#211-multi-stage-transformation-pipelines) @spec(api_file_mgmt_transform_pipelines.md#112-pipeline-validation) @spec(api_file_mgmt_transform_pipelines.md#19-configuration-options)
Feature: ProcessingState Validation

  @REQ-PIPELINE-030 @REQ-FILEMGMT-350 @happy
  Scenario: ValidateSources checks pipeline consistency
    Given a FileEntry with TransformPipeline
    When ValidateSources is called
    Then CurrentSource is validated
    And OriginalSource is validated
    And pipeline structure is validated
    And no error returned if all valid

  @REQ-PIPELINE-031 @happy
  Scenario: Validation checks CurrentSource set for active processing
    Given a FileEntry with active processing
    And CurrentSource is set
    When ValidateSources is called
    Then validation passes
    Given CurrentSource is nil but should be set
    When ValidateSources is called
    Then validation error is returned

  @REQ-PIPELINE-031 @happy
  Scenario: Validation checks pipeline has at least one stage
    Given a TransformPipeline with empty Stages array
    When ValidateSources is called
    Then validation error is returned
    And error indicates pipeline must have at least one stage

  @REQ-PIPELINE-031 @happy
  Scenario: Validation checks CurrentSource matches final stage output
    Given a completed TransformPipeline
    And final stage OutputSource is tempfile.dat
    And CurrentSource points to tempfile.dat
    When ValidateSources is called
    Then validation passes
    And CurrentSource matches final pipeline output

  @REQ-PIPELINE-031 @error
  Scenario: Validation fails when CurrentSource does not match final output
    Given a completed TransformPipeline
    And final stage OutputSource is tempfile.dat
    But CurrentSource points to different file
    When ValidateSources is called
    Then validation error is returned
    And error indicates CurrentSource mismatch

  @REQ-PIPELINE-031 @happy
  Scenario: Validation checks OriginalSource type is valid
    Given a FileEntry with OriginalSource
    When ValidateSources is called
    Then OriginalSource must be IsPackage or IsExternal
    And error returned if neither flag is true
    And validation ensures OriginalSource type is valid

  @REQ-PIPELINE-030 @error
  Scenario: ValidateSources returns structured error on failure
    Given invalid source configuration
    When ValidateSources is called
    Then PackageError is returned
    And error type is ErrTypeValidation
    And error provides clear context about validation failure
