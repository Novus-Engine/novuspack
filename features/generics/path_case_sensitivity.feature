@domain:generics @m2 @REQ-GEN-044 @spec(api_generics.md#1338-case-sensitivity)
Feature: Path Case Sensitivity

  @REQ-GEN-044 @happy
  Scenario: Paths stored case-sensitively
    Given files "file.txt" and "FILE.txt"
    When both files are added to package
    Then both paths are stored as distinct entries
    And "file.txt" is stored with lowercase 'f'
    And "FILE.txt" is stored with uppercase 'F'
    And case is preserved exactly

  @REQ-GEN-044 @happy
  Scenario: Distinct paths for different cases
    Given paths differing only in case
    When paths are stored in package
    Then each path is treated as unique
    And "Config.json" is distinct from "config.json"
    And "README.md" is distinct from "readme.md"
    And separate FileEntry for each path

  @REQ-GEN-044 @happy
  Scenario: Case preservation examples
    Given file "Config.json" with capital 'C'
    When file is added to package
    Then exact case is preserved in storage
    And path retrieved shows "Config.json"
    And capital 'C' is maintained
    And case-sensitive filesystems extract correctly

  @REQ-GEN-045 @happy
  Scenario: No issues on case-sensitive filesystems
    Given package with paths differing only in case
    When package is extracted on Linux
    Then all files extract successfully
    And "file.txt" and "FILE.txt" both created
    And no conflicts occur
    And filesystem handles distinct files

  @REQ-GEN-045 @error
  Scenario: Extraction error on case-insensitive filesystems
    Given package contains "file.txt" and "FILE.txt"
    When package is extracted on Windows
    Then extraction error is returned
    And error message states case-insensitive conflict
    And error message: "Path 'FILE.txt' conflicts with existing path 'file.txt'"
    And extraction is aborted for conflicting files

  @REQ-GEN-045 @error
  Scenario: Extraction error on default macOS
    Given package contains paths differing only in case
    When package is extracted on default macOS volume
    Then extraction error occurs
    And error indicates case-insensitive filesystem conflict
    And clear error message provided
    And extraction does not silently overwrite

  @REQ-GEN-045 @happy
  Scenario: Explicit error with clear message
    Given conflicting case-sensitive paths in package
    When extracted on case-insensitive filesystem
    Then error is explicit not silent
    And error message clearly states the issue
    And error provides both conflicting path names
    And user understands the problem

  @REQ-GEN-044 @happy
  Scenario: Portability warning for package design
    Given package is being designed
    When case-conflicting paths are considered
    Then documentation warns about portability
    And packages with case conflicts not portable
    And recommendation to avoid case-only differences
    And consider target filesystem types

  @REQ-GEN-045 @happy
  Scenario: Comparison with ZIP/TAR behavior
    Given case-conflicting paths in archive
    When comparing NovusPack to ZIP/TAR
    Then ZIP/TAR may use last-write-wins or error (tool dependent)
    And 7-Zip errors on conflicts
    And NovusPack provides explicit clear error
    And NovusPack behavior is predictable and documented
