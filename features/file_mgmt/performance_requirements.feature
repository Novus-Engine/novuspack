@domain:file_mgmt @m2 @REQ-FILEMGMT-075 @spec(api_file_mgmt_addition.md#31-processing-order-requirements)
Feature: Performance Requirements

  @REQ-FILEMGMT-075 @happy
  Scenario: Performance requirements optimize deduplication efficiency
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFile is called
    Then processed size is used as early elimination filter
    And CRC32 is used as early elimination filter
    And SHA-256 is only computed when size and CRC32 match
    And deduplication efficiency is optimized

  @REQ-FILEMGMT-075 @happy
  Scenario: Performance requirements optimize SHA-256 computation
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFile performs deduplication
    Then SHA-256 hash is only computed when necessary
    And expensive SHA-256 computation is avoided when possible
    And performance is optimized through early elimination

  @REQ-FILEMGMT-075 @happy
  Scenario: Performance requirements manage memory efficiently
    Given an open NovusPack package
    And a large file to be added
    And a valid context
    When AddFile is called with large file
    Then memory management is efficient
    And streaming is used when needed
    And large files are handled without excessive memory usage
    And memory footprint is controlled
