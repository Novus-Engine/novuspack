@domain:file_mgmt @REQ-PIPELINE-015 @REQ-PIPELINE-016 @REQ-FILEMGMT-343 @spec(api_file_mgmt_transform_pipelines.md#14-disk-space-management)
Feature: Disk Space Pre-flight Checks

  @REQ-PIPELINE-015 @REQ-FILEMGMT-343 @happy
  Scenario: Pre-flight check estimates required space
    Given a file requiring multi-stage transformation
    And file size is 10GB
    And transformations include decrypt and decompress
    When pre-flight check is performed
    Then system estimates required space for all stages
    And estimation includes source plus intermediate files
    And estimation accounts for transformation expansion

  @REQ-PIPELINE-016 @REQ-FILEMGMT-343 @happy
  Scenario: Pre-flight check verifies available space
    Given estimated space requirement is 50GB
    When pre-flight check is performed
    Then system checks available disk space in temp directory
    And check uses appropriate filesystem API
    And available space is compared to requirement

  @REQ-PIPELINE-016 @REQ-FILEMGMT-343 @error
  Scenario: Insufficient space returns descriptive error
    Given estimated requirement is 100GB
    And available disk space is 50GB
    When pre-flight check is performed
    Then ErrTypeIO error is returned
    And error indicates insufficient disk space
    And error includes required space (100GB)
    And error includes available space (50GB)
    And no transformation is attempted

  @REQ-PIPELINE-015 @happy
  Scenario: Pre-flight prevents partial transformations
    Given insufficient disk space
    When pre-flight check runs before pipeline starts
    Then check fails before any temporary files created
    And no partial transformations occur
    And system state remains clean
    And user gets clear error before wasting time
