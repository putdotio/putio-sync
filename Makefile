all:

build-web:
	@cd web;./node_modules/gulp/bin/gulp.js;cd ..
	@esc -o http/static.go -pkg http -prefix web/build web/build
	@go build -o putio-sync-cli

build-all:
	@GOOS=linux GOARCH=386 go build -o putio-sync.linux-386
	@GOOS=linux GOARCH=amd64 go build -o putio-sync.linux-amd64
	@GOOS=linux GOARCH=arm go build -o putio-sync.linux-arm
	@GOOS=darwin GOARCH=amd64 go build -o putio-sync.darwin-amd64
	@GOOS=windows GOARCH=386 go build -o putio-sync.windows-386
	@GOOS=windows GOARCH=amd64 go build -o putio-sync.windows-amd64

clean:
	@rm putio-sync-cli

.PHONY: all build clean
