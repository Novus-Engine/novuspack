@domain:file_format @m2 @REQ-FILEFMT-034 @spec(package_file_format.md#242-vendorid-appid-combination-examples)
Feature: VendorID + AppID Combination Examples

  @REQ-FILEFMT-034 @happy
  Scenario: VendorID + AppID combination examples demonstrate combined usage
    Given a NovusPack package
    When VendorID + AppID combination examples are examined
    Then examples demonstrate combined usage
    And examples show platform-application identification
    And examples illustrate combined identifier encoding

  @REQ-FILEFMT-034 @happy
  Scenario: Steam CS:GO combination example
    Given a NovusPack package
    When VendorID is set to 0x53544541 (STEAM)
    And AppID is set to 0x00000000000002DA (730)
    Then package is associated with Steam CS:GO
    And VendorID + AppID combination identifies platform application
    And combined identifiers enable platform-specific identification

  @REQ-FILEFMT-034 @happy
  Scenario: Epic Games Fortnite combination example
    Given a NovusPack package
    When VendorID is set to 0x45504943 (EPIC)
    And AppID is set to 0x00000000ABCDEF01
    Then package is associated with Epic Games Fortnite
    And VendorID + AppID combination identifies platform application
    And combined identifiers demonstrate Epic Games Store format

  @REQ-FILEFMT-034 @happy
  Scenario: GOG Witcher 3 combination example
    Given a NovusPack package
    When VendorID is set to 0x474F4720 (GOG)
    And AppID is set to 0x0000000012345678
    Then package is associated with GOG Witcher 3
    And VendorID + AppID combination identifies platform application
    And combined identifiers demonstrate GOG format

  @REQ-FILEFMT-034 @happy
  Scenario: Itch.io Indie Game combination example
    Given a NovusPack package
    When VendorID is set to 0x49544348 (ITCH)
    And AppID is set to 0x0000000056789ABC
    Then package is associated with Itch.io Indie Game
    And VendorID + AppID combination identifies platform application
    And combined identifiers demonstrate Itch.io format

  @REQ-FILEFMT-034 @happy
  Scenario: Unity Asset combination example
    Given a NovusPack package
    When VendorID is set to 0x554E4954 (UNIT)
    And AppID is set to 0x00000000FEDCBA98
    Then package is associated with Unity Asset
    And VendorID + AppID combination identifies platform application
    And combined identifiers demonstrate Unity Asset Store format

  @REQ-FILEFMT-034 @happy
  Scenario: Combined identifiers enable platform-application identification
    Given a NovusPack package
    When VendorID and AppID are set to specific values
    Then combined identifiers uniquely identify platform and application
    And platform-application association is established
    And combined usage enables distribution system integration
