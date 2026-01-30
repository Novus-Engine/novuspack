@domain:file_format @m2 @REQ-FILEFMT-076 @spec(package_file_format.md#611-newfileindex-function)
Feature: FileIndex factory function

  @REQ-FILEFMT-076 @happy
  Scenario: NewFileIndex creates new FileIndex with zero values
    Given NewFileIndex is called
    Then a FileIndex is returned
    And EntryCount is 0
    And Reserved is 0
    And FirstEntryOffset is 0
    And Entries is empty slice
    And FileIndex is in initialized state

  @REQ-FILEFMT-076 @happy
  Scenario: NewFileIndex creates ready-to-use FileIndex
    Given NewFileIndex is called
    When FileIndex is created
    Then FileIndex can have entries added
    And FileIndex can be serialized
    And FileIndex can be validated
