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
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDHelp, bot.MatchTypeExact, helpHandler)
	// for admin, we match "/bindApp <client_id> <app_alias>" in replyHandler
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDBindApp, bot.MatchTypeExact, bindAppHandler)
	// for all, we match "/bind <app_alias>" in replyHandler
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDBind, bot.MatchTypeExact, bindAccountHandler)
	// for all, we match "/reAuth <app_alias> <user_alias>" in replyHandler
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDReAuth, bot.MatchTypeExact, reAuthAccountHandler)
	// for all, we match "/unbind <app_alias> <user_alias>" in replyHandler
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDBind, bot.MatchTypeExact, unBindAccountHandler)
	// for admin, we match "/unbindOther <app_alias> <user_alias>" in replyHandler
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDUnbindOther, bot.MatchTypeExact, unBindAccountHandlerOther)
	// for admin, we match "/delApp <app_alias>" in replyHandler
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDDelApp, bot.MatchTypeExact, delAppHandler)
	// for all
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDListApps, bot.MatchTypeExact, showAPPsHandler)
	// for all
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDListUsers, bot.MatchTypeExact, showBoundUsersHandler)
	// for all
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDStat, bot.MatchTypeExact, statHandler)
	// for admin
	botTelegram.RegisterHandler(bot.HandlerTypeMessageText, CMDStatAll, bot.MatchTypeExact, statAllHandler)
	// init background task
	// this must be called after bot initialized
	InitBackgroundTasks()
	// this is for test only
	// PerformTasks()
	botTelegram.Start(ctx)
	// show logo after boot start
	fmt.Print(logo)
}
