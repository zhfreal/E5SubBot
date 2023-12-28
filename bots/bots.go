package bots

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/logger"
	"github.com/zhfreal/E5SubBot/storage"
)

func Start(conf_file string, show_token bool, account string) {
	var err error

	config.Init(conf_file)
	logger.Init(config.LogIntoFile, config.LogDir, config.LogFile, config.LogLevel, config.MaxSize, config.MaxBackups, config.MaxAge)
	// storage init must be done after logger init, because storage.Init() would using logger
	storage.Init()
	// self Init
	Init()
	if show_token {
		ShowToken(account)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(replyHandler),
	}
	// add socks5 proxy if it is set in config
	// make proxied bot client base on config.ProxyObj
	if config.ProxyObj != nil && config.ProxyObj.Url != nil {
		transport := &http.Transport{
			Proxy: http.ProxyURL(config.ProxyObj.Url),
		}
		opts = append(opts, bot.WithHTTPClient(time.Minute, &http.Client{Transport: transport}))
	}

	botTelegram, err = bot.New(config.BotToken, opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		fmt.Println("failed to create bot", "error", err.Error())
		os.Exit(1)
	}
	// register handlers
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)
	// for admin
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/bindapp", bot.MatchTypeExact, bindAppHandler)
	// for all
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/bind", bot.MatchTypeExact, bindAccountHandler)
	// for all
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/reauth", bot.MatchTypeExact, reAuthAccountHandler)
	// for all
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/unbind", bot.MatchTypeExact, unBindAccountHandler)
	// for admin
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/unbindother", bot.MatchTypeExact, unBindAccountHandlerOther)
	// for admin
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/delapp", bot.MatchTypeExact, delAppHandler)
	// for all
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/stat", bot.MatchTypeExact, statHandler)
	// for admin
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, "/statall", bot.MatchTypeExact, statAllHandler)
	// init background task
	// this must be called after bot initialized
	InitBackgroundTasks()
	// this is for test only
	// PerformTasks()
	botTelegram.Start(ctx)
	// show logo after boot start
	fmt.Print(logo)
}
