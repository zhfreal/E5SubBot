package bots

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-telegram/bot"
)

func Start(config_file string) {
	var err error
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(replyHandler),
	}
	// add socks5 proxy if it is set in config
	// make proxied bot client base on config.ProxyObj
	if ProxyObj != nil && ProxyObj.UrlStr != "" {
		transport := &http.Transport{
			Proxy: http.ProxyURL(ProxyObj.Url),
		}
		opts = append(opts, bot.WithHTTPClient(time.Minute, &http.Client{Transport: transport}))
	}

	botTelegram, err = bot.New(ConfigYamlObj.BotToken, opts...)
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

	// to void potential error, we need run botTelegram.Start() in goroutine
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		botTelegram.Start(ctx)
		wg.Done()
	}()
	// init background task
	// this must be init cronjob after bot start
	// init_background_tasks(ConfigYamlObj.CronConf)
	// debug only
	PerformTasks()
	wg.Wait()
}
