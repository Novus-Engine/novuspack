@domain:compression @m2 @REQ-COMPR-008 @spec(api_package_compression.md#72-compression-status-methods)
Feature: Package compression status methods

  @happy
  Scenario: GetPackageCompressionInfo returns compression information
    Given an open package
    When GetPackageCompressionInfo is called
    Then PackageCompressionInfo is returned
    And compression type is included
    And compression status is included
    And compression ratio is included if compressed

  @REQ-COMPR-009 @happy
  Scenario: IsPackageCompressed checks compression status
    Given an open compressed package
    When IsPackageCompressed is called
    Then true is returned

  @happy
  Scenario: IsPackageCompressed returns false for uncompressed package
    Given an open uncompressed package
    When IsPackageCompressed is called
    Then false is returned

  @REQ-COMPR-010 @happy
  Scenario: GetPackageCompressionType returns compression type
    Given an open compressed package
    When GetPackageCompressionType is called
    Then compression type is returned
    And type matches package compression

  @happy
  Scenario: GetPackageCompressionType returns none for uncompressed
    Given an open uncompressed package
    When GetPackageCompressionType is called
    Then no compression type is returned

  @REQ-COMPR-011 @happy
  Scenario: SetPackageCompressionType sets compression type
    Given an open writable uncompressed package
    When SetPackageCompressionType is called with compression type
    Then compression type is set
    And type is accessible via GetPackageCompressionType

  @REQ-COMPR-012 @happy
  Scenario: CanCompressPackage checks if compression is possible
    Given an open package
    When CanCompressPackage is called
    Then true is returned if compression is possible
    And false is returned if compression is not possible
    And status reflects package state

  @error
  Scenario: SetPackageCompressionType fails for read-only package
    Given a read-only package
    When SetPackageCompressionType is called
    Then structured validation error is returned

  @error
  Scenario: SetPackageCompressionType fails with invalid compression type
    Given an open writable package
    When SetPackageCompressionType is called with invalid type
    Then structured validation error is returned

  @REQ-COMPR-013 @REQ-COMPR-014 @error
  Scenario: Compression status methods validate compression type parameter
    Given an open writable package
    When SetPackageCompressionType is called with invalid compression type
    Then structured validation error is returned
    And error indicates unsupported compression type

  @REQ-COMPR-013 @REQ-COMPR-016 @error
  Scenario: Compression status methods respect context cancellation
    Given an open package
    And a cancelled context
    When compression status method is called
    Then structured context error is returned
    And error type is context cancellation
