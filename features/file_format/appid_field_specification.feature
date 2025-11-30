@domain:file_format @m1 @REQ-FILEFMT-006 @spec(package_file_format.md#24-appid-field-specification)
Feature: AppID field specification

  @happy
  Scenario: AppID defaults to zero for generic packages
    Given a new NovusPack package
    When the package is created
    Then AppID equals 0
    And AppID indicates no specific application association

  @happy
  Scenario: AppID supports 64-bit application identifiers
    Given a NovusPack package
    When AppID is set to any 64-bit value
    Then AppID is an unsigned 64-bit integer
    And AppID is preserved correctly

  @happy
  Scenario Outline: AppID supports platform-specific identifiers
    Given a NovusPack package with VendorID=<VendorID>
    When AppID is set to <AppID>
    Then AppID equals <AppID>
    And AppID represents <Description>

    Examples:
      | VendorID   | AppID                 | Description        |
      | 0x53544541 | 0x00000000000002DA    | Steam CS:GO        |
      | 0x53544541 | 0x00000000000001B8    | Steam TF2          |
      | 0x45504943 | 0x00000000ABCDEF01    | Epic Games App     |
      | 0x00000000 | 0x1234567890ABCDEF    | Generic 64-bit ID  |

  @happy
  Scenario: VendorID and AppID combination identifies platform applications
    Given a NovusPack package
    When VendorID is set to 0x53544541 and AppID is set to 0x00000000000002DA
    Then the package is associated with Steam CS:GO
    When VendorID is set to 0x45504943 and AppID is set to 0x00000000ABCDEF01
    Then the package is associated with Epic Games Store application
