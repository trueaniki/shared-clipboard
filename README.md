# Shared Clipboard
A CLI tool to share clipboard between devices in local area network

Start the application on multiple devices sharing the same network.
Press the combination (by default: `CTRL+SHIFT+A`) to share the clipboard.
Press the combination (by default: `CTRL+SHIFT+D`) on other device to adopt the clipboard.

## Usage
```sh
sharedclipboard -n 192.168.0.0/24
```
Will run the application in foreground mode.
### Daemon mode
```sh
sharedclipboard start -n 192.168.0.0/24
```
Will start the application in daemon mode.

```sh
sharedclipboard stop
```
Will stop the application running in daemon mode.
### Options
- `--network` or `-n` --- Network to scan for peers running shared clipboard in CIDR format.
- `--conf` or `-c` --- Path to hotkeys config file.
### Configuring hotkeys