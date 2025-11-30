@domain:file_format @m2 @REQ-FILEFMT-048 @spec(package_file_format.md#4121-unique-file-identification)
Feature: Unique File Identification

  @REQ-FILEFMT-048 @happy
  Scenario: Unique file identification provides file uniqueness
    Given a NovusPack package
    And multiple file entries
    When FileID fields are examined
    Then each file entry has a unique 64-bit FileID
    And FileID provides stable identification
    And FileID enables efficient file tracking

  @REQ-FILEFMT-048 @happy
  Scenario: FileID is assigned sequentially during file addition
    Given a NovusPack package
    When files are added to the package
    Then FileID is assigned sequentially (1, 2, 3, ...)
    And each new file gets next available FileID
    And FileID assignment follows sequential pattern

  @REQ-FILEFMT-048 @happy
  Scenario: FileID remains constant for lifetime of file entry
    Given a NovusPack package
    And a file entry with FileID
    When file metadata or paths are modified
    Then FileID remains constant
    And FileID provides stable reference across operations
    And FileID persistence enables reliable file tracking

  @REQ-FILEFMT-048 @happy
  Scenario: FileID enables efficient API operations
    Given a NovusPack package
    And file entries with FileID values
    When file operations use FileID
    Then FileID enables efficient file tracking
    And FileID enables API operations without path-based lookup
    And FileID supports future extensibility

  @REQ-FILEFMT-048 @happy
  Scenario: FileID supports up to 18 quintillion files
    Given a NovusPack package
    When FileID range is examined
    Then FileID is 64-bit unsigned integer
    And FileID supports range up to 18446744073709551615
    And FileID future-proofs file identification

  @REQ-FILEFMT-048 @error
  Scenario: FileID value 0 is reserved and must not be used
    Given a NovusPack package
    When FileID equals 0
    Then FileID 0 is reserved
    And FileID 0 must not be used for actual files
    And validation flags reserved FileID usage
