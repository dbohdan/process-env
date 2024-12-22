# process-env

**process-env** prints select environment variables of a process, typically the current user's
desktop session, as shell commands to set those variables or as JSON.
This is useful, e.g., for accessing an [ssh-agent](https://en.wikipedia.org/wiki/Ssh-agent) started for your desktop session in remote sessions.

The default environment variables are as follows:

- `DBUS_SESSION_BUS_ADDRESS`
- `DISPLAY`
- `SSH_AUTH_SOCK`

## Installation

You will need Go 1.22.

```shell
go install github.com/dbohdan.com/process-env@master
```

## Compatibility

process-env is known to **work** on:
- Linux
- Windows 10

process-env (getting the environment variables of a process using [gopsutils](https://github.com/shirou/gopsutil)) is known to **not work** on:
- FreeBSD
- macOS
- NetBSD
- OpenBSD

## Usage

```none
Usage: process-env [options] process-name [var-name ...]
```

### Options:

- `-f`, `--fish`&thinsp;&mdash;&thinsp;output fish shell commands
- `-j`, `--json`&thinsp;&mdash;&thinsp;output JSON
- `-p`, `--posix`&thinsp;&mdash;&thinsp;output POSIX shell commands (default)

## Examples

Get the default environment variables from a [MATE](https://en.wikipedia.org/wiki/MATE_(desktop_environment)) session process using POSIX syntax:

```shell
process-env mate-session
```

Get just the `DISPLAY` environment variable from a [KDE Plasma](https://en.wikipedia.org/wiki/KDE_Plasma) shell process using fish syntax:

```shell
process-env -f plasmashell DISPLAY
```

## License

MIT.
See the file [`LICENSE`](LICENSE).
