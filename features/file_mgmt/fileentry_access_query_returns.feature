@domain:file_mgmt @m2 @REQ-FILEMGMT-273 @spec(api_file_mgmt_queries.md#13-fileentry-access)
Feature: FileEntry access defines query function return types

  @REQ-FILEMGMT-273 @happy
  Scenario: FileEntry access defines return types
    Given an open NovusPack package
    When query functions return FileEntry or slices
    Then FileEntry access defines query function return types as specified
    And the behavior matches the fileentry-access specification
    And return types are FileEntry or []FileEntry and error
    And callers receive complete FileEntry objects
