@domain:compression @m2 @REQ-COMPR-039 @spec(api_package_compression.md#1132-option-2-compress-and-write-in-one-step)
Feature: Compress and Write in One Step

  @REQ-COMPR-039 @happy
  Scenario: CompressPackageFile compresses and writes package in one step
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type
    And an overwrite flag
    When CompressPackageFile is called
    Then package content is compressed
    And compressed package is written to specified path
    And operation completes in single step

  @REQ-COMPR-039 @happy
  Scenario: Compress and write in one step workflow combines compression and writing
    Given an open NovusPack package
    And a valid context
    And compression and writing are needed
    When compress and write in one step workflow is followed
    Then CompressPackageFile performs both operations
    And no intermediate steps are required
    And workflow is simplified compared to separate operations
