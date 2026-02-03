"""
Type and helper-function keyword scoring rules.
"""

from __future__ import annotations

from typing import List, Tuple

from lib.go_defs_index._go_defs_index_scoring_rules_core import ScoringContext


def score_other_types_suffix(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if "other types" not in ctx.section_lower:
        return 0.0, []
    if ctx.name_lower in {
        "addfileoptions",
        "extractpathoptions",
        "removedirectoryoptions",
        "fileinfo",
        "filemetadataupdate",
        "fileindex",
        "indexentry",
        "createoptions",
        "packageconfig",
        "pathhandling",
        "destpathspec",
        "symlinkconvertoptions",
        "transformpipeline",
        "transformstage",
        "transformtype",
        "tag",
        "tagvaluetype",
        "recoveryfileheader",
    }:
        return 0.0, []
    suffixes = [
        "options",
        "config",
        "info",
        "entry",
        "type",
        "spec",
        "handling",
        "pipeline",
        "index",
        "header",
        "rule",
        "strategy",
        "builder",
        "worker",
        "pool",
    ]
    if any(ctx.name_lower.endswith(suffix) for suffix in suffixes):
        return 0.55, ["Type suffix matches Other Types section: +55%"]
    return 0.0, []


def score_generic_type_keywords(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if "generic types" not in ctx.section_lower:
        return 0.0, []
    if ctx.name_lower in {
        "addfileoptions",
        "extractpathoptions",
        "removedirectoryoptions",
        "fileinfo",
        "filemetadataupdate",
        "createoptions",
        "packageconfig",
        "pathhandling",
        "destpathspec",
        "symlinkconvertoptions",
        "fileindex",
        "indexentry",
        "transformpipeline",
        "transformstage",
        "transformtype",
        "tag",
        "tagvaluetype",
        "recoveryfileheader",
    }:
        return 0.0, []
    keywords = [
        "config",
        "builder",
        "option",
        "optional",
        "result",
        "strategy",
        "validator",
        "worker",
        "rule",
        "pool",
        "job",
        "thread",
        "pathentry",
    ]
    if any(keyword in ctx.name_lower for keyword in keywords):
        return 0.30, ["Generic type keyword matches Generic Types: +30%"]
    return 0.0, []


def score_metadata_type_keywords(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if (
        "package metadata types" not in ctx.section_lower
        and "metadata types" not in ctx.section_lower
    ):
        return 0.0, []
    if ctx.name_lower in {
        "pathentry",
        "addfileoptions",
        "extractpathoptions",
        "removedirectoryoptions",
        "fileinfo",
        "filemetadataupdate",
        "tag",
        "tagvaluetype",
        "transformpipeline",
        "transformstage",
        "transformtype",
    }:
        return 0.0, []
    keywords = [
        "metadata",
        "manifest",
        "index",
        "signaturedata",
        "signatureinfo",
        "packageconfig",
        "pathhandling",
        "createoptions",
        "destpathspec",
        "symlinkconvertoptions",
        "fileindex",
        "indexentry",
        "pathmetadata",
        "pathinfo",
        "pathstats",
        "pathnode",
        "pathtree",
        "pathfilesystem",
        "pathinheritance",
        "pathmetadatapatch",
        "pathmetadatatype",
        "pathmetadataentry",
    ]
    if any(keyword in ctx.name_lower for keyword in keywords):
        return 0.30, ["Metadata type keyword matches Package Metadata Types: +30%"]
    return 0.0, []


def score_signature_type_keywords(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if "signature types" not in ctx.section_lower:
        return 0.0, []
    if ctx.name_lower in {"signaturedata", "signatureinfo"}:
        return 0.0, []
    if ctx.name_lower in {"unsupportederrorcontext", "validationerrorcontext"}:
        return 0.50, ["Signature error context matches Signature Types: +50%"]
    if ctx.name_lower == "signingkey":
        return 0.30, ["Signature-adjacent type matches Signature Types: +30%"]
    if "signature" in ctx.name_lower:
        return 0.30, ["Signature type keyword matches Signature Types: +30%"]
    return 0.0, []


def score_error_type_keywords(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if "error types" not in ctx.section_lower:
        return 0.0, []
    allowed = {
        "errortype",
        "packageerror",
        "packageerrorcontext",
        "ioerrorcontext",
        "patternerrorcontext",
        "readonlyerrorcontext",
    }
    if ctx.name_lower in allowed:
        return 0.30, ["Error type keyword matches Error Types: +30%"]
    return 0.0, []


def score_file_entry_type_keywords(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if "fileentry types" not in ctx.section_lower:
        return 0.0, []
    keywords = [
        "fileentry",
        "filesource",
        "hash",
        "optionaldata",
        "processingstate",
        "addfileoptions",
        "extractpathoptions",
        "removedirectoryoptions",
        "fileinfo",
        "filemetadataupdate",
        "tag",
        "tagvaluetype",
        "transformpipeline",
        "transformstage",
        "transformtype",
    ]
    if any(keyword in ctx.name_lower for keyword in keywords):
        return 0.30, ["FileEntry type keyword matches FileEntry Types: +30%"]
    return 0.0, []


def score_file_info_preference(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if ctx.name_lower != "fileinfo":
        return 0.0, []
    if "fileentry types" in ctx.section_lower:
        return 0.40, ["FileInfo prefers FileEntry Types: +40%"]
    if "package interface types" in ctx.section_lower:
        return -0.40, ["FileInfo avoids Package Interface Types: -40%"]
    return 0.0, []


def score_recovery_file_header_preference(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if ctx.name_lower != "recoveryfileheader":
        return 0.0, []
    if "package interface types" in ctx.section_lower:
        return 0.40, ["RecoveryFileHeader prefers Package Interface Types: +40%"]
    return 0.0, []


def score_security_error_context_types(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if ctx.name_lower not in {"securityerrorcontext", "encryptionerrorcontext"}:
        return 0.0, []
    if "encryption and security types" in ctx.section_lower:
        return 0.20, ["Security error context matches Security Types: +20%"]
    return 0.0, []


def score_generic_helper_functions(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if "generic helper functions" not in ctx.section_lower:
        return 0.0, []
    helper_names = {
        "err",
        "ok",
        "processconcurrently",
        "composevalidators",
        "validateall",
        "validatewith",
    }
    if ctx.definition.file == "api_generics.md" or ctx.name_lower in helper_names:
        return 0.50, ["Generic helper function matches Generic Helpers: +50%"]
    return 0.0, []


def score_other_type_helper_functions(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if "other type helper functions" not in ctx.section_lower:
        return 0.0, []
    if ctx.name_lower == "newpackagewithoptions":
        return 0.0, []
    if ctx.name_lower.startswith("new") and "options" in ctx.name_lower:
        return 0.60, ["Options constructor matches Other Type Helpers: +60%"]
    return 0.0, []


def score_generic_core_type_preference(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    generic_core = {
        "config",
        "configbuilder",
        "strategy",
        "validationrule",
        "validator",
        "workerpool",
        "pathentry",
    }
    if ctx.name_lower not in generic_core:
        return 0.0, []
    if "generic types" in ctx.section_lower:
        return 0.30, ["Generic core type prefers Generic Types: +30%"]
    if "other types" in ctx.section_lower:
        return -0.30, ["Generic core type avoids Other Types: -30%"]
    return 0.0, []


def score_package_error_constructor(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if ctx.name_lower != "newpackageerror":
        return 0.0, []
    if "error helper functions" in ctx.section_lower:
        return 0.60, ["NewPackageError prefers Error Helper Functions: +60%"]
    if "package helper functions" in ctx.section_lower:
        return -0.60, ["NewPackageError avoids Package Helper Functions: -60%"]
    return 0.0, []


def score_signature_comment_helpers(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if (
        "package metadata helper functions" not in ctx.section_lower
        and "metadata helper functions" not in ctx.section_lower
    ):
        return 0.0, []
    if "signaturecomment" in ctx.name_lower:
        return 0.20, ["Signature comment helper matches Metadata Helpers: +20%"]
    return 0.0, []


def score_package_open_helpers(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if "package helper functions" not in ctx.section_lower:
        return 0.0, []
    if ctx.name_lower.startswith("open"):
        return 0.15, ["Open helper prefers Package Helper Functions: +15%"]
    return 0.0, []


def score_package_read_header_helpers(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if "package helper functions" not in ctx.section_lower:
        return 0.0, []
    if ctx.name_lower == "readheader":
        return 0.20, ["ReadHeader prefers Package Helper Functions: +20%"]
    return 0.0, []


def score_metadata_header_helpers(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if (
        "package metadata helper functions" not in ctx.section_lower
        and "metadata helper functions" not in ctx.section_lower
    ):
        return 0.0, []
    if ctx.name_lower == "readheaderfrompath":
        return 0.40, ["ReadHeaderFromPath matches Metadata Helpers: +40%"]
    return 0.0, []


def score_metadata_destpath_helpers(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if (
        "package metadata helper functions" not in ctx.section_lower
        and "metadata helper functions" not in ctx.section_lower
    ):
        return 0.0, []
    if ctx.name_lower == "setdestpath":
        return 0.40, ["SetDestPath matches Metadata Helpers: +40%"]
    return 0.0, []


def score_package_helper_overrides(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    if "package helper functions" not in ctx.section_lower:
        return 0.0, []
    if ctx.name_lower == "newpackagewithoptions":
        return 0.40, ["NewPackageWithOptions prefers Package Helper Functions: +40%"]
    if ctx.name_lower == "readheaderfrompath":
        return 0.20, ["ReadHeaderFromPath prefers Package Helper Functions: +20%"]
    return 0.0, []


def score_create_options_preference(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if ctx.name_lower != "createoptions":
        return 0.0, []
    if "package metadata types" in ctx.section_lower:
        return 0.30, ["CreateOptions prefers Package Metadata Types: +30%"]
    if "other types" in ctx.section_lower:
        return -0.30, ["CreateOptions avoids Other Types: -30%"]
    if "generic types" in ctx.section_lower:
        return -0.30, ["CreateOptions avoids Generic Types: -30%"]
    return 0.0, []
