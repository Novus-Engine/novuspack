@domain:file_mgmt @REQ-FILEMGMT-348 @REQ-PIPELINE-028 @REQ-PIPELINE-029 @spec(api_file_mgmt_addition.md#28-addfileoptions-struct) @spec(api_file_mgmt_transform_pipelines.md#19-configuration-options)
Feature: ValidateProcessingState Configuration

  @REQ-PIPELINE-028 @REQ-FILEMGMT-348 @happy
  Scenario: ValidateProcessingState default is false
    Given AddFileOptions without ValidateProcessingState specified
    When file is added with transformations
    Then ValidateProcessingState defaults to false
    And validation is disabled for performance
    And ProcessingState is not validated against actual data state

  @REQ-PIPELINE-028 @happy
  Scenario: ValidateProcessingState configurable via AddFileOptions
    Given AddFileOptions with ValidateProcessingState set to true
    When file is added with transformations
    Then ProcessingState validation is enabled
    And system validates state matches actual transformations
    And useful for debugging and development

  @REQ-PIPELINE-029 @REQ-FILEMGMT-348 @happy
  Scenario: Validation checks ProcessingState matches actual data state
    Given ValidateProcessingState enabled
    And file data is Compressed
    And ProcessingState is set to Compressed
    When validation runs
    Then validation passes
    And no error is returned

  @REQ-PIPELINE-029 @REQ-FILEMGMT-348 @error
  Scenario: Validation fails on ProcessingState mismatch
    Given ValidateProcessingState enabled
    And file data is Compressed
    But ProcessingState is set to Raw
    When validation runs
    Then validation error is returned
    And error type is ErrTypeValidation
    And error indicates ProcessingState mismatch
    And error shows expected (Compressed) vs actual (Raw)

  @REQ-PIPELINE-028 @happy
  Scenario: Validation recommended for testing scenarios
    Given test environment with strict validation needs
    When ValidateProcessingState is enabled
    Then ProcessingState errors are caught early
    And mismatches are detected immediately
    And debugging is simplified with clear validation errors
