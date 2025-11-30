@domain:file_format @m1 @REQ-FILEFMT-007 @spec(package_file_format.md#27-localeid-field-specification)
Feature: LocaleID field specification

  @happy
  Scenario: LocaleID defaults to zero for system default locale
    Given a new NovusPack package
    When the package is created
    Then LocaleID equals 0
    And LocaleID indicates system default locale

  @happy
  Scenario Outline: LocaleID encodes standard locale identifiers
    Given a NovusPack package
    When LocaleID is set to <LocaleID>
    Then LocaleID equals <LocaleID>
    And LocaleID is a 32-bit unsigned integer
    And LocaleID represents <Locale>

    Examples:
      | LocaleID | Locale |
      | 0x0409   | en-US  |
      | 0x0411   | ja-JP  |
      | 0x040C   | fr-FR  |
      | 0x0407   | de-DE  |
      | 0x0410   | it-IT  |
      | 0x0405   | cs-CZ  |

  @happy
  Scenario: LocaleID applies package-wide to all file paths
    Given a NovusPack package with LocaleID set to 0x0411
    When file paths are encoded
    Then all paths use the Japanese locale encoding
    And LocaleID is preserved in the header
