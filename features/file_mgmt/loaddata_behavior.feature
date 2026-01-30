@domain:file_mgmt @m2 @REQ-FILEMGMT-430 @spec(api_file_mgmt_file_entry.md#104-loaddata-behavior)
Feature: LoadData behavior defines file content loading process

  @REQ-FILEMGMT-430 @happy
  Scenario: LoadData behavior defines loading process
    Given a FileEntry with a valid data source
    When LoadData is performed
    Then LoadData behavior defines file content loading process as specified
    And the behavior matches the LoadData behavior specification
    And content is loaded from CurrentSource
    And loading process follows specification
