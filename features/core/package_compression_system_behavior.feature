@domain:core @m2 @REQ-CORE-039 @spec(api_core.md#63-package-compression-behavior)
Feature: Package Compression System Behavior

  @REQ-CORE-039 @happy
  Scenario: Compression scope includes file entries, data, and index
    Given a NovusPack package with files
    When compression is performed
    Then file entries (directory structure) are compressed
    And file data (actual file contents) is compressed
    And package index is compressed
    And compression scope covers package content

  @REQ-CORE-039 @happy
  Scenario: Header, comment, and signatures remain uncompressed
    Given a NovusPack package with header, comment, and signatures
    When compression is performed
    Then package header remains uncompressed for compression type detection
    And package comment remains uncompressed for easy reading without decompression
    And digital signatures remain uncompressed for validation

  @REQ-CORE-039 @happy
  Scenario: Package compression is applied after per-file compression
    Given a NovusPack package with per-file compression
    When package compression is performed
    Then package compression is applied after per-file compression
    And relationship to per-file compression is maintained
    And compression order is correct

  @REQ-CORE-039 @happy
  Scenario: Package decompression must occur before per-file decompression
    Given a compressed NovusPack package
    When decompression is performed
    Then package decompression occurs before per-file decompression
    And decompression order ensures proper processing
    And file decompression follows package decompression

  @REQ-CORE-039 @happy
  Scenario: Compressed packages can be signed but signed packages cannot be compressed
    Given compression and signing operations
    When compression is attempted on signed package
    Then signing compatibility is enforced
    And compressed packages can be signed
    And signed packages cannot be compressed

  @REQ-CORE-039 @happy
  Scenario: Package compression is efficient for small packages
    Given a small NovusPack package
    When package compression is considered
    Then package-level compression is more efficient
    And compression provides benefits for small packages
    And use case is appropriate
