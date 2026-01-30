@domain:file_mgmt @m2 @REQ-FILEMGMT-123 @spec(api_file_mgmt_file_entry.md#4-data-management)
Feature: Generic FileEntry Operations

  @REQ-FILEMGMT-123 @happy
  Scenario: Generic FileEntry operations provide FindFileEntry function
    Given an open NovusPack package
    And a valid context
    And FileEntry objects exist
    When FindFileEntry is called with predicate
    Then first FileEntry matching predicate is returned
    And boolean indicates if found
    And type-safe predicate filtering is used

  @REQ-FILEMGMT-123 @happy
  Scenario: Generic FileEntry operations provide FilterFileEntries function
    Given an open NovusPack package
    And a valid context
    And FileEntry objects exist
    When FilterFileEntries is called with predicate
    Then all FileEntry objects matching predicate are returned
    And slice of FileEntry objects is returned
    And type-safe predicate filtering is used

  @REQ-FILEMGMT-123 @happy
  Scenario: Generic FileEntry operations provide MapFileEntries function
    Given an open NovusPack package
    And a valid context
    And FileEntry objects exist
    When MapFileEntries is called with mapper function
    Then FileEntry objects are transformed using mapper
    And slice of transformed values is returned
    And type-safe mapping is used

  @REQ-FILEMGMT-123 @happy
  Scenario: Generic FileEntry operations support reusable operation patterns
    Given an open NovusPack package
    And a valid context
    When generic FileEntry operations are used
    Then type-safe predicates enable filtering
    And type-safe mappers enable transformation
    And reusable operation patterns are supported
