//go:build !freebsd

package main

import "github.com/shirou/gopsutil/v3/process"

func environ(process *process.Process) ([]string, error) {
	env, err := process.Environ()
	if err != nil {
		return nil, err
	}

	return env, nil
}
