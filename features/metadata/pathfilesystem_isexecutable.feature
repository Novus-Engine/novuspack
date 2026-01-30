@domain:metadata @m2 @REQ-META-124 @spec(api_metadata.md#842-pathmetadataentry-filesystem-properties)
Feature: PathFileSystem IsExecutable field tracks file execute permissions

  @REQ-META-124 @happy
  Scenario: PathFileSystem IsExecutable tracks execute permissions
    Given a PathMetadataEntry with filesystem properties
    When PathFileSystem properties are accessed
    Then IsExecutable field tracks file execute permissions as specified
    And the behavior matches the path metadata entry filesystem properties specification
    And execute permission is persisted and retrievable
    And permission is consistent with source file or override
