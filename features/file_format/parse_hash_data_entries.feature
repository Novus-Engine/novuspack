@domain:file_format @m1 @REQ-FILEFMT-015 @spec(package_file_format.md#4143-hash-data)
Feature: Parse Hash Data Entries

  @REQ-FILEFMT-015 @happy
  Scenario: Hash data structure contains count and entries array
    Given a file entry with HashCount > 0
    When hash data is parsed
    Then HashCount (1 byte) indicates number of hash entries
    And hash entries array follows
    And entries are located at HashDataOffset from start of variable-length data

  @REQ-FILEFMT-015 @happy
  Scenario: Each hash entry has HashType, HashPurpose, HashLength, and HashData
    Given a file entry with hash entries
    When hash entries are parsed
    Then each entry has HashType (1 byte)
    And each entry has HashPurpose (1 byte)
    And each entry has HashLength (2 bytes)
    And each entry has variable-length HashData
    And HashLength matches the actual hash data length

  @REQ-FILEFMT-015 @happy
  Scenario: Supported hash types parse correctly
    Given a file entry with hash entries
    And entries have supported HashType values
    When hash data is parsed
    Then SHA-256 (0x00) entries parse correctly
    And SHA-512 (0x01) entries parse correctly
    And BLAKE3 (0x02) entries parse correctly
    And XXH3 (0x03) entries parse correctly

  @REQ-FILEFMT-015 @happy
  Scenario: Supported hash purposes parse correctly
    Given a file entry with hash entries
    And entries have supported HashPurpose values
    When hash data is parsed
    Then content verification (0x00) entries parse correctly
    And deduplication (0x01) entries parse correctly
    And integrity (0x02) entries parse correctly

  @REQ-FILEFMT-015 @happy
  Scenario: Reserved hash types and purposes are accepted
    Given a file entry with hash entries
    And entries have reserved HashType (0x04-0xFF) or HashPurpose (0x03-0xFF) values
    When hash data is parsed
    Then reserved types and purposes are accepted
    And entries are parsed with unknown type handling
    And hash data payload is preserved

  @REQ-FILEFMT-015 @happy
  Scenario Outline: Hash entries parse and validate types and lengths
    Given a file entry with HashCount=<Count>
    And hash entries with types <Types> purposes <Purposes> lengths <Lengths>
    When hash data is parsed
    Then there are <Count> hash entries
    And each entry has a supported HashType and HashPurpose
    And each entry's data length matches its HashLength

    Examples:
      | Count | Types        | Purposes   | Lengths    |
      | 2     | 0x00,0x02    | 0x00,0x01  | 32,32      |
      | 1     | 0x03         | 0x03       | 8          |

  @REQ-FILEFMT-015 @error
  Scenario: Hash length mismatch is rejected
    Given a file entry with hash data where HashLength does not match actual data
    When hash data is parsed
    Then a structured invalid hash data error is returned
    And error indicates length mismatch
    And error follows structured error format

  @REQ-FILEFMT-015 @error
  Scenario: Invalid HashDataOffset is rejected
    Given a file entry with invalid HashDataOffset
    When hash data is parsed
    Then a structured invalid hash data error is returned
    And error indicates invalid offset
    And error follows structured error format

  @happy
  Scenario: WriteTo serializes hash entry to binary format
    Given a HashEntry with values
    When hash entry WriteTo is called with writer
    Then hash entry is written to writer
    And written data matches hash entry content

  @happy
  Scenario: ReadFrom deserializes hash entry from binary format
    Given a reader with valid hash entry data
    When hash entry ReadFrom is called with reader
    Then hash entry is read from reader
    And hash entry fields match reader data
    And hash entry is valid

  @happy
  Scenario: Hash entry round-trip serialization preserves all fields
    Given a HashEntry with all fields set
    When hash entry WriteTo is called with writer
    And ReadFrom is called with written data
    Then all hash entry fields are preserved
    And hash entry is valid

  @error
  Scenario: ReadFrom fails with incomplete data
    Given a reader with incomplete hash entry data
    When hash entry ReadFrom is called with reader
    Then structured IO error is returned
    And error indicates read failure
