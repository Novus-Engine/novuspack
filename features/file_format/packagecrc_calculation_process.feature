@domain:file_format @m2 @REQ-FILEFMT-029 @spec(package_file_format.md#224-packagecrc-calculation-process)
Feature: PackageCRC Calculation Process

  @REQ-FILEFMT-029 @happy
  Scenario: PackageCRC calculation process defines checksum computation
    Given a NovusPack package
    And PackageCRC calculation is needed
    When PackageCRC calculation process is examined
    Then checksum computation is defined
    And calculation order is specified
    And calculation scope includes required sections

  @REQ-FILEFMT-029 @happy
  Scenario: PackageCRC is calculated over data in order
    Given a NovusPack package
    And PackageCRC calculation is performed
    When calculation process is examined
    Then calculation follows specified order
    And order is: file entries and data, file index, package comment
    And order ensures consistent checksum computation

  @REQ-FILEFMT-029 @happy
  Scenario: File entries and data are calculated first
    Given a NovusPack package
    And PackageCRC calculation is performed
    When calculation order is examined
    Then file entries and data are calculated first
    And all file entries and associated data are included
    And compressed/encrypted content is included in calculation

  @REQ-FILEFMT-029 @happy
  Scenario: File index is calculated second
    Given a NovusPack package
    And PackageCRC calculation is performed
    When calculation order is examined
    Then file index is calculated second
    And complete file index section is included
    And file index contributes to CRC32 checksum

  @REQ-FILEFMT-029 @happy
  Scenario: Package comment is calculated third if present
    Given a NovusPack package
    And package has a comment
    And PackageCRC calculation is performed
    When calculation order is examined
    Then package comment is calculated third
    And package comment section is included if present
    And comment contributes to CRC32 checksum

  @REQ-FILEFMT-029 @happy
  Scenario: PackageCRC uses CRC32 algorithm
    Given a NovusPack package
    And PackageCRC calculation is performed
    When calculation algorithm is examined
    Then PackageCRC uses CRC32 algorithm
    And algorithm matches file-level checksums for consistency
    And CRC32 provides integrity validation
