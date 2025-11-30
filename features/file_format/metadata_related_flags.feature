@domain:file_format @m2 @REQ-FILEFMT-036 @spec(package_file_format.md#252-metadata-related-flags)
Feature: Metadata-Related Flags

  @REQ-FILEFMT-036 @happy
  Scenario: Metadata-related flags define metadata flag bits for package features
    Given a NovusPack package
    And package header is examined
    When metadata-related flags are inspected
    Then Bit 7 indicates Metadata-only package
    And Bit 6 indicates Has special metadata files
    And Bit 5 indicates Has per-file tags
    And Bit 4 indicates Has package comment

  @REQ-FILEFMT-036 @happy
  Scenario: Bit 7 indicates metadata-only package
    Given a NovusPack package
    And package contains only special metadata files
    When package header flags are examined
    Then Bit 7 is set to 1
    And flag indicates metadata-only package
    And flag corresponds to special metadata file detection

  @REQ-FILEFMT-036 @happy
  Scenario: Bit 6 indicates package has special metadata files
    Given a NovusPack package
    And package contains special metadata files
    When package header flags are examined
    Then Bit 6 is set to 1
    And flag indicates special metadata files are present
    And flag corresponds to special metadata file detection

  @REQ-FILEFMT-036 @happy
  Scenario: Bit 5 indicates package has per-file tags
    Given a NovusPack package
    And package has files with per-file tags
    When package header flags are examined
    Then Bit 5 is set to 1
    And flag indicates per-file tag system is used
    And flag corresponds to per-file tag system usage

  @REQ-FILEFMT-036 @happy
  Scenario: Bit 4 indicates package has comment
    Given a NovusPack package
    And package has a comment
    When package header flags are examined
    Then Bit 4 is set to 1
    And flag indicates comment is present
    And flag corresponds to HasComment in PackageInfo
