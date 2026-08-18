# GOOS=windows GOARCH=amd64 go build -o output_name.exe main.go
.PHONY: build front build-windows gen run

gen:
	sqlc generate

front:
	cd frontend && npm install
	cd frontend && npm run build
	cp -a frontend/build/. internal/web/build/

build: gen front
	go build -o bin/myapp ./cmd/app

run: build
	./bin/myapp

# Windows GUI subsystem prevents a console window from appearing when the
# packaged application is launched by the user.
build-windows: gen front
	GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o bin/card-platform.exe ./cmd/app
