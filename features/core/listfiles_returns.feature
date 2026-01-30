@domain:core @m2 @REQ-CORE-087 @spec(api_core.md#1163-listfiles-returns) @spec(api_core.md#packagereaderlistfiles-returns)
Feature: ListFiles returns define sorted file information slice

  @REQ-CORE-087 @happy
  Scenario: ListFiles returns a sorted file information slice
    Given an opened package
    When ListFiles is called
    Then a slice of file information is returned
    And the slice is sorted by PrimaryPath
    And the return type matches the PackageReader ListFiles contract
