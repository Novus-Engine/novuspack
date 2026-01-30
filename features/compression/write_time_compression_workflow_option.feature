@domain:compression @m2 @REQ-COMPR-040 @spec(api_package_compression.md#11331-process-option-3)
Feature: Write-time compression workflow option

  @REQ-COMPR-040 @happy
  Scenario: Write with compression applies compression during write operations
    Given an uncompressed in-memory package
    And a write workflow that performs compression at write time
    When the package is written using the write-with-compression option
    Then compression is applied during the write process
    And the output file is a compressed package representation
    And compression behavior follows the documented process option
    And large payloads are handled without requiring pre-compression of the entire package file
    And the resulting package remains readable by compatible readers

