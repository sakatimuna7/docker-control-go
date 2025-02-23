dev:
	air

build:
	go build -o dist/main

run:
	./dist/main

# seed:
seed-all:
	go run src/database/seeding/main.go

# migrate fresh:
migrate-fresh:
	go run src/database/migrations/main.go