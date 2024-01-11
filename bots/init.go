package bots

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	cron "github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/logger"
	"github.com/zhfreal/E5SubBot/storage"
)

func read_config(file_path string) (*config.ConfigYaml, error) {
	var config_yaml config.ConfigYaml
	if len(file_path) == 0 {
		file_path = "./"
	}
	info, err := os.Stat(file_path)
	if err != nil {
		err := fmt.Errorf("can't access %v", file_path)
		return nil, err
	}
	// check the mode
	if info.IsDir() {
		viper.AddConfigPath(file_path)
		viper.SetConfigName("config")
	} else {
		viper.SetConfigFile(file_path)
	}
	// read config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	// unmarshal using viper
	if err := viper.Unmarshal(&config_yaml); err != nil {
		return nil, err
	}

	// print deprecated warning info
	if len(config_yaml.Socks5) > 0 {
		fmt.Println("<Init> WARNING: socks5 is deprecated, use \"Proxy: socks5://127.0.0.1:1080\" or \"Proxy: http://127.0.0.1:8080\" or \"Proxy: https://example.com:8443\" instead")
		if len(config_yaml.Proxy) == 0 {
			config_yaml.Proxy = config_yaml.Socks5
			if !strings.HasPrefix(config_yaml.Proxy, "socks5://") {
				config_yaml.Proxy = "socks5://" + config_yaml.Proxy
			}
		}
	}

	// proxy
	ProxyObj, err = config.NewProxyValid(config_yaml.Proxy)
	if err != nil {
		return nil, err
	}

	// check db settings
	switch config_yaml.DB.DBType {
	case "mysql":
		if config_yaml.DB.Mysql == nil {
			err = fmt.Errorf("<Init> FATAL: mysql settings is empty, please set \"DB.mysql\"")
			return nil, err
		}
		if config_yaml.DB.Mysql.Host == "" ||
			config_yaml.DB.Mysql.Port > 65535 ||
			config_yaml.DB.Mysql.Port < 1 ||
			config_yaml.DB.Mysql.User == "" ||
			config_yaml.DB.Mysql.Password == "" ||
			config_yaml.DB.Mysql.Database == "" {
			err = fmt.Errorf("<Init> FATAL: mysql settings is invalid, please check \"DB.mysql\"")
			return nil, err
		}
	case "sqlite":
		// detect sqlite.db db file
		if config_yaml.DB.Sqlite == nil {
			err = fmt.Errorf("<Init> FATAL: sqlite settings is empty, please set \"DB.sqlite\"")
			return nil, err
		}
		if config_yaml.DB.Sqlite.DBFile == "" {
			err = fmt.Errorf("<Init> FATAL: sqlite settings is invalid, please check \"DB.sqlite\"")
			return nil, err
		}
	}

	return &config_yaml, nil
}

