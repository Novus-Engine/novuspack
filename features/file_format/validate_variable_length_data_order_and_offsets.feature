@domain:file_format @m1 @REQ-FILEFMT-015 @spec(package_file_format.md#414-variable-length-data-follows-fixed-structure)
Feature: Validate variable-length data order and offsets

  # Covers 4.1.4.1 order and offset rules
  @happy
  Scenario: Variable-length sections appear in required order
    Given a file entry with paths, hashes, and optional data
    When the variable-length data region is parsed
    Then paths start at offset 0
    And hash data begins at HashDataOffset
    And optional data begins at OptionalDataOffset
    And all sections fit within the variable-length region

  @error
  Scenario Outline: Misordered or overlapping variable-length sections are rejected
    Given a file entry with HashDataOffset=<HOff> HashDataLen=<HLen> OptionalDataOffset=<OOff> OptionalDataLen=<OLen>
    When the variable-length data region is validated
    Then validation result is <Result>

    Examples:
      | HOff | HLen | OOff | OLen | Result  |
      | 32   | 16   | 64   | 16   | valid   |
      | 64   | 32   | 48   | 16   | invalid |
      | 0    | 0    | 0    | 0    | valid   |
