@domain:metadata @m2 @REQ-META-062 @spec(api_metadata.md#5-special-metadata-file-types)
Feature: Special metadata file types define special file classifications

  @REQ-META-062 @happy
  Scenario: Special metadata file types define classifications
    Given a package with special metadata files
    When special file types are queried or used
    Then special metadata file types define classifications as specified
    And file type 65000, 65001, 65002, 65003 are recognized
    And the behavior matches the special metadata file types specification
    And special files are handled correctly
