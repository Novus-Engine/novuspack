@domain:file_mgmt @m2 @REQ-FILEMGMT-047 @spec(api_file_mgmt_file_entry.md#6-marshaling)
Feature: FileEntry Serialization

  @REQ-FILEMGMT-047 @happy
  Scenario: Serialization supports parsing fileentry from binary data
    Given binary data containing file entry
    When ParseFileEntry is called with the data
    Then FileEntry is parsed from binary data
    And FileEntry structure is created
    And parsing operation completes successfully

  @REQ-FILEMGMT-047 @happy
  Scenario: ParseFileEntry converts binary format to FileEntry
    Given valid binary file entry data
    When ParseFileEntry is called
    Then FileEntry structure is created
    And all fields are populated from binary data
    And structure matches binary format

  @REQ-FILEMGMT-047 @happy
  Scenario: ParseFileEntry handles variable-length data sections
    Given binary data with variable-length sections (paths, hashes, optional data)
    When ParseFileEntry is called
    Then paths are parsed correctly
    And hash data is parsed correctly
    And optional data is parsed correctly
    And variable-length sections are handled

  @REQ-FILEMGMT-047 @happy
  Scenario: ParseFileEntry validates binary format structure
    Given binary data representing file entry
    When ParseFileEntry is called
    Then binary format structure is validated
    And field alignment is verified
    And structure requirements are checked

  @REQ-FILEMGMT-047 @error
  Scenario: ParseFileEntry returns error for invalid binary data
    Given invalid or malformed binary data
    When ParseFileEntry is called
    Then structured parsing error is returned
    And error indicates invalid binary format
    And error follows structured error format

  @REQ-FILEMGMT-047 @error
  Scenario: ParseFileEntry returns error for incomplete data
    Given incomplete binary data (missing fields or sections)
    When ParseFileEntry is called
    Then structured parsing error is returned
    And error indicates incomplete data
    And error follows structured error format
