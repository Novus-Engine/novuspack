@domain:file_mgmt @m2 @REQ-FILEMGMT-419 @spec(api_file_mgmt_file_entry.md#15-processingstate-type)
Feature: ProcessingState type defines the current state of file data transformations

  @REQ-FILEMGMT-419 @happy
  Scenario: ProcessingState type defines transformation state
    Given a FileEntry with file data transformations
    When ProcessingState is used for state tracking
    Then ProcessingState type defines the current state of file data transformations as specified
    And the behavior matches the ProcessingState type specification
    And state values are Raw, Compressed, Encrypted, etc.
    And state transitions are defined
