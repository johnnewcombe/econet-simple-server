version := 0.0.1

build:
	# This will fail if the version already exists, this is by
	# design and prevents overwriting an existing version
	#mkdir ./releases/$(version)

	#env GOOS=linux GOARCH=amd64  go build -o /usr/local/bin/PiconetFS ./src/main.go
	env GOOS=darwin GOARCH=amd64  go build -o ./releases/$(version)/PiconetFS-macos-amd64 ./src/main.go
	env GOOS=darwin GOARCH=arm64  go build -o ./releases/$(version)/PiconetFS-macos-arm64 ./src/main.go
	env GOOS=windows GOARCH=amd64 go build -o ./releases/$(version)/PiconetFS-windows-amd64.exe ./src/main.go
	env GOOS=linux GOARCH=amd64   go build -o ./releases/$(version)/PiconetFS-linux-amd64 ./src/main.go
	env GOOS=linux GOARCH=arm64   go build -o ./releases/$(version)/PiconetFS-linux-arm64 ./src/main.go
