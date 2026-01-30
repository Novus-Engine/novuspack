@domain:file_mgmt @metadata @REQ-META-128 @REQ-META-129 @REQ-META-130 @REQ-META-131 @REQ-META-132 @spec(api_metadata.md#8216-packagesetdestpath-method)
Feature: SetDestPath destination override management

  @REQ-META-128 @REQ-META-129 @happy
  Scenario: SetDestPath sets destination override with DestPathOverride struct
    Given an open NovusPack package
    And file "data.txt" exists
    When SetDestPath is called with DestPathOverride{DestPath: stringPtr("/custom/path/data.txt")}
    Then PathMetadataEntry.DestPath is set to "/custom/path/data.txt"
    And destination override is stored

  @REQ-META-128 @REQ-META-129 @happy
  Scenario: SetDestPath sets Windows-specific destination
    Given an open NovusPack package
    And file "data.txt" exists
    When SetDestPath is called with DestPathOverride{DestPathWin: stringPtr("C:\\custom\\data.txt")}
    Then PathMetadataEntry.DestPathWin is set to "C:\\custom\\data.txt"
    And destination override is stored

  @REQ-META-128 @REQ-META-130 @happy
  Scenario: SetDestPathTyped accepts string input
    Given an open NovusPack package
    And file "data.txt" exists
    When SetDestPathTyped is called with string "/custom/path/data.txt"
    Then PathMetadataEntry.DestPath is set to "/custom/path/data.txt"

  @REQ-META-130 @REQ-META-132 @happy
  Scenario: SetDestPathTyped detects Windows path from string
    Given an open NovusPack package
    And file "data.txt" exists
    When SetDestPathTyped is called with string "C:\\custom\\data.txt"
    Then PathMetadataEntry.DestPathWin is set to "C:\\custom\\data.txt"
    And PathMetadataEntry.DestPath is not set

  @REQ-META-130 @happy
  Scenario: SetDestPathTyped accepts map input
    Given an open NovusPack package
    And file "data.txt" exists
    When SetDestPathTyped is called with map{"DestPath": "/unix/path", "DestPathWin": "C:\\win\\path"}
    Then PathMetadataEntry.DestPath is set to "/unix/path"
    And PathMetadataEntry.DestPathWin is set to "C:\\win\\path"

  @REQ-META-128 @happy
  Scenario: SetDestPath creates PathMetadataEntry if missing
    Given an open NovusPack package
    And file "data.txt" exists
    And PathMetadataEntry does not exist for "data.txt"
    When SetDestPath is called with destination override
    Then PathMetadataEntry is created
    And destination override is stored

  @REQ-META-131 @happy
  Scenario: SetDestPath normalizes storedPath by prefixing leading slash
    Given an open NovusPack package
    And file "data.txt" exists
    When SetDestPath is called with storedPath "data.txt" (no leading slash)
    Then PathMetadataEntry is found or created for "/data.txt"
    And destination override is stored

  @REQ-META-132 @happy
  Scenario: SetDestPath detects UNC paths as Windows-only
    Given an open NovusPack package
    And file "data.txt" exists
    When SetDestPathTyped is called with string "\\\\server\\share\\data.txt"
    Then PathMetadataEntry.DestPathWin is set to "\\\\server\\share\\data.txt"
    And PathMetadataEntry.DestPath is not set

  @REQ-META-130 @error
  Scenario: SetDestPathTyped rejects invalid map keys
    Given an open NovusPack package
    And file "data.txt" exists
    When SetDestPathTyped is called with map{"InvalidKey": "value"}
    Then ErrTypeValidation error is returned
    And error indicates invalid key

  @REQ-META-128 @happy
  Scenario: SetDestPath updates existing PathMetadataEntry
    Given an open NovusPack package
    And file "data.txt" exists with PathMetadataEntry
    And PathMetadataEntry.DestPath is "/old/path"
    When SetDestPath is called with DestPathOverride{DestPath: stringPtr("/new/path")}
    Then PathMetadataEntry.DestPath is updated to "/new/path"
    And old destination is replaced
