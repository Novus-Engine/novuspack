# NovusPack Technical Specifications - File Types System

---

## 0. Overview

This document defines the comprehensive file type system, detection algorithms, and special file handling for the NovusPack system.

### 0.1 Cross-References

- [Main Index](_main.md) - Central navigation for all NovusPack specifications
- [Testing Requirements](testing.md) - Comprehensive testing requirements and validation
- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Package File Format](package_file_format.md) - .npk format and file entry structure
- [Metadata System](metadata.md) - Package metadata and tags system

---

## 1. File Type System Specification

The NovusPack file type system provides comprehensive semantic categorization of files within packages using a range-based architecture. This enables specialized handling, validation, and processing based on file content and purpose.

### 1.1 File Type Range Architecture

| Range           | Category          | Type IDs      | Description                             | Special Handling                        |
| --------------- | ----------------- | ------------- | --------------------------------------- | --------------------------------------- |
| **0-999**       | **Binary Files**  | 0-999         | Binary executables, libraries, archives | Security scanning, execution validation |
| **1000-1999**   | **Text Files**    | 1000-1999     | Text documents, markup, data files      | Text processing, encoding validation    |
| **2000-3999**   | **Script Files**  | 2000-3999     | Programming scripts and code            | Syntax validation, security analysis    |
| **4000-4999**   | **Config Files**  | 4000-4999     | Configuration and settings files        | Schema validation, config parsing       |
| **5000-6999**   | **Image Files**   | 5000-6999     | Image and graphics files                | Format validation, image processing     |
| **7000-7999**   | **Audio Files**   | 7000-7999     | Audio and sound files                   | Format validation, audio processing     |
| **8000-9999**   | **Video Files**   | 8000-9999     | Video and multimedia files              | Format validation, video processing     |
| **10000-10999** | **System Files**  | 10000-10999   | System files and directories            | System validation, path handling        |
| **11000-64999** | **Reserved**      | 11000-64999   | Reserved for future use                 | Reserved                                |
| **65000-65535** | **Special Files** | 65000-65535   | Special package files and metadata      | Special processing, reserved handling   |

**Note:** This is the authoritative definition of file type ranges. All other references to file types should link to this document.

### 1.2 Special File Naming Strategy

Special package files use a systematic naming convention to ensure uniqueness:

- **Prefix**: `__NPK_` - Clearly identifies NovusPack special files
- **Type Code**: `META`, `MAN`, `IDX`, `SIG` - Abbreviated type identifier
- **Type ID**: `240`, `241`, `242`, `243` - Numeric type identifier
- **Suffix**: `__` - Delimiter for consistency
- **Extension**: `.npk*` - Unique extension for each type

#### 1.3.1 Unique Extensions

- **`.npkmeta`**: Package metadata files (YAML content)
- **`.npkman`**: Package manifest files (YAML content)
- **`.npkidx`**: Package index files (YAML content)
- **`.npksig`**: Digital signature files (binary content)

## 2. Range-Based Category Queries

### 2.1 Category Checking Functions

```go
// IsBinaryFile returns true if file type is within binary file range (0-999)
func IsBinaryFile(fileType FileType) bool

// IsTextFile returns true if file type is within text file range (1000-1999)
func IsTextFile(fileType FileType) bool

// IsScriptFile returns true if file type is within script file range (2000-3999)
func IsScriptFile(fileType FileType) bool

// IsConfigFile returns true if file type is within config file range (4000-4999)
func IsConfigFile(fileType FileType) bool

// IsImageFile returns true if file type is within image file range (5000-6999)
func IsImageFile(fileType FileType) bool

// IsAudioFile returns true if file type is within audio file range (7000-7999)
func IsAudioFile(fileType FileType) bool

// IsVideoFile returns true if file type is within video file range (8000-9999)
func IsVideoFile(fileType FileType) bool

// IsSystemFile returns true if file type is within system file range (10000-10999)
func IsSystemFile(fileType FileType) bool

// IsSpecialFile returns true if file type is within special file range (65000-65535)
func IsSpecialFile(fileType FileType) bool
```

