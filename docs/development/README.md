# Development
## Useful Commands
Build:
```zsh
go build ./cmd/raspilive
```

Run tests:
```zsh
go test ./... -v
go test ./... -cover
```

Run tests with code coverage:
```zsh
go test ./... -cover
```

Copy files from localhost to the Raspberry Pi:
```zsh
scp <file_path> pi@raspberrypi:<remote_dir>
scp -r <local_dir> pi@raspberrypi:<remote_dir>
```

Copy files from the Raspberry Pi to localhost:
```zsh
scp pi@raspberrypi:/home/pi/image.jpg .
scp -r pi@raspberrypi:<remote_dir> <local_dir>
```
