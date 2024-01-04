package bots

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/robfig/cron/v3"
	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/logger"
	ms "github.com/zhfreal/E5SubBot/microsoft"
	"github.com/zhfreal/E5SubBot/storage"
	"github.com/zhfreal/E5SubBot/utils"
)

func NotifyStats() {
	// get all Stats
	tg_list, e := storage.GetAllTgId()
	if e != nil {
		logger.Errorf("<NotifyStats> failed to get tg_id from table users, failed with: %v\n", e.Error())
		return
	}
	if len(tg_list) == 0 {
		return
	}
	// loop and send stats
	for _, tg_id := range tg_list {
		stats, e := storage.GetAppsStatsByTgId(tg_id)
		if e != nil {
			logger.Errorf("<NotifyStats> failed to perform storage.GetAppsStatsByTgId by tg_id %v, failed with: %v\n", tg_id, e.Error())
			continue
		}
		handleSendStats(context.Background(), botTelegram, tg_id, stats)
	}
}

func PerformTasks() {
	// record this task
	t_start_time := time.Now()
	// get all users config
	all_users_config, e := storage.GetAllUsersEnabled()
	if e != nil {
		logger.Error("<PerformTasks> failed to get all users config, failed with: %v\n", e.Error())
		return
	}
	t_len_users_config := len(all_users_config)
	// check if there is any users config
	if t_len_users_config == 0 {
		return
	}
	app_map := make(map[uint]string, 0)
	all_apps, _ := storage.GetAllAppEnabled()
	for _, v := range all_apps {
		app_map[v.ID] = v.Alias
	}
	var wg sync.WaitGroup
	var in chan ms.Args
	var out chan ms.ApiResult
	var done chan bool
	tasks_count := 0
	// init all chan
	thread_count := utils.MinInt(t_len_users_config, config.MaxGoroutines)
	in = make(chan ms.Args, t_len_users_config)
	out = make(chan ms.ApiResult, t_len_users_config)
	done = make(chan bool, t_len_users_config)
	// put task
	// add to wg
	wg.Add(1)
	go func() {
		for _, uc := range all_users_config {
			// mails
			user_id := uc.ID
			tg_id := uc.TgId
			t_token, is_refresh_err, e := GetToken(uc)
			if e != nil {
				logger.Errorf("<PerformTasks> failed to get token by id %v, failed with: %v\n", user_id, e.Error())
				// notice user refresh token failed, and give a option to unbind this account
				if is_refresh_err {
					t_refresh_failure_count := UsersConfigCacheObj.GetFailCount(user_id)
					ms_username := uc.MsUsername
					app_name := app_map[uc.AppId]
					// send notice to unbind, this message is not the one which bounded the unbound account action
					// and cached in BindCachedObj
					t_msg := fmt.Sprintf("Failed to get token for account %v which bind in %v and this happened for %v times recently. You may need unbind it and re-authorize it.", ms_username, app_name, t_refresh_failure_count)
					_, e := botTelegram.SendMessage(context.Background(), &bot.SendMessageParams{
						ChatID: tg_id,
						Text:   t_msg,
					})
					if e != nil {
						logger.Errorf("<PerformTasks> failed to send account unbind option message to %v, failed with: %v\n", tg_id, e.Error())
						continue
					}
					// clean old unbind message to void too many messages in user side
					t_key_list := BindCachedObj.FindMsgKeyByChatIDAndReplyTypeExtraDataKey1(tg_id, ReplyForUnbindAccountS2, ms_username)
					for _, v := range t_key_list {
						CleanTGMsgAndBindCached(context.Background(), botTelegram, v)
					}
					// send option message to user to unbound this account, and cached in BindCachedObj
					handleUnbindAccountS1Helper(context.Background(), botTelegram, tg_id, uc, "PerformTasks")
				}
				continue
			}
			args := ms.Args{
				Func:        ms.WorkingOnMails,
				ID:          user_id,
				AccessToken: t_token,
			}
			in <- args
			tasks_count += 1
		}
		// release from wg
		wg.Done()
	}()

	for i := 0; i < thread_count; i++ {
		wg.Add(1)
		go ms.WorkingOnMsFromChan(in, out, done, &wg, config.ProxyObj.UrlStr)
	}
	// handle results

	var stats []*storage.Stats
	var details []*storage.OpDetails

	//
	wg.Add(1)
	go func() {
		t_count := 0
		for r := range out {
			t_s := &storage.Stats{
				UserID:   r.ID,
				OpID:     r.OpID,
				Success:  r.S,
				Failure:  r.F,
				LastTime: r.StartTime,
			}
			stats = append(stats, t_s)
			t_count++
			// record the details
			d := &storage.OpDetails{
				UserID:    r.ID,
				OpID:      r.OpID,
				StartTime: r.StartTime,
				EndTime:   r.EndTime,
				Success:   r.S,
				Failure:   r.F,
				Duration:  r.Duration,
			}
			details = append(details, d)
			// print debug log about this result
			logger.Debugf("<PerformTasks> got result for user - %v, op - %v, s/f - %v/%v, finish at - %v, duration - %v\n", r.ID, r.OpID, r.S, r.F, time.Unix(r.EndTime, 0).Format("15:04:05"), r.Duration)
			logger.Debugf("<PerformTasks> %v/%v, %v\n", t_count, tasks_count, r)
			// send done message to goroutine
			if t_count <= tasks_count {
				done <- true
			}
			if t_count >= tasks_count {
				break
			}
			// time.Sleep(500 * time.Microsecond)
		}
		wg.Done()
	}()
	// wait for all tasks
	wg.Wait()
	// task end here
	t_end_time := time.Now()
	t_duration := t_end_time.Sub(t_start_time).Milliseconds()
	logger.Debugf("<PerformTasks> all tasks done, duration: %vms\n", t_duration)
	// record task record
	t_record := &storage.TaskRecords{
		StartTime: t_start_time.Unix(),
		EndTime:   t_end_time.Unix(),
		Duration:  t_duration,
	}
	var task_records []*storage.TaskRecords
	task_records = append(task_records, t_record)
	// close all chan
	close(in)
	close(out)
	close(done)
	// handle stats
	t_stats_map := make(map[storage.TypeUserIDOpID]*storage.Stats, 0)
	for _, v := range stats {
		t_key := storage.TypeUserIDOpID{
			UserId: v.UserID,
			OpId:   v.OpID,
		}
		t_stats_map[t_key] = v
	}
	// update storage
	storage.UpdateStatsByStats(t_stats_map)
	// just store op_details and task_records when debug on
	if config.LogLevel == "debug" {
		storage.SaveOpDetails(details)
		storage.SaveTaskRecords(task_records)
	}
}

// initialize background cron tasks
// this must be called after botTelegram initialized
func InitBackgroundTasks() {
	c := cron.New()
	c.AddFunc(config.Cron, PerformTasks)
	c.AddFunc(config.CronNotice, NotifyStats)
	c.Start()
}