### 2.2 Compression Integration

#### 2.2.1 SelectCompressionType

```go
    //  SelectCompressionType selects the appropriate compression algorithm based on file type
    //  Skip compression for already compressed formats: Returns CompressionNone for JPEG, PNG, GIF, MP3, MP4, OGG, FLAC files
    //  Special file handling: Uses IsSpecialFile() to check for special file types
    //    - FileTypeSignature: Never compress signature files (returns CompressionNone)
    //    - FileTypeMetadata, FileTypeManifest, FileTypeIndex: Always compress YAML special files (returns CompressionZstd)
    //    - Other special files: Default compression (returns CompressionZstd)
    //  Text-based files: Returns CompressionZstd for text, script, and config files (good compression for text)
    //  Binary media files: Returns CompressionLZ4 for image, audio, and video files (fast compression for binary data)
    //  Default: Returns CompressionZstd as default compression method
func SelectCompressionType(data []byte, fileType FileType) uint8
```

## 3. File Type API

### 3.1 File Type Management

#### 3.1.1 FileType Definition

```go
    //  FileType represents a file type identifier
    //  Note: This is the authoritative definition. All other references should link to this document.
type FileType uint16
```

##### 3.1.1.1 File Type Range Constants

```go
    //  File type range constants (2-byte: 0-65535)
const (
    //  Binary Files (0-999)
    FileTypeBinaryStart = 0
    FileTypeBinaryEnd   = 999

    //  Text Files (1000-1999)
    FileTypeTextStart = 1000
    FileTypeTextEnd   = 1999

    //  Script Files (2000-3999)
    FileTypeScriptStart = 2000
    FileTypeScriptEnd   = 3999

    //  Config Files (4000-4999)
    FileTypeConfigStart = 4000
    FileTypeConfigEnd   = 4999

    //  Image Files (5000-6999)
    FileTypeImageStart = 5000
    FileTypeImageEnd   = 6999

    //  Audio Files (7000-7999)
    FileTypeAudioStart = 7000
    FileTypeAudioEnd   = 7999

    //  Video Files (8000-9999)
    FileTypeVideoStart = 8000
    FileTypeVideoEnd   = 9999

    //  System Files (10000-10999)
    FileTypeSystemStart = 10000
    FileTypeSystemEnd   = 10999

    //  Special Files (65000-65535)
    FileTypeSpecialStart = 65000
    FileTypeSpecialEnd   = 65535
)
```

##### 3.1.1.2 Specific File Type Constants

###### 3.1.1.2.1 Binary File Types (0-999)

```go
const (
    FileTypeBinary    FileType = 0 // Generic binary data
    FileTypeExecutable FileType = 1 // Binary executables
    FileTypeLibrary   FileType = 2 // Shared libraries
    FileTypeArchive   FileType = 3 // Compressed archives
)
```

###### 3.1.1.2.2 Text File Types (1000-1999)

```go
const (
    FileTypeText     FileType = 1000 // Plain text files
    FileTypeMarkdown FileType = 1001 // Markdown documents
    FileTypeHTML     FileType = 1002 // HTML documents
    FileTypeCSS      FileType = 1003 // CSS stylesheets
    FileTypeCSV      FileType = 1004 // Comma-separated values
    FileTypeSQL      FileType = 1005 // SQL scripts
    FileTypeDiff     FileType = 1006 // Diff/patch files
    FileTypeTeX      FileType = 1007 // TeX/LaTeX documents
)
```

###### 3.1.1.2.3 Script File Types (2000-3999)

