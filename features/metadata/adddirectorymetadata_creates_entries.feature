@domain:metadata @m2 @REQ-META-121 @spec(api_metadata.md#82-pathmetadata-management-methods)
Feature: AddDirectoryMetadata creates directory metadata entries without adding file content

  @REQ-META-121 @happy
  Scenario: AddDirectoryMetadata creates directory metadata entries
    Given a package and a directory path
    When AddDirectoryMetadata is invoked
    Then directory metadata entries are created without adding file content
    And path metadata is created as specified
    And the behavior matches the path metadata management methods specification
    And file content is not added for directory entries
