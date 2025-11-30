@domain:compression @m2 @REQ-COMPR-136 @spec(api_package_compression.md#71-compression-information-structure)
Feature: Compression: Compression Information Structure (Core)

  @REQ-COMPR-136 @happy
  Scenario: Compression information structure provides compression type field
    Given an open NovusPack package
    And a valid context
    When GetPackageCompressionInfo is called
    Then PackageCompressionInfo structure is returned
    And structure contains Type field
    And Type field stores compression type value 0-3

  @REQ-COMPR-136 @happy
  Scenario: Compression information structure provides compression state field
    Given an open NovusPack package
    And a valid context
    When GetPackageCompressionInfo is called
    Then PackageCompressionInfo structure is returned
    And structure contains IsCompressed boolean field
    And IsCompressed indicates compression state

  @REQ-COMPR-136 @happy
  Scenario: Compression information structure provides original size field
    Given an open NovusPack package
    And a valid context
    When GetPackageCompressionInfo is called
    Then PackageCompressionInfo structure is returned
    And structure contains OriginalSize int64 field
    And OriginalSize stores size before compression in bytes

  @REQ-COMPR-136 @happy
  Scenario: Compression information structure provides compressed size field
    Given an open NovusPack package
    And a valid context
    When GetPackageCompressionInfo is called
    Then PackageCompressionInfo structure is returned
    And structure contains CompressedSize int64 field
    And CompressedSize stores size after compression in bytes

  @REQ-COMPR-136 @happy
  Scenario: Compression information structure provides compression ratio field
    Given an open NovusPack package
    And a valid context
    When GetPackageCompressionInfo is called
    Then PackageCompressionInfo structure is returned
    And structure contains Ratio float64 field
    And Ratio stores compression ratio between 0.0 and 1.0
    And Ratio represents compression efficiency

  @REQ-COMPR-136 @happy
  Scenario: Compression information structure provides comprehensive compression details
    Given an open NovusPack package
    And a valid context
    When GetPackageCompressionInfo is called
    Then structure provides all compression information
    And information enables compression analysis
    And information supports compression decisions

  @REQ-COMPR-136 @happy
  Scenario: Compression information structure is accessible via GetPackageCompressionInfo method
    Given an open NovusPack package
    And a valid context
    When GetPackageCompressionInfo is called
    Then method returns PackageCompressionInfo structure
    And structure is complete and accurate
    And structure reflects current compression state
