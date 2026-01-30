@domain:generics @m2 @REQ-GEN-030 @spec(api_generics.md#13-pathentry-type)
Feature: PathEntry Type
    As a package developer
    I want a PathEntry type for path representation
    So that file and directory paths are stored consistently with leading slash format

    Background:
        Given the PathEntry type is defined in the generics package
        And PathEntry contains Path string and PathLength uint32 fields

    @REQ-GEN-030 @happy
    Scenario: PathEntry stores file path with leading slash
        When a PathEntry is created with path "/path/to/file.txt"
        Then the Path field MUST be "/path/to/file.txt"
        And the PathLength field MUST equal the byte length of the Path field

    @REQ-GEN-030 @happy
    Scenario: PathEntry stores directory path with leading slash and trailing slash
        When a PathEntry is created with path "/assets/"
        Then the Path field MUST be "/assets/"
        And the PathLength field MUST equal the byte length of the Path field

    @REQ-GEN-030 @happy
    Scenario: PathEntry validates path has leading slash
        When a PathEntry is validated
        Then it MUST verify the Path begins with "/"
        And it MUST return an error if the Path does not begin with "/"

    @REQ-GEN-030 @happy
    Scenario: PathEntry enforces forward slash separators
        When a PathEntry is created with path "/path\to\file.txt"
        Then the Path field MUST use forward slashes "/path/to/file.txt"
        And backslash separators MUST NOT be stored

    @REQ-GEN-030 @happy
    Scenario: PathEntry validates UTF-8 encoding
        When a PathEntry is validated
        Then the Path field MUST be valid UTF-8
        And the validation MUST reject invalid UTF-8 sequences

    @REQ-GEN-030 @happy
    Scenario: PathEntry validates path length consistency
        When a PathEntry is validated
        Then the PathLength field MUST exactly match the byte length of the Path field
        And validation MUST fail if PathLength is inconsistent
