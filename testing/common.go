// Copyright 2018 Celer Network

package testing

import (
	"os"
	"os/exec"

	"github.com/celer-network/sgn/testing/log"
)

func StartProcess(name string, args ...string) *os.Process {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Infoln(err)
	}
	return cmd.Process
}

func KillProcess(process *os.Process) {
	process.Kill()
	process.Release()
}
