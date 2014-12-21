package config

import (
	"fmt"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ConfigTest struct{}

var _ = Suite(&ConfigTest{})

func (s *ConfigTest) TestValidConfig(c *C) {
	configFile := readFile("../tests/config1.conf")
	r := strings.NewReader(configFile)
	config, err := NewConfig(r)
	c.Assert(err, Equals, nil)

	program, ok := config.Programs["program_name"]

	c.Assert(ok, Equals, true)
	c.Assert(program.Name, Equals, "program_name")
	c.Assert(program.Command, Equals, "/bin/ls")
	c.Assert(program.Directory, Equals, "/etc")
	c.Assert(program.Priority, Equals, int64(999))
}

func (s *ConfigTest) TestInvalidEntry(c *C) {
	configFile := readFile("../tests/config2.conf")
	r := strings.NewReader(configFile)
	config, err := NewConfig(r)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)

}
func readFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot open %s", filename))
	}
	return string(data)
}
