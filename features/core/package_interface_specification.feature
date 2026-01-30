@domain:core @m2 @REQ-CORE-061 @spec(api_core.md#13-package-interface)
Feature: Package Interface Specification
    As a package developer
    I want a unified Package interface
    So that both reader and writer capabilities are accessible through a single interface

    Background:
        Given the Package interface combines PackageReader and PackageWriter
        And the interface provides unified access to all package operations

    @REQ-CORE-061 @happy
    Scenario: Package interface provides read operations
        When a Package instance is created
        Then it MUST support all PackageReader operations
        And read operations include ReadFile, ListFiles, GetInfo, GetMetadata, and Validate

    @REQ-CORE-061 @happy
    Scenario: Package interface provides write operations
        When a Package instance is created
        Then it MUST support all PackageWriter operations
        And write operations include Write, SafeWrite, and FastWrite

    @REQ-CORE-061 @happy
    Scenario: Package interface maintains single source of truth
        When package state is modified through write operations
        Then subsequent read operations MUST reflect the updated state
        And there MUST be no inconsistency between reader and writer views

    @REQ-CORE-061 @happy
    Scenario: Package interface enforces write protection after signing
        When a package has been signed
        Then write operations MUST be blocked
        And attempts to modify MUST return structured errors

    @REQ-CORE-061 @happy
    Scenario: Package interface provides lifecycle operations
        When a Package instance is managed
        Then it MUST support Create and CreateWithOptions for package creation
        And it MUST support Close and CloseWithCleanup for resource management
        And it MUST support IsOpen, IsReadOnly, and GetPath for state queries
        And it MUST support Defragment for package optimization

    @REQ-CORE-061 @happy
    Scenario: Package interface supports Create method
        Given a Package instance
        When Create is called with a valid path
        Then a new package MUST be created at the specified path
        And the package MUST be ready for write operations
        And the package MUST be in an open state

    @REQ-CORE-061 @happy
    Scenario: Package interface supports CreateWithOptions method
        Given a Package instance
        When CreateWithOptions is called with path and options
        Then a new package MUST be created with specified options
        And options MUST include compression, encryption, and metadata settings
        And the package MUST reflect the configured options

    @REQ-CORE-061 @happy
    Scenario: Package interface supports Close method
        Given an open Package instance
        When Close is called
        Then all resources MUST be released
        And file handles MUST be closed
        And the package MUST transition to closed state

    @REQ-CORE-061 @happy
    Scenario: Package interface supports CloseWithCleanup method
        Given an open Package instance
        When CloseWithCleanup is called with context
        Then resources MUST be released with proper cleanup
        And temporary files MUST be removed if applicable
        And context cancellation MUST be respected

    @REQ-CORE-061 @happy
    Scenario: Package interface supports IsOpen method
        Given a Package instance
        When IsOpen is called
        Then it MUST return true if package is currently open
        And it MUST return false if package is closed
        And the method MUST not require context parameter

    @REQ-CORE-061 @happy
    Scenario: Package interface supports IsReadOnly method
        Given an open Package instance
        When IsReadOnly is called
        Then it MUST return true if package was opened read-only
        And it MUST return false if package is writable
        And the method MUST reflect the package's access mode

    @REQ-CORE-061 @happy
    Scenario: Package interface supports GetPath method
        Given an open Package instance
        When GetPath is called
        Then it MUST return the filesystem path to the package file
        And the path MUST be absolute or relative as originally specified
        And the method MUST not require context parameter

    @REQ-CORE-061 @happy
    Scenario: Package interface supports Defragment method
        Given an open writable Package instance
        When Defragment is called with context
        Then the package structure MUST be optimized
        And unused space MUST be reclaimed
        And file order MUST be optimized for access patterns
        And context cancellation MUST be respected
