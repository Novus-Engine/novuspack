@domain:file_mgmt @m2 @REQ-FILEMGMT-269 @spec(api_file_mgmt_queries.md#11-package-methods)
Feature: Package methods define FileExists and ListFiles methods

  @REQ-FILEMGMT-269 @happy
  Scenario: Package methods define FileExists and ListFiles
    Given an open NovusPack package
    When Package methods are used for file queries
    Then FileExists and ListFiles methods are defined as specified
    And the behavior matches the package-methods specification
    And FileExists checks path existence
    And ListFiles returns all file entries
