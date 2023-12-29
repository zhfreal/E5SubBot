package config

import (
	"fmt"
	"net/url"

	"github.com/zhfreal/E5SubBot/utils"
)

var (
	BotToken           string
	Socks5             string
	Proxy              string
	ProxyHTTPSInsecure bool
	BindMaxNum         int
	MaxGoroutines      int
	MaxErrTimes        int
	Cron               string
	CronNotice         string
	Notice             string
	Admins             []int64
	AdminSet           *AdminList
	DB                 string
	Mysql              *mysqlConfig
	Sqlite             *sqliteConfig
	ProxyObj           *ProxyType
	// this is available for other package to search and delete mails
	// it will be modified by config.yaml during initialization
	// ms.mail.auto-delete.enabled
	MailAutoDeleteEnabled bool = true
	// ms.mail.auto-delete.keyword
	MailAutoDeleteKeyWord string = "George Best quote"
	// ms.mail.auto-delete.quantity
	MailAutoDeleteQuantity int = 20
	// TelegramBot        *bot.Bot

	// for logging
	LogIntoFile bool   = true
	LogFile     string = "latest.log"
	LogDir      string = "./log"
	MaxSize     int    = 5      // single file max size in MiB
	MaxBackups  int    = 20     // max quantity of log files to keep
	MaxAge      int    = 7      // max days of log files to keep
	LogLevel    string = "warn" // debug, info, warn, error, fatal, panic
)

type sqliteConfig struct {
	DB string `json:"db,omitempty"`
}
type mysqlConfig struct {
	Host                string `json:"host,omitempty"`
	Port                int    `json:"port,omitempty"`
	User                string `json:"user,omitempty"`
	Password            string `json:"password,omitempty"`
	DB                  string `json:"db,omitempty"`
	SSLMode             string `json:"ssl_mode,omitempty"`
	EnabledTLSProtocols string `json:"enabled_tls_protocols,omitempty"`
}

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

type ProxyType struct {
	Url    *url.URL
	UrlStr string
}

func NewProxyType(proxy string) (*ProxyType, error) {
	t_url, t_e := utils.ParseProxy(proxy)
	if t_e != nil {
		e := fmt.Errorf("invalid proxy %v", proxy)
		return nil, e
	}
	return &ProxyType{
		Url:    t_url,
		UrlStr: proxy,
	}, nil
}
