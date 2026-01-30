@domain:basic_ops @m2 @REQ-API_BASIC-207 @spec(api_basic_operations.md#3352-specialfiles-mapping)
Feature: SpecialFiles mapping

  @REQ-API_BASIC-207 @happy
  Scenario: SpecialFiles mapping links special file types to FileEntry records
    Given a package containing special metadata files
    When special files are loaded
    Then SpecialFiles mapping defines special file type to FileEntry mapping
    And mapping supports efficient retrieval of special metadata files
    And mapping remains consistent with the file index and entries
    And mapping supports optional and required special files as specified
    And mapping aligns with special file naming and file type rules

