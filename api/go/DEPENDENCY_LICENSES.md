# Go Dependencies License Report

This document lists all dependencies (direct and transitive) for the Go v1 implementation and their licenses.

Generated: 2026-01-11

## Direct Dependencies

| Package                     | Version  | License | Status        |
| --------------------------- | -------- | ------- | ------------- |
| `github.com/cucumber/godog` | v0.15.1  | MIT     | ✅ Compatible |
| `github.com/goccy/go-yaml`  | v1.15.13 | MIT     | ✅ Compatible |
| `github.com/samber/lo`      | v1.52.0  | MIT     | ✅ Compatible |

## Transitive Dependencies

| Package                                   | Version                            | License                 | Status           |
| ----------------------------------------- | ---------------------------------- | ----------------------- | ---------------- |
| `github.com/cpuguy83/go-md2man/v2`        | v2.0.2                             | MIT                     | ✅ Compatible    |
| `github.com/cucumber/gherkin/go/v26`      | v26.2.0                            | MIT                     | ✅ Compatible    |
| `github.com/cucumber/messages/go/v21`     | v21.0.1                            | MIT                     | ✅ Compatible    |
| `github.com/cucumber/messages/go/v22`     | v22.0.0                            | MIT                     | ✅ Compatible    |
| `github.com/davecgh/go-spew`              | v1.1.1                             | ISC                     | ✅ Compatible    |
| `github.com/gofrs/uuid`                   | v4.4.0+incompatible                | MIT                     | ✅ Compatible    |
| `github.com/hashicorp/go-immutable-radix` | v1.3.1                             | MPL-2.0                 | ✅ Compatible    |
| `github.com/hashicorp/go-memdb`           | v1.3.5                             | MPL-2.0                 | ✅ Compatible    |
| `github.com/hashicorp/go-uuid`            | v1.0.2                             | MPL-2.0                 | ✅ Compatible    |
| `github.com/hashicorp/golang-lru`         | v1.0.2                             | MPL-2.0                 | ✅ Compatible    |
| `github.com/inconshreveable/mousetrap`    | v1.1.0                             | Apache-2.0              | ✅ Compatible    |
| `github.com/kr/pretty`                    | v0.2.1                             | MIT                     | ✅ Compatible    |
| `github.com/kr/pty`                       | v1.1.1                             | MIT                     | ✅ Compatible    |
| `github.com/kr/text`                      | v0.1.0                             | MIT                     | ✅ Compatible    |
| `github.com/pmezard/go-difflib`           | v1.0.0                             | BSD-3-Clause            | ✅ Compatible    |
| `github.com/russross/blackfriday/v2`      | v2.1.0                             | BSD-2-Clause            | ✅ Compatible    |
| `github.com/spf13/cobra`                  | v1.7.0                             | Apache-2.0              | ✅ Compatible    |
| `github.com/spf13/pflag`                  | v1.0.10                            | BSD-style               | ✅ Compatible    |
| `github.com/stretchr/objx`                | v0.5.0                             | MIT                     | ✅ Compatible    |
| `github.com/stretchr/testify`             | v1.8.2                             | MIT                     | ✅ Compatible    |
| `golang.org/x/mod`                        | v0.31.0                            | BSD-3-Clause            | ✅ Compatible    |
| `golang.org/x/sync`                       | v0.19.0                            | BSD-3-Clause            | ✅ Compatible    |
| `golang.org/x/text`                       | v0.33.0                            | BSD-3-Clause            | ✅ Compatible    |
| `golang.org/x/tools`                      | v0.40.0                            | BSD-3-Clause            | ✅ Compatible    |
| `gopkg.in/check.v1`                       | v1.0.0-20201130134442-10cb98267c6c | BSD-3-Clause            | ✅ Compatible    |
| `gopkg.in/yaml.v3`                        | v3.0.1                             | MIT / Apache-2.0 (dual) | ✅ Compatible    |

## License Summary

- **MIT**: 15 packages
- **BSD-3-Clause**: 6 packages
- **MPL-2.0**: 4 packages (HashiCorp packages)
- **Apache-2.0**: 2 packages
- **BSD-2-Clause**: 1 package
- **ISC**: 1 package
- **BSD-style**: 1 package

## Notes

### MPL-2.0 Dependencies

The following HashiCorp packages use MPL-2.0 (Mozilla Public License 2.0):

- `github.com/hashicorp/go-immutable-radix`
- `github.com/hashicorp/go-memdb`
- `github.com/hashicorp/go-uuid`
- `github.com/hashicorp/golang-lru`

**MPL-2.0 Compatibility**: MPL-2.0 is a weak copyleft license (file-level) that is compatible with Apache 2.0.
MPL-2.0 files must remain under MPL-2.0, but can be combined with other licensed code in a larger work.
Since these are transitive dependencies that are not modified, they can be used as-is.

### All Other Dependencies

All other dependencies use permissive licenses (MIT, BSD, Apache-2.0, ISC) that are fully compatible with the project's license requirements.

## Verification

To regenerate this list:

```bash
cd api/go
go list -m all
```

Then verify licenses using:

- `go-licenses` tool (if installed)
- Manual verification via GitHub repositories
- Module cache inspection
