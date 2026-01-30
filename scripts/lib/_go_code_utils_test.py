#!/usr/bin/env python3
"""Test script for _go_code_utils normalization logic."""

import sys
from pathlib import Path

# Add scripts directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))

from lib._go_code_utils import (  # noqa: E402
    normalize_go_signature, normalize_generic_name, extract_receiver_type,
    parse_go_def_signature,
    find_go_code_blocks, is_in_go_code_block, is_example_code,
    is_example_signature_name,
    is_definition_start_line, is_continuation_line, is_signature_only_code_block
)


def test_normalize_generic_name():
    """Test generic name normalization."""
    print("Testing normalize_generic_name:")
    test_cases = [
        # Basic generics
        ("Option[T]", "Option"),
        ("BufferPool[T]", "BufferPool"),
        ("ConfigBuilder[T any]", "ConfigBuilder"),
        ("Option", "Option"),
        ("Result[T, E]", "Result"),
        # Multiple type parameters
        ("Map[K, V]", "Map"),
        ("Pair[A, B, C]", "Pair"),
        # With constraints
        ("Option[T comparable]", "Option"),
        ("Result[T any, E error]", "Result"),
        # Nested generics (should remove all)
        ("Container[Option[T]]", "Container"),
        # Edge cases
        ("Type[]", "Type"),  # Empty generic params
        ("Type[T]", "Type"),  # Single param
        ("Type[T, U, V]", "Type"),  # Multiple params
    ]
    for input_name, expected in test_cases:
        result = normalize_generic_name(input_name)
        status = "✓" if result == expected else "✗"
        print(f"  {status} {input_name} -> {result} (expected: {expected})")
    print()


def test_package_qualified_types():
    """Test package-qualified type name normalization."""
    print("Testing package-qualified type normalization:")
    test_cases = [
        # Basic package types
        ("generics.Tag", "Tag"),
        ("metadata.PackageMetadata", "PackageMetadata"),
        ("fileformat.PackageHeader", "PackageHeader"),
        ("pkgerrors.ErrorType", "ErrorType"),
        ("signatures.Signature", "Signature"),

        # Already normalized (no package)
        ("Tag", "Tag"),
        ("PackageMetadata", "PackageMetadata"),

        # In function signatures
        ("func (p *Package) AddFile(ctx context.Context, path string) error",
         "func (p *Package) AddFile(ctx context.Context, path string) error"),
        ("func CreatePackage(ctx context.Context, config generics.Tag) (*Package, error)",
         "func CreatePackage(ctx context.Context, config Tag) (*Package, error)"),

        # Multiple package types
        ("func Process(data metadata.PackageInfo, sig signatures.Signature) pkgerrors.PackageError",
         "func Process(data PackageInfo, sig Signature) PackageError"),

        # With pointers and slices
        ("func GetInfo() (*metadata.PackageMetadata, error)",
         "func GetInfo() (*PackageMetadata, error)"),
        ("func GetFiles() []fileformat.FileEntry",
         "func GetFiles() []FileEntry"),

        # In method receivers
        ("func (p *Package) SetMetadata(info metadata.PackageInfo)",
         "func (p *Package) SetMetadata(info PackageInfo)"),

        # Nested packages (if any)
        ("internal.helper.Utility", "Utility"),

        # Edge cases
        ("func Test() string", "func Test() string"),  # No package types
        ("type MyType struct", "type MyType struct"),  # Type definition
    ]

    for input_sig, expected in test_cases:
        result = normalize_go_signature(input_sig)
        # For comparison, we'll check if package qualifiers were removed
        packages = ['generics.', 'metadata.', 'fileformat.', 'pkgerrors.', 'signatures.']
        has_package_qualifier = any(pkg in input_sig for pkg in packages)
        if has_package_qualifier:
            # Should have removed package qualifiers
            has_removed = not any(pkg in result for pkg in packages)
            status = "✓" if has_removed else "✗"
            print(f"  {status} Package qualifiers removed: {input_sig[:60]}...")
        else:
            # Should remain unchanged (no package qualifiers to remove)
            status = "✓" if result == expected or (not has_package_qualifier) else "✗"
            print(f"  {status} No change needed: {input_sig[:60]}...")
    print()


