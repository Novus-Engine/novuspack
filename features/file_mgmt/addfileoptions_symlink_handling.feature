@domain:file_mgmt @m2 @REQ-FILEMGMT-310 @spec(api_file_mgmt_addition.md#2827-symlink-options)
Feature: AddFileOptions symlink handling defines default follow behavior

  @REQ-FILEMGMT-310 @happy
  Scenario: AddFileOptions symlink handling defines follow behavior
    Given AddFileOptions with symlink handling configuration
    When a symlink path is added
    Then symlink handling defines default follow behavior as specified
    And the behavior matches the symlink-options specification
    And symlinks are followed or not per options
    And target content or link is added per configuration
