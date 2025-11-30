@domain:file_format @m2 @REQ-FILEFMT-058 @spec(package_file_format.md#422-hashdatalen-field)
Feature: HashDataLen Field

  @REQ-FILEFMT-058 @happy
  Scenario: HashDataLen field stores total length of all hash data
    Given a NovusPack package
    And file entry has hash data
    When HashDataLen field is examined
    Then HashDataLen field stores total hash data length
    And length includes all hash entries
    And length is stored as 2 bytes value

  @REQ-FILEFMT-058 @happy
  Scenario: HashDataLen of 0 indicates no hash data
    Given a NovusPack package
    And file entry has no hash data
    When HashDataLen is examined
    Then HashDataLen is 0
    And no hash data is present
    And hash data section is empty

  @REQ-FILEFMT-058 @happy
  Scenario: HashDataLen enables hash data buffer allocation
    Given a NovusPack package
    And file entry has hash data
    When HashDataLen is read
    Then hash data buffer size can be determined
    And buffer can be allocated appropriately
    And hash data can be read efficiently
