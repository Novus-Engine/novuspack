@domain:validation @m2 @REQ-FILEMGMT-038 @spec(file_validation.md#11-file-name-validation)
Feature: File Name Validation

  @REQ-FILEMGMT-038 @happy
  Scenario: File name validation enforces name requirements
    Given a NovusPack package
    And an open NovusPack package
    When file name validation is performed
    Then empty names are prohibited (files with empty names are rejected)
    And whitespace-only names are prohibited (names containing only whitespace are rejected)
    And minimum name requirements are enforced (names must contain at least one non-whitespace character)
    And validation error handling provides clear error messages

  @REQ-FILEMGMT-038 @happy
  Scenario: File name validation rejects empty names
    Given a NovusPack package
    And an open NovusPack package
    When file with empty name is added
    Then file is rejected
    And error message indicates empty name rejection
    And error message indicates which file was rejected

  @REQ-FILEMGMT-038 @happy
  Scenario: File name validation rejects whitespace-only names
    Given a NovusPack package
    And an open NovusPack package
    When file with whitespace-only name is added
    Then file is rejected
    And error message indicates whitespace-only name rejection
    And error message indicates which file was rejected

  @REQ-FILEMGMT-038 @error
  Scenario: File name validation provides clear error messages
    Given a NovusPack package
    And an open NovusPack package
    When file name validation fails
    Then clear error message is provided
    And error message indicates which files were rejected
    And error message indicates why files were rejected
    And error follows structured error format
