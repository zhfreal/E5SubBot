package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zhfreal/E5SubBot/bots"
	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/logger"
	"github.com/zhfreal/E5SubBot/storage"
)

const (
	logo = `
  ______ _____ _____       _     ____        _   
 |  ____| ____/ ____|     | |   |  _ \      | |  
 | |__  | |__| (___  _   _| |__ | |_) | ___ | |_ 
 |  __| |___ \\___ \| | | | '_ \|  _ < / _ \| __|
 | |____ ___) |___) | |_| | |_) | |_) | (_) | |_ 
 |______|____/_____/ \__,_|_.__/|____/ \___/ \__|
`
)

var (
	APPName = "E5SubBot"
	Version = "dev"
)

func print_version() {
	fmt.Printf("%v %v\n", APPName, Version)
}

func main() {
	var conf string
	var show_version bool
	var show_token bool
	var accounts_for_show string
	help := APPName + " " + Version + `
    Usage: E5SubBot [options]
    options:
        -c|-conf          config file path or folder contain "config.yaml" or "config.yml". "/etc/e5bot/" as default.
        -S|-show-token    show all bounded accounts's token.
        -a|-account       specific the account bounded for show, work with "-S|-show-token".
        -v|-version       show version.
        -h|-help          show help.
`
	flag.StringVar(&conf, "conf", "/etc/e5bot/", "config file path or folder contain \"config.yaml\" or \"config.yml\".")
	flag.StringVar(&conf, "c", "/etc/e5bot/", "config file path or folder contain \"config.yaml\" or \"config.yml\".")
	flag.BoolVar(&show_token, "show-token", false, "show all bounded accounts's token.")
	flag.BoolVar(&show_token, "S", false, "show all bounded accounts's token.")
	flag.StringVar(&accounts_for_show, "account", "", "specific the account bounded for show, work with -S|-show-token.")
	flag.StringVar(&accounts_for_show, "a", "", "specific the account bounded for show, work with -S|-show-token.")
	flag.BoolVar(&show_version, "version", false, "Show version.")
	flag.BoolVar(&show_version, "v", false, "Show version.")
	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()
	if show_version {
		print_version()
		os.Exit(0)
	}

	// do initialization
	config.Init(conf)
	logger.Init(config.LogIntoFile, config.LogDir, config.LogFile, config.LogLevel, config.MaxSize, config.MaxBackups, config.MaxAge)
	// storage init must be done after logger init, because storage.Init() would using logger
	storage.Init()
	// self Init
	bots.Init()
	if show_token {
		bots.ShowToken(accounts_for_show)
		os.Exit(0)
	}
	// this is for test only
	// PerformTasks()
	// show logo
	fmt.Print(logo)
	bots.Start()
}
