build:
	go build -o result/bin/spotify-search main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/spotify-search-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/spotify-search-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/spotify-search-freebsd-386 main.go