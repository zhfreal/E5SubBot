package bots

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	cron "github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/logger"
	ms "github.com/zhfreal/E5SubBot/microsoft"
	"github.com/zhfreal/E5SubBot/storage"
)

var vp *viper.Viper

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
		vp.AddConfigPath(file_path)
		vp.SetConfigName("config")
	} else {
		vp.SetConfigFile(file_path)
	}
	// read config
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	// unmarshal using vp
	if err := vp.Unmarshal(&config_yaml); err != nil {
		return nil, err
	}

	// print deprecated warning info
	if len(config_yaml.Socks5) > 0 {
		logger.Warnln("<Init-read-config> WARNING: socks5 is deprecated, use \"Proxy: socks5://127.0.0.1:1080\" or \"Proxy: http://127.0.0.1:8080\" or \"Proxy: https://example.com:8443\" instead")
		if len(config_yaml.Proxy) == 0 {
			config_yaml.Proxy = config_yaml.Socks5
			if !strings.HasPrefix(config_yaml.Proxy, "socks5://") {
				config_yaml.Proxy = "socks5://" + config_yaml.Proxy
			}
		}
	}

	// proxy
	proxyObj, err = config.NewProxyValid(config_yaml.Proxy)
	if err != nil {
		return nil, err
	}

	// check db settings
	switch config_yaml.DB.DBType {
	case "mysql":
		if config_yaml.DB.Mysql == nil {
			err = fmt.Errorf("<Init-read-config> FATAL: mysql settings is empty, please set \"DB.mysql\"")
			return nil, err
		}
		if config_yaml.DB.Mysql.Host == "" ||
			config_yaml.DB.Mysql.Port > 65535 ||
			config_yaml.DB.Mysql.Port < 1 ||
			config_yaml.DB.Mysql.User == "" ||
			config_yaml.DB.Mysql.Password == "" ||
			config_yaml.DB.Mysql.Database == "" {
			err = fmt.Errorf("<Init-read-config> FATAL: mysql settings is invalid, please check \"DB.mysql\"")
			return nil, err
		}
	case "sqlite":
		// detect sqlite.db db file
		if config_yaml.DB.Sqlite == nil {
			err = fmt.Errorf("<Init-read-config> FATAL: sqlite settings is empty, please set \"DB.sqlite\"")
			return nil, err
		}
		if config_yaml.DB.Sqlite.DBFile == "" {
			err = fmt.Errorf("<Init-read-config> FATAL: sqlite settings is invalid, please check \"DB.sqlite\"")
			return nil, err
		}
	}

	// lower-case and check ConfigYaml.Log.LogLevel, if it's not in debug, info, warn, error
	// then send a warn and set it to LogLvlWarn
	config_yaml.Log.LogLevel = strings.ToLower(config_yaml.Log.LogLevel)
	if config_yaml.Log.LogLevel != logger.LogLvlDebug &&
		config_yaml.Log.LogLevel != logger.LogLvlInfo &&
		config_yaml.Log.LogLevel != logger.LogLvlWarn &&
		config_yaml.Log.LogLevel != logger.LogLvlError {
		logger.Warnf("<Init-read-config> WARNING log.loglevel %v is invalid, must be in debug, info, warn, error.", config_yaml.Log.LogLevel)
		logger.Warnln(" Set it to default as level warn!")
		config_yaml.Log.LogLevel = logger.LogLvlWarn
	}

	return &config_yaml, nil
}

// read ms.mail.autosendmails.template and ms.mail.autosendmails.templatecontent
func read_mail_template() {
	template := vp.GetString("ms.mail.autosendmails.template")
	template_content := ""
	if len(template) > 0 {
		t_byts, err := os.ReadFile(template)
		if err != nil {
			t_path := filepath.Join(configYamlObj.Workspace, template)
			t_byts, _ = os.ReadFile(t_path)
		}
		t_template_content := string(t_byts)
		if len(t_template_content) > 0 {
			template_content = t_template_content
		}
	}
	if len(template_content) == 0 {
		template_content = vp.GetString("ms.mail.autosendmails.templatecontent")
	}
	if len(template_content) == 0 {
		template_content = ms.MailTemplate
	}
	configYamlObj.MS.Mail.AutoSendMails.Template = template
	configYamlObj.MS.Mail.AutoSendMails.TemplateContent = template_content
}

