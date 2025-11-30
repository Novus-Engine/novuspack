@domain:file_format @m2 @REQ-FILEFMT-060 @spec(package_file_format.md#424-optionaldatalen-field)
Feature: OptionalDataLen Field

  @REQ-FILEFMT-060 @happy
  Scenario: OptionalDataLen field stores total length of optional data
    Given a NovusPack package
    And file entry has optional data
    When OptionalDataLen field is examined
    Then OptionalDataLen field stores optional data length
    And length includes all optional data entries
    And length is stored as 2 bytes value

  @REQ-FILEFMT-060 @happy
  Scenario: OptionalDataLen of 0 indicates no optional data
    Given a NovusPack package
    And file entry has no optional data
    When OptionalDataLen is examined
    Then OptionalDataLen is 0
    And no optional data is present
    And optional data section is empty

  @REQ-FILEFMT-060 @happy
  Scenario: OptionalDataLen enables optional data buffer allocation
    Given a NovusPack package
    And file entry has optional data
    When OptionalDataLen is read
    Then optional data buffer size can be determined
    And buffer can be allocated appropriately
    And optional data can be read efficiently
