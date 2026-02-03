"""Models for Go specification reference validation."""

import re
from typing import Optional


class SpecReference:
    """Represents a specification reference from a Go file."""

    def __init__(self, file_path, line_num: int, raw_ref: str):
        self.file_path = file_path
        self.line_num = line_num
        self.raw_ref = raw_ref
        self.spec_file: Optional[str] = None
        self.section: Optional[str] = None
        self.heading: Optional[str] = None
        self.is_valid_format = False
        self.function_name: Optional[str] = None
        self.suggested_ref: Optional[str] = None
        self._parse()

    def _parse(self):
        """Parse the raw reference into components."""
        ref = self.raw_ref.strip()
        ref = re.sub(r'^\.\./', '', ref)
        ref = re.sub(r'^\.\.\\', '', ref)
        if '..' in ref:
            return

        pattern = r'^([a-zA-Z0-9_\-]+\.md):\s+(\d+(?:\.\d+)*)\.?\s+(.+)$'
        match = re.match(pattern, ref)
        if match:
            self.is_valid_format = True
            self.spec_file = match.group(1)
            self.section = match.group(2)
            self.heading = match.group(3).strip()
        else:
            anchor_pattern = r'^([a-zA-Z0-9_\-]+\.md)#([^#\s]+)$'
            anchor_match = re.match(anchor_pattern, ref)
            if anchor_match:
                self.spec_file = anchor_match.group(1)
                anchor = anchor_match.group(2)
                section_match = re.match(r'^(\d+)(?:-|$)', anchor)
                if section_match:
                    digits_str = section_match.group(1)
                    if len(digits_str) == 1:
                        self.section = digits_str
                    elif len(digits_str) == 2:
                        self.section = f"{digits_str[0]}.{digits_str[1]}"
                    elif len(digits_str) == 3:
                        self.section = f"{digits_str[0]}.{digits_str[1]}.{digits_str[2]}"
                    elif len(digits_str) == 4:
                        self.section = (
                            f"{digits_str[0]}.{digits_str[1]}."
                            f"{digits_str[2]}.{digits_str[3]}"
                        )
                return
            section_pattern = r'^([a-zA-Z0-9_\-]+\.md)\s+Section\s+(\d+(?:\.\d+)*)(?:\s*-\s*(.+))?$'
            section_match = re.match(section_pattern, ref)
            if section_match:
                self.spec_file = section_match.group(1)
                self.section = section_match.group(2)
                if section_match.group(3):
                    self.heading = section_match.group(3).strip()
                return
            if ':' in ref:
                parts = ref.split(':', 1)
                self.spec_file = parts[0].strip()
            elif ref.endswith('.md'):
                self.spec_file = ref.strip()
            else:
                md_match = re.search(r'([a-zA-Z0-9_\-]+\.md)', ref)
                if md_match:
                    self.spec_file = md_match.group(1)

    def __repr__(self):
        if self.is_valid_format:
            return f"{self.spec_file}: {self.section} {self.heading}"
        return self.raw_ref
