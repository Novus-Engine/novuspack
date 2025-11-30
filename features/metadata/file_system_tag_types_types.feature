@domain:metadata @m2 @REQ-META-023 @REQ-META-050 @REQ-META-092 @spec(metadata.md#126-file-system)
Feature: File System Tag Types

  @REQ-META-023 @happy
  Scenario: File system tag value types support path values
    Given an open NovusPack package
    And a file entry
    When tags with file system types are set
    Then Path type (0x0D) is supported
    And path values are stored as UTF-8 strings
    And file system paths are accessible

  @REQ-META-023 @happy
  Scenario: File system tag value types support MIME type values
    Given an open NovusPack package
    And a file entry
    When tags with file system types are set
    Then MimeType type (0x0E) is supported
    And MIME type values are stored as UTF-8 strings
    And MIME type information is accessible

  @REQ-META-050 @happy
  Scenario: File search by tags finds files by category tag
    Given an open NovusPack package
    And files with category tags
    When GetFilesByTag is called with category tag
    Then all files with matching category are found
    And search results are accurate
    And tag-based search works correctly

  @REQ-META-050 @happy
  Scenario: File search by tags finds files by type tag
    Given an open NovusPack package
    And files with type tags
    When GetFilesByTag is called with type tag
    Then all files with matching type are found
    And search results match type criteria
    And type-based search works correctly

  @REQ-META-050 @happy
  Scenario: File search by tags finds files by priority level
    Given an open NovusPack package
    And files with priority tags
    When GetFilesByTag is called with priority tag
    Then all high-priority files are found
    And search results match priority criteria
    And priority-based search works correctly

  @REQ-META-092 @happy
  Scenario: File type requirements define special file type rules
    Given an open NovusPack package
    And a special metadata file
    When file type requirements are examined
    Then special files must use types in range 65000-65535
    And special files must have reserved names
    And special files must be uncompressed
    And file type rules are enforced

  @REQ-META-011 @error
  Scenario: File search by tags fails with invalid tag key
    Given an open NovusPack package
    And files with tags
    When GetFilesByTag is called with invalid tag key
    Then structured validation error is returned
    And error indicates invalid tag key

  @REQ-META-011 @error
  Scenario: File type requirements validation fails with invalid type
    Given an open NovusPack package
    And file entry with invalid special file type
    When file type is validated
    Then structured validation error is returned
    And error indicates invalid file type
