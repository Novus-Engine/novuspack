@domain:file_mgmt @m2 @REQ-FILEMGMT-317 @spec(api_file_mgmt_addition.md#264-session-base-package-level-automatic-basepath)
Feature: Session Base is established automatically on first absolute path and persists

  @REQ-FILEMGMT-317 @happy
  Scenario: Session Base persists for package construction session
    Given a package construction session
    When first absolute path is added
    Then Session Base is established automatically
    And Session Base persists for the package construction session
    And the behavior matches the session-base specification
    And subsequent absolute paths use the same base
