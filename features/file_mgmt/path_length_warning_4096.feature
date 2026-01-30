@domain:file_mgmt @m2 @REQ-FILEMGMT-320 @spec(api_file_mgmt_addition.md#2672-platform-extraction-limits)
Feature: Path length validation emits warnings for paths over 4,096 bytes

  @REQ-FILEMGMT-320 @happy
  Scenario: Path length warning for paths over 4096 bytes
    Given path input over 4,096 bytes
    When path is used for addition or extraction
    Then warnings are emitted for paths over 4,096 bytes due to platform extraction limits
    And the behavior matches the platform-extraction-limits specification
    And operation may still succeed with warning
    And extraction on some platforms may fail for long paths
