@domain:compression @m2 @REQ-COMPR-041 @spec(api_package_compression.md#1133-option-3-write-with-compression)
Feature: Write with Compression

  @REQ-COMPR-041 @happy
  Scenario: Write applies compression during write operation
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type
    And an overwrite flag
    When Write is called with compression type
    Then compression is applied during write process
    And package is written to disk with compression
    And operation handles compression and writing together

  @REQ-COMPR-041 @happy
  Scenario: Write with compression workflow writes package with compression applied
    Given an open NovusPack package
    And a valid context
    And compression should be applied during write
    When write with compression workflow is followed
    Then Write applies compression during write
    And single operation handles both compression and writing
    And workflow integrates compression with write operation
