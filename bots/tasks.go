package bots

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
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
		stats, e := GetAppsStatsByTgId(tg_id)
		if e != nil {
			logger.Errorf("<NotifyStats> failed to perform storage.GetAppsStatsByTgId by tg_id %v, failed with: %v\n", tg_id, e.Error())
			continue
		}
		handleSendStats(context.Background(), botTelegram, tg_id, stats)
	}
}

func PerformTasks() {
	if ok := JobLock.TryLock(); !ok {
		logger.Info("<PerformTasks> job is running, skip this time")
		return
	}
	defer JobLock.Unlock()
	var t_start_time time.Time
	var t_end_time time.Time
	// record this task
	if ConfigYamlObj.Log.SaveTaskRecords {
		t_start_time = time.Now()
	}
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
	var wg_prd WaitGroupCount  // for task producer
	var wg_con WaitGroupCount  // for task consumer
	var wg_task WaitGroupCount // for tasks themselves
	var in chan *ms.Task
	var out chan *ms.ApiResult
	var done chan bool
	// tasks_count := 0
	// init all chan
	thread_count := utils.MinInt(t_len_users_config*3, ConfigYamlObj.Goroutine)
	// make more rooms for chan, consider we split the tasks into smaller tasks.
	chan_num := t_len_users_config * 10
	in = make(chan *ms.Task, chan_num)
	out = make(chan *ms.ApiResult, chan_num)
	done = make(chan bool, thread_count)
	// loop all_users_config, get user's mail
	mails_list := []*string{}
	for _, uc := range all_users_config {
		mails_list = append(mails_list, &uc.MsUsername)
	}
	// shuffle mails_list
	utils.ShuffleSlice(mails_list)
	// put task
	// add to wg
	for i, uc := range all_users_config {
		wg_prd.Add(1)
		go func(uc *storage.Users, to *string) {
			// mails
			// release from wg
			defer wg_prd.Done()
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
						return
					}
					// clean old unbind message to void too many messages in user side
					t_key_list := BindCachedObj.FindMsgKeyByChatIDAndReplyTypeExtraDataKey1(tg_id, ReplyForUnbindAccountS2, ms_username)
					for _, v := range t_key_list {
						CleanTGMsgAndBindCached(context.Background(), botTelegram, v)
					}
					// send option message to user to unbound this account, and cached in BindCachedObj
					handleUnbindAccountS1Helper(context.Background(), botTelegram, tg_id, uc, "PerformTasks")
				}
				return
			}

			args := map[string]interface{}{
				ms.ArgUserID:          user_id,
				ms.ArgAccessToken:     &t_token,
				ms.ArgTo:              to,
				ms.ArgReadAttachments: false,
			}
			// read mails from root
			if ConfigYamlObj.MS.Mail.ReadMails.Enabled {
				wg_task.Add(1)
				args[ms.ArgReadAttachments] = ConfigYamlObj.MS.Mail.ReadMails.ReadAttachments
				in <- &ms.Task{
					Func: ms.MailListMailsFrom,
					Args: args,
				}
			}
			// list mail folders
			if ConfigYamlObj.MS.Mail.ReadMailFolders.Enabled {
				wg_task.Add(1)
				args[ms.ArgReadAttachments] = ConfigYamlObj.MS.Mail.ReadMailFolders.ReadAttachments
				in <- &ms.Task{
					Func: ms.MailListMailFolders,
					Args: args,
				}
			}
			// search mails
			if ConfigYamlObj.MS.Mail.SearchMails.Enabled {
				args[ms.ArgReadAttachments] = ConfigYamlObj.MS.Mail.SearchMails.ReadAttachments
				for _, t_keyword := range ConfigYamlObj.MS.Mail.SearchMails.Keywords {
					wg_task.Add(1)
					in <- &ms.Task{
						Func: ms.MailsSearch,
						Args: map[string]interface{}{
							ms.ArgUserID:      user_id,
							ms.ArgAccessToken: &t_token,
							ms.ArgTo:          to,
							ms.ArgKeyword:     &t_keyword,
						},
					}
				}
			}
			// delete mails
			if ConfigYamlObj.MS.Mail.AutoDeleteMails.Enabled {
				wg_task.Add(1)
				args[ms.ArgReadAttachments] = ConfigYamlObj.MS.Mail.AutoDeleteMails.ReadAttachments
				in <- &ms.Task{
					Func: ms.MailsDelListSpecificMailFolders,
					Args: args,
				}

			}
			// send mail
			if ConfigYamlObj.MS.Mail.AutoSendMails.Enabled {
				wg_task.Add(1)
				in <- &ms.Task{
					Func: ms.MailsSend,
					Args: args,
				}

			}
		}(uc, mails_list[i])
	}
	// we need read the template as the content for sending emails
	// read template before we start the task
	if ConfigYamlObj.MS.Mail.AutoSendMails.Enabled {
		read_mail_template()
	}
	// to reduce the cache size, we need clean Template and TemplateContent after background task done
	defer func() {
		ConfigYamlObj.MS.Mail.AutoSendMails.Template = ""
		ConfigYamlObj.MS.Mail.AutoSendMails.TemplateContent = ""
	}()

	// start multiple goroutines to perform WorkingOnMsFromChan
	for i := 0; i < thread_count; i++ {
		wg_con.Add(1)
		go WorkingOnMsFromChan(in, out, done, &wg_con, ConfigYamlObj.Proxy, ConfigYamlObj.MS)
	}
	// handle results

	var stats []*storage.Stats
	var details []*storage.OpDetails

	t_count := 0
	var cache []*ms.Task
