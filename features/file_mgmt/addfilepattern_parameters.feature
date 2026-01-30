@domain:file_mgmt @m2 @REQ-FILEMGMT-056 @spec(api_file_mgmt_addition.md#242-addfilepattern-parameters)
Feature: AddFilePattern parameters include context, pattern, and options

  @REQ-FILEMGMT-056 @happy
  Scenario: AddFilePattern accepts context, pattern, and options
    Given a package and a pattern for file addition
    When AddFilePattern is called
    Then parameters include context, pattern, and options
    And the parameter contract matches the specification
    And the behavior matches the AddFilePattern parameters specification
