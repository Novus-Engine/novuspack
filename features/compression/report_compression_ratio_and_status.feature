@domain:compression @m2 @REQ-COMPR-003 @spec(api_package_compression.md#7-compression-information-and-status)
Feature: Report compression ratio and status

  @happy
  Scenario: Compression status is queryable
    Given a package with some compressed content
    When I query compression status
    Then I should see the current algorithm and compression ratio

  @happy
  Scenario: GetPackageCompressionInfo returns complete information
    Given a package
    When GetPackageCompressionInfo is called
    Then PackageCompressionInfo is returned
    And Type field indicates compression type
    And IsCompressed field indicates compression status
    And OriginalSize and CompressedSize are included
    And Ratio is calculated correctly

  @happy
  Scenario: IsPackageCompressed indicates compression status
    Given a package
    When IsPackageCompressed is called
    Then true is returned if package is compressed
    And false is returned if package is not compressed
    And status matches header flags

  @happy
  Scenario: GetPackageCompressionType returns compression type
    Given a package
    When GetPackageCompressionType is called
    Then compression type is returned
    And type matches header flags bits 15-8

  @happy
  Scenario: Compression ratio is calculated correctly
    Given a compressed package
    When compression information is retrieved
    Then compression ratio equals CompressedSize / OriginalSize
    And ratio is between 0.0 and 1.0
    And ratio indicates compression effectiveness
