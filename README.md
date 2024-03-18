# Shared Clipboard
A CLI tool to share clipboard between devices in local area network

Start the application on multiple devices sharing the same network.
Press the combination (by default: `CTRL+SHIFT+S`) to share the clipboard.
Press the combination (by default: `CTRL+SHIFT+A`) on other device to adopt the clipboard.

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
- `--conf` or `-c`    --- Path to hotkeys config file.
### Configuring hotkeys
To set custom hotkeys, create a textfile with contents like folowing:
```
Share=ctrl+s
Adopt=ctrl+a
```
Pass the file path to the application using `--conf` flag:
```sh
sharedclipboard start -n 192.168.0.0/24 -c ~/.sharedclipboard.conf
```

When writing hotkeys, follow the rules of [github.com/trueaniki/go-parse-hotkeys](https://github.com/trueaniki/go-parse-hotkeys). Separator is `+`.

### Init command
Use `init` command to generate config file at `~/.shared-clipboard.conf`:
```sh
sharedclipboard init
```
This will generate file with the following contents:
```
# Share=Ctrl+Shift+S
# Adopt=Ctrl+Shift+A
```
You can uncomment the lines and set hotkeys to whatever you want according to [github.com/trueaniki/go-parse-hotkeys](https://github.com/trueaniki/go-parse-hotkeys) using `+` separator.

If you have `~/.shared-clipboard.conf` file set up, no need to pass it exteranlly using `--conf` flag, it will be handled automatically.