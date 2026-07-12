// Package version holds build-time metadata injected via ldflags.
package version

import "runtime"

// These variables are set at build time via -ldflags.
// When building without ldflags (e.g. during development),
// sensible defaults are used based on compile-time constants.
var (
	Version   = "0.0.1"
	Arch      = runtime.GOARCH // compile-time constant, reflects target arch
	BuildTime = "dev"           // set via ldflags in release builds
)
