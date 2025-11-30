@domain:testing @m2 @REQ-TEST-010 @spec(testing.md#22-path-normalization-testing)
Feature: Testing Path Validation

  @REQ-TEST-010 @happy
  Scenario: Path normalization testing validates path normalization
    Given a NovusPack package
    And path normalization testing configuration
    When path normalization testing is performed
    Then tar-like path normalization is tested ("dir//file.txt" becomes "dir/file.txt")
    And relative reference resolution is tested ("dir/./file.txt" becomes "dir/file.txt")
    And parent directory resolution is tested ("dir/../file.txt" becomes "file.txt")
    And multiple normalization is tested ("dir/./../subdir//file.txt" becomes "subdir/file.txt")

  @REQ-TEST-010 @happy
  Scenario: Tar-like path normalization testing validates redundant separators
    Given a NovusPack package
    And path normalization testing configuration
    When tar-like path normalization testing is performed
    Then paths with double separators are normalized ("dir//file.txt" becomes "dir/file.txt")
    And redundant separators are removed
    And normalized paths are consistent

  @REQ-TEST-010 @happy
  Scenario: Relative reference resolution testing validates dot references
    Given a NovusPack package
    And path normalization testing configuration
    When relative reference resolution testing is performed
    Then dot references are resolved ("dir/./file.txt" becomes "dir/file.txt")
    And relative references are normalized correctly

  @REQ-TEST-010 @happy
  Scenario: Parent directory resolution testing validates parent references
    Given a NovusPack package
    And path normalization testing configuration
    When parent directory resolution testing is performed
    Then parent directory references are resolved ("dir/../file.txt" becomes "file.txt")
    And parent references are normalized correctly

  @REQ-TEST-010 @happy
  Scenario: Multiple normalization testing validates complex paths
    Given a NovusPack package
    And path normalization testing configuration
    When multiple normalization testing is performed
    Then complex paths are normalized ("dir/./../subdir//file.txt" becomes "subdir/file.txt")
    And multiple normalization steps are applied correctly
    And complex path normalization is verified
