build:
	go build -o snip.exe main.go
build-windows:
	set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=1&& go build -o snip.exe main.go
test:
	go test -v ./internal/test/...

bench:
	go test -run='^$$' -bench=. ./internal/test/...