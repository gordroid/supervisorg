package config

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-shellwords"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type Program struct {
	Name      string
	Command   string
	Directory string
	Priority  int64
}

func (p *Program) isValid() bool {
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

type Config struct {
	Programs map[string]*Program
}

func (c *Config) addProgram(p *Program) {
	if p != nil && p.isValid() {
		c.Programs[p.Name] = p
	}
}

func NewConfig(f io.Reader) (*Config, error) {
	programSectionPrefix := "[program:"
	config := &Config{
		Programs: make(map[string]*Program),
	}

	r := bufio.NewReader(f)
	inProgram := false
	var program *Program = nil
	lineNum := 0

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if inProgram {
				config.addProgram(program)
			}
			break
		}

		lineNum += 1
		line = strings.TrimSpace(line)

		if line == "" {
			if inProgram {
				config.addProgram(program)
			}
			inProgram = false
			continue
		}

		if strings.HasPrefix(line, programSectionPrefix) {
			name := strings.TrimSuffix(strings.TrimPrefix(line, programSectionPrefix), "]")
			inProgram = true
			program = &Program{
				Name:     name,
				Priority: 999,
			}
			continue
		}

		if inProgram {
			e := parseProgramLine(line, program)
			if e != nil {
				return nil, e
			}
		}
	}
	return config, nil
}

func parseProgramLine(line string, program *Program) error {
	key, value, err := getKeyVal(line)
	if err != nil {
		return err
	}

	if key == "command" {
		program.Command = value
	}
	if key == "directory" {
		program.Directory = value
	}
	if key == "priority" {
		program.Priority, err = strconv.ParseInt(value, 10, 64)
	}
	return err
}

func getKeyVal(line string) (key, value string, err error) {
	parts := strings.Split(line, "=")
	if len(parts) != 2 {
		err = fmt.Errorf("Invalid program entry: %s", line)
		return
	}

	key = strings.TrimSpace(strings.ToLower(parts[0]))
	value = strings.TrimSpace(parts[1])
	return
}
