@domain:core @m2 @REQ-CORE-105 @spec(api_core.md#1176-does-not-include)
Feature: GetInfo does not include individual FileEntry metadata or special metadata file contents

  @REQ-CORE-105 @happy
  Scenario: GetInfo excludes per-file and special metadata
    Given an opened package with file entries and metadata
    When GetInfo is called
    Then individual FileEntry metadata is not included
    And special metadata file contents are not included
    And the result is lightweight package-level information only
