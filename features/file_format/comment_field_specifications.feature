@domain:file_format @m2 @REQ-FILEFMT-064 @spec(package_file_format.md#6111-field-specifications)
Feature: Comment Field Specifications

  @REQ-FILEFMT-064 @happy
  Scenario: Field specifications define comment field format
    Given a NovusPack package
    And package has a comment
    When comment field specifications are examined
    Then comment field format is defined
    And comment field follows UTF-8 encoding
    And comment field is null-terminated

  @REQ-FILEFMT-064 @happy
  Scenario: Comment field uses UTF-8 encoding
    Given a NovusPack package
    And package has a comment
    When comment field is examined
    Then comment uses UTF-8 encoding
    And comment supports international characters
    And comment encoding is standardized

  @REQ-FILEFMT-064 @happy
  Scenario: Comment field is null-terminated
    Given a NovusPack package
    And package has a comment
    When comment field is examined
    Then comment is null-terminated
    And null byte indicates comment end
    And comment length includes null terminator
