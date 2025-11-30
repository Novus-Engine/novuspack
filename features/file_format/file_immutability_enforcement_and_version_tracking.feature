@domain:file_format @m2 @REQ-FILEFMT-040 @REQ-FILEFMT-049 @spec(package_file_format.md#291-file-immutability-enforcement)
Feature: File Immutability Enforcement and Version Tracking

  @REQ-FILEFMT-040 @happy
  Scenario: File immutability enforcement prevents file modification
    Given a signed NovusPack package
    And SignatureOffset is non-zero
    When file modification is attempted
    Then modification is prohibited
    And structured immutability error is returned
    And file content cannot be changed on signed packages

  @REQ-FILEFMT-040 @happy
  Scenario: Write operations check SignatureOffset before modification
    Given a NovusPack package
    And write operation is needed
    When write operation is performed
    Then SignatureOffset is checked first
    And if SignatureOffset > 0, modification is prohibited
    And immutability is enforced before any changes

  @REQ-FILEFMT-040 @happy
  Scenario: Only reads and signature addition are allowed on signed packages
    Given a signed NovusPack package
    And SignatureOffset is non-zero
    When read operations are performed
    Then read operations are allowed
    And signature addition operations are allowed
    And all other modifications are prohibited

  @REQ-FILEFMT-049 @happy
  Scenario: File version tracking provides version management
    Given a file entry
    When file version fields are examined
    Then FileVersion tracks file content changes
    And MetadataVersion tracks file metadata changes
    And dual versioning enables granular change detection

  @REQ-FILEFMT-049 @happy
  Scenario: FileVersion increments when file data is modified
    Given a file entry
    And FileVersion has current value
    When file content is modified
    Then FileVersion is incremented
    And version change indicates file data modification
    And version enables change detection

  @REQ-FILEFMT-049 @happy
  Scenario: MetadataVersion increments when file metadata is modified
    Given a file entry
    And MetadataVersion has current value
    When file metadata is modified (paths, tags, compression, encryption)
    Then MetadataVersion is incremented
    And version change indicates metadata modification
    And version enables metadata change detection

  @REQ-FILEFMT-049 @happy
  Scenario: File version fields have initial value of 1 for new files
    Given a new file entry
    When file version fields are examined
    Then FileVersion is set to 1
    And MetadataVersion is set to 1
    And initial versions indicate new file entry

  @REQ-FILEFMT-049 @happy
  Scenario: File version tracking supports incremental operations
    Given a file entry
    And file has version history
    When incremental operations use version fields
    Then version fields enable efficient incremental operations
    And change detection uses version comparison
    And conflict resolution uses version information

  @REQ-FILEFMT-040 @error
  Scenario: File modification on signed package returns immutability error
    Given a signed NovusPack package
    And SignatureOffset > 0
    When file modification is attempted
    Then structured immutability error is returned
    And error indicates signed package cannot be modified
    And modification fails with appropriate error
