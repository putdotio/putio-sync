all:

build-web:
	@cd web;./node_modules/gulp/bin/gulp.js;cd ..
	@esc -o http/static.go -pkg http -prefix web/build web/build

build-all:
	@mkdir build/
	@GOOS=linux GOARCH=386 go build -o build/putio-sync.linux-386
	@GOOS=linux GOARCH=amd64 go build -o build/putio-sync.linux-amd64
	@GOOS=linux GOARCH=arm go build -o build/putio-sync.linux-arm
	@GOOS=darwin GOARCH=amd64 go build -o build/putio-sync.darwin-amd64
	@GOOS=windows GOARCH=386 go build -o build/putio-sync.windows-386
	@GOOS=windows GOARCH=amd64 go build -o build/putio-sync.windows-amd64

clean:
	@rm -rf build/

.PHONY: all build clean
