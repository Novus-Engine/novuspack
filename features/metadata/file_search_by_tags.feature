@domain:metadata @m2 @REQ-META-050 @spec(metadata.md#file-search-by-tags)
Feature: File Search by Tags

  @REQ-META-050 @happy
  Scenario: File search by tags searches files by specific tag values
    Given a NovusPack package
    When file search by tags is used
    Then GetFilesByTag searches for files by specific tag values
    And search enables finding files by tag criteria
    And search supports tag-based file discovery

  @REQ-META-050 @happy
  Scenario: File search finds files by category tag
    Given a NovusPack package
    And files with category tags
    When GetFilesByTag is called with category value
    Then all files with matching category tag are returned
    And search finds texture files by category
    And search enables category-based organization

  @REQ-META-050 @happy
  Scenario: File search finds files by type tag
    Given a NovusPack package
    And files with type tags
    When GetFilesByTag is called with type value
    Then all files with matching type tag are returned
    And search finds UI files by type
    And search enables type-based filtering

  @REQ-META-050 @happy
  Scenario: File search finds files by priority level
    Given a NovusPack package
    And files with priority tags
    When GetFilesByTag is called with priority value
    Then all files with matching priority tag are returned
    And search finds high-priority files by priority level
    And search enables priority-based file management

  @REQ-META-050 @error
  Scenario: File search handles invalid tag queries
    Given a NovusPack package
    When invalid tag queries are provided
    Then search validation detects invalid queries
    And appropriate errors are returned
