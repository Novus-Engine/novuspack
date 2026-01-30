@domain:file_mgmt @m2 @REQ-FILEMGMT-270 @spec(api_file_mgmt_queries.md#111-package-fileexists-method)
Feature: Package FileExists method checks if a file with the given path exists

  @REQ-FILEMGMT-270 @happy
  Scenario: FileExists checks if file exists at path
    Given an open NovusPack package with file entries
    When FileExists is called with a path
    Then the method checks if a file with the given path exists in the package
    And the behavior matches the FileExists method specification
    And true is returned when file exists
    And false is returned when file does not exist