```go
const (
    FileTypeScript       FileType = 2000 // Generic scripts
    FileTypePython       FileType = 2001 // Python scripts
    FileTypeJavaScript   FileType = 2002 // JavaScript code
    FileTypeTypeScript   FileType = 2003 // TypeScript code
    FileTypeShell        FileType = 2004 // Shell scripts
    FileTypeLua          FileType = 2005 // Lua scripts
    FileTypePerl         FileType = 2006 // Perl scripts
    FileTypeRuby         FileType = 2007 // Ruby scripts
    FileTypePHP          FileType = 2008 // PHP scripts
    FileTypeJava         FileType = 2009 // Java source
    FileTypeCSharp       FileType = 2010 // C# source
    FileTypeGo           FileType = 2011 // Go source
    FileTypeRust         FileType = 2012 // Rust source
    FileTypeScala        FileType = 2013 // Scala source
    FileTypeKotlin       FileType = 2014 // Kotlin source
    FileTypeSwift        FileType = 2015 // Swift source
    FileTypeCoffeeScript FileType = 2016 // CoffeeScript
    FileTypeDart         FileType = 2017 // Dart scripts
    FileTypeElixir       FileType = 2018 // Elixir scripts
    FileTypeErlang       FileType = 2019 // Erlang scripts
    FileTypeHaskell      FileType = 2020 // Haskell scripts
    FileTypeClojure      FileType = 2021 // Clojure scripts
    FileTypeFSharp       FileType = 2022 // F# scripts
    FileTypeOCaml        FileType = 2023 // OCaml scripts
    FileTypeR            FileType = 2024 // R scripts
    FileTypeMATLAB       FileType = 2025 // MATLAB scripts
    FileTypeJulia        FileType = 2026 // Julia scripts
    FileTypePowerShell   FileType = 2027 // PowerShell scripts
    FileTypeBatch        FileType = 2028 // Windows batch files
    FileTypeVBScript     FileType = 2029 // VBScript
    FileTypeAppleScript  FileType = 2030 // AppleScript
    FileTypeAutoHotkey   FileType = 2031 // AutoHotkey scripts
    FileTypeGroovy       FileType = 2032 // Groovy scripts
    FileTypeCrystal      FileType = 2033 // Crystal scripts
    FileTypeNim          FileType = 2034 // Nim scripts
    FileTypeZig          FileType = 2035 // Zig scripts
    FileTypeV            FileType = 2036 // V scripts
    FileTypeD            FileType = 2037 // D scripts
    FileTypeAda          FileType = 2038 // Ada scripts
    FileTypeFortran      FileType = 2039 // Fortran scripts
)
```

###### 3.1.1.2.4 Config File Types (4000-4999)

```go
const (
    FileTypeYAML       FileType = 4000 // YAML configuration
    FileTypeJSON       FileType = 4001 // JSON configuration
    FileTypeXML        FileType = 4002 // XML configuration
    FileTypeTOML       FileType = 4003 // TOML configuration
    FileTypeHOCON      FileType = 4004 // HOCON configuration
    FileTypeEDN        FileType = 4005 // EDN configuration
    FileTypeCUE        FileType = 4006 // CUE configuration
    FileTypeProperties FileType = 4007 // Properties files
    FileTypeINI        FileType = 4008 // INI configuration
)
```

###### 3.1.1.2.5 Image File Types (5000-6999)

