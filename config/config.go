package config

import (
	"bufio"
	"fmt"
	"github.com/cfsalguero/supervisorg/program"
	"io"
	"strconv"
	"strings"
)

type Config struct {
	Programs map[string]*program.Program
}

func (c *Config) addProgram(p *program.Program) {
	if p != nil && p.IsValid() {
		c.Programs[p.Name] = p
	}
}

func NewConfig(f io.Reader) (*Config, error) {
	programSectionPrefix := "[program:"
	config := &Config{
		Programs: make(map[string]*program.Program),
	}

	r := bufio.NewReader(f)
	inProgram := false
	var prg *program.Program = nil
	lineNum := 0

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if inProgram {
				config.addProgram(prg)
			}
			break
		}

		lineNum += 1
		line = strings.TrimSpace(line)

		if line == "" {
			if inProgram {
				config.addProgram(prg)
			}
			inProgram = false
			continue
		}

		if strings.HasPrefix(line, programSectionPrefix) {
			name := strings.TrimSuffix(strings.TrimPrefix(line, programSectionPrefix), "]")
			inProgram = true
			prg = &program.Program{
				Name:     name,
				Priority: 999,
			}
			continue
		}

		if inProgram {
			e := parseProgramLine(line, prg)
			if e != nil {
				return nil, e
			}
		}
	}
	return config, nil
}

func parseProgramLine(line string, program *program.Program) error {
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
