putio-sync
==========

Command-line program to sync a folder between put.io and your computer.

**WARNING: The project is still in development and features are subject to change!**

[![Build Status](https://travis-ci.org/putdotio/putio-sync.svg?branch=v2)](https://travis-ci.org/putdotio/putio-sync)
[![GitHub Release](https://img.shields.io/github/release/putdotio/putio-sync.svg)](https://github.com/putdotio/putio-sync/releases)

Installing
----------

If you are on MacOS you can install from [brew](https://brew.sh/):
```sh
brew tap putdotio/putio-sync
brew install putio-sync
```

Otherwise, get the latest binary from [releases page](https://github.com/putdotio/putio-sync/releases).

Usage
-----

Run the program with your account credentials:
```sh
putio-sync -username <username> -password <password>
```

Then program is going to sync the contents of these folders:
- **$HOME/putio-sync** folder in your computer
- **/putio-sync** folder in your Put.io account

The folders are created if they don't exist.
The program exists after sync.
You can run it again whenever you want to sync.
