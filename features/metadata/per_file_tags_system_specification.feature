@domain:metadata @m2 @REQ-META-122 @spec(metadata.md#1-per-file-tags-system-specification)
Feature: Per-File Tags System Specification
    As a package developer
    I want a tag-based metadata architecture for file entries
    So that flexible metadata can be attached to files without schema constraints

    Background:
        Given the per-file tags system provides tag-based metadata
        And tags are stored in FileEntry structures

    @REQ-META-122 @happy
    Scenario: Tags support multiple value types
        When tags are added to a FileEntry
        Then tags MUST support string, integer, float, boolean, and structured types
        And each tag MUST have a TagValueType identifier

    @REQ-META-122 @happy
    Scenario: Tags are stored with key-value pairs
        When a tag is added to a FileEntry
        Then it MUST have a unique key within that FileEntry
        And it MUST have a value and value type

    @REQ-META-122 @happy
    Scenario: Tag storage format enables efficient serialization
        When FileEntry tags are serialized
        Then the binary format MUST be compact and efficient
        And the format MUST support deserialization without data loss

    @REQ-META-122 @happy
    Scenario: Tags support type-safe operations
        When tags are accessed using Tag[T] generic operations
        Then type safety MUST be enforced at compile time
        And runtime type validation MUST prevent type mismatches

    @REQ-META-122 @happy
    Scenario: Tag synchronization keeps TaggedData and Tags in sync
        When tags are added, updated, or removed
        Then TaggedData and Tags fields MUST remain synchronized
        And inconsistencies MUST be detected and prevented

    @REQ-META-122 @happy
    Scenario: Tags enable file search and filtering
        When packages contain many files
        Then files MUST be searchable by tag keys and values
        And tag-based queries MUST be efficient