def test_extract_receiver_type():
    """Test receiver type extraction."""
    print("Testing extract_receiver_type:")
    test_cases = [
        # Basic receivers
        ("(p *Package)", "Package"),
        ("(fe *FileEntry)", "FileEntry"),
        ("(o *Option[T])", "Option"),
        ("(r Result[T])", "Result"),
        ("(p Package)", "Package"),
        ("Package", "Package"),  # Already extracted
        # With generics (normalized)
        ("(o *Option[T])", "Option"),
        ("(r Result[T, E])", "Result"),
        # Without generics normalization
        ("(o *Option[T])", "Option[T]", False),  # normalize_generics=False
        ("(r Result[T, E])", "Result[T, E]", False),
        # Edge cases
        ("(*Package)", "Package"),  # No variable name
        ("(Package)", "Package"),  # No pointer, no variable
        ("pkg *Package", "Package"),  # No parentheses
        ("*Package", "Package"),  # Just pointer type
    ]
    for test_case in test_cases:
        if len(test_case) == 3:
            input_receiver, expected, normalize = test_case
        else:
            input_receiver, expected = test_case
            normalize = True
        result = extract_receiver_type(input_receiver, normalize_generics=normalize)
        status = "✓" if result == expected else "✗"
        norm_str = f", normalize={normalize}" if len(test_case) == 3 else ""
        print(f"  {status} {input_receiver}{norm_str} -> {result} (expected: {expected})")
    print()


def test_full_signature_normalization():
    """Test full signature normalization with package types."""
    print("Testing full signature normalization:")
    test_cases = [
        # Methods with single return
        (
            ("func (p *Package) AddFile(ctx context.Context, path string, "
             "info metadata.PackageInfo) error"),
            "func (Package) AddFile(context.Context, string, PackageInfo) error"
        ),
        # Functions with multiple returns
        (
            "func CreatePackage(ctx context.Context, config generics.Tag) (*Package, error)",
            "func CreatePackage(context.Context, Tag) (*Package, error)"
        ),
        (
            "func (p *Package) GetMetadata() (*metadata.PackageMetadata, error)",
            "func (Package) GetMetadata() (*PackageMetadata, error)"
        ),
        # Single return with package type
        (
            "func Validate(sig signatures.Signature) pkgerrors.PackageError",
            "func Validate(Signature) PackageError"
        ),
        # No return value
        (
            "func (p *Package) Close()",
            "func (Package) Close()"
        ),
        # Multiple parameters with mixed package types
        (
            ("func Process(data metadata.PackageInfo, sig signatures.Signature, "
             "err pkgerrors.PackageError) error"),
            "func Process(PackageInfo, Signature, PackageError) error"
        ),
        # With slices and pointers
        (
            "func GetFiles() []fileformat.FileEntry",
            "func GetFiles() []FileEntry"
        ),
        (
            "func GetInfo() *metadata.PackageMetadata",
            "func GetInfo() *PackageMetadata"
        ),
        # Standard library types preserved
        (
            "func Read(r io.Reader) ([]byte, error)",
            "func Read(io.Reader) ([]byte, error)"
        ),
        (
            "func Write(ctx context.Context, data []byte) (int, error)",
            "func Write(context.Context, []byte) (int, error)"
        ),
        # Complex return types
        (
            "func Get() (*metadata.PackageMetadata, []fileformat.FileEntry, error)",
            "func Get() (*PackageMetadata, []FileEntry, error)"
        ),
        # Generic types in signatures
        (
            "func Create[T any](ctx context.Context, config generics.Tag[T]) (*Package, error)",
            "func Create(context.Context, Tag) (*Package, error)"
        ),
        # Methods with no parameters
        (
            "func (p *Package) String() string",
            "func (Package) String() string"
        ),
        # Functions with variadic parameters
        (
            "func Join(items ...string) string",
            "func Join(...string) string"
        ),
        # Interface methods (no receiver in interface definition)
        (
            "func (p *Package) Open(ctx context.Context) error",
            "func (Package) Open(context.Context) error"
        ),
    ]

    for input_sig, expected in test_cases:
        result = normalize_go_signature(input_sig)
        # Compare normalized versions (may differ in whitespace)
        result_clean = ' '.join(result.split())
        expected_clean = ' '.join(expected.split())
        status = "✓" if result_clean == expected_clean else "✗"
        print(f"  {status}")
        print(f"    Input:  {input_sig}")
        print(f"    Result: {result}")
        print(f"    Expected: {expected}")
        if result_clean != expected_clean:
            print("    MISMATCH!")
        print()


