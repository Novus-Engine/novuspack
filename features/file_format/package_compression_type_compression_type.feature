@domain:file_format @m2 @REQ-FILEFMT-038 @spec(package_file_format.md#254-package-compression-type)
Feature: Package Compression Type

  @REQ-FILEFMT-038 @happy
  Scenario: Package compression type defines compression flag encoding
    Given a NovusPack package header
    When package compression type is examined
    Then compression flag encoding is defined
    And compression type is encoded in flags bits 15-8
    And compression type determines package compression behavior

  @REQ-FILEFMT-038 @happy
  Scenario: Compression type 0 indicates no package compression
    Given a NovusPack package header
    And flags bits 15-8 equal 0
    When package compression type is examined
    Then compression type is 0 (no compression)
    And package has no package-level compression
    And package content is uncompressed

  @REQ-FILEFMT-038 @happy
  Scenario: Compression type 1 indicates Zstd compression
    Given a NovusPack package header
    And flags bits 15-8 equal 1
    When package compression type is examined
    Then compression type is 1 (Zstd)
    And package uses Zstd compression algorithm
    And package content is compressed with Zstd

  @REQ-FILEFMT-038 @happy
  Scenario: Compression type 2 indicates LZ4 compression
    Given a NovusPack package header
    And flags bits 15-8 equal 2
    When package compression type is examined
    Then compression type is 2 (LZ4)
    And package uses LZ4 compression algorithm
    And package content is compressed with LZ4

  @REQ-FILEFMT-038 @happy
  Scenario: Compression type 3 indicates LZMA compression
    Given a NovusPack package header
    And flags bits 15-8 equal 3
    When package compression type is examined
    Then compression type is 3 (LZMA)
    And package uses LZMA compression algorithm
    And package content is compressed with LZMA

  @REQ-FILEFMT-038 @happy
  Scenario: Compression types 4-255 are reserved for future algorithms
    Given a NovusPack package header
    And flags bits 15-8 are in range 4-255
    When package compression type is examined
    Then compression types 4-255 are reserved
    And reserved values enable future algorithm support
    And extensibility is supported

  @REQ-FILEFMT-038 @error
  Scenario: Invalid compression type values are rejected
    Given a NovusPack package header
    And compression type exceeds valid range
    When package compression type is validated
    Then validation may flag invalid compression type
    And invalid compression type violations are detected
