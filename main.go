package main

import (
	"fmt"
	"github.com/cfsalguero/supervisorg/config"
	"log"
	"os"
)

func main() {
	configFiles := []string{"/etc/supervisorg/supervisor.conf", "supervisor.conf"}
	f, err := chooseConfiFile(configFiles)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	conf, err := config.NewConfig(f)
	for name, program := range conf.Programs {
		fmt.Printf("Name: %s, Config: %+v\n", name, program)
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
