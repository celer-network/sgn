// Copyright 2018 Celer Network

package testing

import (
	"os"
	"os/exec"

	"github.com/celer-network/sgn/testing/log"
)

// Killable is object that has Kill() func
type Killable interface {
	Kill() error
}

func TearDown(tokill []Killable) {
	log.Info("Tear down Killables ing...")
	for _, p := range tokill {
		p.Kill()
	}
}

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
