@domain:file_types @m2 @REQ-FILETYPES-028 @spec(file_type_system.md#4112-extension-fallback-mapping)
Feature: Extension Fallback Mapping

  @REQ-FILETYPES-028 @happy
  Scenario: Extension fallback mapping provides extension-based type detection
    Given a NovusPack package
    And a file with extension
    When content-based detection fails
    Then extension fallback mapping checks file extension
    And file type is determined from extension mapping

  @REQ-FILETYPES-028 @happy
  Scenario: Extension fallback mapping handles text file extensions
    Given a NovusPack package
    And a file with extension ".txt"
    When extension fallback mapping processes the file
    Then FileTypeText is returned
    When a file with extension ".text" is processed
    Then FileTypeText is returned

  @REQ-FILETYPES-028 @happy
  Scenario: Extension fallback mapping handles YAML file extensions
    Given a NovusPack package
    And a file with extension ".yaml"
    When extension fallback mapping processes the file
    Then FileTypeYAML is returned
    When a file with extension ".yml" is processed
    Then FileTypeYAML is returned

  @REQ-FILETYPES-028 @happy
  Scenario: Extension fallback mapping handles script and config file extensions
    Given a NovusPack package
    And a file with extension ".lua"
    When extension fallback mapping processes the file
    Then FileTypeLua is returned
    When a file with extension ".ini" is processed
    Then FileTypeINI is returned
    When a file with extension ".cfg" is processed
    Then FileTypeINI is returned
    When a file with extension ".js" is processed
    Then FileTypeJavaScript is returned

  @REQ-FILETYPES-028 @error
  Scenario: Extension fallback mapping handles unknown extensions
    Given a NovusPack package
    And a file with unknown extension ".xyz"
    When extension fallback mapping processes the file
    Then extension mapping fails
    And detection process continues to next stage
