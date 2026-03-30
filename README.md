putio-sync
==========

Command-line program to sync a folder between put.io and your computer.

**WARNING: The project is still in development and features are subject to change!**

If you are using MacOS or Windows, you can install desktop version: [putio-sync-desktop](https://github.com/putdotio/putio-sync-desktop).

Installing
----------

If you are on macOS you can install with [Homebrew](https://brew.sh/):
```sh
brew install putdotio/tap/putio-sync
```

> [!NOTE]
> If you previously installed via the old tap (`putdotio/putio-sync`), you'll need to uninstall first and reinstall with the new formula name:
> ```sh
> brew uninstall putio-sync
> brew untap putdotio/putio-sync
> brew install putdotio/tap/putio-sync
> ```

Otherwise, get the latest binary from [releases page](https://github.com/putdotio/putio-sync/releases).

Usage
-----

Run the program with your account credentials:
```sh
PUTIO_Username=<username> PUTIO_Password=<password> putio-sync
```

Then program is going to sync the contents of these folders:
- **$HOME/putio-sync** in your computer
- **/putio-sync** in your Put.io account

The folders are created if they don't exist.
Files will be synced periodically or when a change has been detected.
