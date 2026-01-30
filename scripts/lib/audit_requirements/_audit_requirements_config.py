"""
Constants for the requirements coverage audit.

Used by audit_requirements_scan and audit_requirements_classify.
"""

import re

# Maximum lines of prose for a heading to be considered organizational
MAX_ORGANIZATIONAL_PROSE_LINES = 5

# Minimum thresholds for meaningful prose
MIN_MEANINGFUL_SENTENCES = 1
MIN_PROSE_LINES = 2

# Exceptions to catch when we want to log and continue (avoid bare Exception)
_SCRIPT_ERROR_EXCEPTIONS = (
    TypeError, ValueError, KeyError, IndexError, AttributeError, RuntimeError, OSError
)

# Headings matching these patterns (case-insensitive) are excluded from coverage
HEADING_EXCLUSION_PATTERNS = [
    "example usage",
    "usage examples",
    "best practices",
    "cross-reference",
    "cross-references",
    "implementation details",
    "implementation structure",
    "internal",
    "table of contents",
    "overview",
]

# Keywords that indicate functional/testable behavior (case-insensitive)
FUNCTIONAL_BEHAVIOR_KEYWORDS = [
    "behavior",
    "behaviour",
    "process",
    "operation",
    "operations",
    "method",
    "function",
    "creates",
    "validates",
    "returns",
    "error conditions",
    "error handling",
    "parameters",
    "parameter",
]

# Architectural keyword weights (case-insensitive). Format: (keyword, weight)
ARCHITECTURAL_KEYWORD_WEIGHTS = [
    ("system architecture", 3),
    ("software architecture", 3),
    ("type definition", 3),
    ("type definitions", 3),
    ("interface definition", 3),
    ("interface definitions", 3),
    ("struct definition", 3),
    ("struct definitions", 3),
    ("architectural design", 3),
    ("architectural pattern", 3),
    ("architectural patterns", 3),
    ("system design", 3),
    ("architecture design", 3),
    ("design pattern", 3),
    ("design patterns", 3),
    ("implementation structure", 3),
    ("data structure", 3),
    ("structure", 2),
    ("system", 2),
    ("interface", 2),
    ("architecture", 2),
    ("organization", 1),
    ("organisation", 1),
    ("struct", 1),
]

# Non-architectural phrase weights (negative). Format: (phrase, negative_weight)
NON_ARCHITECTURAL_PHRASE_WEIGHTS = [
    ("purpose", -2),
    ("usage notes", -2),
    ("usage", -2),
    ("management", -2),
    ("effects", -2),
    ("effect", -2),
    ("derivation", -2),
    ("validation", -2),
    ("execution", -2),
    ("policy", -2),
    ("features", -2),
    ("feature", -2),
    ("examples", -2),
    ("example", -2),
    ("helpers", -2),
    ("helper", -2),
    ("use cases", -2),
    ("use case", -2),
    ("specification", -2),
    ("computing", -2),
    ("options", -1),
    ("option", -1),
    ("details", -1),
    ("file", -1),
]

ARCHITECTURAL_SCORE_THRESHOLD = 2

# Content analysis scoring weights
FUNCTION_SIGNATURE_WEIGHT = 2
TYPE_DEFINITION_WEIGHT = 2
EXAMPLE_CODE_WEIGHT = -1

# Compiled regex patterns
FUNCTIONAL_KEYWORD_PATTERNS = [
    re.compile(rf'\b{re.escape(kw)}\b', re.IGNORECASE)
    for kw in FUNCTIONAL_BEHAVIOR_KEYWORDS
]
ARCHITECTURAL_KEYWORD_PATTERNS = [
    (re.compile(rf'\b{re.escape(kw)}\b', re.IGNORECASE), weight)
    for kw, weight in ARCHITECTURAL_KEYWORD_WEIGHTS
]
NON_ARCHITECTURAL_PHRASE_PATTERNS = [
    (re.compile(rf'\b{re.escape(phrase)}\b', re.IGNORECASE), weight)
    for phrase, weight in NON_ARCHITECTURAL_PHRASE_WEIGHTS
]
IMPLEMENTATION_PATTERN = re.compile(r'\bimplementation\b', re.IGNORECASE)
LINK_ONLY_PATTERN = re.compile(r'^\s*(?:[-*+]|\d+\.)?\s*\[[^\]]+\]\([^)]+\)\s*$')
LABEL_ONLY_PATTERN = re.compile(r'^\s*(?:\*\*[^*]+\*\*|[A-Za-z0-9][A-Za-z0-9 \-/]+):\s*$')
URL_ONLY_PATTERN = re.compile(r'^\s*(https?://\S+|www\.\S+|mailto:\S+)\s*$', re.IGNORECASE)