def test_parse_go_def_signature():
    """Test parsing Go definition signatures (unified function)."""
    print("Testing parse_go_def_signature:")
    test_cases = [
        # Methods
        ("func (p *Package) AddFile(ctx context.Context, path string) error",
         ("AddFile", "method", "Package", "ctx context.Context, path string", "error")),
        # Functions
        ("func CreatePackage(ctx context.Context) (*Package, error)",
         ("CreatePackage", "func", None, "ctx context.Context", "(*Package, error)")),
        # No parameters
        ("func String() string",
         ("String", "func", None, "", "string")),
        # No return
        ("func Close()",
         ("Close", "func", None, "", "")),
        # Generic function
        ("func Create[T any](config T) T",
         ("Create", "func", None, "config T", "T")),
        # Structs
        ("type Package struct", ("Package", "struct", None, None, None)),
        ("type FileEntry struct {", ("FileEntry", "struct", None, None, None)),
        # Interfaces
        ("type Reader interface", ("Reader", "interface", None, None, None)),
        ("type Writer interface {", ("Writer", "interface", None, None, None)),
        # Type aliases
        ("type MyInt int", ("MyInt", "type", None, None, None)),
        ("type Handler func()", ("Handler", "type", None, None, None)),
        # Generic types
        # Note: generic_params stored in Signature object, not in tuple
        ("type Option[T any] struct", ("Option", "struct", None, None, None)),
        ("type Result[T, E] struct", ("Result", "struct", None, None, None)),
    ]
    for input_line, expected in test_cases:
        result = parse_go_def_signature(input_line)
        if result is None:
            status = "✗"
            print(f"  {status} {input_line[:50]}... -> None (expected signature)")
            continue
        # Compare: (name, kind, receiver, params, returns)
        actual = (result.name, result.kind, result.receiver, result.params, result.returns)
        status = "✓" if actual == expected else "✗"
        print(f"  {status} {input_line[:50]}...")
        if actual != expected:
            print(f"    Got: {actual}")
            print(f"    Expected: {expected}")
    print()


def test_find_go_code_blocks():
    """Test finding Go code blocks in markdown."""
    print("Testing find_go_code_blocks:")
    content = """# Title

Some text.

```go
type Package struct {
    Name string
}
```

More text.

```go
func Create() *Package {
    return &Package{}
}
```

```rust
// This should be ignored
```

```go
// Another block
type File struct {}
```
"""
    blocks = find_go_code_blocks(content)
    expected_count = 3
    status = "✓" if len(blocks) == expected_count else "✗"
    print(f"  {status} Found {len(blocks)} code blocks (expected: {expected_count})")
    for i, (start, end, code) in enumerate(blocks, 1):
        print(f"    Block {i}: lines {start}-{end}, {len(code.split(chr(10)))} lines")
    print()


def test_is_in_go_code_block():
    """Test checking if line is in Go code block."""
    print("Testing is_in_go_code_block:")
    content = """Line 1
Line 2
```go
Line 4 (in block)
Line 5 (in block)
```
Line 7 (not in block)
```go
Line 9 (in block)
```
Line 11 (not in block)
"""
    test_cases = [
        (1, False),  # Before any block
        (4, True),   # In first block
        (5, True),   # In first block
        (7, False),  # Between blocks
        (9, True),   # In second block
        (11, False),  # After blocks
    ]
    for line_num, expected in test_cases:
        result = is_in_go_code_block(content, line_num)
        status = "✓" if result == expected else "✗"
        print(f"  {status} Line {line_num}: {result} (expected: {expected})")
    print()


def test_is_example_code():
    """Test example code detection."""
    print("Testing is_example_code:")

    # Test case 1: ExampleType with "hypothetical" marker
    code1 = "// This is a hypothetical example\ntype ExampleType struct {\n    Field string\n}"
    lines1 = code1.split('\n')
    result1 = is_example_code(code1, 1, lines=lines1, check_single_line=1)
    status1 = "✓" if result1 is True else "✗"
    print(f"  {status1} ExampleType (hypothetical marker): {result1} (expected: True)")

    # Test case 2: Package - no example markers
    code2 = "// Regular type\ntype Package struct {\n    Name string\n}"
    lines2 = code2.split('\n')
    result2 = is_example_code(code2, 1, lines=lines2, check_single_line=1)
    status2 = "✓" if result2 is False else "✗"
    print(f"  {status2} Package (no markers): {result2} (expected: False)")

    # Test case 3: HypotheticalStruct with "not the actual" marker
    code3 = "// This is not the actual type\ntype HypotheticalStruct struct {}"
    lines3 = code3.split('\n')
    result3 = is_example_code(code3, 1, lines=lines3, check_single_line=1)
    status3 = "✓" if result3 is True else "✗"
    print(f"  {status3} HypotheticalStruct (not the actual): {result3} (expected: True)")

    # Test case 4: MockService - starts with Mock
    code4 = "type MockService struct {}"
    lines4 = code4.split('\n')
    result4 = is_example_code(code4, 1, lines=lines4, check_single_line=0)
    status4 = "✓" if result4 is True else "✗"
    print(f"  {status4} MockService (starts with Mock): {result4} (expected: True)")

    # Test case 5: TestHelper - starts with Test
    code5 = "type TestHelper struct {}"
    lines5 = code5.split('\n')
    result5 = is_example_code(code5, 1, lines=lines5, check_single_line=0)
    status5 = "✓" if result5 is True else "✗"
    print(f"  {status5} TestHelper (starts with Test): {result5} (expected: True)")

    # Test case 6: Code block checking (multiple lines)
    code6 = "// Example code\n// This is an example\ntype ExampleStruct struct {}"
    lines6 = code6.split('\n')
    result6 = is_example_code(code6, 1, lines=lines6)  # Check entire block
    status6 = "✓" if result6 is True else "✗"
    print(f"  {status6} Code block check (example in comments): {result6} (expected: True)")

    print()


