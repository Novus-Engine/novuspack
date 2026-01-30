@domain:file_mgmt @m2 @REQ-FILEMGMT-319 @spec(api_file_mgmt_addition.md#2671-on-disk-format-limit)
Feature: Path length validation enforces 65,535 byte maximum from PathEntry.PathLength uint16

  @REQ-FILEMGMT-319 @happy
  Scenario: Path length validation enforces 65535 byte maximum
    Given path input for file addition or storage
    When path length validation is applied
    Then 65,535 byte maximum from PathEntry.PathLength uint16 is enforced
    And the behavior matches the on-disk-format-limit specification
    And paths over limit are rejected with validation error
    And PathLength uint16 constraint is respected
