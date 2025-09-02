version := 0.0.1

build: test
	# This will fail if the version already exists, this is by
	# design and prevents overwriting an existing version
	#mkdir ./releases/$(version)


	env GOOS=darwin GOARCH=amd64  go build -o ./releases/$(version)/PiconetSFS-macos-amd64 ./src/main.go
	env GOOS=darwin GOARCH=arm64  go build -o ./releases/$(version)/PiconetSFS-macos-arm64 ./src/main.go
	env GOOS=windows GOARCH=amd64 go build -o ./releases/$(version)/PiconetSFS-windows-amd64.exe ./src/main.go
	env GOOS=linux GOARCH=amd64   go build -o ./releases/$(version)/PiconetSFS-linux-amd64 ./src/main.go
	env GOOS=linux GOARCH=arm64   go build -o ./releases/$(version)/PiconetSFS-linux-arm64 ./src/main.go

test:
	go test ./...

install: test

	###################################################################################
	# SORT THIS OUT BEFORE REPO IS PUBLIC
	###################################################################################
	cp releases/0.0.1/PiconetSFS-linux-amd64 /opt/piconetSFS/piconetSFS
	#scp releases/0.0.1/PiconetSFS-linux-amd64 john@192.168.1.200:~/piconetSFS
	scp releases/0.0.1/PiconetSFS-linux-amd64 john@s1.glasstty.local:~/piconetSFS
