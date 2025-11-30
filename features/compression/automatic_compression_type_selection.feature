@domain:compression @m2 @REQ-COMPR-150 @REQ-COMPR-151 @REQ-COMPR-152 @REQ-COMPR-153 @spec(api_package_compression.md#1122-automatic-compression-type-selection)
Feature: Automatic Compression Type Selection

  @REQ-COMPR-150 @REQ-COMPR-151 @happy
  Scenario: Automatic selection analyzes package properties when compression type is not specified
    Given a NovusPack package
    And compression type is not specified (compressionType = 0)
    When CompressPackage is called with compressionType 0
    Then package properties are analyzed
    And total package size is calculated
    And file count is determined
    And file type distribution is analyzed
    And average file size is calculated
    And content compressibility is estimated
    And optimal compression type is selected based on analysis

  @REQ-COMPR-150 @REQ-COMPR-152 @happy
  Scenario: Automatic selection selects LZ4 for packages with already-compressed content
    Given a NovusPack package
    And package contains >50% already-compressed formats (JPEG, PNG, GIF, MP3, MP4, OGG, FLAC)
    When CompressPackage is called with compressionType 0
    Then LZ4 compression is automatically selected
    And selection prioritizes speed over compression ratio
    And rationale is minimal benefit from heavy compression on already-compressed content

  @REQ-COMPR-150 @REQ-COMPR-152 @happy
  Scenario: Automatic selection selects LZ4 for small packages
    Given a NovusPack package
    And total package size is less than 10MB
    When CompressPackage is called with compressionType 0
    Then LZ4 compression is automatically selected
    And selection prioritizes speed over compression ratio
    And rationale is compression overhead outweighs benefits for small packages

  @REQ-COMPR-150 @REQ-COMPR-152 @happy
  Scenario: Automatic selection selects LZ4 for packages with many small files
    Given a NovusPack package
    And file count is greater than 100
    And average file size is less than 10KB
    When CompressPackage is called with compressionType 0
    Then LZ4 compression is automatically selected
    And selection prioritizes speed
    And rationale is package structure overhead makes compression ratio less important than speed

  @REQ-COMPR-150 @REQ-COMPR-152 @happy
  Scenario: Automatic selection selects LZMA for large packages with text-heavy content
    Given a NovusPack package
    And total package size is greater than 100MB
    And text-based files (text, scripts, configs) represent >60% of content
    When CompressPackage is called with compressionType 0
    Then LZMA compression is automatically selected
    And selection prioritizes size reduction
    And rationale is text compresses well and large size justifies CPU cost

  @REQ-COMPR-150 @REQ-COMPR-152 @happy
  Scenario: Automatic selection selects Zstandard for large packages with mixed content
    Given a NovusPack package
    And total package size is greater than 100MB
    And text-based files represent 30-60% of content
    When CompressPackage is called with compressionType 0
    Then Zstandard compression is automatically selected
    And selection provides balanced compression
    And rationale is mixed content benefits from balanced approach

  @REQ-COMPR-150 @REQ-COMPR-152 @happy
  Scenario: Automatic selection selects Zstandard for large packages with binary-heavy content
    Given a NovusPack package
    And total package size is greater than 100MB
    And binary files represent >60% of content
    When CompressPackage is called with compressionType 0
    Then Zstandard compression is automatically selected
    And selection provides good compression for binary with reasonable speed
    And rationale is binary doesn't compress as well, balanced approach optimal

  @REQ-COMPR-150 @REQ-COMPR-152 @happy
  Scenario: Automatic selection selects Zstandard for medium packages
    Given a NovusPack package
    And total package size is between 10MB and 100MB
    When CompressPackage is called with compressionType 0
    Then Zstandard compression is automatically selected
    And selection provides balanced performance
    And rationale is default balanced approach for moderate sizes

  @REQ-COMPR-150 @REQ-COMPR-152 @happy
  Scenario: Automatic selection defaults to Zstandard as fallback
    Given a NovusPack package
    And package properties do not match any specific selection rules
    When CompressPackage is called with compressionType 0
    Then Zstandard compression is automatically selected as default
    And selection provides safe balanced default
    And rationale is Zstandard provides good balance of speed and compression for most scenarios

  @REQ-COMPR-153 @happy
  Scenario: File type classification supports automatic selection
    Given a NovusPack package
    And files of various types are present
    When automatic compression selection analyzes package
    Then files are classified as text-based, binary, already-compressed, or media
    And classification uses SelectCompressionType logic where applicable
    And classification informs compression type selection

  @REQ-COMPR-150 @happy
  Scenario: Automatic selection works with CompressPackageFile
    Given a NovusPack package
    And compression type is not specified
    When CompressPackageFile is called with compressionType 0
    Then automatic selection algorithm is applied
    And selected compression type is used for compression
    And compressed package is written to file

  @REQ-COMPR-150 @happy
  Scenario: Automatic selection works with CompressPackageStream
    Given a NovusPack package
    And compression type is not specified
    When CompressPackageStream is called with compressionType 0
    Then automatic selection algorithm is applied
    And selected compression type is used for streaming compression

  @REQ-COMPR-150 @happy
  Scenario: Automatic selection works with Write method
    Given a NovusPack package
    And compression type is not specified
    When Write is called with compressionType 0
    Then automatic selection algorithm is applied
    And selected compression type is used for write operation

  @REQ-COMPR-150 @happy
  Scenario: Automatic selection can be overridden by explicit compression type
    Given a NovusPack package
    And automatic selection would choose LZ4
    When CompressPackage is called with compressionType 3 (LZMA)
    Then LZMA compression is used
    And automatic selection is bypassed
    And explicit compression type takes precedence

  @REQ-COMPR-150 @happy
  Scenario: Automatic selection is consistent for same package properties
    Given a NovusPack package with specific properties
    When CompressPackage is called with compressionType 0
    Then compression type A is selected
    When same package properties are analyzed again
    And CompressPackage is called with compressionType 0
    Then same compression type A is selected
    And selection is consistent and deterministic
