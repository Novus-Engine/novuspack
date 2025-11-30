@domain:file_types @m2 @REQ-FILETYPES-019 @spec(file_type_system.md#31124-config-file-types-4000-4999)
Feature: Config File Type Constants

  @REQ-FILETYPES-019 @happy
  Scenario: Config file type constants define config file type range
    Given a NovusPack package
    When config file type constants are examined
    Then FileTypeConfigStart is 4000
    And FileTypeConfigEnd is 4999
    And config file types are within range 4000-4999

  @REQ-FILETYPES-019 @happy
  Scenario: Specific config file type constants are defined
    Given a NovusPack package
    When config file type constants are examined
    Then FileTypeYAML is 4000
    And FileTypeJSON is 4001
    And FileTypeXML is 4002
    And FileTypeTOML is 4003
    And FileTypeHOCON is 4004
    And FileTypeEDN is 4005
    And FileTypeCUE is 4006
    And FileTypeProperties is 4007
    And FileTypeINI is 4008

  @REQ-FILETYPES-019 @happy
  Scenario: Config file types are recognized by IsConfigFile
    Given a NovusPack package
    When FileTypeYAML is checked with IsConfigFile
    Then IsConfigFile returns true
    When FileTypeJSON is checked with IsConfigFile
    Then IsConfigFile returns true
    When FileTypeXML is checked with IsConfigFile
    Then IsConfigFile returns true
    When FileTypeINI is checked with IsConfigFile
    Then IsConfigFile returns true

  @REQ-FILETYPES-019 @error
  Scenario: Non-config file types are not recognized by IsConfigFile
    Given a NovusPack package
    When FileTypeText is checked with IsConfigFile
    Then IsConfigFile returns false
    When FileTypeImage is checked with IsConfigFile
    Then IsConfigFile returns false
    When FileTypeBinary is checked with IsConfigFile
    Then IsConfigFile returns false

  @REQ-FILETYPES-019 @happy
  Scenario: Config file types support schema validation
    Given a NovusPack package
    And a config file with type FileTypeYAML
    When the config file is validated
    Then schema validation is performed
    And config parsing is appropriate for config files
