.PHONY: build front gen run

gen:
	sqlc generate

front:
	cd frontend && pnpm install
	cd frontend && pnpm build
	rm -rf internal/web/dist
	cp -r frontend/dist internal/web/dist

build: gen front
	go build -o bin/myapp ./cmd/app

run: build
	./bin/myapp
