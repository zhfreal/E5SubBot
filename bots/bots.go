package bots

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/db"
	"github.com/zhfreal/E5SubBot/logger"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
)

var bot *tb.Bot

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

func Start(work_dir string) {
	var err error
	fmt.Print(logo)

	config.Init(work_dir)
	logger.Init(work_dir)
	db.Init()
	InitTask()

	poller := &tb.LongPoller{Timeout: 15 * time.Second}
	midPoller := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		if upd.Message == nil {
			return true
		}
		if !upd.Message.Private() {
			return false
		}
		return true
	})
	setting := tb.Settings{
		Token:  config.BotToken,
		Poller: midPoller,
	}

	if config.Socks5 != "" {
		dialer, err := proxy.SOCKS5("tcp", config.Socks5, nil, proxy.Direct)
		if err != nil {
			zap.S().Fatalw("failed to get dialer",
				"error", err, "proxy", config.Socks5)
		}
		transport := &http.Transport{}
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
		setting.Client = &http.Client{Transport: transport}
	}

	bot, err = tb.NewBot(setting)
	if err != nil {
		zap.S().Fatalw("failed to create bot", "error", err)
	}
	fmt.Printf("Bot: %d %s\n", bot.Me.ID, bot.Me.Username)

	makeHandlers()
	fmt.Println("Bot Start")
	fmt.Println("------------")
	bot.Start()
}

func makeHandlers() {
	// 所有用户
	bot.Handle("/start", bStart)
	bot.Handle("/my", bMy)
	bot.Handle("/bind", bBind)
	bot.Handle("/unbind", bUnBind)
	bot.Handle("/export", bExport)
	bot.Handle("/help", bHelp)
	bot.Handle(tb.OnText, bOnText)
	// 管理员
	bot.Handle("/task", bTask)
	bot.Handle("/log", bLog)
}
