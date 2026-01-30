@domain:file_mgmt @m2 @REQ-FILEMGMT-395 @spec(api_file_mgmt_extraction.md#25-extractpathoptions-concurrency)
Feature: ExtractOptions provides configurable security limits and concurrency settings

  @REQ-FILEMGMT-395 @happy
  Scenario: ExtractOptions provides security limits and concurrency
    Given extraction operations with ExtractOptions
    When ExtractOptions is used for extraction
    Then ExtractOptions provides configurable security limits and concurrency settings with safe defaults
    And the behavior matches the ExtractPathOptions concurrency specification
    And safe defaults are applied when options are nil
    And concurrency and security limits are configurable
