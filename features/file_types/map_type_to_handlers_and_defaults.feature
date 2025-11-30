@domain:file_types @m2 @spec(file_type_system.md#31-file-type-management)
Feature: Map Type to Handlers and Defaults

  @REQ-FILETYPES-002 @happy
  Scenario: Select handler based on declared type
    Given a file with declared type "image"
    When I select a handler for the file
    Then the selected handler should be the image handler
    And handler selection uses declared type priority

  @REQ-FILETYPES-002 @happy
  Scenario: Select handler based on detected type
    Given a file with an undeclared type and a detectable type "text"
    When I select a handler for the file
    Then the selected handler should be the text handler
    And handler selection uses detected type when declared type is missing

  @REQ-FILETYPES-002 @happy
  Scenario: Handler selection prioritizes declared type over detected type
    Given a file with declared type "binary"
    And a detectable type "text"
    When I select a handler for the file
    Then the selected handler should be the binary handler
    And declared type takes precedence over detected type

  @REQ-FILETYPES-002 @happy
  Scenario: Handler selection uses probe result when no declared type
    Given a file with no declared type
    And a probe result indicating type "script"
    When I select a handler for the file
    Then the selected handler should be the script handler
    And handler selection uses probe result as fallback

  @REQ-FILETYPES-002 @error
  Scenario: Handler selection handles unknown file types
    Given a file with undeclared type
    And no detectable type from probe
    When I select a handler for the file
    Then default handler is selected
    And appropriate fallback handler is used
