@domain:file_types @m2 @REQ-FILETYPES-006 @spec(file_type_system.md#12-special-file-naming-strategy)
Feature: Special File Naming Strategy

  @REQ-FILETYPES-006 @happy
  Scenario: Special file naming strategy uses systematic naming convention
    Given a NovusPack package
    When special file naming strategy is examined
    Then prefix "__NVPK_" clearly identifies NovusPack special files
    And type code provides abbreviated type identifier
    And type ID provides numeric type identifier
    And suffix "__" provides delimiter for consistency
    And extension provides unique extension for each type

  @REQ-FILETYPES-006 @happy
  Scenario: Special file naming strategy components are structured
    Given a NovusPack package
    When special file names are examined
    Then prefix is "__NVPK_"
    And type codes include "META", "MAN", "IDX", "SIG"
    And type IDs include 240, 241, 242, 243
    And suffix is "__"
    And extensions include ".nvpkmeta", ".nvpkman", ".nvpkidx", ".nvpksig"

  @REQ-FILETYPES-006 @happy
  Scenario: Special file naming strategy ensures uniqueness
    Given a NovusPack package
    When special file names are examined
    Then each special file type has unique extension
    And naming convention prevents conflicts with regular files
    And systematic naming enables easy identification

  @REQ-FILETYPES-006 @error
  Scenario: Special file naming strategy handles invalid names
    Given a NovusPack package
    When invalid special file name is provided
    Then naming strategy validation detects invalid format
    And appropriate error is returned
