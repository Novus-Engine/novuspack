@domain:metadata @m2 @REQ-META-065 @spec(api_metadata.md#53-package-index-file-type-65002)
Feature: Package Index File Type 65002

  @REQ-META-065 @happy
  Scenario: Package index file type 65002 provides index file operations
    Given a NovusPack package
    When package index file operations are used
    Then AddIndexFile adds index file
    And GetIndexFile retrieves index
    And UpdateIndexFile updates index
    And RemoveIndexFile removes index file
    And HasIndexFile checks for index file

  @REQ-META-065 @happy
  Scenario: AddIndexFile adds package index file
    Given a NovusPack package
    And IndexData
    When AddIndexFile is called
    Then index file is added to package
    And file type is set to 65002
    And file contains IndexData

  @REQ-META-065 @happy
  Scenario: Package index file provides file navigation and indexing
    Given a NovusPack package
    And a package index file
    When index file is examined
    Then file provides file location mappings
    And file provides content-based indexing
    And file provides search and navigation data
    And file provides file relationship mappings

  @REQ-META-065 @happy
  Scenario: GetIndexFile retrieves index
    Given a NovusPack package
    And a package with index file
    When GetIndexFile is called
    Then IndexData is returned
    And index contains all navigation and indexing information

  @REQ-META-065 @error
  Scenario: Package index file operations handle errors
    Given a NovusPack package
    When invalid index or file operations fail
    Then appropriate errors are returned
    And errors follow structured error format
