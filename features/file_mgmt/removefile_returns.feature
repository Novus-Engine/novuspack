@domain:file_mgmt @m2 @REQ-FILEMGMT-138 @spec(api_file_mgmt_removal.md#24-removefile-returns)
Feature: RemoveFile returns error on failure

  @REQ-FILEMGMT-138 @happy
  Scenario: RemoveFile returns an error on failure
    Given a package with file entries
    When RemoveFile is called and a failure occurs
    Then an error is returned
    And the error is structured
    And the behavior matches the RemoveFile returns specification
