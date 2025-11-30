@domain:compression @m2 @REQ-COMPR-020 @REQ-COMPR-136 @spec(api_package_compression.md#13-compression-information-structure)
Feature: Compression: Compression Information Structure (PackageCompressionInfo)

  @REQ-COMPR-020 @REQ-COMPR-136 @happy
  Scenario: PackageCompressionInfo structure contains compression type
    Given a NovusPack package
    And compression information is retrieved
    When PackageCompressionInfo structure is examined
    Then structure contains Type field with compression type value
    And Type field stores compression type 0-3
    And Type indicates compression algorithm used

  @REQ-COMPR-020 @REQ-COMPR-136 @happy
  Scenario: PackageCompressionInfo structure indicates if package is compressed
    Given a NovusPack package
    And compression information is retrieved
    When PackageCompressionInfo structure is examined
    Then structure contains IsCompressed boolean field
    And IsCompressed indicates whether package is compressed
    And IsCompressed is true for compressed packages

  @REQ-COMPR-020 @REQ-COMPR-136 @happy
  Scenario: PackageCompressionInfo structure contains original size
    Given a NovusPack package
    And compression information is retrieved
    When PackageCompressionInfo structure is examined
    Then structure contains OriginalSize field
    And OriginalSize stores package size before compression
    And OriginalSize is int64 value in bytes

  @REQ-COMPR-020 @REQ-COMPR-136 @happy
  Scenario: PackageCompressionInfo structure contains compressed size
    Given a NovusPack package
    And compression information is retrieved
    When PackageCompressionInfo structure is examined
    Then structure contains CompressedSize field
    And CompressedSize stores package size after compression
    And CompressedSize is int64 value in bytes

  @REQ-COMPR-020 @REQ-COMPR-136 @happy
  Scenario: PackageCompressionInfo structure contains compression ratio
    Given a NovusPack package
    And compression information is retrieved
    When PackageCompressionInfo structure is examined
    Then structure contains Ratio field
    And Ratio stores compression ratio as float64
    And Ratio is between 0.0 and 1.0
    And Ratio represents compression efficiency

  @REQ-COMPR-020 @REQ-COMPR-136 @happy
  Scenario: Compression information structure provides comprehensive compression details
    Given a NovusPack package
    And compression information is retrieved
    When PackageCompressionInfo structure is examined
    Then structure provides all compression details
    And details include type, state, sizes, and ratio
    And details enable compression analysis
