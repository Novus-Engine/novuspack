@domain:file_types @m2 @REQ-FILETYPES-008 @spec(file_type_system.md#2-range-based-category-queries)
Feature: Range-Based Category Queries

  @REQ-FILETYPES-008 @happy
  Scenario: Range-based category queries provide category checking
    Given a NovusPack package
    And a file type value
    When range-based category queries are used
    Then category checking functions check if file type is within range
    And range-based queries support all file type categories

  @REQ-FILETYPES-008 @happy
  Scenario: Range-based category queries check binary file range
    Given a NovusPack package
    When IsBinaryFile is called with file type in range 0-999
    Then IsBinaryFile returns true
    When IsBinaryFile is called with file type outside range 0-999
    Then IsBinaryFile returns false

  @REQ-FILETYPES-008 @happy
  Scenario: Range-based category queries check text file range
    Given a NovusPack package
    When IsTextFile is called with file type in range 1000-1999
    Then IsTextFile returns true
    When IsTextFile is called with file type outside range 1000-1999
    Then IsTextFile returns false

  @REQ-FILETYPES-008 @happy
  Scenario: Range-based category queries handle all file type ranges
    Given a NovusPack package
    When category checking functions are used
    Then binary file range 0-999 is supported
    And text file range 1000-1999 is supported
    And script file range 2000-3999 is supported
    And config file range 4000-4999 is supported
    And image file range 5000-6999 is supported
    And audio file range 7000-7999 is supported
    And video file range 8000-9999 is supported
    And system file range 10000-10999 is supported
    And special file range 65000-65535 is supported

  @REQ-FILETYPES-008 @error
  Scenario: Range-based category queries handle values outside all ranges
    Given a NovusPack package
    When file type value is outside all defined ranges
    Then all category checking functions return false
    And reserved range 11000-64999 is not recognized by any category function
