@domain:file_format @m2 @REQ-FILEFMT-043 @spec(package_file_format.md#32-compression-behavior)
Feature: Compression Format Behavior

  @REQ-FILEFMT-043 @happy
  Scenario: Compression behavior defines compression process
    Given a NovusPack package
    And package compression is enabled
    When compression behavior is examined
    Then compression process is defined
    And behavior is determined by compression type in header flags
    And compression type is encoded in bits 15-8

  @REQ-FILEFMT-043 @happy
  Scenario: Compression behavior is determined by header flags
    Given a NovusPack package
    And package compression type is specified in header flags (bits 15-8)
    When compression behavior is examined
    Then compression type 0 indicates no compression
    And compression type 1 indicates Zstd compression
    And compression type 2 indicates LZ4 compression
    And compression type 3 indicates LZMA compression

  @REQ-FILEFMT-043 @happy
  Scenario: Package compression is applied after per-file compression
    Given a NovusPack package
    And files have per-file compression
    When package compression is performed
    Then package compression is applied after per-file compression
    And relationship to per-file compression is maintained
    And compression order is correct

  @REQ-FILEFMT-043 @happy
  Scenario: Package decompression must occur before per-file decompression
    Given a compressed NovusPack package
    When decompression is performed
    Then package decompression occurs before per-file decompression
    And decompression order ensures proper processing
    And file decompression follows package decompression

  @REQ-FILEFMT-043 @happy
  Scenario: Compressed packages can be signed
    Given a compressed NovusPack package
    When signing is attempted
    Then compressed packages can be signed
    And signing succeeds after compression
    And signature validation works with compressed packages

  @REQ-FILEFMT-043 @error
  Scenario: Signed packages cannot be compressed
    Given a signed NovusPack package
    And SignatureOffset > 0
    When compression is attempted
    Then compression fails
    And structured immutability error is returned
    And signed packages cannot be compressed
