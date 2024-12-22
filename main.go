package main

import (
	"bytes"
	"encoding/json"
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
	jsonObject
	posixShell
)

var defaultEnvVarNames = []string{"DBUS_SESSION_BUS_ADDRESS", "DISPLAY", "SSH_AUTH_SOCK"}

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
	jsonObj := false
	processName := ""

	pflag.BoolVarP(&fish, "fish", "f", false, "output fish shell commands")
	pflag.BoolVarP(&jsonObj, "json", "j", false, "output JSON")
	pflag.BoolVarP(&posix, "posix", "p", false, "output POSIX shell commands (default)")

	defaultEnvVarList := "  - " + strings.Join(defaultEnvVarNames, "\n  - ") + "\n"

	pflag.Usage = func() {
		message := fmt.Sprintf(
			"Usage: %s [options] process-name [var-name ...]\n\n"+
				"Print select environment variables of a process, "+
				"typically the current user's desktop session, "+
				"as shell commands to set those variables or as JSON.\n\n"+
				"Default variables:\n%s\nOptions:",
			filepath.Base(os.Args[0]),
			defaultEnvVarList,
		)

		fmt.Fprintln(os.Stderr, wrapForTerm(message))
		pflag.PrintDefaults()
	}

	pflag.Parse()

	flags := 0
	if fish {
		flags++
	}
	if jsonObj {
		flags++
	}
	if posix {
		flags++
	}
	if flags > 1 {
		fmt.Fprintln(os.Stderr, "can only specify one output format")
		os.Exit(2)
	}

	args := pflag.Args()

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "process name is required")
		os.Exit(2)
	}

	processName = args[0]
	envVarNames := args[1:]
	if len(envVarNames) == 0 {
		envVarNames = defaultEnvVarNames
	}
	format := posixShell
	if fish {
		format = fishShell
	} else if jsonObj {
		format = jsonObject
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

	outputMap := make(map[string]string)
	for _, name := range envVarNames {
		if value, ok := envMap[name]; ok {
			switch format {
			case fishShell:
				fmt.Printf("set -x %s %s\n", fishQuote(name), fishQuote(value))

			case jsonObject:
				outputMap[name] = value

			case posixShell:
				fmt.Printf("export %s=%s\n", posixQuote(name), posixQuote(value))
			default:
				fmt.Fprintf(os.Stderr, "unknown format")
				os.Exit(1)
			}
		}
	}

	if format == jsonObject {
		b, err := json.Marshal(outputMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to marshal to JSON: %v", err)
			os.Exit(1)
		}

		var out bytes.Buffer
		err = json.Indent(&out, b, "", "    ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to ident JSON: %v", err)
			os.Exit(1)
		}
		_, err = out.WriteTo(os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to stdout: %v", err)
			os.Exit(1)
		}

		fmt.Println()
	}
}
