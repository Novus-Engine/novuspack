@domain:file_format @m2 @REQ-FILEFMT-032 @REQ-FILEFMT-033 @REQ-FILEFMT-034 @spec(package_file_format.md#231-vendorid-example-mappings)
Feature: VendorID and AppID Usage Examples

  @REQ-FILEFMT-032 @happy
  Scenario: VendorID example mappings demonstrate vendor identifier usage
    Given a NovusPack package
    When VendorID example mappings are examined
    Then vendor identifier usage is demonstrated
    And examples show common storefront and platform identifiers
    And examples illustrate VendorID encoding

  @REQ-FILEFMT-032 @happy
  Scenario: VendorID examples include Steam, Epic, GOG, and other platforms
    Given a NovusPack package
    When VendorID examples are examined
    Then Steam VendorID is 0x53544541 (STEAM)
    And Epic Games Store VendorID is 0x45504943 (EPIC)
    And GOG VendorID is 0x474F4720 (GOG )
    And Itch.io VendorID is 0x49544348 (ITCH)
    And examples cover major distribution platforms

  @REQ-FILEFMT-033 @happy
  Scenario: AppID examples demonstrate application identifier usage
    Given a NovusPack package
    When AppID examples are examined
    Then application identifier usage is demonstrated
    And examples show platform-specific AppID formats
    And examples illustrate AppID encoding for different platforms

  @REQ-FILEFMT-033 @happy
  Scenario: AppID examples include Steam CS:GO, TF2, and custom formats
    Given a NovusPack package
    When AppID examples are examined
    Then Steam CS:GO AppID is 0x00000000000002DA (730)
    And Steam TF2 AppID is 0x00000000000001B8 (440)
    And Itch.io Game ID format is demonstrated
    And Epic Games AppID format is demonstrated
    And examples cover various platform formats

  @REQ-FILEFMT-034 @happy
  Scenario: VendorID + AppID combination examples demonstrate combined usage
    Given a NovusPack package
    When VendorID + AppID combination examples are examined
    Then combined usage is demonstrated
    And examples show platform-application identification
    And examples illustrate combined identifier encoding

  @REQ-FILEFMT-034 @happy
  Scenario: Combination examples include Steam CS:GO, Epic Fortnite, GOG Witcher 3
    Given a NovusPack package
    When combination examples are examined
    Then Steam CS:GO combination is demonstrated (VendorID=0x53544541, AppID=0x00000000000002DA)
    And Epic Games Fortnite combination is demonstrated
    And GOG Witcher 3 combination is demonstrated
    And examples show platform-application pairing
