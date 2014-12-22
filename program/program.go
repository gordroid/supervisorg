package program

import (
	"github.com/mattn/go-shellwords"
	"os/exec"
	"syscall"
)

type Program struct {
	Name       string
	Command    string
	Directory  string
	Priority   int64
	Instances  int
	StopSignal syscall.Signal
	//
	running  bool
	stopChan chan bool
}

func (p *Program) IsValid() bool {
	if p.Name != "" && p.Command != "" {
		return true
	}
	return false
}

func (p *Program) GetCmd() (cmd *exec.Cmd, err error) {
	parser := shellwords.NewParser()
	parser.ParseEnv = true
	args, err := parser.Parse(p.Command)
	if err != nil {
		return nil, err
	}
	cmd = exec.Command(args[0], args[1:]...)
	if p.Directory != "" {
		cmd.Dir = p.Directory
	}
	return cmd, nil
}

func NewProgram(name string) *Program {
	p := &Program{
		Name:       name,
		stopChan:   make(chan bool),
		Instances:  1,
		StopSignal: syscall.SYS_KILL,
	}
	return p
}

func (p *Program) Run() {

}

func (p *Program) Stop() {

}
