@domain:file_format @m1 @REQ-FILEFMT-011 @spec(package_file_format.md#42-fileentry-field-specifications)
Feature: FileEntry field specifications and validation

  @happy
  Scenario: HashCount field defaults to zero
    Given a new file entry
    When the file entry is created
    Then HashCount equals 0
    And HashCount is an 8-bit unsigned integer

  @happy
  Scenario: HashCount supports up to 255 hash entries
    Given a file entry
    When HashCount is set to 255
    Then HashCount equals 255
    And HashCount validation passes

  @happy
  Scenario: HashDataLen defaults to zero when no hashes
    Given a file entry with HashCount equals 0
    When the file entry is created
    Then HashDataLen equals 0
    And HashDataLen is a 16-bit unsigned integer

  @happy
  Scenario: HashDataLen equals total length of all hash data
    Given a file entry with multiple hash entries
    When hash data is added
    Then HashDataLen equals the sum of all hash data lengths
    And HashDataLen does not exceed 65535 bytes

  @error
  Scenario: HashDataLen mismatch with actual hash data is invalid
    Given a file entry with HashCount > 0
    When HashDataLen does not match the actual hash data length
    Then a structured invalid file entry error is returned

  @happy
  Scenario: HashDataOffset defaults to zero when no hashes
    Given a file entry with HashCount equals 0
    When the file entry is created
    Then HashDataOffset equals 0
    And HashDataOffset is a 32-bit unsigned integer

  @happy
  Scenario: HashDataOffset points to hash data in variable-length section
    Given a file entry with hash entries
    When variable-length data is structured
    Then HashDataOffset points to the start of hash data
    And HashDataOffset is relative to start of variable-length data

  @error
  Scenario: HashDataOffset out of bounds is invalid
    Given a file entry with HashCount > 0
    When HashDataOffset points beyond variable-length data section
    Then a structured invalid file entry error is returned

  @happy
  Scenario: OptionalDataLen defaults to zero when no optional data
    Given a new file entry
    When the file entry is created
    Then OptionalDataLen equals 0
    And OptionalDataLen is a 16-bit unsigned integer

  @happy
  Scenario: OptionalDataLen equals total length of all optional data
    Given a file entry with optional data entries
    When optional data is added
    Then OptionalDataLen equals the sum of all optional data lengths
    And OptionalDataLen does not exceed 65535 bytes

  @error
  Scenario: OptionalDataLen mismatch with actual optional data is invalid
    Given a file entry with optional data
    When OptionalDataLen does not match the actual optional data length
    Then a structured invalid file entry error is returned

  @happy
  Scenario: OptionalDataOffset defaults to zero when no optional data
    Given a file entry with OptionalDataLen equals 0
    When the file entry is created
    Then OptionalDataOffset equals 0
    And OptionalDataOffset is a 32-bit unsigned integer

  @happy
  Scenario: OptionalDataOffset points to optional data in variable-length section
    Given a file entry with optional data entries
    When variable-length data is structured
    Then OptionalDataOffset points to the start of optional data
    And OptionalDataOffset is relative to start of variable-length data

  @error
  Scenario: OptionalDataOffset out of bounds is invalid
    Given a file entry with OptionalDataLen > 0
    When OptionalDataOffset points beyond variable-length data section
    Then a structured invalid file entry error is returned

  @happy
  Scenario: Variable-length data ordering is correct
    Given a file entry with paths, hashes, and optional data
    When variable-length data is structured
    Then path entries come first at offset 0
    And hash data comes after paths at HashDataOffset
    And optional data comes after hash data at OptionalDataOffset
    And all offsets are non-overlapping

  @error
  Scenario: Overlapping variable-length data sections are invalid
    Given a file entry with paths, hashes, and optional data
    When HashDataOffset overlaps with path data
    Then a structured invalid file entry error is returned
    When OptionalDataOffset overlaps with hash data
    Then a structured invalid file entry error is returned
