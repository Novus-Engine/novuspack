@domain:file_format @m2 @REQ-FILEFMT-059 @spec(package_file_format.md#423-hashdataoffset-field)
Feature: HashDataOffset Field

  @REQ-FILEFMT-059 @happy
  Scenario: HashDataOffset field stores offset to hash data
    Given a NovusPack package
    And file entry has hash data
    When HashDataOffset field is examined
    Then HashDataOffset field stores hash data offset
    And offset is relative to start of variable-length data
    And offset is stored as 4 bytes value

  @REQ-FILEFMT-059 @happy
  Scenario: HashDataOffset enables hash data location
    Given a NovusPack package
    And file entry has hash data
    When HashDataOffset is read
    Then hash data location can be determined
    And hash data can be accessed directly
    And hash data section is navigable

  @REQ-FILEFMT-059 @happy
  Scenario: HashDataOffset of 0 indicates no hash data offset
    Given a NovusPack package
    And file entry has no hash data
    When HashDataOffset is examined
    Then HashDataOffset may be 0 or undefined
    And no hash data offset is specified
