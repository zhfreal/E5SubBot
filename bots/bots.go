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
	// we do all binding in replyHandler
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
	// bot.Start() will block the current running thread,
	// so we need run botTelegram.Start() in goroutine, and do other things
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		botTelegram.Start(ctx)
		wg.Done()
	}()
	// init background task, th
	// this must be init cronjob after bot start
	init_background_tasks(ConfigYamlObj.CronConf)
	// debug only
	// PerformTasks()
	wg.Wait()
}
