@domain:writing @m2 @REQ-WRITE-066 @spec(api_writing.md#533-packagewrite-method)
Feature: Package.Write method behavior

  @REQ-WRITE-066 @happy
  Scenario: Package.Write uses internal compression methods based on package state
    Given an open writable package
    When Package.Write is called
    Then internal compression handling is applied based on header and FileEntry state
    And compression is not provided as a direct parameter

  @REQ-WRITE-066 @happy
  Scenario: Package.Write selects SafeWrite vs FastWrite based on compression state
    Given an open writable package
    When Package.Write is called
    Then SafeWrite is selected for package-compressed packages
    And FastWrite is selected for uncompressed packages when criteria are met

