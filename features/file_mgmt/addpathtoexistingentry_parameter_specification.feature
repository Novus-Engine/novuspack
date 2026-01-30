@domain:file_mgmt @m2 @REQ-FILEMGMT-199 @spec(api_deduplication.md#314-addpathtoexistingentry-parameters)
Feature: AddPathToExistingEntry Parameter Specification

  @REQ-FILEMGMT-199 @happy
  Scenario: AddPathToExistingEntry parameters include existingEntry and newPath
    Given an open NovusPack package
    And a valid context
    And an existing FileEntry
    And a new path string
    When AddPathToExistingEntry is called
    Then existingEntry parameter references existing file entry
    And newPath parameter specifies new path to add
    And parameters enable path addition to existing entry

  @REQ-FILEMGMT-199 @happy
  Scenario: AddPathToExistingEntry supports deduplication path addition
    Given an open NovusPack package
    And a valid context
    And duplicate content is detected
    And existing FileEntry is found
    When AddPathToExistingEntry is called with new path
    Then new path is added to existing entry
    And content is shared between paths
    And storage space is saved

  @REQ-FILEMGMT-199 @error
  Scenario: AddPathToExistingEntry validates existingEntry parameter
    Given an open NovusPack package
    And a valid context
    And an invalid or nil FileEntry
    When AddPathToExistingEntry is called with invalid entry
    Then appropriate error is returned
    And error indicates invalid entry
    And error follows structured error format

  @REQ-FILEMGMT-199 @error
  Scenario: AddPathToExistingEntry validates newPath parameter
    Given an open NovusPack package
    And a valid context
    And an existing FileEntry
    And an invalid path string
    When AddPathToExistingEntry is called with invalid path
    Then appropriate error is returned
    And error indicates invalid path
    And error follows structured error format
