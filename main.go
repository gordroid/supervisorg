package main

import (
	"flag"
	"fmt"
	"github.com/cfsalguero/supervisorg/config"
	"log"
	"os"
	"strings"
)

func main() {
	var configFile string
	configFiles := []string{"/etc/supervisorg/supervisor.conf", "supervisor.conf"}

	flag.StringVar(&configFile, "configuration", "", "Default: "+strings.Join(configFiles, ", "))
	flag.StringVar(&configFile, "c", "", "Default: "+strings.Join(configFiles, ", "))
	flag.Parse()

	if configFile != "" {
		configFiles = []string{configFile}
	}

	f, err := chooseConfiFile(configFiles)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	conf, err := config.NewConfig(f)
	for name, program := range conf.Programs {
		cmd, err := program.GetCmd()
		if err == nil {
			fmt.Printf("Name: %s, Config: %+v\n", name, cmd)
		}
	}
}

func chooseConfiFile(files []string) (*os.File, error) {
	for _, file := range files {
		f, err := os.Open(file)
		if err == nil {
			return f, nil
		}
	}
	return nil, fmt.Errorf("Cannot find config file in: %v", files)
}
