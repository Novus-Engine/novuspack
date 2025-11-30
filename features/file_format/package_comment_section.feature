@domain:file_format @m2 @REQ-FILEFMT-062 @spec(package_file_format.md#6-package-comment-section-optional)
Feature: Package Comment Section

  @REQ-FILEFMT-062 @happy
  Scenario: Package comment section provides comment storage
    Given a NovusPack package
    And package has a comment
    When package comment section is examined
    Then comment storage is provided
    And comment can be stored in package
    And comment is optional section

  @REQ-FILEFMT-062 @happy
  Scenario: Package comment section stores human-readable description
    Given a NovusPack package
    And package has a comment
    When package comment section is used
    Then human-readable description is stored
    And comment provides package information
    And comment is accessible without decompression

  @REQ-FILEFMT-062 @happy
  Scenario: Package comment section is optional
    Given a NovusPack package
    And package has no comment
    When package comment section is examined
    Then comment section is optional
    And CommentLength is 0
    And comment section is omitted
