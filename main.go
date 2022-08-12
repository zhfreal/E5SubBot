package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/zhfreal/E5SubBot/bots"
)

func main() {
	var work_dir string
	flag.StringVarP(&work_dir, "conf-dir", "d", "", "Write result to csv file, disabled by default.")
	flag.Parse()
	if len(work_dir) == 0 {
		fmt.Println("No config directory, please use \"-d|--conf-dir\"")
		os.Exit(1)
	}
	if t_dir_list, t_err := os.ReadDir(work_dir); t_err != nil {
		fmt.Printf("Config directory \"%v\"does not exist!", work_dir)
		os.Exit(1)
	} else {
		t_config_yaml_found := false
		for _, t_file := range t_dir_list {
			if t_file.Name() == "config.yml" {
				t_config_yaml_found = true
				break
			}
		}
		if !t_config_yaml_found {
			fmt.Printf("\"config.yml\" is not in Config directory \"%v\"!", work_dir)
			os.Exit(1)
		}
	}
	bots.Start(work_dir)
}
