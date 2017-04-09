all: build

format:
	gofmt -w .

build: format
	cd client && go build
	cd server && go build

clean:
	rm server/server client/client
