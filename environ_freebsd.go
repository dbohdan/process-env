//go:build freebsd

package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/shirou/gopsutil/v3/process"
)

type procstatOutput struct {
	Version  string `json:"__version"`
	Procstat struct {
		Env map[string]struct {
			ProcessID int      `json:"process_id"`
			Command   string   `json:"command"`
			Envp      []string `json:"envp"`
		} `json:"env"`
	} `json:"procstat"`
}

func environ(process *process.Process) ([]string, error) {
	if process == nil {
		return nil, fmt.Errorf("process is nil")
	}

	// Convert the PID to string for both the command and the later map lookup.
	// Do *not* use `string(process.Pid)`.
	// With an int32, this produces a Unicode character.
	pidStr := strconv.Itoa(int(process.Pid))

	cmd := exec.Command("procstat", "--libxo", "json", "penv", pidStr)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run procstat(1): %v", err)
	}

	var result procstatOutput
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse procstat(1) output: %v", err)
	}

	procEnv, ok := result.Procstat.Env[pidStr]
	if !ok {
		return nil, fmt.Errorf("no environment variables found for PID %v", pidStr)
	}

	if len(procEnv.Envp) == 0 {
		return nil, fmt.Errorf("empty array of environment variables for PID %v", pidStr)
	}

	return procEnv.Envp, nil
}
