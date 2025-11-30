@domain:metadata @m2 @REQ-META-094 @spec(api_metadata.md#8313-package-header-flags)
Feature: Package Header Flags

  @REQ-META-094 @happy
  Scenario: Package header flags define special metadata file flags
    Given a NovusPack package
    When package header flags are examined
    Then Bit 6 indicates has special metadata files
    And Bit 5 indicates has per-file tags
    And flags are set when special files exist

  @REQ-META-094 @happy
  Scenario: Bit 6 flag indicates special metadata files
    Given a NovusPack package
    And special metadata files in package
    When package header flags are examined
    Then Bit 6 is set to 1 when special files exist
    And flag accurately reflects special file presence

  @REQ-META-094 @happy
  Scenario: Bit 5 flag indicates per-file tags
    Given a NovusPack package
    And directory metadata providing inheritance
    When package header flags are examined
    Then Bit 5 is set to 1 if directory metadata provides inheritance
    And flag accurately reflects per-file tag support

  @REQ-META-094 @happy
  Scenario: Package header flags are updated when special files change
    Given a NovusPack package
    And a valid context
    When special metadata files are added or removed
    Then UpdateSpecialMetadataFlags updates package header flags
    And flags accurately reflect current state
    And context supports cancellation

  @REQ-META-094 @error
  Scenario: Package header flags handle invalid flag states
    Given a NovusPack package
    When invalid flag states are detected
    Then flag validation detects inconsistencies
    And appropriate errors are returned
