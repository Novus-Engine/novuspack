@domain:file_format @m2 @REQ-FILEFMT-041 @spec(package_file_format.md#31-compression-scope)
Feature: Compression Scope

  @REQ-FILEFMT-041 @happy
  Scenario: Compression scope defines compression boundaries
    Given a NovusPack package
    And package compression is enabled
    When compression scope is examined
    Then compression boundaries are defined
    And compressed content is within scope
    And uncompressed content is excluded from scope

  @REQ-FILEFMT-041 @happy
  Scenario: Compression scope includes file entries
    Given a NovusPack package
    And package compression is enabled
    When compression scope is examined
    Then file entries are within compression scope
    And file entries directory structure is compressed
    And file entries are part of compressed content

  @REQ-FILEFMT-041 @happy
  Scenario: Compression scope includes file data
    Given a NovusPack package
    And package compression is enabled
    When compression scope is examined
    Then file data sections are within compression scope
    And actual file contents are compressed
    And file data is part of compressed content

  @REQ-FILEFMT-041 @happy
  Scenario: Compression scope includes package index
    Given a NovusPack package
    And package compression is enabled
    When compression scope is examined
    Then file index is within compression scope
    And package index is compressed
    And file index is part of compressed content

  @REQ-FILEFMT-041 @happy
  Scenario: Compression scope excludes package header
    Given a NovusPack package
    And package compression is enabled
    When compression scope is examined
    Then package header is excluded from compression scope
    And header remains uncompressed for compression type detection
    And header is directly accessible

  @REQ-FILEFMT-041 @happy
  Scenario: Compression scope excludes package comment
    Given a NovusPack package
    And package compression is enabled
    And package has a comment
    When compression scope is examined
    Then package comment is excluded from compression scope
    And comment remains uncompressed for easy reading
    And comment can be read without decompression

  @REQ-FILEFMT-041 @happy
  Scenario: Compression scope excludes digital signatures
    Given a signed NovusPack package
    And package compression is enabled
    When compression scope is examined
    Then digital signatures are excluded from compression scope
    And signatures remain uncompressed for validation
    And signatures can be validated without decompression
