@domain:file_mgmt @m2 @REQ-FILEMGMT-191 @spec(api_file_mgmt_compression.md#33-getfilecompressioninfo-parameters) @spec(api_file_mgmt_file_entry.md#714-getfilecompressioninfo-parameters)
Feature: GetFileCompressionInfo Parameter Specification

  @REQ-FILEMGMT-191 @happy
  Scenario: GetFileCompressionInfo parameters include context and path
    Given an open NovusPack package
    And a valid context
    And a file path in the package
    When GetFileCompressionInfo is called
    Then context parameter supports cancellation and timeout handling
    And path parameter specifies virtual file path to inspect
    And parameters are validated

  @REQ-FILEMGMT-191 @happy
  Scenario: GetFileCompressionInfo returns compression information
    Given an open NovusPack package
    And a valid context
    And a file path in the package
    When GetFileCompressionInfo is called
    Then FileCompressionInfo structure is returned
    And compression information includes compression status
    And compression information includes compression type
    And compression information includes size and ratio

  @REQ-FILEMGMT-191 @error
  Scenario: GetFileCompressionInfo handles package not open errors
    Given a closed NovusPack package
    And a valid context
    And a file path
    When GetFileCompressionInfo is called
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-191 @error
  Scenario: GetFileCompressionInfo handles file not found errors
    Given an open NovusPack package
    And a valid context
    And a non-existent file path
    When GetFileCompressionInfo is called with non-existent path
    Then a structured error is returned
    And error indicates file not found
    And error follows structured error format
