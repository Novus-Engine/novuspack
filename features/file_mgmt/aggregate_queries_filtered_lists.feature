@domain:file_mgmt @m2 @REQ-FILEMGMT-261 @spec(api_file_mgmt_queries.md#4-aggregate-queries-and-filtered-lists)
Feature: Aggregate queries and filtered lists return non-FileEntry values or filtered lists

  @REQ-FILEMGMT-261 @happy
  Scenario: Aggregate queries and filtered lists return filtered results
    Given an open NovusPack package with file entries
    When aggregate queries or filtered lists are used
    Then non-FileEntry values or filtered lists for common criteria are returned
    And the behavior matches the aggregate-queries specification
    And ListCompressedFiles and ListEncryptedFiles are available
    And results match the filter criteria
