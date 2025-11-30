@domain:file_format @m2 @REQ-FILEFMT-033 @spec(package_file_format.md#241-appid-examples)
Feature: AppID Examples

  @REQ-FILEFMT-033 @happy
  Scenario: AppID examples demonstrate application identifier usage
    Given a NovusPack package
    When AppID examples are examined
    Then examples demonstrate application identifier usage
    And examples show platform-specific AppID formats
    And examples illustrate AppID encoding

  @REQ-FILEFMT-033 @happy
  Scenario: Steam AppID CS:GO example
    Given a NovusPack package
    When AppID is set to 0x00000000000002DA
    Then AppID represents Steam CS:GO (730)
    And AppID is 64-bit unsigned integer
    And Steam AppID is stored in lower 32 bits

  @REQ-FILEFMT-033 @happy
  Scenario: Steam AppID TF2 example
    Given a NovusPack package
    When AppID is set to 0x00000000000001B8
    Then AppID represents Steam TF2 (440)
    And AppID is 64-bit unsigned integer
    And Steam AppID format is demonstrated

  @REQ-FILEFMT-033 @happy
  Scenario: Itch.io Game ID example
    Given a NovusPack package
    When AppID is set to 0x0000000012345678
    Then AppID represents Itch.io Game ID
    And AppID demonstrates custom format
    And AppID supports numeric game IDs

  @REQ-FILEFMT-033 @happy
  Scenario: Epic Games AppID example
    Given a NovusPack package
    When AppID is set to 0x00000000ABCDEF01
    Then AppID represents Epic Games AppID
    And AppID demonstrates custom format
    And AppID supports Epic Games Store identifiers

  @REQ-FILEFMT-033 @happy
  Scenario: Generic 64-bit AppID example
    Given a NovusPack package
    When AppID is set to 0x1234567890ABCDEF
    Then AppID represents generic 64-bit identifier
    And AppID demonstrates custom 64-bit ID format
    And AppID supports proprietary systems

  @REQ-FILEFMT-033 @happy
  Scenario: AppID defaults to zero for generic packages
    Given a new NovusPack package
    When AppID is examined
    Then AppID equals 0x0000000000000000
    And AppID indicates no specific application association
    And default value demonstrates no association
