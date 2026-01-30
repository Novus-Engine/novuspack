@domain:file_mgmt @m2 @REQ-FILEMGMT-058 @spec(api_file_mgmt_addition.md#244-addfilepattern-behavior)
Feature: AddFilePattern behavior includes pattern scanning and bulk file addition

  @REQ-FILEMGMT-058 @happy
  Scenario: AddFilePattern performs pattern scanning and bulk file addition
    Given a package and a filesystem pattern
    When AddFilePattern is called
    Then pattern scanning is performed as specified
    And bulk file addition follows the behavior specification
    And the behavior matches the AddFilePattern behavior specification
