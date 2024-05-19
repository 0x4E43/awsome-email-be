init:
	@touch db/email.db

clean:
	@rm db/email.db

run:
	@go run main.go

build:
	@go build -o bin/