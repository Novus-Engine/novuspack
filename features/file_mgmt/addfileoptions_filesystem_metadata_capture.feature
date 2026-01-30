@domain:file_mgmt @m2 @REQ-FILEMGMT-311 @spec(api_file_mgmt_addition.md#2829-filesystem-metadata-capture-options)
Feature: AddFileOptions filesystem metadata capture options define opt-in capture behavior

  @REQ-FILEMGMT-311 @happy
  Scenario: AddFileOptions filesystem metadata capture is opt-in
    Given AddFileOptions with filesystem metadata capture configuration
    When a file is added from the filesystem
    Then filesystem metadata capture options define opt-in capture behavior
    And the behavior matches the filesystem-metadata-capture-options specification
    And capture is opt-in per options
    And permissions and timestamps are captured when enabled
