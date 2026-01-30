@domain:file_mgmt @m2 @REQ-FILEMGMT-431 @spec(api_file_mgmt_file_entry.md#105-processdata-behavior)
Feature: ProcessData behavior defines compression and encryption application

  @REQ-FILEMGMT-431 @happy
  Scenario: ProcessData behavior defines compression and encryption
    Given a FileEntry with loaded file content
    When ProcessData is performed
    Then ProcessData behavior defines compression and encryption application as specified
    And the behavior matches the ProcessData behavior specification
    And compression or encryption is applied per configuration
    And processing follows specification
