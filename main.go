package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cornfeedhobo/pflag"
	tsize "github.com/kopoli/go-terminal-size"
	"github.com/mitchellh/go-wordwrap"
	"github.com/shirou/gopsutil/v3/process"
)

type outputFormat int

const (
	fishShell outputFormat = iota
	posixShell
)

var defaultEnvVars = []string{"DBUS_SESSION_BUS_ADDRESS", "DISPLAY", "SSH_AUTH_SOCK"}

func pgrep(user string, procName string) ([]*process.Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var matches []*process.Process
	for _, p := range processes {
		username, err := p.Username()
		if err != nil {
			continue
		}

		name, err := p.Name()
		if err != nil {
			continue
		}

		if username == user && name == procName {
			matches = append(matches, p)
		}
	}

	return matches, nil
}

func shellSafe(s string) bool {
	re := regexp.MustCompile("^[A-Za-z0-9%+,-./:=@_]+$")
	return re.MatchString(s)
}

func fishQuote(s string) string {
	if shellSafe(s) {
		return s
	}

	return "'" + strings.ReplaceAll(s, "'", `\'`) + "'"
}

func posixQuote(s string) string {
	if shellSafe(s) {
		return s
	}

	// Simple POSIX shell quoting: wrap in single quotes and escape single quotes.
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

func setVarCommand(name, value string, format outputFormat) (string, error) {
	switch format {
	case fishShell:
		return fmt.Sprintf("set -x %s %s", fishQuote(name), fishQuote(value)), nil

	case posixShell:
		return fmt.Sprintf("export %s=%s", posixQuote(name), posixQuote(value)), nil

	default:
		return "", fmt.Errorf("unknown format")
	}
}

func wrapForTerm(s string) string {
	size, err := tsize.GetSize()
	if err != nil {
		return s
	}

	return wordwrap.WrapString(s, uint(size.Width))
}

func main() {
	fish := false
	posix := false
	processName := ""

	pflag.BoolVarP(&fish, "fish", "f", false, "output fish shell commands")
	pflag.BoolVarP(&posix, "posix", "p", false, "output POSIX shell commands (default)")

	defaultEnvVarList := "  - " + strings.Join(defaultEnvVars, "\n  - ") + "\n"

	pflag.Usage = func() {
		message := fmt.Sprintf(
			"Usage: %s [options] process-name [var-name ...]\n\n"+
				"Print shell commands to set environment variables to those of another process, "+
				"typically the current user's desktop session.\n\n"+
				"Default variables:\n%s\nOptions:",
			filepath.Base(os.Args[0]),
			defaultEnvVarList,
		)

		fmt.Fprintln(os.Stderr, wrapForTerm(message))
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if fish && posix {
		fmt.Fprintln(os.Stderr, "can't specify both `--fish` and `--posix`")
		os.Exit(2)
	}

	args := pflag.Args()

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "process name is required")
		os.Exit(2)
	}

	processName = args[0]
	envVars := args[1:]
	if len(envVars) == 0 {
		envVars = defaultEnvVars
	}
	format := posixShell
	if fish {
		format = fishShell
	}

	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting username: %v\n", err)
		os.Exit(1)
	}

	processes, err := pgrep(user.Username, processName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding processes: %v\n", err)
		os.Exit(1)
	}

	if len(processes) == 0 {
		fmt.Fprintln(os.Stderr, "no process found")
		os.Exit(1)
	}

	if len(processes) > 1 {
		fmt.Fprintln(os.Stderr, "more than one process found")
		os.Exit(1)
	}

	process := processes[0]
	env, err := process.Environ()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting process environment: %v\n", err)
		os.Exit(1)
	}

	// Convert a slice of `KEY=VALUE` strings to a map of keys to values.
	envMap := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)

		if len(parts) != 2 {
			continue
		}

		envMap[parts[0]] = parts[1]
	}

	for _, v := range envVars {
		if val, ok := envMap[v]; ok {
			cmd, err := setVarCommand(v, val, format)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating command: %v\n", err)
				os.Exit(1)
			}

			fmt.Println(cmd)
		}
	}
}
