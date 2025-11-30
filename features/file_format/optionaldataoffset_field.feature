@domain:file_format @m2 @REQ-FILEFMT-061 @spec(package_file_format.md#425-optionaldataoffset-field)
Feature: OptionalDataOffset Field

  @REQ-FILEFMT-061 @happy
  Scenario: OptionalDataOffset field stores offset to optional data
    Given a NovusPack package
    And file entry has optional data
    When OptionalDataOffset field is examined
    Then OptionalDataOffset field stores optional data offset
    And offset is relative to start of variable-length data
    And offset is stored as 4 bytes value

  @REQ-FILEFMT-061 @happy
  Scenario: OptionalDataOffset enables optional data location
    Given a NovusPack package
    And file entry has optional data
    When OptionalDataOffset is read
    Then optional data location can be determined
    And optional data can be accessed directly
    And optional data section is navigable

  @REQ-FILEFMT-061 @happy
  Scenario: OptionalDataOffset of 0 indicates no optional data offset
    Given a NovusPack package
    And file entry has no optional data
    When OptionalDataOffset is examined
    Then OptionalDataOffset may be 0 or undefined
    And no optional data offset is specified
