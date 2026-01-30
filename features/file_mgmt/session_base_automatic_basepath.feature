@domain:file_mgmt @m2 @REQ-FILEMGMT-316 @spec(api_file_mgmt_addition.md#264-session-base-package-level-automatic-basepath)
Feature: Session Base provides package-level automatic BasePath for absolute filesystem paths

  @REQ-FILEMGMT-316 @happy
  Scenario: Session Base provides automatic BasePath
    Given a package construction session and absolute filesystem paths
    When Session Base is used for path derivation
    Then package-level automatic BasePath for absolute paths is provided
    And the behavior matches the session-base specification
    And first absolute path establishes Session Base
    And subsequent paths are relative to Session Base