```go
const (
    FileTypeImage FileType = 5000 // Generic image
    FileTypePNG   FileType = 5001 // PNG images
    FileTypeJPEG  FileType = 5002 // JPEG images
    FileTypeGIF   FileType = 5003 // GIF images
    FileTypeBMP   FileType = 5004 // BMP images
    FileTypeWebP  FileType = 5005 // WebP images
    FileTypeTIFF  FileType = 5006 // TIFF images
    FileTypeSVG   FileType = 5007 // SVG images
    FileTypeRAW   FileType = 5008 // RAW images
    FileTypeHEIC  FileType = 5009 // HEIC images
    FileTypeAVIF  FileType = 5010 // AVIF images
    FileTypeICO   FileType = 5011 // ICO images
    FileTypeTGA   FileType = 5012 // TGA images
    FileTypePSD   FileType = 5013 // PSD images
    FileTypeXCF   FileType = 5014 // XCF images
    FileTypeDDS   FileType = 5015 // DDS images
    FileTypeEXR   FileType = 5016 // EXR images
    FileTypeHDR   FileType = 5017 // HDR images
    FileTypePPM   FileType = 5018 // PPM images
    FileTypePGM   FileType = 5019 // PGM images
    FileTypePBM   FileType = 5020 // PBM images
    FileTypeXBM   FileType = 5021 // XBM images
    FileTypeXPM   FileType = 5022 // XPM images
    FileTypePCX   FileType = 5023 // PCX images
    FileTypeTIF   FileType = 5024 // TIF images
    FileTypeJPG   FileType = 5025 // JPG images
    FileTypeJP2   FileType = 5026 // JP2 images
    FileTypeJ2K   FileType = 5027 // J2K images
    FileTypeJXR   FileType = 5028 // JXR images
    FileTypeWDP   FileType = 5029 // WDP images
)
```

###### 3.1.1.2.6 Audio File Types (7000-7999)

```go
const (
    FileTypeAudio FileType = 7000 // Generic audio
    FileTypeMP3   FileType = 7001 // MP3 audio
    FileTypeWAV   FileType = 7002 // WAV audio
    FileTypeOGG   FileType = 7003 // OGG audio
    FileTypeFLAC  FileType = 7004 // FLAC audio
    FileTypeAAC   FileType = 7005 // AAC audio
    FileTypeWMA   FileType = 7006 // WMA audio
    FileTypeAIFF  FileType = 7007 // AIFF audio
    FileTypeALAC  FileType = 7008 // ALAC audio
    FileTypeAPE   FileType = 7009 // APE audio
    FileTypeOpus  FileType = 7010 // Opus audio
    FileTypeM4A   FileType = 7011 // M4A audio
    FileTypeOGV   FileType = 7012 // OGV audio
    FileTypeVorbis FileType = 7013 // Vorbis audio
    FileTypeSpeex FileType = 7014 // Speex audio
    FileTypeAMR   FileType = 7015 // AMR audio
    FileType3GP   FileType = 7016 // 3GP audio
    FileTypeAC3   FileType = 7017 // AC3 audio
    FileTypeDTS   FileType = 7018 // DTS audio
    FileTypeMKA   FileType = 7019 // MKA audio
)
```

###### 3.1.1.2.7 Video File Types (8000-9999)

```go
const (
    FileTypeVideo FileType = 8000 // Generic video
    FileTypeMP4   FileType = 8001 // MP4 video
    FileTypeMKV   FileType = 8002 // MKV video
    FileTypeAVI   FileType = 8003 // AVI video
    FileTypeMOV   FileType = 8004 // MOV video
    FileTypeWebM  FileType = 8005 // WebM video
    FileTypeWMV   FileType = 8006 // WMV video
    FileTypeFLV   FileType = 8007 // FLV video
    FileTypeM4V   FileType = 8008 // M4V video
    FileTypeMPEG  FileType = 8009 // MPEG video
    FileType3GP   FileType = 8010 // 3GP video
    FileTypeHEVC  FileType = 8011 // HEVC video
    FileTypeAV1   FileType = 8012 // AV1 video
    FileTypeDivX  FileType = 8013 // DivX video
    FileTypeXvid  FileType = 8014 // Xvid video
    FileTypeVP8   FileType = 8015 // VP8 video
    FileTypeVP9   FileType = 8016 // VP9 video
    FileTypeH264  FileType = 8017 // H.264 video
    FileTypeH265  FileType = 8018 // H.265 video
    FileTypeOGV   FileType = 8019 // OGV video
    FileTypeASF   FileType = 8020 // ASF video
    FileTypeRM    FileType = 8021 // RM video
    FileTypeRMVB  FileType = 8022 // RMVB video
    FileTypeVOB   FileType = 8023 // VOB video
    FileTypeTS    FileType = 8024 // TS video
    FileTypeM2TS  FileType = 8025 // M2TS video
    FileTypeMTS   FileType = 8026 // MTS video
    FileTypeM2V   FileType = 8027 // M2V video
    FileTypeM1V   FileType = 8028 // M1V video
    FileTypeMPG   FileType = 8029 // MPG video
)
```

