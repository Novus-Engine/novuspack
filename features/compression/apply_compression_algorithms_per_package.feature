@domain:compression @m2 @REQ-COMPR-001 @spec(api_package_compression.md#12-compression-types)
Feature: Apply compression algorithms per package

  @happy
  Scenario: Algorithms negotiated and applied per spec
    Given a package with compressible content
    When I apply the default compression settings
    Then the package should report a supported algorithm in use

  @happy
  Scenario: CompressionNone applies no compression
    Given a package
    When CompressionNone is applied
    Then package remains uncompressed
    And IsCompressed returns false
    And compression ratio is 1.0

  @happy
  Scenario: CompressionZstd applies Zstd compression
    Given a package with compressible content
    When CompressionZstd is applied
    Then package is compressed with Zstd algorithm
    And IsCompressed returns true
    And compression ratio is less than 1.0

  @happy
  Scenario: CompressionLZ4 applies LZ4 compression
    Given a package with compressible content
    When CompressionLZ4 is applied
    Then package is compressed with LZ4 algorithm
    And IsCompressed returns true
    And compression is faster than Zstd

  @happy
  Scenario: CompressionLZMA applies LZMA compression
    Given a package with compressible content
    When CompressionLZMA is applied
    Then package is compressed with LZMA algorithm
    And IsCompressed returns true
    And compression ratio is typically highest

  @happy
  Scenario: Compression scope excludes header, comment, and signatures
    Given a package with header, comment, and signatures
    When compression is applied
    Then header remains uncompressed
    And comment remains uncompressed
    And signatures remain uncompressed
    And file entries, data, and index are compressed

  @error
  Scenario: Compression fails on signed packages
    Given a signed package with SignatureOffset > 0
    When compression is attempted
    Then structured validation error is returned
    And error indicates signed package cannot be compressed

  @error
  Scenario: Compression operations respect context cancellation
    Given a package
    And a cancelled context
    When compression operation is called
    Then structured context error is returned
