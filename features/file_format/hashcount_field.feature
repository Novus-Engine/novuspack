@domain:file_format @m2 @REQ-FILEFMT-057 @spec(package_file_format.md#421-hashcount-field)
Feature: HashCount Field

  @REQ-FILEFMT-057 @happy
  Scenario: HashCount field stores number of hash entries
    Given a NovusPack package
    And file entry has hash data
    When HashCount field is examined
    Then HashCount field stores hash count
    And count indicates number of hash entries
    And count is stored as 1 byte value

  @REQ-FILEFMT-057 @happy
  Scenario: HashCount of 0 indicates no hashes
    Given a NovusPack package
    And file entry has no hash data
    When HashCount is examined
    Then HashCount is 0
    And no hash entries are present
    And hash data section is empty

  @REQ-FILEFMT-057 @happy
  Scenario: HashCount enables hash array processing
    Given a NovusPack package
    And file entry has multiple hash entries
    When HashCount is read
    Then hash array size is determined
    And hash entries can be processed iteratively
    And hash data structure is navigable
