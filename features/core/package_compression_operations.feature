@domain:core @m1 @REQ-CORE-007 @spec(api_core.md#6-package-compression-operations)
Feature: Package compression operations

  @happy
  Scenario: Package compression types are defined
    Given package compression functionality
    When compression types are examined
    Then compression type 0 represents no compression
    And compression type 1 represents Zstd compression
    And compression type 2 represents LZ4 compression
    And compression type 3 represents LZMA compression

  @happy
  Scenario: Package compression scope excludes header, comment, and signatures
    Given a NovusPack package
    When package compression is applied
    Then package header remains uncompressed
    And package comment remains uncompressed
    And digital signatures remain uncompressed
    And file entries, file data, and index are compressed

  @happy
  Scenario: Package compression is applied after per-file compression
    Given a NovusPack package with per-file compressed files
    When package compression is applied
    Then per-file compression is preserved
    And package compression is applied to compressed content
    And decompression order is: package first, then per-file

  @happy
  Scenario: Compressed packages can be signed
    Given a compressed NovusPack package
    When a signature is added
    Then signature operation succeeds
    And signatures validate compressed content
    And header, comment, and signatures remain accessible

  @error
  Scenario: Signed packages cannot be compressed
    Given a signed NovusPack package
    When package compression is attempted
    Then a structured error is returned
    And error type indicates CompressSignedPackageError
    And compression fails
    And package remains unchanged

  @happy
  Scenario: Package decompression occurs before per-file decompression
    Given a compressed NovusPack package with compressed files
    When package is accessed
    Then package is decompressed first
    Then per-file decompression occurs
    And decompression order is correct

  @happy
  Scenario: Package compression improves efficiency for small packages
    Given a small NovusPack package with many small files
    When package compression is applied
    Then overall package size is reduced
    And compression efficiency is improved

  @error
  Scenario: Package compression fails if package is signed
    Given a signed NovusPack package with SignatureOffset > 0
    When CompressPackage is called
    Then CompressSignedPackageError is returned
    And compression is prevented
