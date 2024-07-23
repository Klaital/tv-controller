
.PHONY: build
build: 
	@go generate ./...
	@go build -tags raspberrypi ./cmd/server

run-vlc:
	@"C:\Program Files\VideoLAN\VLC\vlc.exe" --http-host=127.0.0.1 --http-port=8090 --extraintf=http --http-password=bedroomtv123