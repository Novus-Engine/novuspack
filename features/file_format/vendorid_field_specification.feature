@domain:file_format @m1 @REQ-FILEFMT-005 @spec(package_file_format.md#23-vendorid-field-specification)
Feature: VendorID field specification

  @happy
  Scenario: VendorID defaults to zero for generic packages
    Given a new NovusPack package
    When the package is created
    Then VendorID equals 0
    And VendorID indicates no specific vendor association

  @happy
  Scenario Outline: VendorID encodes common storefront identifiers
    Given a NovusPack package
    When VendorID is set to <VendorID>
    Then VendorID equals <VendorID>
    And VendorID is a 32-bit unsigned integer
    And VendorID represents <Platform>

    Examples:
      | VendorID   | Platform           |
      | 0x53544541 | Steam             |
      | 0x45504943 | Epic Games Store  |
      | 0x474F4720 | GOG               |
      | 0x49544348 | Itch.io           |
      | 0x48554D42 | Humble Bundle     |
      | 0x4D494352 | Microsoft Store   |
      | 0x50534E59 | PlayStation Store |
      | 0x58424F58 | Xbox Store        |
      | 0x4E54444F | Nintendo eShop    |
      | 0x554E4954 | Unity Asset Store |
      | 0x554E5245 | Unreal Marketplace|
      | 0x47495448 | GitHub            |
      | 0x4749544C | GitLab            |
