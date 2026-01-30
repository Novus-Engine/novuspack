@domain:metadata @m2 @REQ-META-118 @spec(metadata.md#132-path-metadata-entry-structure)
Feature: PathMetadataEntry paths stored with leading slash for package root reference

  @REQ-META-118 @happy
  Scenario: PathMetadataEntry paths use leading slash
    Given PathMetadataEntry structures for path metadata
    When paths are stored in PathMetadataEntry
    Then paths are stored with leading slash for package root reference
    And storage format is consistent with the specification
    And the behavior matches the path metadata entry structure specification
    And display conversion strips leading slash for user output
