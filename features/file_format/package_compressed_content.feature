@domain:file_format @m2 @REQ-FILEFMT-042 @spec(package_file_format.md#311-compressed-content)
Feature: Package Compressed Content

  @REQ-FILEFMT-042 @happy
  Scenario: Compressed content defines what is compressed
    Given a NovusPack package
    And package compression is enabled
    When compressed content is examined
    Then compressed content is defined
    And what is compressed is specified
    And compression scope includes specified sections

  @REQ-FILEFMT-042 @happy
  Scenario: File entries are compressed content
    Given a NovusPack package
    And package compression is enabled
    When compressed content is examined
    Then file entries (directory structure) are compressed
    And file entries are part of compressed content
    And file entries contribute to compressed package size

  @REQ-FILEFMT-042 @happy
  Scenario: File data is compressed content
    Given a NovusPack package
    And package compression is enabled
    When compressed content is examined
    Then file data (actual file contents) is compressed
    And file data is part of compressed content
    And file data contributes to compressed package size

  @REQ-FILEFMT-042 @happy
  Scenario: Package index is compressed content
    Given a NovusPack package
    And package compression is enabled
    When compressed content is examined
    Then package index is compressed
    And file index is part of compressed content
    And file index contributes to compressed package size

  @REQ-FILEFMT-042 @happy
  Scenario: Compressed content excludes header, comment, and signatures
    Given a NovusPack package
    And package compression is enabled
    When compressed content is examined
    Then package header is not compressed
    And package comment is not compressed
    And digital signatures are not compressed
    And excluded sections remain directly accessible
