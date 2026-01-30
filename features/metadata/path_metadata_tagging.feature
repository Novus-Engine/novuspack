@domain:metadata @m2 @REQ-META-049 @spec(metadata.md#path-tagging)
Feature: Metadata Path Tagging

  @REQ-META-049 @happy
  Scenario: Path tagging sets tags on paths
    Given a NovusPack package
    And a PathMetadataEntry
    When path tagging is used
    Then tags are set on paths
    And tags are inherited by child paths via ParentPath
    And path tags support inheritance

  @REQ-META-049 @happy
  Scenario: Path tags are inherited by child paths
    Given a NovusPack package
    And a PathMetadataEntry with tags
    And child paths in the hierarchy
    When path tagging is used
    Then tags set on paths are inherited by child paths
    And inheritance follows tag inheritance rules via ParentPath
    And child paths receive parent path tags

  @REQ-META-049 @happy
  Scenario: Path tagging example demonstrates usage
    Given a NovusPack package
    And a textures path
    When path is tagged
    Then category tag is set to "texture"
    And compression tag is set to "lossless"
    And mipmaps tag is set to true
    And all files associated with that path inherit these tags via PathMetadataEntry

  @REQ-META-049 @error
  Scenario: Path tagging validates tag values
    Given a NovusPack package
    When invalid path tag values are provided
    Then tag validation detects invalid values
    And appropriate errors are returned
