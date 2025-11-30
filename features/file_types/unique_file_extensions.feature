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
    Then ".npkmeta" extension is used for package metadata files
    And ".npkmeta" extension is unique to metadata files
    And metadata files contain YAML content

  @REQ-FILETYPES-007 @happy
  Scenario: Package manifest files use unique extension
    Given a NovusPack package
    When special file naming strategy is examined
    Then ".npkman" extension is used for package manifest files
    And ".npkman" extension is unique to manifest files
    And manifest files contain YAML content

  @REQ-FILETYPES-007 @happy
  Scenario: Package index files use unique extension
    Given a NovusPack package
    When special file naming strategy is examined
    Then ".npkidx" extension is used for package index files
    And ".npkidx" extension is unique to index files
    And index files contain YAML content

  @REQ-FILETYPES-007 @happy
  Scenario: Digital signature files use unique extension
    Given a NovusPack package
    When special file naming strategy is examined
    Then ".npksig" extension is used for digital signature files
    And ".npksig" extension is unique to signature files
    And signature files contain binary content

  @REQ-FILETYPES-007 @error
  Scenario: Unique extensions prevent conflicts with regular file extensions
    Given a NovusPack package
    When special file extensions are examined
    Then ".npkmeta" does not conflict with regular file extensions
    And ".npkman" does not conflict with regular file extensions
    And ".npkidx" does not conflict with regular file extensions
    And ".npksig" does not conflict with regular file extensions
