@domain:file_format @m2 @REQ-FILEFMT-030 @spec(package_file_format.md#2241-excluded-from-calculation)
Feature: PackageCRC Excluded from Calculation

  @REQ-FILEFMT-030 @happy
  Scenario: Excluded from calculation defines checksum exclusions
    Given a NovusPack package
    And PackageCRC calculation is needed
    When excluded from calculation is examined
    Then checksum exclusions are defined
    And excluded sections are identified
    And exclusion enables proper CRC calculation

  @REQ-FILEFMT-030 @happy
  Scenario: Package header is excluded from PackageCRC calculation
    Given a NovusPack package
    And PackageCRC calculation is performed
    When excluded from calculation is examined
    Then package header is excluded from calculation
    And header exclusion avoids circular dependency
    And header is not included in CRC32 checksum

  @REQ-FILEFMT-030 @happy
  Scenario: Digital signatures are excluded from PackageCRC calculation
    Given a signed NovusPack package
    And PackageCRC calculation is performed
    When excluded from calculation is examined
    Then digital signatures are excluded from calculation
    And signature exclusion allows signature addition without recalculating CRC
    And signatures are not included in CRC32 checksum

  @REQ-FILEFMT-030 @happy
  Scenario: File entries are included in PackageCRC calculation
    Given a NovusPack package
    And PackageCRC calculation is performed
    When calculation scope is examined
    Then file entries are included in calculation
    And file entries contribute to CRC32 checksum
    And file entries are part of integrity validation

  @REQ-FILEFMT-030 @happy
  Scenario: File data is included in PackageCRC calculation
    Given a NovusPack package
    And PackageCRC calculation is performed
    When calculation scope is examined
    Then file data is included in calculation
    And file data contributes to CRC32 checksum
    And file data is part of integrity validation

  @REQ-FILEFMT-030 @happy
  Scenario: File index is included in PackageCRC calculation
    Given a NovusPack package
    And PackageCRC calculation is performed
    When calculation scope is examined
    Then file index is included in calculation
    And file index contributes to CRC32 checksum
    And file index is part of integrity validation

  @REQ-FILEFMT-030 @happy
  Scenario: Package comment is included in PackageCRC calculation
    Given a NovusPack package
    And package has a comment
    And PackageCRC calculation is performed
    When calculation scope is examined
    Then package comment is included in calculation
    And comment contributes to CRC32 checksum
    And comment is part of integrity validation
