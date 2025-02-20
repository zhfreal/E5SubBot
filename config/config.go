package config

import (
	"os"
	"strconv"
	"strings"

	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	LogBasePath    string = "./log/"
	WelcomeContent string = "欢迎使用E5SubBot!"
	HelpContent    string = `
	命令：
	/my 查看已绑定账户信息
	/bind  绑定新账户
	/unbind 解绑账户
	/export 导出账户信息(JSON)
	/help 帮助
	源码及使用方法：https://github.com/iyear/E5SubBot
`
)

func Init(work_dir string) {

	viper.SetConfigName("config")
	viper.AddConfigPath(work_dir)

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatalw("failed to read config", "error", err)
	}
	BotToken = viper.GetString("bot_token")
	Cron = viper.GetString("cron")
	Socks5 = viper.GetString("socks5")

	viper.SetDefault("errlimit", 5)
	viper.SetDefault("bindmax", 5)
	viper.SetDefault("goroutine", 10)

	BindMaxNum = viper.GetInt("bindmax")
	MaxErrTimes = viper.GetInt("errlimit")
	Notice = viper.GetString("notice")

	MaxGoroutines = viper.GetInt("goroutine")
	Admins = getAdmins()
	DB = viper.GetString("db")
	Table = viper.GetString("table")

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
		sqlite_db_new := filepath.Join(work_dir, sqlite_db_org)
		if _, err := os.Stat(sqlite_db_org); err == nil {
			sqlite_db_new = sqlite_db_org
		}
		Sqlite = &sqliteConfig{
			DB: sqlite_db_new,
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		MaxGoroutines = viper.GetInt("goroutine")
		BindMaxNum = viper.GetInt("bindmax")
		MaxErrTimes = viper.GetInt("errlimit")
		Notice = viper.GetString("notice")
		Admins = getAdmins()
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
