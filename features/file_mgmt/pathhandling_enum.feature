@domain:file_mgmt @m2 @REQ-FILEMGMT-359 @spec(api_file_mgmt_addition.md#27-pathhandling-type)
Feature: PathHandling enum defines multiple path handling strategies

  @REQ-FILEMGMT-359 @happy
  Scenario: PathHandling enum defines path handling strategies
    Given file addition or extraction with path handling options
    When PathHandling is used
    Then PathHandling enum defines multiple path handling strategies (hard links, symlinks, preserve)
    And the behavior matches the PathHandling type specification
    And strategies are hard links, symlinks, preserve as specified
    And options are applied per PathHandling value
