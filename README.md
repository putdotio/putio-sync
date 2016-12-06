# putio-sync

putio-sync is a command line program to download your preferred Put.io directory to your computer.

It periodically watches for changes and downloads new content when they arrive.

## features

* Synchronize upstream folder to make upstream directory hiearchy identical
* Checks for file hash equality
* Pause/resume support
* HTTP API for external clients
* Web UI

# installation

```sh
go get github.com/putdotio/putio-sync
```

# usage

run `putio-sync -server` and visit `http://127.0.0.1:3000`