def test_is_example_signature_name():
    """Test example signature name detection."""
    print("Testing is_example_signature_name:")
    test_cases = [
        ("ExampleType", True),
        ("HypotheticalStruct", True),
        ("MockService", True),
        ("TestHelper", True),
        ("Package", False),
        ("FileEntry", False),
        ("Example", True),  # Exact match
        ("ExamplePackage", True),  # Starts with Example
    ]
    for name, expected in test_cases:
        result = is_example_signature_name(name)
        status = "✓" if result == expected else "✗"
        print(f"  {status} {name}: {result} (expected: {expected})")
    print()


def test_package_normalization_edge_cases():
    """Test edge cases for package name normalization."""
    print("Testing package normalization edge cases:")
    test_cases = [
        # Standard library should be preserved
        ("func Read(r io.Reader) error", "func Read(io.Reader) error"),
        ("func Write(ctx context.Context, data []byte) (int, error)",
         "func Write(context.Context, []byte) (int, error)"),
        ("func Format(fmt string, args ...interface{}) string",
         "func Format(string, ...interface{}) string"),
        # Internal packages should be normalized
        ("func Process(data metadata.PackageInfo) error",
         "func Process(PackageInfo) error"),
        # Mixed standard and internal
        ("func Process(ctx context.Context, info metadata.PackageInfo) error",
         "func Process(context.Context, PackageInfo) error"),
        # Nested package paths (should still normalize)
        ("func Use(util internal.helper.Utility) error",
         "func Use(Utility) error"),
        # Package names that look like types (lowercase package, uppercase type)
        ("func Get() custompackage.CustomType",
         "func Get() CustomType"),
    ]
    for input_sig, expected in test_cases:
        result = normalize_go_signature(input_sig)
        result_clean = ' '.join(result.split())
        expected_clean = ' '.join(expected.split())
        status = "✓" if result_clean == expected_clean else "✗"
        print(f"  {status} {input_sig[:55]}...")
        if result_clean != expected_clean:
            print(f"    Got: {result}")
            print(f"    Expected: {expected}")
    print()


def test_is_signature_only_code_block():
    """Test is_signature_only_code_block and related helpers."""
    print("Testing is_signature_only_code_block / is_definition_start_line:")
    # Struct with doc comments and fields (like SymlinkConvertOptions)
    struct_block = '''// SymlinkConvertOptions configures path-to-symlink conversion.
type SymlinkConvertOptions struct {
    // PrimaryPath explicitly specifies which path should be the primary
    PrimaryPath Option[string]
    // PrimaryPathSelector chooses which path becomes the primary path
    PrimaryPathSelector Option[func(paths []string) string]
    ValidateTargetExists Option[bool]
}
'''
    result = is_signature_only_code_block(struct_block)
    status = "✓" if result else "✗"
    print(f"  {status} struct with comments and fields -> {result} (expected True)")
    # Single func signature
    func_block = (
        "func (p *Package) ConvertPathsToSymlinks(ctx context.Context, entry *FileEntry, "
        "options *SymlinkConvertOptions) (*FileEntry, []SymlinkEntry, error)"
    )
    result = is_signature_only_code_block(func_block)
    status = "✓" if result else "✗"
    print(f"  {status} single func signature -> {result} (expected True)")
    # Mix of definition and non-continuation (e.g. standalone statement)
    mixed_block = '''type X struct {}
var y = 1
'''
    result = is_signature_only_code_block(mixed_block)
    status = "✓" if not result else "✗"
    print(f"  {status} type + var -> {result} (expected False)")
    # Definition start line detection
    assert is_definition_start_line("type SymlinkConvertOptions struct {")
    assert is_definition_start_line("func (p *Package) Foo() error")
    assert not is_definition_start_line("    PrimaryPath Option[string]")
    assert is_continuation_line("    PrimaryPath Option[string]")
    assert is_continuation_line("}")
    print("  ✓ is_definition_start_line / is_continuation_line edge cases")
    print()


if __name__ == '__main__':
    print("=" * 70)
    print("Testing _go_code_utils normalization logic")
    print("=" * 70)
    print()

    test_normalize_generic_name()
    test_package_qualified_types()
    test_extract_receiver_type()
    test_full_signature_normalization()
    test_parse_go_def_signature()
    test_find_go_code_blocks()
    test_is_in_go_code_block()
    test_is_example_code()
    test_is_example_signature_name()
    test_package_normalization_edge_cases()
    test_is_signature_only_code_block()

    print("=" * 70)
    print("Tests complete")
    print("=" * 70)
