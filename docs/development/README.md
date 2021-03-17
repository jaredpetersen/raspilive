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

## Debian Package & Release
Building Debian package:
```zsh
dpkg-deb --build raspilive-1.0-1_armhf.deb
```

Installing Debian package:
```zsh
sudo apt install ./raspilive-1.0-1_armhf.deb
```

Structure (`raspilive-1.0-1_armhf`)
```
.
├── DEBIAN
│   └── control
└── usr
    └── bin
        └── raspilive
```

Control file:
```
Package: raspilive
Version: 1.0-0
Section: video
Priority: optional
Architecture: armhf
Depends: ffmpeg (>= 7:4.1.6-1~deb10u1)
Maintainer: Jared Petersen <jaredtoddpetersen@gmail.com>
Homepage: https://github.com/jaredpetersen/raspilive
Description: Stream video from the Raspberry Pi Camera module to the web
```
