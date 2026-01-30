@domain:file_mgmt @m2 @REQ-FILEMGMT-420 @spec(api_file_mgmt_file_entry.md#17-optionaldata-structure)
Feature: OptionalData structure represents structured optional data for a file entry

  @REQ-FILEMGMT-420 @happy
  Scenario: OptionalData structure represents optional data
    Given a FileEntry with optional metadata or data
    When OptionalData structure is used
    Then OptionalData structure represents structured optional data for a file entry as specified
    And the behavior matches the OptionalData structure specification
    And optional data is present or absent per field
    And type safety is preserved
