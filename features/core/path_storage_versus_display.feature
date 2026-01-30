@domain:core @m1 @REQ-CORE-057 @REQ-CORE-060 @spec(api_core.md#112-package-path-semantics)
Feature: Path Storage Versus Display Format

  Background:
    Given an open NovusPack package
    And a valid context

  @REQ-CORE-057 @happy
  Scenario: Paths are stored internally with leading slash (duplicate removed)
    Given user input path "documents/file.txt"
    When path is normalized and stored
    Then stored path is "/documents/file.txt"
    And stored paths always include a leading slash

  @REQ-CORE-060 @happy
  Scenario: Paths are displayed to users without leading slash
    Given paths stored with leading slash
    When paths are displayed to end users
    Then leading slash is stripped from display
    And users see relative paths
    And display format matches user expectations

  @REQ-CORE-057 @REQ-CORE-060 @happy
  Scenario: Round-trip path handling maintains consistency
    Given user input path "documents/file.txt"
    When path is normalized and stored
    Then stored path is "/documents/file.txt"
    When path is retrieved for display
    Then displayed path is "documents/file.txt"
    And leading slash is only used internally

  @REQ-CORE-060 @happy
  Scenario: File listings never show leading slash to users
    Given package contains multiple files with various paths
    When file listing is generated
    Then all displayed paths lack leading slash
    And paths appear as relative paths
    And users do not see internal storage format

  @REQ-CORE-058 @REQ-CORE-060 @happy
  Scenario: Input paths without leading slash are normalized
    Given user provides path "assets/textures/button.png"
    When path is processed for storage
    Then path is normalized to "/assets/textures/button.png"
    But when shown to user
    Then path is displayed as "assets/textures/button.png"

  @REQ-CORE-059 @happy
  Scenario: Leading slash indicates package root not filesystem root
    Given stored path "/documents/file.txt"
    Then leading slash refers to package root
    And does not refer to filesystem root
    When extracted or displayed
    Then path is relative to extraction directory or current context

  @REQ-CORE-060 @error
  Scenario: Storage format should not leak into user-facing operations
    Given any user-facing operation
    When operation produces output with paths
    Then output paths must not have leading slash
    And internal storage format is not exposed
    And user experience is consistent

  @REQ-CORE-057 @REQ-CORE-060 @happy
  Scenario: Root path special handling
    Given the root path "/" (package root)
    When displayed to users
    Then it may be shown as "." or "" (empty) or root indicator
    But never as "/" to avoid confusion with filesystem root

  @REQ-CORE-057 @error
  Scenario: Dot segments are not permitted in stored paths
    Given user input path with dot segments "documents/./file.txt"
    When path is normalized for storage
    Then dot segments MUST be resolved to canonical form
    And stored path MUST be "/documents/file.txt"
    And dot segments MUST NOT appear in stored paths

  @REQ-CORE-057 @error
  Scenario: Path traversal attempts are rejected
    Given user input path "../../../etc/passwd"
    When path is validated for storage
    Then the path MUST be rejected with validation error
    And error type MUST be ErrTypeValidation
    And paths escaping package root MUST NOT be allowed

  @REQ-CORE-057 @error
  Scenario: Empty paths are not allowed
    Given user input path is empty string
    When path is validated
    Then the path MUST be rejected with validation error
    And error type MUST be ErrTypeValidation
    And error message MUST indicate invalid path

  @REQ-CORE-057 @happy
  Scenario: Trailing slashes distinguish files from directories
    Given stored path "/documents/" with trailing slash
    Then trailing slash MUST indicate a directory entry
    And path "/documents" without trailing slash indicates a file
    And trailing slashes MUST be preserved in storage
    And trailing slash semantics MUST be meaningful

  @REQ-CORE-057 @happy
  Scenario: Path separators are normalized to forward slash
    Given user input path with backslashes "documents\subdir\file.txt"
    When path is normalized for storage
    Then all separators MUST be converted to forward slash
    And stored path MUST be "/documents/subdir/file.txt"
    And Windows-style backslashes MUST be normalized

  @REQ-CORE-057 @happy
  Scenario: Paths are case sensitive
    Given two paths "/documents/File.txt" and "/documents/file.txt"
    When paths are stored and compared
    Then they MUST be treated as different paths
    And case sensitivity MUST be preserved
    And "/File.txt" and "/file.txt" MUST refer to different files
