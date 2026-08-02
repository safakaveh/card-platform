# GOOS=windows GOARCH=amd64 go build -o output_name.exe main.go
.PHONY: build front gen run

gen:
	sqlc generate

front:
	cd frontend && pnpm install
	cd frontend && pnpm build
	rm -rf internal/web/dist
	cp -r frontend/build internal/web/build

build: gen front
	go build -o bin/myapp ./cmd/app

run: build
	./bin/myapp