###### 3.1.1.2.8 System File Types (10000-10999)

```go
const (
    FileTypeRegular   FileType = 10000 // Regular files
    FileTypeDirectory FileType = 10001 // Directories
    FileTypeSymlink   FileType = 10002 // Symbolic links
)
```

###### 3.1.1.2.9 Special File Types (65000-65535)

```go
const (
    FileTypeMetadata  FileType = 65000 // Package metadata
    FileTypeManifest  FileType = 65001 // Package manifest
    FileTypeIndex     FileType = 65002 // Package index
    FileTypeSignature FileType = 65003 // Package signature
)
```

### 3.2 File Type Detection Functions

```go
    //  DetermineFileType detects file type from name and content
func DetermineFileType(name string, data []byte) FileType

    //  SelectCompressionType selects compression based on file type
func SelectCompressionType(data []byte, fileType FileType) uint8
```

## 4. File Type Detection Algorithm

The `DetermineFileType` function uses a sophisticated multi-stage detection process:

1. **Extension-Based Detection**: First checks file extensions for specific formats (OGG, FLAC, ZIP)
2. **Content-Based Detection**: Uses the `mimetype` library to analyze file content
3. **MIME Type Mapping**: Maps MIME types to specific file type constants
4. **Extension Fallback**: Falls back to extension-based detection if content detection fails
5. **Text Analysis**: Performs text file analysis for unknown files
6. **Default Classification**: Defaults to binary for unrecognized files

### 4.1 Detection Process

#### 4.1.1 DetermineFileType

```go
    //  DetermineFileType uses a sophisticated multi-stage detection process to identify file types
func DetermineFileType(name string, data []byte) FileType
```

- **Stage 1**: Extension-based detection for specific formats
  - Checks file extensions: .ogg => FileTypeOGG, .flac => FileTypeFLAC, .zip => FileTypeArchive
- **Stage 2**: Content-based detection using mimetype library
  - Uses mimetype.Detect() to analyze file content
  - Maps MIME types to specific file type constants

##### 4.1.1.1 MIME Type Mapping

- **Image MIME types**: Maps "image/\*" MIME types to specific image file types
- **Audio MIME types**: Maps "audio/\*" MIME types to specific audio file types
- **Video MIME types**: Maps "video/\*" MIME types to specific video file types
- **Text MIME types**: Maps "text/\*" MIME types to specific text file types
- **Additional mappings**: Supports other MIME type categories

##### 4.1.1.2 Extension Fallback Mapping

- **Text files**: ".txt", ".text" => FileTypeText
- **YAML files**: ".yaml", ".yml" => FileTypeYAML
- **Lua files**: ".lua" => FileTypeLua
- **Config files**: ".ini", ".cfg" => FileTypeINI
- **JavaScript files**: ".js" => FileTypeJavaScript
- **Additional extensions**: Supports other file extensions

##### 4.1.1.3 Text File Analysis

- **Data validation**: Checks if data length is greater than 0
- **Text detection**: Analyzes each byte to determine if content is text
- **Control character handling**: Allows newlines, carriage returns, and tabs
- **Printable character check**: Ensures content contains printable ASCII characters
- **Result**: Returns FileTypeText if content is valid text, otherwise continues to binary detection

##### 4.1.1.4 Default Classification

- **Binary fallback**: Returns FileTypeBinary if no other classification matches