func monitor_config_change(file_path string) {
	vp.OnConfigChange(func(e fsnotify.Event) {
		logger.Warnf("Config file changed:\n", e.Name)
		new_config, err := read_config(file_path)
		if err != nil {
			logger.Warnf("failed to reload config, failed with: %v\n", err)
			return
		}
		// bot_token changed, warning to restart daemon to take effect
		if new_config.BotToken != configYamlObj.BotToken {
			logger.Warnf("bot_token changed, please restart daemon to take effect\n")
		}
		// proxy changed, replace ConfigYamlInstance.Proxy to new, and re-create ProxyObj
		if new_config.Proxy != configYamlObj.Proxy {
			configYamlObj.Proxy = new_config.Proxy
			proxyObj, err = config.NewProxyValid(configYamlObj.Proxy)
			if err != nil {
				logger.Warnf("proxy section is invalid, please do double check\n")
			}
		}
		// Goroutine, BindMax, ErrLimit, Notice changed, just copy to ConfigYamlInstance
		configYamlObj.BindMax = new_config.BindMax
		configYamlObj.Goroutine = new_config.Goroutine
		configYamlObj.ErrLimit = new_config.ErrLimit
		configYamlObj.Notice = new_config.Notice
		// re-create AdminSet, if admin list changed
		if len(new_config.Admin) > 0 {
			adminSet = config.NewAdminList(getAdmins())
		}
		// cron changed, need cronjob to stop to change the cron setting
		if new_config.CronConf.Enabled != configYamlObj.CronConf.Enabled ||
			new_config.CronConf.Task != configYamlObj.CronConf.Task ||
			new_config.CronConf.Notice != configYamlObj.CronConf.Notice {
			configYamlObj.CronConf.Enabled = new_config.CronConf.Enabled
			configYamlObj.CronConf.Task = new_config.CronConf.Task
			configYamlObj.CronConf.Notice = new_config.CronConf.Notice
			//  we need wait for the original cronjob to finish
			if cronObj != nil {
				ctx := cronObj.Stop()
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
			init_background_tasks(configYamlObj.CronConf)
		}
		// DB changed, need to restart daemon to take effect
		if new_config.DB.DBType != configYamlObj.DB.DBType {
			logger.Warnf("DB changed, please restart daemon to take effect\n")
		} else if new_config.DB.DBType == "mysql" {
			if new_config.DB.Mysql.Host != configYamlObj.DB.Mysql.Host ||
				new_config.DB.Mysql.Port != configYamlObj.DB.Mysql.Port ||
				new_config.DB.Mysql.User != configYamlObj.DB.Mysql.User ||
				new_config.DB.Mysql.Password != configYamlObj.DB.Mysql.Password ||
				new_config.DB.Mysql.Database != configYamlObj.DB.Mysql.Database ||
				new_config.DB.Mysql.TLS != configYamlObj.DB.Mysql.TLS {
				logger.Warnf("DB changed, please restart daemon to take effect\n")
			}
		} else if new_config.DB.DBType == "sqlite" {
			if new_config.DB.Sqlite.DBFile != configYamlObj.DB.Sqlite.DBFile {
				logger.Warnf("DB changed, please restart daemon to take effect\n")
			}
		}
		// if log settings changed, need to restart daemon to take effect
		if new_config.Log.LogIntoFile != configYamlObj.Log.LogIntoFile ||
			new_config.Log.LogFile != configYamlObj.Log.LogFile ||
			new_config.Log.LogLevel != configYamlObj.Log.LogLevel ||
			new_config.Log.MaxSize != configYamlObj.Log.MaxSize ||
			new_config.Log.MaxBackups != configYamlObj.Log.MaxBackups ||
			new_config.Log.MaxAge != configYamlObj.Log.MaxAge ||
			new_config.Workspace != configYamlObj.Workspace {
			logger.Warnf("log settings changed, please restart daemon to take effect\n")
		}
		// read new setting for SaveOpDetails, SaveOpDetails
		configYamlObj.Log.SaveOpDetails = new_config.Log.SaveOpDetails
		configYamlObj.Log.SaveTaskRecords = new_config.Log.SaveOpDetails
		// if Workspace changed, need to restart daemon to take effect
		if new_config.Workspace != configYamlObj.Workspace {
			logger.Warnf("Workspace changed, please restart daemon to take effect\n")
		}
		// copy new_config.MS to ConfigYamlObj.MS
		*(configYamlObj.MS.Mail.ReadMails) = *(new_config.MS.Mail.ReadMails)
		*(configYamlObj.MS.Mail.SearchMails) = *(new_config.MS.Mail.SearchMails)
		*(configYamlObj.MS.Mail.AutoSendMails) = *(new_config.MS.Mail.AutoSendMails)
		*(configYamlObj.MS.Mail.AutoDeleteMails) = *(new_config.MS.Mail.AutoDeleteMails)
		*(configYamlObj.MS.File.ListFiles) = *(new_config.MS.File.ListFiles)
		*(configYamlObj.MS.Calendar.ListCalendars) = *(new_config.MS.Calendar.ListCalendars)
		*(configYamlObj.MS.Calendar.ListEvents) = *(new_config.MS.Calendar.ListEvents)
		*(configYamlObj.MS.Calendar.ListReminders) = *(new_config.MS.Calendar.ListReminders)
		*(configYamlObj.MS.Calendar.GetSchedule) = *(new_config.MS.Calendar.GetSchedule)
		// resolve template and templatecontent
		read_mail_template()
	})

	vp.WatchConfig()
}

func getAdmins() []int64 {
	var result []int64
	for _, v := range configYamlObj.Admin {
		id, _ := strconv.ParseInt(v, 10, 64)
		result = append(result, id)
	}

	return result
}

// initialize background cron tasks
// this must be called after botTelegram initialized
func init_background_tasks(cron_conf *config.ConfigCron) {
	if cron_conf.Enabled {
		cronObj = cron.New()
		cronObj.AddFunc(cron_conf.Task, PerformTasks)
		cronObj.AddFunc(cron_conf.Notice, NotifyStats)
		cronObj.Start()
	}
}

func Init(conf string) {
	// init vp
	vp = viper.New()
	// read config
	var err error
	configYamlObj, err = read_config(conf)
	if err != nil {
		logger.Errorf("read_config error: %v\n", err.Error())
		os.Exit(1)
	}
	// init logger
	logger.Init(configYamlObj.Log.LogIntoFile,
		configYamlObj.Workspace,
		configYamlObj.Log.LogFile,
		configYamlObj.Log.LogLevel,
		configYamlObj.Log.MaxSize,
		configYamlObj.Log.MaxBackups,
		configYamlObj.Log.MaxAge)
	// storage init must be done after logger init, because storage.Init() would using logger
	storage.Init(configYamlObj.Workspace, configYamlObj.DB)
	// do cache initialization
	authCachedObj = NewAuthCache()
	bindCachedObj = NewBindCache()
	usersConfigCacheObj = NewUsersConfigCache()
	jobLock = &sync.Mutex{}
	// AdminList
	adminSet = config.NewAdminList(getAdmins())
	// setup monitor to monitor the change of config file
	// this should be run in goroutine, because vp.WatchConfig() will block the main goroutine
	go monitor_config_change(conf)
}
