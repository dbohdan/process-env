# process-env

**process-env** prints select environment variables of a process, typically the current user's desktop session, as shell commands to set those variables or as JSON.
This is useful if you want to access an [ssh-agent](https://en.wikipedia.org/wiki/Ssh-agent) started for your desktop session in a remote session and in similar scenarios.

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
- FreeBSD 14
- Linux
- Windows 10

It is known to **not work** on:
- macOS
- NetBSD
- OpenBSD

## Usage

**process-env** [_options_] (_pid_|_process-name_) [_var-name_ ...]

### Options:

- **-f**, **--fish**&thinsp;&mdash;&thinsp;output fish shell commands
- **-j**, **--json**&thinsp;&mdash;&thinsp;output JSON
- **-p**, **--posix**&thinsp;&mdash;&thinsp;output POSIX shell commands (default)

## Examples

Get the default environment variables from a [MATE](https://en.wikipedia.org/wiki/MATE_(desktop_environment)) session process using POSIX syntax:

```shell
process-env mate-session
```

Get only the `DISPLAY` environment variable from a [KDE Plasma](https://en.wikipedia.org/wiki/KDE_Plasma) shell process using fish syntax:

```shell
process-env -f plasmashell DISPLAY
```

## License

MIT.
See the file [`LICENSE`](LICENSE).
