package config

import (
	"fmt"
	"net/url"

	"github.com/zhfreal/E5SubBot/utils"
)

var (
// BotToken           string
// Socks5             string
// Proxy              string
// ProxyHTTPSInsecure bool
// BindMaxNum    int
// MaxGoroutines int
// MaxErrTimes int
// Cron       string
// CronNotice string
// Notice   string
// Admins   []int64
// AdminSet *AdminList
// DB       string
// Mysql    *mysqlConfig
// Sqlite   *sqliteConfig
// ProxyObj *ProxyType
// this is available for other package to search and delete mails
// it will be modified by config.yaml during initialization
// ms.mail.auto-delete.enabled
// 	MailAutoDeleteEnabled bool = true
// 	// ms.mail.auto-delete.keyword
// 	MailAutoDeleteKeyWords []string = []string{"George Best quote"}
// 	// ms.mail.auto-delete.quantity
// 	MailAutoDeleteQuantity int = 20
// 	// TelegramBot        *bot.Bot
// 	MailReadUnread bool = false // read unread mails

// // for logging
// LogIntoFile bool   = true
// LogFile     string = "latest.log"
// Workspace   string = "/var/lib/e5bot/"
// MaxSize     int    = 5      // single file max size in MiB
// MaxBackups  int    = 20     // max quantity of log files to keep
// MaxAge      int    = 7      // max days of log files to keep
// LogLevel    string = "warn" // debug, info, warn, error, fatal, panic
)

type AdminList struct {
	admins map[int64]bool
}

func NewAdminList(ids []int64) *AdminList {
	a := &AdminList{}
	a.admins = make(map[int64]bool)
	a.AddMore(ids)
	return a
}

func (a *AdminList) Add(id int64) {
	a.admins[id] = true
}

func (a *AdminList) AddMore(ids []int64) {
	for _, id := range ids {
		a.admins[id] = true
	}
}

func (a *AdminList) Has(id int64) bool {
	return a.admins[id]
}

type ProxyValid struct {
	Url    *url.URL
	UrlStr string
}

func NewProxyValid(proxy string) (*ProxyValid, error) {
	t_url, t_e := utils.ParseProxy(proxy)
	if t_e != nil {
		e := fmt.Errorf("invalid proxy %v", proxy)
		return nil, e
	}
	return &ProxyValid{
		Url:    t_url,
		UrlStr: proxy,
	}, nil
}

// define a struct based on config.yaml, to unmarshal config.yaml into this struct.
// config.yaml has content like:
// # bot-token: 6478297263:AAFeqBKfNlf5hw1qjogC7KLP22aXRFMMInY
// bot-token: 6616723465:AAFGrzuUG2Eal_RGLx_URJxj2ZMGlz0njg8
// #bot-token: 5366662643:AAH4vtyBl5Vp8sXEATDHSofjrBBEMRx7hrw
// #proxy: socks5://192.168.2.1:1084
// proxy: ""
// bind-max: 999
// goroutine: 20
// admin:
//   - 358920093
// err-limit: 9
// notice: |-
//   welcome!
// cron:
//   task: "*/5 * * * *"
//   notice: "*/30 * * * *"
// workspace: "./"
// db:
//   type: sqlite
//   mysql:
//     host: 127.0.0.1
//     port: 3306
//     user: root
//     password: pwd
//     database: e5sub
//     charset: utf8mb4
//     tls: true
//   sqlite:
//     db-file: data-test.db
// ms:
//   mail:
//     auto-delete:
//       enabled: true
//       keyword: ""
//       keywords:
//        - "As promised in the last email"
//        - "George Best quote"
//        - "Delivery has failed to these recipients or groups"
//       quantity: 20
//     read-mails:
//       enabled: true
//       quantity: 20
//       read-unread: true
//     search-mails:
//       enabled: true
//       quantity: 20
//       read-unread: true
//     auto-send:
//       enabled: true
//       template: "send_mail.html"
//       template_content: ""
//       template_type: "html"
// log:
//     log-into-file: true
//     log-file: "logs/latest.log"
//     log-level: "debug"
//     # in MiB
//     max-size: 5
//     # in days
//     max-age: 7
//     # quantity
//     max-backups: 20

