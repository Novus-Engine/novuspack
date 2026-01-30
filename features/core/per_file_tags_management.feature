@domain:core @m1 @REQ-CORE-012 @spec(api_core.md#8-per-file-tags-management) @spec(api_file_mgmt_file_entry.md#3-tag-management)
Feature: Per-file tags management

  @happy
  Scenario: SetFileTags sets tags for a specific file
    Given an open writable NovusPack package
    When SetFileTags is called with path and tags map
    Then tags are set for the file
    And tags are stored in file entry
    And tags are accessible

  @happy
  Scenario: GetFileTags retrieves tags for a specific file
    Given an open NovusPack package with tagged file
    When GetFileTags is called with file path
    Then tags map is returned
    And all file tags are included
    And tag values are correct

  @happy
  Scenario: UpdateFileTags updates existing tags
    Given an open writable NovusPack package with tagged file
    When UpdateFileTags is called with updates map
    Then existing tags are updated
    And new tags are added
    And unchanged tags remain

  @happy
  Scenario: RemoveFileTags removes specific tag keys
    Given an open writable NovusPack package with tagged file
    When RemoveFileTags is called with keys list
    Then specified tag keys are removed
    And other tags remain unchanged

  @happy
  Scenario: ClearFileTags removes all tags from file
    Given an open writable NovusPack package with tagged file
    When ClearFileTags is called
    Then all tags are removed from file
    And file has no tags

  @happy
  Scenario: GetFilesByTag searches files by tag key-value pair
    Given an open NovusPack package with multiple tagged files
    When GetFilesByTag is called with key and value
    Then matching files are returned
    And file paths match tag criteria
    And all matching files are included

  @happy
  Scenario: GetInheritedTags retrieves tags with path inheritance
    Given an open NovusPack package with directory tags
    When GetInheritedTags is called with file path
    Then file tags are included
    And parent directory tags are included
    And tag inheritance hierarchy is respected

  @error
  Scenario: SetFileTags fails for non-existent file
    Given an open writable NovusPack package
    When SetFileTags is called with non-existent path
    Then a structured validation error is returned

  @error
  Scenario: Tag operations fail if package is read-only
    Given a read-only open NovusPack package
    When SetFileTags is called
    Then a structured validation error is returned
