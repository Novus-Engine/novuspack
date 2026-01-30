@domain:file_mgmt @m2 @REQ-FILEMGMT-082 @spec(api_core.md#1114-fileerrorcontext-structure)
Feature: FileErrorContext Structure
    As a package API consumer
    I want typed error context for file operations
    So that error debugging and handling are improved with operation-specific details

    Background:
        Given FileErrorContext provides typed error context
        And the structure contains operation-specific fields for file errors

    @REQ-FILEMGMT-082 @happy
    Scenario: FileErrorContext includes file path information
        When a file operation error occurs
        Then the FileErrorContext MUST include the Path field
        And the path identifies the file involved in the error

    @REQ-FILEMGMT-082 @happy
    Scenario: FileErrorContext includes operation identifier
        When a file operation error occurs
        Then the FileErrorContext MUST include the Operation field
        And the operation identifies what action was being performed

    @REQ-FILEMGMT-082 @happy
    Scenario: FileErrorContext includes file size information
        When a file operation error occurs
        Then the FileErrorContext MAY include the Size field
        And size provides context about the file dimensions

    @REQ-FILEMGMT-082 @happy
    Scenario: FileErrorContext attaches to PackageError
        When a file operation error is created
        Then FileErrorContext MUST be attachable as typed context
        And context MUST be retrievable using GetErrorContext helper

    @REQ-FILEMGMT-082 @happy
    Scenario: Multiple context types can coexist
        When errors involve multiple subsystems
        Then multiple typed contexts MAY be attached to a single error
        And each context type MUST be independently retrievable

    @REQ-FILEMGMT-082 @happy
    Scenario: Error context supports serialization
        When errors need to be logged or transmitted
        Then FileErrorContext MUST be serializable to structured formats
        And serialization MUST preserve all context fields
