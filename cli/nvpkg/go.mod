module github.com/novus-engine/novuspack/cli/nvpkg

go 1.25

replace github.com/novus-engine/novuspack => ../..

replace github.com/novus-engine/novuspack/api/go => ../../api/go

require (
	github.com/creack/pty v1.1.24
	github.com/novus-engine/novuspack/api/go v0.0.0
	github.com/peterh/liner v1.2.2
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.5.0 // indirect
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/godog v0.15.1 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.5 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/samber/lo v1.52.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
)
