@domain:file_mgmt @m2 @REQ-FILEMGMT-271 @spec(api_file_mgmt_queries.md#112-package-listfiles-method)
Feature: Package ListFiles method returns all file entries in the package

  @REQ-FILEMGMT-271 @happy
  Scenario: ListFiles returns all file entries
    Given an open NovusPack package with file entries
    When ListFiles is called
    Then all file entries in the package are returned
    And the behavior matches the ListFiles method specification
    And results are sorted as specified
    And error is returned when package is not open
