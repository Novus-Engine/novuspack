@domain:metadata @m2 @REQ-META-148 @spec(api_metadata.md#646-write-operation-requirements)
Feature: Write operations detect metadata-only status and set header Bit 7

  @REQ-META-148 @happy
  Scenario: Write operations detect metadata-only and sync PackageInfo
    Given a package that may be metadata-only (FileCount = 0)
    When write operations are performed
    Then metadata-only status is automatically detected
    And header Bit 7 is set when FileCount = 0
    And PackageInfo.IsMetadataOnly is synchronized before writing
    And the behavior matches the write operation requirements specification
