@domain:file_format @m1 @REQ-FILEFMT-015 @spec(package_file_format.md#4144-optional-data)
Feature: Parse Optional Data Entries

  @REQ-FILEFMT-015 @happy
  Scenario: Optional data structure contains count and entries
    Given a file entry with OptionalDataLen > 0
    When optional data is parsed
    Then OptionalDataCount (2 bytes) indicates number of entries
    And optional data entries array follows the count
    And entries are located at OptionalDataOffset from start of variable-length data

  @REQ-FILEFMT-015 @happy
  Scenario: Each optional data entry has DataType, DataLength, and Data
    Given a file entry with optional data entries
    When optional data entries are parsed
    Then each entry has DataType (1 byte)
    And each entry has DataLength (2 bytes)
    And each entry has variable-length Data payload
    And DataLength matches the actual data payload length

  @REQ-FILEFMT-015 @happy
  Scenario: Supported optional data types parse correctly
    Given a file entry with optional data entries
    And entries have supported DataType values (0x00-0x08)
    When optional data is parsed
    Then TagsData (0x00) entries parse correctly
    And PathEncoding (0x01) entries parse correctly
    And PathFlags (0x02) entries parse correctly
    And CompressionDictionaryID (0x03) entries parse correctly
    And SolidGroupID (0x04) entries parse correctly
    And FileSystemFlags (0x05) entries parse correctly
    And WindowsAttributes (0x06) entries parse correctly
    And ExtendedAttributes (0x07) entries parse correctly
    And ACLData (0x08) entries parse correctly

  @REQ-FILEFMT-015 @happy
  Scenario: Reserved optional data types are accepted
    Given a file entry with optional data entries
    And entries have reserved DataType values (0x09-0xFF)
    When optional data is parsed
    Then reserved types are accepted
    And entries are parsed with unknown type handling
    And data payload is preserved

  @REQ-FILEFMT-015 @happy
  Scenario Outline: Optional data entries parse with supported types
    Given a file entry with OptionalDataLen > 0
    And optional data entries of types <Types> with lengths <Lengths>
    When optional data is parsed
    Then each entry's DataType is supported or reserved
    And each entry's DataLength matches the data payload

    Examples:
      | Types           | Lengths |
      | 0x00,0x01,0x04  | 12,1,4  |
      | 0x06            | 4       |

  @REQ-FILEFMT-015 @error
  Scenario: Optional data exceeds region bounds is rejected
    Given optional data where DataLength exceeds available region
    When optional data is parsed
    Then a structured invalid optional data error is returned
    And error indicates bounds violation
    And error follows structured error format

  @REQ-FILEFMT-015 @error
  Scenario: Invalid OptionalDataOffset is rejected
    Given a file entry with invalid OptionalDataOffset
    When optional data is parsed
    Then a structured invalid optional data error is returned
    And error indicates invalid offset
    And error follows structured error format

  @happy
  Scenario: WriteTo serializes optional data entry to binary format
    Given an OptionalDataEntry with values
    When optional data entry WriteTo is called with writer
    Then optional data entry is written to writer
    And written data matches optional data entry content

  @happy
  Scenario: ReadFrom deserializes optional data entry from binary format
    Given a reader with valid optional data entry data
    When optional data entry ReadFrom is called with reader
    Then optional data entry is read from reader
    And optional data entry fields match reader data
    And optional data entry is valid

  @happy
  Scenario: Optional data entry round-trip serialization preserves all fields
    Given an OptionalDataEntry with all fields set
    When optional data entry WriteTo is called with writer
    And ReadFrom is called with written data
    Then all optional data entry fields are preserved
    And optional data entry is valid

  @error
  Scenario: ReadFrom fails with incomplete data
    Given a reader with incomplete optional data entry data
    When optional data entry ReadFrom is called with reader
    Then structured IO error is returned
    And error indicates read failure
