all: build

build: client server
	cd client && go build
	cd server && go build

clean:
	rm server/server client/client
