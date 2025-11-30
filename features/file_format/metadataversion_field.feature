@domain:file_format @m2 @REQ-FILEFMT-027 @spec(package_file_format.md#222-metadataversion-field)
Feature: MetadataVersion Field

  @REQ-FILEFMT-027 @happy
  Scenario: MetadataVersion field tracks package metadata version
    Given a NovusPack package
    And package has metadata
    When MetadataVersion field is examined
    Then MetadataVersion field tracks package metadata version
    And version increments on metadata changes
    And version enables metadata change detection

  @REQ-FILEFMT-027 @happy
  Scenario: MetadataVersion increments when metadata changes
    Given a NovusPack package
    And package metadata is modified
    When MetadataVersion is examined
    Then MetadataVersion is incremented
    And version change indicates metadata modification
    And version tracks metadata evolution

  @REQ-FILEFMT-027 @happy
  Scenario: MetadataVersion initial value is 1 for new packages
    Given a new NovusPack package
    When MetadataVersion is examined
    Then MetadataVersion is set to 1
    And initial version indicates new package
    And version starts at 1 for new packages

  @REQ-FILEFMT-027 @happy
  Scenario: MetadataVersion tracks package-level metadata changes
    Given a NovusPack package
    And package-level metadata is modified
    When MetadataVersion is examined
    Then MetadataVersion reflects package metadata changes
    And version tracks package comment changes
    And version tracks package header metadata changes
