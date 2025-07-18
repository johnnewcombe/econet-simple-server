version := 0.0.1

build:
	# This will fail if the version already exists, this is by
	# design and prevents overwriting an existing version
	#mkdir ./releases/$(version)

	env GOOS=linux GOARCH=amd64  go build -o EconetSimpleServer main.go
	env GOOS=darwin GOARCH=arm64  go build -o ./releases/$(version)/EconetSFS-macos-arm64 main.go
	env GOOS=windows GOARCH=amd64 go build -o ./releases/$(version)/EconetSFS-windows-amd64.exe main.go
	env GOOS=linux GOARCH=amd64   go build -o ./releases/$(version)/EconetSFS-linux-amd64 main.go
	env GOOS=linux GOARCH=arm64   go build -o ./releases/$(version)/EconetSFS-linux-arm64 main.go
