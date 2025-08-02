version := 0.0.1

build:
	# This will fail if the version already exists, this is by
	# design and prevents overwriting an existing version
	#mkdir ./releases/$(version)

	# REMOVE THIS ONCE REPO IS PUBLIC
	env GOOS=linux GOARCH=amd64  go build -o /opt/piconetSFS/piconetSFS ./src/main.go

	env GOOS=darwin GOARCH=amd64  go build -o ./releases/$(version)/PiconetSFS-macos-amd64 ./src/main.go
	env GOOS=darwin GOARCH=arm64  go build -o ./releases/$(version)/PiconetSFS-macos-arm64 ./src/main.go
	env GOOS=windows GOARCH=amd64 go build -o ./releases/$(version)/PiconetSFS-windows-amd64.exe ./src/main.go
	env GOOS=linux GOARCH=amd64   go build -o ./releases/$(version)/PiconetSFS-linux-amd64 ./src/main.go
	env GOOS=linux GOARCH=arm64   go build -o ./releases/$(version)/PiconetSFS-linux-arm64 ./src/main.go

