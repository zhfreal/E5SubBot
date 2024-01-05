package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// const (
//     LogBasePath    string = "./log/"
// )

func Init(file_path string) {
	if len(file_path) == 0 {
		file_path = "./"
	}
	info, err := os.Stat(file_path)
	if err != nil {
		fmt.Printf("Can't access %v\n", file_path)
		os.Exit(1)
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
		fmt.Println("failed to read config, ", "error:", err)
		os.Exit(1)
	}
	BotToken = viper.GetString("bot_token")
	Cron = viper.GetString("cron")
	CronNotice = viper.GetString("cron_notice")
	Socks5 = viper.GetString("socks5")
	// print deprecated warning info
	if len(Socks5) > 0 {
		fmt.Println("<Init> WARNING: socks5 is deprecated, use \"Proxy: socks5://127.0.0.1:1080\" or \"Proxy: http://127.0.0.1:8080\" or \"Proxy: https://example.com:8443\" instead")
	}
	Proxy = viper.GetString("proxy")
	if len(Proxy) > 0 {
		var e error
		ProxyObj, e = NewProxyType(Proxy)
		if e != nil {
			fmt.Println("<Init> WARNING: invalid proxy settings, use \"Proxy: socks5://127.0.0.1:1080\" or \"Proxy: http://127.0.0.1:8080\" or \"Proxy: https://example.com:8443\".")
		}
	} else if len(Socks5) > 0 {
		var e error
		ProxyObj, e = NewProxyType("socks5://" + Socks5)
		if e != nil {
			fmt.Println("<Init> WARNING: invalid socks5 proxy settings, use \"Socks5: 127.0.0.1:1080\", or just \"Proxy: socks5://127.0.0.1:1080\" or \"Proxy: http://127.0.0.1:8080\" or \"Proxy: https://example.com:8443\" instead.")
		}
	} else {
		ProxyObj = &ProxyType{Url: nil, UrlStr: ""}
	}

	viper.SetDefault("errlimit", 5)
	viper.SetDefault("bindmax", 5)
	viper.SetDefault("goroutine", 10)
	viper.SetDefault("ms.mail.auto-delete.enabled", MailAutoDeleteEnabled)
	viper.SetDefault("ms.mail.auto-delete.keyword", MailAutoDeleteKeyWords[0])
	viper.SetDefault("ms.mail.auto-delete.keywords", MailAutoDeleteKeyWords)
	viper.SetDefault("ms.mail.auto-delete.quantity", MailAutoDeleteQuantity)
	viper.SetDefault("ms.mail.read-unread", MailReadUnread)
	// logging
	viper.SetDefault("log.log-into-file", LogIntoFile)
	viper.SetDefault("log.log-dir", LogDir)
	viper.SetDefault("log.log-file", LogFile)
	viper.SetDefault("log.log-level", LogLevel)
	viper.SetDefault("log.max-size", MaxSize)
	viper.SetDefault("log.max-backups", MaxBackups)
	viper.SetDefault("log.max-age", MaxAge)

	BindMaxNum = viper.GetInt("bindmax")
	MaxErrTimes = viper.GetInt("errlimit")
	Notice = viper.GetString("notice")

	// read from config.yaml
	// mail deletion settings
	readMsSection()

	// logging settings
	LogIntoFile = viper.GetBool("log.log-into-file")
	LogDir = viper.GetString("log.log-dir")
	LogFile = viper.GetString("log.log-file")
	LogLevel = strings.ToLower(viper.GetString("log.log-level"))
	MaxSize = viper.GetInt("log.max-size")
	MaxBackups = viper.GetInt("log.max-backups")
	MaxAge = viper.GetInt("log.max-age")

	MaxGoroutines = viper.GetInt("goroutine")
	Admins = getAdmins()
	AdminSet = NewAdminList(Admins)
	DB = viper.GetString("db")

	switch DB {
	case "mysql":
		Mysql = &mysqlConfig{
			Host:                viper.GetString("mysql.host"),
			Port:                viper.GetInt("mysql.port"),
			User:                viper.GetString("mysql.user"),
			Password:            viper.GetString("mysql.password"),
			DB:                  viper.GetString("mysql.database"),
			SSLMode:             viper.GetString("mysql.ssl_mode"),
			EnabledTLSProtocols: viper.GetString("mysql.enabled_tls_protocols"),
		}
	case "sqlite":
		// detect sqlite.db db file
		sqlite_db_org := viper.GetString("sqlite.db")
		Sqlite = &sqliteConfig{
			DB: sqlite_db_org,
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		MaxGoroutines = viper.GetInt("goroutine")
		BindMaxNum = viper.GetInt("bindmax")
		MaxErrTimes = viper.GetInt("errlimit")
		Notice = viper.GetString("notice")
		Admins = getAdmins()
		// reload ms section
		readMsSection()
	})
}

func getAdmins() []int64 {
	var result []int64
	admins := strings.Split(viper.GetString("admin"), ",")
	for _, v := range admins {
		id, _ := strconv.ParseInt(v, 10, 64)
		result = append(result, id)
	}
	return result
}

func readMsSection() {
	// ms.mail.auto-delete.enabled
	MailAutoDeleteEnabled = viper.GetBool("ms.mail.auto-delete.enabled")
	// ms.mail.auto-delete.keywords
	MailAutoDeleteKeyWords = viper.GetStringSlice("ms.mail.auto-delete.keywords")
	// ms.mail.auto-delete.keyword
	MailAutoDeleteKeyWords = append(MailAutoDeleteKeyWords, viper.GetString("ms.mail.auto-delete.keyword"))
	// ms.mail.auto-delete.quantity
	MailAutoDeleteQuantity = viper.GetInt("ms.mail.auto-delete.quantity")
	// make a map to remove duplicates item in MailAutoDeleteKeyWords
	// make sure each item in MailAutoDeleteKeyWords is trim and not empty
	MailAutoDeleteKeyWordsMap := make(map[string]bool)
	for _, v := range MailAutoDeleteKeyWords {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			MailAutoDeleteKeyWordsMap[v] = true
		}
	}
	MailAutoDeleteKeyWords = make([]string, 0, len(MailAutoDeleteKeyWordsMap))
	for k := range MailAutoDeleteKeyWordsMap {
		MailAutoDeleteKeyWords = append(MailAutoDeleteKeyWords, k)
	}
	MailReadUnread = viper.GetBool("ms.mail.read-unread")
}
