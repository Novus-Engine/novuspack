@domain:file_format @m1 @REQ-FILEFMT-018 @spec(package_file_format.md#312-uncompressed-content)
Feature: Header, comment, and signatures remain uncompressed

  @REQ-FILEFMT-018 @happy
  Scenario: Package header remains uncompressed when compression is enabled
    Given a package with compression type Zstd (value 1)
    When the package layout is examined
    Then the package header region is uncompressed
    And header can be accessed directly without decompression
    And header provides direct access for validation

  @REQ-FILEFMT-018 @happy
  Scenario: Package comment remains uncompressed when compression is enabled
    Given a package with compression type LZ4 (value 2)
    And package has a comment
    When the package layout is examined
    Then the package comment region is uncompressed
    And comment can be read directly without decompression
    And comment provides easy reading access

  @REQ-FILEFMT-018 @happy
  Scenario: Digital signatures remain uncompressed when compression is enabled
    Given a signed package with compression type LZMA (value 3)
    When the package layout is examined
    Then the digital signatures region is uncompressed
    And signatures can be validated directly without decompression
    And signatures provide direct access for validation

  @REQ-FILEFMT-018 @happy
  Scenario: Uncompressed sections enable direct access for all compression types
    Given a compressed package
    When header, comment, or signatures are accessed
    Then direct access is enabled without decompression
    And performance is optimized for these elements
    And access patterns are efficient

  @REQ-FILEFMT-018 @happy
  Scenario: Uncompressed sections preserve package structure integrity
    Given a compressed package
    When package layout is examined
    Then header region maintains fixed structure
    And comment region maintains format structure
    And signatures region maintains signature structure
    And package structure integrity is preserved

  @REQ-FILEFMT-018 @happy
  Scenario Outline: Uncompressed sections when package compression is enabled
    Given a package with Flags compression type = <Compression>
    When the package layout is examined
    Then the header region is uncompressed
    And the package comment region is uncompressed
    And the signatures region is uncompressed

    Examples:
      | Compression |
      | 1           |
      | 2           |
      | 3           |
