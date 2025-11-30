@domain:compression @m2 @REQ-COMPR-102 @spec(api_package_compression.md#31-compression-information-interface)
Feature: Compression Information Interface

  @REQ-COMPR-102 @happy
  Scenario: CompressionInfo interface provides GetCompressionInfo method
    Given an open NovusPack package
    And a valid context
    When GetCompressionInfo is called
    Then PackageCompressionInfo structure is returned
    And structure contains comprehensive compression details
    And details include type, state, sizes, and ratio

  @REQ-COMPR-102 @happy
  Scenario: CompressionInfo interface provides IsCompressed method
    Given an open NovusPack package
    When IsCompressed is called
    Then boolean value is returned
    And value indicates whether package is compressed
    And value is true for compressed packages

  @REQ-COMPR-102 @happy
  Scenario: CompressionInfo interface provides GetCompressionType method
    Given an open NovusPack package
    When GetCompressionType is called
    Then compression type and isSet flag are returned
    And type is uint8 value 0-3
    And isSet indicates if compression type is set

  @REQ-COMPR-102 @happy
  Scenario: CompressionInfo interface provides GetCompressionRatio method
    Given an open NovusPack package
    When GetCompressionRatio is called
    Then compression ratio and isSet flag are returned
    And ratio is float64 value between 0.0 and 1.0
    And isSet indicates if compression ratio is available

  @REQ-COMPR-102 @happy
  Scenario: CompressionInfo interface provides CanCompress method
    Given an open NovusPack package
    When CanCompress is called
    Then boolean value is returned
    And value indicates if package can be compressed
    And value is false if package is signed

  @REQ-COMPR-102 @happy
  Scenario: CompressionInfo interface provides read-only access
    Given an open NovusPack package
    When CompressionInfo interface methods are used
    Then all methods provide read-only access
    And no methods modify package state
    And information is retrieved without side effects

  @REQ-COMPR-102 @happy
  Scenario: CompressionInfo interface accepts context parameter
    Given an open NovusPack package
    And a valid context
    When GetCompressionInfo is called with context
    Then context is accepted as parameter
    And context supports cancellation and timeout
    And context follows standard Go patterns