RESULT_LOPPER:
	for {
		select {
		case r := <-out:
			wg_task.Done()
			t_s := &storage.Stats{
				UserID:   r.ID,
				OpID:     r.OpID,
				Success:  r.S,
				Failure:  r.F,
				LastTime: r.EndTime.Unix(),
			}
			stats = append(stats, t_s)
			t_count++
			// record the details
			if ConfigYamlObj.Log.SaveOpDetails {
				d := &storage.OpDetails{
					UserID:    r.ID,
					OpID:      r.OpID,
					StartTime: r.StartTime.Unix(),
					EndTime:   r.EndTime.Unix(),
					Success:   r.S,
					Failure:   r.F,
					Duration:  r.Duration,
				}
				details = append(details, d)
				if len(details) >= MaxOpDetailsRecordCached {
					e := storage.SaveOpDetails(details)
					if e != nil {
						logger.Errorf("<PerformTasks> failed to store details, failed with: %v\n", e.Error())
					} else {
						// clean details while storage succeed
						details = []*storage.OpDetails{}
					}
				}
			}
			// print debug log about this result
			logger.Debugf("<PerformTasks> got result for user - %v, op - %v, s/f - %v/%v, finish at - %v, duration - %v\n", r.ID, r.OpID, r.S, r.F, r.EndTime.Format("15:04:05"), r.Duration)
			logger.Debugf("<PerformTasks> busy-%v, done-%v, todo-%v\n", wg_task.GetCount(), t_count, len(in)+len(cache))
			t_task_list := r.Tasks
			// put task into cache
			if len(t_task_list) > 0 {
				cache = append(cache, t_task_list...)
			}
			// check length of stats, store back into db if stats has too many records
			if len(stats) >= MaxStatsRecordCached {
				e := storage.UpdateStatsByStatsNew(stats)
				if e != nil {
					logger.Errorf("<PerformTasks> failed to store stats, failed with: %v\n", e.Error())
				} else {
					// clean stats while storage succeed
					stats = []*storage.Stats{}
				}
			}
		default:
			// put task from cache into in
			spared := chan_num - len(in)
			for spared > 0 && len(cache) > 0 {
				wg_task.Add(1)
				t_task := cache[0]
				in <- t_task
				if len(cache) > 1 {
					cache = cache[1:]
				} else {
					cache = []*ms.Task{}
				}
				spared = chan_num - len(in)
			}
			// when we get enough results and task producer all finished, then break the looper, we are finishing
			if wg_task.GetCount() == 0 && wg_prd.GetCount() == 0 && len(in) == 0 && len(cache) == 0 {
				break RESULT_LOPPER
			}
			time.Sleep(50 * time.Millisecond)
		}
	}

	// put done for all threads after we receive all results,
	for i := 0; i < thread_count; i++ {
		done <- true
	}

	// wait for all tasks
	wg_prd.Wait()
	wg_con.Wait()
	wg_task.Wait()

	// close all chan
	close(in)
	close(out)
	close(done)
	// update storage
	e = storage.UpdateStatsByStatsNew(stats)
	if e != nil {
		logger.Errorf("<PerformTasks> failed to store stats, failed with: %v\n", e.Error())
	}
	// just store op_details and task_records when setting is on
	if ConfigYamlObj.Log.SaveOpDetails {
		storage.SaveOpDetails(details)
	}
	if ConfigYamlObj.Log.SaveTaskRecords {
		// task end here
		t_end_time = time.Now()
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
		storage.SaveTaskRecords(task_records)
	}
}

func WorkingOnMsFromChan(in chan *ms.Task, out chan *ms.ApiResult, done chan bool, wg *WaitGroupCount, proxy string, ms_conf *config.ConfigMs) {
	for {
		select {
		case task := <-in:
			task.Func(out, &proxy, ConfigYamlObj.MS, task.Args)
		case ok := <-done:
			if ok {
				wg.Done()
				return
			}
		default:
			time.Sleep(ms.APIInterval)
		}
	}
}
