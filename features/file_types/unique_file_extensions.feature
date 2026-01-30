@domain:file_types @m2 @REQ-FILETYPES-007 @spec(file_type_system.md#131-unique-extensions)
Feature: Unique File Extensions

  @REQ-FILETYPES-007 @happy
  Scenario: Unique extensions provide extension uniqueness
    Given a NovusPack package
    When special file extensions are examined
    Then each special file type has a unique extension
    And extensions ensure uniqueness for special package files

  @REQ-FILETYPES-007 @happy
  Scenario: Package metadata files use unique extension
    Given a NovusPack package
    When special file naming strategy is examined
    Then ".nvpkmeta" extension is used for package metadata files
    And ".nvpkmeta" extension is unique to metadata files
    And metadata files contain YAML content

  @REQ-FILETYPES-007 @happy
  Scenario: Package manifest files use unique extension
    Given a NovusPack package
    When special file naming strategy is examined
    Then ".nvpkman" extension is used for package manifest files
    And ".nvpkman" extension is unique to manifest files
    And manifest files contain YAML content

  @REQ-FILETYPES-007 @happy
  Scenario: Package index files use unique extension
    Given a NovusPack package
    When special file naming strategy is examined
    Then ".nvpkidx" extension is used for package index files
    And ".nvpkidx" extension is unique to index files
    And index files contain YAML content

  @REQ-FILETYPES-007 @happy
  Scenario: Digital signature files use unique extension
    Given a NovusPack package
    When special file naming strategy is examined
    Then ".nvpksig" extension is used for digital signature files
    And ".nvpksig" extension is unique to signature files
    And signature files contain binary content

  @REQ-FILETYPES-007 @error
  Scenario: Unique extensions prevent conflicts with regular file extensions
    Given a NovusPack package
    When special file extensions are examined
    Then ".nvpkmeta" does not conflict with regular file extensions
    And ".nvpkman" does not conflict with regular file extensions
    And ".nvpkidx" does not conflict with regular file extensions
    And ".nvpksig" does not conflict with regular file extensions
