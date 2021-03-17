# Development
## Useful Commands
Build for local:
```zsh
go build ./cmd/raspilive
```

Build for Raspberry Pi:
```zsh
env GOOS=linux GOARCH=arm GOARM=6 go build ./cmd/raspilive
```

Run tests:
```zsh
go test ./...
go test ./... -v
go test ./... -cover
```

Copy binary from localhost to the Raspberry Pi:
```zsh
scp raspilive pi@raspberrypi:/home/pi
```