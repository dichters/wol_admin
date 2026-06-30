module github.com/mydev/mydev-service

go 1.25.6

require (
	github.com/gofiber/fiber/v2 v2.50.0
	github.com/google/uuid v1.6.0
	github.com/mydev/mydev-api v0.0.0
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.50.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/mydev/mydev-api => ../mydev-api
