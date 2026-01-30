@domain:file_mgmt @m2 @REQ-FILEMGMT-066 @spec(api_file_mgmt_addition.md#281-addfileoptions-purpose)
Feature: AddFileOptions purpose configures file processing behavior

  @REQ-FILEMGMT-066 @happy
  Scenario: AddFileOptions purpose configures file processing behavior
    Given file addition operations that accept options
    When AddFileOptions is used to configure behavior
    Then file processing behavior is configured as specified
    And the purpose matches the AddFileOptions specification
    And options control compression, encryption, and path behavior
