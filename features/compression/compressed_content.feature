@domain:compression @m2 @REQ-COMPR-018 @spec(api_package_compression.md#111-compressed-content)
Feature: Compressed content

  @REQ-COMPR-018 @happy
  Scenario: Compressed content includes file entries
    Given a NovusPack package compression operation
    When package is compressed
    Then file entries are included in compressed content
    And directory structure is compressed
    And file entry metadata is compressed

  @REQ-COMPR-018 @happy
  Scenario: Compressed content includes file data
    Given a NovusPack package compression operation
    When package is compressed
    Then file data is included in compressed content
    And actual file contents are compressed
    And file data is compressed as part of package content

  @REQ-COMPR-018 @happy
  Scenario: Compressed content includes package index
    Given a NovusPack package compression operation
    When package is compressed
    Then package index is included in compressed content
    And file index is compressed
    And index references are compressed

  @REQ-COMPR-018 @happy
  Scenario: Header, comment, and signatures remain uncompressed
    Given a NovusPack package compression operation
    When package is compressed
    Then package header remains uncompressed
    And package comment remains uncompressed
    And digital signatures remain uncompressed
    And these components are directly accessible without decompression
