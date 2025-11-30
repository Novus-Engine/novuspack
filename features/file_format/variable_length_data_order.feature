@domain:file_format @m2 @REQ-FILEFMT-055 @spec(package_file_format.md#4141-variable-length-data-order)
Feature: Variable-Length Data Order

  @REQ-FILEFMT-055 @happy
  Scenario: Variable-length data order defines data arrangement
    Given a file entry
    And file entry has variable-length data
    When variable-length data order is examined
    Then data arrangement follows specified order
    And order is: path entries, hash data, optional data
    And data order enables efficient parsing

  @REQ-FILEFMT-055 @happy
  Scenario: Path entries are at offset 0
    Given a file entry
    And PathCount > 0
    When variable-length data order is examined
    Then path entries start at offset 0
    And path entries are first in variable-length data
    And path entries are located immediately after fixed structure

  @REQ-FILEFMT-055 @happy
  Scenario: Hash data is located at HashDataOffset
    Given a file entry
    And HashCount > 0
    When variable-length data order is examined
    Then hash data starts at HashDataOffset
    And HashDataOffset is from start of variable-length data
    And hash data follows path entries

  @REQ-FILEFMT-055 @happy
  Scenario: Optional data is located at OptionalDataOffset
    Given a file entry
    And OptionalDataLen > 0
    When variable-length data order is examined
    Then optional data starts at OptionalDataOffset
    And OptionalDataOffset is from start of variable-length data
    And optional data follows hash data

  @REQ-FILEFMT-055 @happy
  Scenario: Variable-length data order is fixed structure, paths, hashes, optional
    Given a file entry
    And file entry has all variable-length sections
    When variable-length data order is examined
    Then fixed structure comes first (64 bytes)
    And path entries follow at offset 0
    And hash data follows at HashDataOffset
    And optional data follows at OptionalDataOffset

  @REQ-FILEFMT-055 @error
  Scenario: Invalid offset values cause validation failure
    Given a file entry
    And HashDataOffset or OptionalDataOffset points outside variable-length section
    When variable-length data is validated
    Then validation fails
    And structured invalid file entry error is returned
    And offset violations are detected
