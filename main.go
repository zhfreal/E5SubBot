package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/zhfreal/E5SubBot/bots"
)

var (
	runable = "E5SubBot"
	version = "dev"
)

func print_version() {
	fmt.Printf("%v %v\n", runable, version)
}

func main() {
	var work_dir string
	var show_version bool
	var help = "\n" + runable + " " + version + `
    Usage: E5SubBot [options]
    options:
        -d|--conf-dir  work directory for E5SubBot.
        -v|--version   show version.
        -h|--help      show help.
    `
	flag.StringVarP(&work_dir, "conf-dir", "d", ".", "work directory for E5SubBot.")
	flag.BoolVarP(&show_version, "version", "v", false, "Show version.")
	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()
	if show_version {
		print_version()
		os.Exit(0)
	}
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
			if t_file.Name() == "config.yml" || t_file.Name() == "config.yaml" {
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
