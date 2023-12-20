# Development
## Useful Commands
Build for local:
```zsh
go build ./cmd/raspi-live-libcamera
```

Build for Raspberry Pi:
```zsh
env GOOS=linux GOARCH=arm GOARM=6 go build ./cmd/raspi-live-libcamera
```

Go install command:
```zsh
go install github.com/amd940/rraspi-live-libcamera/cmd/raspi-live-libcamera
```

Run tests:
```zsh
go test ./...
go test ./... -v
go test ./... -cover
```

Copy binary from localhost to the Raspberry Pi:
```zsh
scp raspi-live-libcamera pi@raspberrypi:/home/pi
```
