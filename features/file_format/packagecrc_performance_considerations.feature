@domain:file_format @m2 @REQ-FILEFMT-031 @spec(package_file_format.md#2242-performance-considerations)
Feature: PackageCRC Performance Considerations

  @REQ-FILEFMT-031 @happy
  Scenario: PackageCRC calculation can be skipped for performance
    Given a NovusPack package
    And package write operation is needed
    When PackageCRC is set to 0
    Then CRC calculation is skipped
    And write operation performance is improved
    And PackageCRC zero value indicates calculation was skipped

  @REQ-FILEFMT-031 @happy
  Scenario: PackageCRC can be calculated after write operation
    Given a NovusPack package
    And PackageCRC was skipped during write (set to 0)
    When PackageCRC is calculated using API methods
    Then PackageCRC is updated with calculated value
    And integrity validation is enabled after write
    And CRC can be updated post-write for performance

  @REQ-FILEFMT-031 @happy
  Scenario: PackageCRC calculation is recommended for production packages
    Given a NovusPack package
    And package is for production use
    When PackageCRC is considered
    Then PackageCRC should be calculated
    And integrity validation is recommended
    And production packages benefit from checksum validation

  @REQ-FILEFMT-031 @happy
  Scenario: PackageCRC calculation can be computationally expensive for large packages
    Given a large NovusPack package
    When PackageCRC calculation is performed
    Then calculation may be computationally expensive
    And performance impact increases with package size
    And skipping calculation is an option for large packages

  @REQ-FILEFMT-031 @happy
  Scenario: PackageCRC is recommended for integrity-critical scenarios
    Given a NovusPack package
    And integrity validation is critical
    When PackageCRC is considered
    Then PackageCRC should be calculated
    And checksum provides integrity validation
    And integrity-critical scenarios benefit from validation
