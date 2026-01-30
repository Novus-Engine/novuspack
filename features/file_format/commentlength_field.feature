@domain:file_format @m2 @v2 @REQ-FILEFMT-069 @spec(package_file_format.md#724-commentlength-field)
Feature: CommentLength Field

  @REQ-FILEFMT-069 @happy
  Scenario: CommentLength field stores comment length
    Given a NovusPack package
    And package has a comment
    When package header is examined
    Then CommentLength field stores comment length
    And length is stored in appropriate format
    And length indicates comment size in bytes

  @REQ-FILEFMT-069 @happy
  Scenario: CommentLength field enables comment size determination
    Given a NovusPack package
    And package has a comment
    When CommentLength field is read
    Then comment size can be determined
    And comment reading can use appropriate buffer size
    And comment boundaries are known

  @REQ-FILEFMT-069 @happy
  Scenario: CommentLength of 0 indicates no comment
    Given a NovusPack package
    And package has no comment
    When CommentLength field is read
    Then CommentLength is 0
    And Comment field is skipped during read
    And no comment is stored