type ConfigYaml struct {
	BotToken  string      `yaml:"bottoken"`
	Proxy     string      `yaml:"proxy"`
	Socks5    string      `yaml:"socks5"`
	BindMax   int         `yaml:"bindmax" default:"9999"`
	Goroutine int         `yaml:"goroutine" default:"20"`
	Admin     []string    `yaml:"admin"`
	ErrLimit  int         `yaml:"errlimit" default:"6"`
	Notice    string      `yaml:"notice" default:"welcome!"`
	CronConf  *ConfigCron `yaml:"cronconf"`
	Workspace string      `yaml:"workspace" default:"./"`
	DB        *ConfigDb   `yaml:"db"`
	MS        *ConfigMs   `yaml:"ms"`
	Log       *ConfigLog  `yaml:"log"`
}
type ConfigCron struct {
	Task    string `yaml:"task" default:"*/10 * * * *"`
	Notice  string `yaml:"notice" default:"*/30 * * * *"`
	Enabled bool   `yaml:"enabled" default:"false"`
}
type ConfigDb struct {
	DBType string        `yaml:"dbtype" default:"sqlite3"`
	Mysql  *MySqlConfig  `yaml:"mysql"`
	Sqlite *SqliteConfig `yaml:"sqlite"`
}
type SqliteConfig struct {
	DBFile string `yaml:"dbfile" default:"data.db"`
}
type MySqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset" default:"utf8mb4"`
	TLS      string `yaml:"tls" default:"false"`
}
type ConfigMs struct {
	Mail *ConfigMail `yaml:"mail"`
}
type ConfigMail struct {
	AutoDeleteMails *ConfigAutoDelete      `yaml:"autodeletemails"`
	ReadMails       *ConfigReadMails       `yaml:"readmails"`
	ReadMailFolders *ConfigReadMailFolders `yaml:"readmailfolders"`
	SearchMails     *ConfigSearchMails     `yaml:"searchmails"`
	AutoSendMails   *ConfigAutoSend        `yaml:"autosendmails"`
}
type ConfigAutoDelete struct {
	Enabled         bool     `yaml:"enabled" default:"false"`
	Keyword         string   `yaml:"keyword"`
	Keywords        []string `yaml:"keywords"`
	Quantity        int      `yaml:"quantity" default:"20"`
	ReadUnread      bool     `yaml:"readunread" default:"true"`
	ReadAttachments bool     `yaml:"readattachments" default:"false"`
}
type ConfigReadMails struct {
	Enabled         bool `yaml:"enabled" default:"true"`
	Quantity        int  `yaml:"quantity" default:"20"`
	ReadUnread      bool `yaml:"readunread" default:"true"`
	ReadAttachments bool `yaml:"readattachments" default:"false"`
}
type ConfigReadMailFolders struct {
	Enabled         bool `yaml:"enabled" default:"true"`
	Quantity        int  `yaml:"quantity" default:"20"`
	ReadUnread      bool `yaml:"readunread" default:"true"`
	ReadAttachments bool `yaml:"readattachments" default:"false"`
}
type ConfigSearchMails struct {
	Enabled         bool     `yaml:"enabled" default:"true"`
	Keyword         string   `yaml:"keyword"`
	Keywords        []string `yaml:"keywords"`
	Quantity        int      `yaml:"quantity" default:"20"`
	ReadUnread      bool     `yaml:"readunread" default:"true"`
	ReadAttachments bool     `yaml:"readattachments" default:"false"`
}
type ConfigAutoSend struct {
	Enabled         bool   `yaml:"enabled" default:"false"`
	Template        string `yaml:"template"`
	TemplateContent string `yaml:"templatecontent"`
	TemplateType    string `yaml:"templatetype"`
	Subject         string `yaml:"subject"`
}
type ConfigLog struct {
	LogIntoFile     bool   `yaml:"logintofile" default:"true"`
	LogFile         string `yaml:"logfile" default:"logs/latest.log"`
	LogLevel        string `yaml:"loglevel" default:"warn"`
	MaxSize         int    `yaml:"maxsize" default:"5"`
	MaxAge          int    `yaml:"maxage" default:"7"`
	MaxBackups      int    `yaml:"maxbackups" default:"20"`
	SaveOpDetails   bool   `yaml:"saveopdetails" default:"false"`
	SaveTaskRecords bool   `yaml:"savetaskrecords" default:"false"`
}