func monitor_config_change(file_path string) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Warnf("Config file changed:\n", e.Name)
		new_config, err := read_config(file_path)
		if err != nil {
			logger.Warnf("failed to reload config error, failed with: %v\n", err)
			return
		}
		// bot_token changed, warning to restart daemon to take effect
		if new_config.BotToken != ConfigYamlObj.BotToken {
			logger.Warnf("bot_token changed, please restart daemon to take effect\n")
		}
		// proxy changed, replace ConfigYamlInstance.Proxy to new, and re-create ProxyObj
		if new_config.Proxy != ConfigYamlObj.Proxy {
			ConfigYamlObj.Proxy = new_config.Proxy
			ProxyObj, err = config.NewProxyValid(ConfigYamlObj.Proxy)
			if err != nil {
				logger.Warnf("proxy section is invalid, please do double check\n")
			}
		}
		// Goroutine, BindMax, ErrLimit, Notice changed, just copy to ConfigYamlInstance
		ConfigYamlObj.BindMax = new_config.BindMax
		ConfigYamlObj.Goroutine = new_config.Goroutine
		ConfigYamlObj.ErrLimit = new_config.ErrLimit
		ConfigYamlObj.Notice = new_config.Notice
		// re-create AdminSet, if admin list changed
		if len(new_config.Admin) > 0 {
			AdminSet = config.NewAdminList(getAdmins())
		}
		// cron changed, need cronjob to stop to change the cron setting
		if new_config.CronConf.Enabled != ConfigYamlObj.CronConf.Enabled ||
			new_config.CronConf.Task != ConfigYamlObj.CronConf.Task ||
			new_config.CronConf.Notice != ConfigYamlObj.CronConf.Notice {
			ConfigYamlObj.CronConf.Enabled = new_config.CronConf.Enabled
			ConfigYamlObj.CronConf.Task = new_config.CronConf.Task
			ConfigYamlObj.CronConf.Notice = new_config.CronConf.Notice
			//  we need wait for the original cronjob to finish
			if CronObj != nil {
				ctx := CronObj.Stop()
				// cron done
				done_ch := ctx.Done()
			THIS_LOOP:
				for {
					select {
					case <-done_ch:
						break THIS_LOOP
					default:
						time.Sleep(time.Millisecond * 50)
					}
				}
			}
			init_background_tasks(ConfigYamlObj.CronConf)
		}
		// DB changed, need to restart daemon to take effect
		if new_config.DB.DBType != ConfigYamlObj.DB.DBType {
			logger.Warnf("DB changed, please restart daemon to take effect\n")
		} else if new_config.DB.DBType == "mysql" {
			if new_config.DB.Mysql.Host != ConfigYamlObj.DB.Mysql.Host ||
				new_config.DB.Mysql.Port != ConfigYamlObj.DB.Mysql.Port ||
				new_config.DB.Mysql.User != ConfigYamlObj.DB.Mysql.User ||
				new_config.DB.Mysql.Password != ConfigYamlObj.DB.Mysql.Password ||
				new_config.DB.Mysql.Database != ConfigYamlObj.DB.Mysql.Database ||
				new_config.DB.Mysql.TLS != ConfigYamlObj.DB.Mysql.TLS {
				logger.Warnf("DB changed, please restart daemon to take effect\n")
			}
		} else if new_config.DB.DBType == "sqlite" {
			if new_config.DB.Sqlite.DBFile != ConfigYamlObj.DB.Sqlite.DBFile {
				logger.Warnf("DB changed, please restart daemon to take effect\n")
			}
		}
		// if log settings changed, need to restart daemon to take effect
		if new_config.Log.LogIntoFile != ConfigYamlObj.Log.LogIntoFile ||
			new_config.Log.LogFile != ConfigYamlObj.Log.LogFile ||
			new_config.Log.LogLevel != ConfigYamlObj.Log.LogLevel ||
			new_config.Log.MaxSize != ConfigYamlObj.Log.MaxSize ||
			new_config.Log.MaxBackups != ConfigYamlObj.Log.MaxBackups ||
			new_config.Log.MaxAge != ConfigYamlObj.Log.MaxAge ||
			new_config.Workspace != ConfigYamlObj.Workspace {
			logger.Warnf("log settings changed, please restart daemon to take effect\n")
		}
		// if Workspace changed, need to restart daemon to take effect
		if new_config.Workspace != ConfigYamlObj.Workspace {
			logger.Warnf("Workspace changed, please restart daemon to take effect\n")
		}
	})
	viper.WatchConfig()
}

func getAdmins() []int64 {
	var result []int64
	for _, v := range ConfigYamlObj.Admin {
		id, _ := strconv.ParseInt(v, 10, 64)
		result = append(result, id)
	}

	return result
}

// initialize background cron tasks
// this must be called after botTelegram initialized
func init_background_tasks(cron_conf *config.ConfigCron) {
	if cron_conf.Enabled {
		CronObj := cron.New()
		CronObj.AddFunc(cron_conf.Task, PerformTasks)
		CronObj.AddFunc(cron_conf.Notice, NotifyStats)
		CronObj.Start()
	}
}

func Init(conf string) {
	// read config
	var err error
	ConfigYamlObj, err = read_config(conf)
	if err != nil {
		fmt.Printf("read_config error: %v\n", err)
		os.Exit(1)
	}
	// init logger
	logger.Init(ConfigYamlObj.Log.LogIntoFile,
		ConfigYamlObj.Workspace,
		ConfigYamlObj.Log.LogFile,
		ConfigYamlObj.Log.LogLevel,
		ConfigYamlObj.Log.MaxSize,
		ConfigYamlObj.Log.MaxBackups,
		ConfigYamlObj.Log.MaxAge)
	// storage init must be done after logger init, because storage.Init() would using logger
	storage.Init(ConfigYamlObj.Workspace, ConfigYamlObj.DB)
	// do cache initialization
	AuthCachedObj = NewAuthCache()
	BindCachedObj = NewBindCache()
	UsersConfigCacheObj = NewUsersConfigCache()
	// AdminList
	AdminSet = config.NewAdminList(getAdmins())
}
