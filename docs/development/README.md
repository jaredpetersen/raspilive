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
scp pi@raspberrypi:<file_path> .
scp -r pi@raspberrypi:<remote_dir> <local_dir>
```

Building Debian package:
```zsh
dpkg-deb --build raspilive-1.0-0_armhf.deb
```

Installing Debian package:
```zsh
sudo apt install ./raspilive-1.0-0_armhf.deb
```