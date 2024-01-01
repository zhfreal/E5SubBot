package bots

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/logger"
	ms "github.com/zhfreal/E5SubBot/microsoft"
	"github.com/zhfreal/E5SubBot/storage"
	"github.com/zhfreal/E5SubBot/utils"
)

// get expiredAt time
// expired_time in minutes
func GetExpiredTimeFromNowAfter(seconds int) time.Time {
	return GetExpiredTimeFromATimeAfter(time.Now(), seconds)
}

func GetExpiredTimeFromATimeAfter(from_time time.Time, seconds int) time.Time {
	return from_time.Add(time.Duration(seconds) * time.Second)
}

// delete specific message
func DeleteMsg(ctx context.Context, b *bot.Bot, chat_id int64, msg_id int) {
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chat_id,
		MessageID: msg_id,
	})
}

func CleanTGMsgAndBindCached(ctx context.Context, b *bot.Bot, key *MsgKey) {
	// delete msg
	DeleteMsg(ctx, b, key.ChatID, key.MsgID)
	// delete cached
	BindCachedObj.Del(key)
}

func CleanMsgBindCachedAndUnlockAuthCache(ctx context.Context, b *bot.Bot, key *MsgKey) {
	CleanTGMsgAndBindCached(ctx, b, key)
	// unlock AuthCachedObj
	AuthCachedObj.Unlock(key.ChatID)
}

func CleanMsgBindCachedAndCancelAuthCache(ctx context.Context, b *bot.Bot, key *MsgKey) {
	CleanTGMsgAndBindCached(ctx, b, key)
	// delete cached
	AuthCachedObj.Cancel(key.ChatID)
}

func CleanMsgBindCachedUnlockAndCancelAuthCache(ctx context.Context, b *bot.Bot, key *MsgKey) {
	CleanTGMsgAndBindCached(ctx, b, key)
	// delete cached
	AuthCachedObj.Cancel(key.ChatID)
	// unlock AuthCachedObj
	AuthCachedObj.Unlock(key.ChatID)
}

// handle reply messages
func replyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// msg_id, chat_id represent the previous message stored in BindCachedObj
	var msg_id int
	var chat_id int64
	var received_msg string
	var received_action string
	var t_key *MsgKey
	var t_value *MsgValue
	this_msg_type := -1
	// it's a callback reply
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		msg_id = update.CallbackQuery.Message.ID
		chat_id = update.CallbackQuery.Message.Chat.ID
		received_action = update.CallbackQuery.Data
		this_msg_type = ReplyFromCallBack
	} else if update.Message != nil {
		chat_id = update.Message.Chat.ID
		if update.Message.ReplyToMessage != nil {
			// pure text reply
			msg_id = update.Message.ReplyToMessage.ID
			chat_id = update.Message.ReplyToMessage.Chat.ID
			this_msg_type = ReplyWithPureMsg
		} else {
			// pure text message initialize from user side
			this_msg_type = InitialMsg
		}
		received_msg = update.Message.Text
	} else {
		// unsupported msg, send help message to user
		helpHandler(ctx, b, update)
		return
	}
	// get the previous message stored in BindCachedObj, check it expired or not
	// if it's expired, clean it and send msg to user
	if this_msg_type == ReplyFromCallBack || this_msg_type == ReplyWithPureMsg {
		t_key = &MsgKey{MsgID: msg_id, ChatID: chat_id}
		// don't handle non-cached msg
		if !BindCachedObj.Has(t_key) {
			// clean the callback
			DeleteMsg(ctx, b, chat_id, msg_id)
			// send msg to user
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chat_id,
				Text:   "This message has expired, please start over again.",
			})
			return
		}
		t_value = BindCachedObj.Get(t_key)
		// don't handle expired msg
		if t_value.IsExpired() {
			// send msg to user
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chat_id,
				Text:   "This message has expired, please start over again.",
			})
			// clean
			if t_value.MsgType == ReplyForWaitingAuth {
				// authentication go routing will clean itself, just do ctx cancel
				AuthCachedObj.Cancel(t_key.ChatID)
			} else {
				CleanTGMsgAndBindCached(ctx, b, t_key)
			}
			return
		}
	}

	// handle reply, this is for "/BindApp"
	if this_msg_type == ReplyWithPureMsg {
		// got an incorrect response message
		if t_value.MsgType == ReplyForBindAPP {
			handleAPPBound(ctx, b, t_key, received_msg)
		} else {
			// b.SendMessage(ctx, &bot.SendMessageParams{
			// 	ChatID: chat_id,
			// 	Text:   "Unsupported reply message.",
			// })
			// print help
			helpHandler(ctx, b, update)
			return
		}
	} else if this_msg_type == ReplyFromCallBack {
		if t_value.MsgType == ReplyForBindAccount {
			// user request to bind an account
			if received_action == ActionYes {
				// run binding in a goroutine
				go handleAccountBinding(ctx, b, t_key, false)
			}
		} else if t_value.MsgType == ReplyForReAuth {
			// user request to bind an account
			if received_action == ActionYes {
				// run binding in a goroutine
				go handleAccountBinding(ctx, b, t_key, true)
			}
		} else if t_value.MsgType == ReplyForWaitingAuth {
			// user request to cancel auth
			if received_action == ActionCancel {
				CleanMsgBindCachedAndCancelAuthCache(ctx, b, t_key)
			}
		} else if t_value.MsgType == ReplyForUnbindAccountS1 {
			// hand unbind account step 1: handle user's choose (it's an app)
			if received_action == ActionYes {
				handleUnbindAccountS1(ctx, b, t_key)
			}
		} else if t_value.MsgType == ReplyForUnbindAccountS2 {
			// hand unbind account step 2: handle unbind action based on user's choose
			if received_action == ActionYes {
				handleUnbindAccountS2(ctx, b, t_key)
			}
		} else if t_value.MsgType == ReplyForDeleteAPP {
			if received_action == ActionYes {
				handleAPPDeletion(ctx, b, t_key)
			}
		} else {
			helpHandler(ctx, b, update)
		}
	} else if this_msg_type == InitialMsg {
		full_message := update.Message.Text
		t_list := utils.SplitString(full_message)
		// empty message, print help
		if len(t_list) == 0 {
			helpHandler(ctx, b, update)
		} else if len(t_list) == 1 {
			// just single cmd without any parameters
			t_cmd := t_list[0]
			switch t_cmd {
			case CMDBindApp:
				bindAppHandler(ctx, b, update)
			case CMDDelApp:
				delAppHandler(ctx, b, update)
			case CMDBind:
				bindAccountHandler(ctx, b, update)
			case CMDUnbind:
				unBindAccountHandler(ctx, b, update)
			case CMDReAuth:
				reAuthAccountHandler(ctx, b, update)
			case CMDUnbindOther:
				unBindAccountHandlerOther(ctx, b, update)
			case CMDListApps:
				showAPPsHandler(ctx, b, update)
			case CMDListUsers:
				showBoundUsersHandler(ctx, b, update)
			case CMDStat:
				statHandler(ctx, b, update)
			case CMDStatAll:
				statAllHandler(ctx, b, update)
			default:
				// print help
				helpHandler(ctx, b, update)
			}
		} else if len(t_list) > 1 {
			// cmd with parameters
			t_cmd := t_list[0]
			switch t_cmd {
			case CMDBindApp:
				bindAppHandlerFromCMD(ctx, b, chat_id, t_list)
			case CMDDelApp:
				delAppByFromCMD(ctx, b, chat_id, t_list)
			case CMDBind:
				bindAccountFromCMD(ctx, b, chat_id, t_list)
			case CMDUnbind:
				unbindUserFromCMD(ctx, b, chat_id, t_list)
			case CMDReAuth:
				reAuthUserByAppAliasUserAliasFromCMD(ctx, b, chat_id, t_list)
			case CMDUnbindOther:
				unbindUserOtherFromCMD(ctx, b, chat_id, t_list)
			default:
				// print help
				helpHandler(ctx, b, update)
			}
		}

	} else {
		// print help
		helpHandler(ctx, b, update)
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Bind APP
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// handle "/bindApp" and "/bindApp <client_id> <app_alias>"
// Limitation: just admin and bind app
func bindAppHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chat_id := update.Message.Chat.ID
	// check if user is admin by it's chat_id
	// non-admin just return
	if !config.AdminSet.Has(chat_id) {
		// b.SendMessage(ctx, &bot.SendMessageParams{
		//     ChatID: chat_id,
		//     Text:   "You are not admin, please contact admin to bind an APP.",
		// })
		return
	}
	// do bind
	bindAPPHandlerHelper(ctx, b, chat_id)
}

// handle "/bindApp <client_id> <app_alias>"
// Limitation: just admin and bind app
func bindAppHandlerFromCMD(ctx context.Context, b *bot.Bot, chat_id int64, msg_list []string) {
	// check if user is admin by it's chat_id
	// non-admin just return
	if !config.AdminSet.Has(chat_id) {
		// b.SendMessage(ctx, &bot.SendMessageParams{
		//     ChatID: chat_id,
		//     Text:   "You are not admin, please contact admin to bind an APP.",
		// })
		return
	}
	t_list := msg_list
	// single /binApp
	if len(t_list) == 1 {
		bindAPPHandlerHelper(ctx, b, chat_id)
	} else if len(t_list) == 3 {
		// /bindApp <client_id> <app_alias>
		saveAppBoundInfo(t_list[1], t_list[2], ctx, b, chat_id)
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Incorrect format, please reply as \"/bindApp <client_id> <app_alias>\" or just \"/bindApp\"",
		})
	}
}

// send tg message "Please reply as \"client id\" \"alias\"
// and store key MsgKey{chat_id, msg_id} in BindCached
func bindAPPHandlerHelper(ctx context.Context, b *bot.Bot, chat_id int64) {
	m, e := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		// MessageThreadID: myUSeq.Get(),
		Text:        "Please reply as \"client id\" \"alias\"\n",
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: &models.ForceReply{ForceReply: true, Selective: true},
	})
	if nil != e {
		logger.Errorf("<bindAPPHandlerHelper> send bind message failed, failed with: %v\n", e.Error())
	} else {
		// store key MsgKey{chat_id, msg_id} in BindCached
		msg_id := m.ID
		msgTime := time.Unix(int64(m.Date), 0)
		expireAt := GetExpiredTimeFromATimeAfter(msgTime, BindCacheTimeInSeconds)
		key := &MsgKey{chat_id, msg_id}
		// set msg_cached.BindCached
		v := &MsgValue{MsgType: ReplyForBindAPP, ExpiredAt: expireAt, Extra: nil}
		BindCachedObj.Add(key, v)
	}
}

// handle APP bound message, it should be a reply to "Please reply as \"client id\" \"alias\"\n"
// store client id and alias as user's response
// and send user a message about successfully APP binding and an option to bind an account
// "msg" format: "<client_id> <alias>"
func handleAPPBound(ctx context.Context, b *bot.Bot, key *MsgKey, msg string) {
	t_list := strings.Split(msg, " ")
	chat_id := key.ChatID
	// clean
	defer CleanTGMsgAndBindCached(ctx, b, key)
	// incorrect format
	if len(t_list) != 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Incorrect format, please reply as \"<client id> <alias>\"",
		})
		bindAPPHandlerHelper(ctx, b, chat_id)
		return
	}
	client_id := t_list[0]
	alias := t_list[1]
	// handle app bound: store app info and send tg message to offer an option for binding an account
	saveAppBoundInfo(client_id, alias, ctx, b, chat_id)
}

// save app bound info
func saveAppBoundInfo(client_id string, alias string, ctx context.Context, b *bot.Bot, chat_id int64) {
	// check alias exists or not
	if storage.AppsHasAlias(alias) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   fmt.Sprintf("The alias %s already used", alias),
		})
		bindAPPHandlerHelper(ctx, b, chat_id)
		return
	}
	// check client_id exists or not
	if storage.AppsHasClientID(client_id) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   fmt.Sprintf("The client id %s already exists in APP config", client_id),
		})
		bindAPPHandlerHelper(ctx, b, chat_id)
		return
	}
	// store APP info
	e := storage.AddAppConfig(&storage.AppConfig{
		ClientId:     client_id,
		Alias:        alias,
		ClientSecret: "",
		FromLegacy:   false,
	})
	// handle database error
	if e != nil {
		logger.Errorf("<HandleAPPBound> failed to store app config, failed with: %v\n", e.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "failed to store APP bound info, please contact admin",
		})
		return
	}
	// send t_msg to user and
	// give an option to bind an account
	app_id, e := storage.GetAppIDByClientID(client_id)
	if e != nil {
		logger.Errorf("<HandleAPPBound> failed to get app id by client id %v, failed with: %v\n", client_id, e.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "failed to store APP bound info, please contact admin",
		})
		return
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Bind an account", CallbackData: ActionYes},
			},
		},
	}
	t_msg := fmt.Sprintf("Success to bind an app:  %s", alias)
	m, e := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chat_id,
		Text:        t_msg,
		ReplyMarkup: kb,
	})
	// handle telegram bot api error
	if e != nil {
		logger.Errorf("<HandleAPPBound> failed to response to user for bind an account, failed with: %v\n", e.Error())
		return
	}
	msg_id := m.ID
	expireAt := GetExpiredTimeFromNowAfter(BindCacheTimeInSeconds)
	// store to BindCached
	t_key := &MsgKey{m.Chat.ID, msg_id}
	t_value := &MsgValue{
		MsgType:   ReplyForBindAccount,
		ExpiredAt: expireAt,
		Extra: &ExtraData{
			ExtraData1String: client_id,
			ExtraData1Uint:   app_id,
			ExtraData2String: alias,
		},
	}
	BindCachedObj.Add(t_key, t_value)
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// bind or re-auth account
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// handle "/bind"
func bindAccountHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// get chat id
	chat_id := update.Message.Chat.ID
	// check if user has a pending device code for authentication
	// if yes, send message to user and return
	if AuthCachedObj.IsLocked(chat_id) {
		l_dc := AuthCachedObj.GetDeviceCode(chat_id)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   fmt.Sprintf("Your have pending device code for authenticate, please finish it! %v", l_dc.Msg),
		})
		return
	}
	// get app list
	app_list, e := storage.GetAllApp()
	if e != nil {
		logger.Errorf("<BindAccountHandler> failed to get app list from database, failed with: %v\n", e.Error())
		return
	}
	// check if app list is empty
	if len(app_list) == 0 {
		msg := "There is no app bound."
		if config.AdminSet.Has(chat_id) {
			// admin
			msg = fmt.Sprintf("%v.Please use /bindApp to bind an app", msg)
		} else {
			// non-admin
			msg = fmt.Sprintf("%v.Please contact admin to bind an app", msg)
		}
		// send message and returnSleep
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   msg,
		})
		return
	}
	// send tips
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   "Choose an app to bind an account",
	})

	// list all app alias and offer an option to user to choose to bind an account
	for _, app := range app_list {
		time.Sleep(BotSendMsgInterval)
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Bind", CallbackData: ActionYes},
				},
			},
		}
		m, e := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chat_id,
			ReplyMarkup: kb,
			Text:        fmt.Sprintf("Bind an account to %s", app.Alias),
		})
		if e != nil {
			logger.Errorf("<BindAccountHandler> failed to send bind option message to %v, failed with: %v\n", chat_id, e.Error())
			continue
		}
		msg_id := m.ID
		expireAt := GetExpiredTimeFromNowAfter(BindCacheTimeInSeconds)
		// store to BindCached
		t_key := &MsgKey{m.Chat.ID, msg_id}
		t_value := &MsgValue{
			MsgType:   ReplyForBindAccount,
			ExpiredAt: expireAt,
			Extra: &ExtraData{
				ExtraData1Uint: app.ID,
			},
		}
		BindCachedObj.Add(t_key, t_value)
	}
}

// after user click on button "Bind an account"
// request a device code, notice user to authorize it,
// and check the authorization result in time.
// re_auth: indicate this an re-auth procession
// When do /bind, ExtraData like:
//
//	  &ExtraData{
//	                ExtraData1Uint: app.ID,
//				}
//
// When do /reAuth, ExtraData like:
//
//	  &ExtraData{
//	                ExtraData1Uint: app.ID,
//	                ExtraData2Uint: users.ID,
//				},
func handleAccountBinding(ctx context.Context, b *bot.Bot, key *MsgKey, re_auth bool) {
	chat_id := key.ChatID
	v := BindCachedObj.Get(key)
	app_id := v.Extra.ExtraData1Uint
	// get app config from database
	app_conf, e := storage.GetAppByID(app_id)
	if e != nil || app_conf == nil {
		logger.Errorf("<handleAccountBinding> failed to get app config from database, failed with: %v\n", e.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to process your request, please contact admin",
		})
		return
	}
	// clean original cached msg
	CleanTGMsgAndBindCached(context.Background(), b, key)
	var user_id uint
	if re_auth {
		user_id = v.Extra.ExtraData2Uint
	}
	handleAccountAuth(ctx, b, chat_id, app_conf.ID, app_conf.ClientId, app_conf.Alias, user_id, re_auth)
}

// request device code and check it's authorization status
// user_id work with re_auth, while we do /reAuth, user_id means the specific user's config ID we chose to reAuth
// it used for double check of new authorization account with old one
func handleAccountAuth(ctx context.Context, b *bot.Bot, chat_id int64, app_id uint, client_id, app_alias string, user_id uint, re_auth bool) {
	token_cache, e := ms.GetDeviceCode(context.Background(), client_id)
	if e != nil {
		logger.Errorf("<handleAccountAuth> get device code failed with %v\n", e.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get device code, please contact admin",
		})
		return
	}
	t_msg := token_cache.Message
	t_code_expired_at := time.Unix(token_cache.ExpireTime, 0)
	t_msg = fmt.Sprintf("Binding account into %v. %v. This code will expire at %v",
		app_alias, t_msg, t_code_expired_at.Local().Format("2006-01-02 15:04:05 MST"))
	// send message to auth device code, and offer an option to cancel this authorization
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Cancel auth", CallbackData: ActionCancel},
			},
		},
	}
	m, e := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chat_id,
		Text:        t_msg,
		ReplyMarkup: kb,
	})
	if e != nil {
		logger.Errorf("<handleAccountAuth> failed to send device code message to %v with error: %v\n", chat_id, e.Error())
		return
	}
	// add this message to BindCached
	this_msg_id := m.ID
	t_key := &MsgKey{chat_id, this_msg_id}
	t_value := &MsgValue{
		MsgType:   ReplyForWaitingAuth,
		ExpiredAt: t_code_expired_at,
		Extra: &ExtraData{
			ExtraData1String: token_cache.DeviceCode,
		},
	}
	BindCachedObj.Add(t_key, t_value)
	// lock AuthCached
	dc_for_cached := &PendingDeviceCode{ClientID: client_id, DeviceCode: token_cache.DeviceCode, Msg: t_msg}
	AuthCachedObj.Lock(chat_id, dc_for_cached, nil)
	t_expired_durations := time.Until(t_code_expired_at)
	ctx, cancel := context.WithTimeout(ctx, t_expired_durations)
	AuthCachedObj.AddCancelFunc(chat_id, cancel)
	// clean
	defer CleanMsgBindCachedAndUnlockAuthCache(ctx, b, t_key)
	e = ms.CheckAuthStatusOfDeviceCode(ctx, token_cache)
	if e != nil {
		logger.Errorf("<handleAccountAuth> failed to auth device code with error: %v\n", e.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to auth device code!",
		})
		return
	}
	// authenticate successfully
	// store token
	userConfig := storage.NewUsersConfig(app_id, chat_id, token_cache)
	// check usersConfig is exist or not
	// if exist, send message to user and return
	found, old_userConfig := storage.FindUserByAppIDTgIDMsUsername(userConfig.AppId, userConfig.TgId, userConfig.MsUsername)
	if re_auth {
		// check if userConfig is exist or not
		if !found || user_id != old_userConfig.ID {
			// not exists
			logger.Errorf("<handleAccountAuth> new authorization info is not match with old one, key: \appId - %v, tgId - %v, ms-username -%v\n",
				userConfig.AppId, userConfig.TgId, userConfig.MsUsername)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chat_id,
				Text:   "Failed to re-auth, please contact admin",
			})
		} else {
			// exists the original user's config
			// try to update storage
			storage.UpdateUserByAppIDTgIDMsUsername(userConfig)
		}
	} else {
		// not re-auth, means we do /bind
		// no same ms account in table users
		if !found {
			// add into storage
			e = storage.AddUsers(userConfig)
			// failed to add to storage
			if e != nil {
				logger.Errorf("<handleAccountAuth> failed to save binding info to db with error: %v\n", e.Error())
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chat_id,
					Text:   "Failed to bind an account, please contact admin",
				})
			} else {
				// successfully add to storage
				// add stats
				storage.InitStatsByUserID(userConfig.ID)
				// send message to user
				t_msg = fmt.Sprintf("Succeed to bind %v into APP %v!", userConfig.MsUsername, app_alias)
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chat_id,
					Text:   t_msg,
				})
			}
		} else {
			// we found the same ms account in table users
			// do update
			storage.UpdateUserByAppIDTgIDMsUsername(userConfig)
			logger.Debugf("<handleAccountAuth> found this account %v already bound in app %v, we do authorization update.\n", userConfig.MsUsername, app_alias)
			// send message to user
			t_msg = fmt.Sprintf("This account %v already bound in APP %v. We update authorization info !", userConfig.MsUsername, app_alias)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chat_id,
				Text:   t_msg,
			})
		}
	}
	// this bind or re-auth process will be considered a successful process for this account, reset failure count
	UsersConfigCacheObj.InitFailCount(userConfig.ID)
}

// handle "/reAuth"
func reAuthAccountHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// get chat id
	chat_id := update.Message.Chat.ID
	// check if user has a pending device code for authentication
	// if yes, send message to user and return
	if AuthCachedObj.IsLocked(chat_id) {
		l_dc := AuthCachedObj.GetDeviceCode(chat_id)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   fmt.Sprintf("Your have pending device code for authenticate, please finish it! %v", l_dc.Msg),
		})
		return
	}
	// the following code handle single "/reAuth" message
	// get all UsersConfig
	uc_list, e := storage.GetUsersByTgID(chat_id)
	// handle database error
	if e != nil || len(uc_list) == 0 {
		var t_msg string
		if e != nil {
			// got error while get data from database
			t_msg = "Failed to get your binding accounts, please contact admin"
			logger.Errorf("<reAuthAccountHandler> failed to get users config from db with error: %v\n", e.Error())
		} else {
			// send different message for empty user's config
			t_msg = "No binding accounts found! Please binding one with message \"/bind\"."
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   t_msg,
		})
		return
	}
	// get apps data for their alias
	apps_config, e := storage.GetAllApp()
	if e != nil || len(apps_config) == 0 {
		t_msg := "Failed to get apps-config from db, please contact admin"
		if e != nil {
			// got error while get data from database
			logger.Errorf("<reAuthAccountHandler> failed to get apps-config from db with error: %v\n", e.Error())
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   t_msg,
		})
		return
	}
	// make map to store {id, alias} for apps
	app_alias_map := make(map[uint]string)
	for _, app := range apps_config {
		app_alias_map[app.ID] = app.Alias
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   "Choose an account to re-auth.",
	})
	// loop through uc_list, and send option to user to choose to unbind it's account
	for _, uc := range uc_list {
		time.Sleep(BotSendMsgInterval)
		app_alias := app_alias_map[uc.AppId]
		t_msg := fmt.Sprintf("Re-auth: %v from APP - %v", uc.MsUsername, app_alias)
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Re-auth", CallbackData: ActionYes},
				},
			},
		}
		m, e := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chat_id,
			ReplyMarkup: kb,
			Text:        t_msg,
		})
		if e != nil {
			logger.Errorf("<reAuthAccountHandler> failed to send re-auth message to %v with error: %v\n", chat_id, e.Error())
			return
		}
		msg_id := m.ID
		// add to BindCached
		t_key := &MsgKey{m.Chat.ID, msg_id}
		t_value := &MsgValue{
			MsgType:   ReplyForReAuth,
			ExpiredAt: GetExpiredTimeFromNowAfter(BindCacheTimeInSeconds),
			Extra: &ExtraData{
				ExtraData1Uint: uc.AppId,
				ExtraData2Uint: uc.ID,
			},
		}
		BindCachedObj.Add(t_key, t_value)
	}
}

// handle "/bind <app_alias>", msg_list = []string{"/bind", "<app_alias>"}
func bindAccountFromCMD(ctx context.Context, b *bot.Bot, chat_id int64, msg_list []string) {
	// no app_alias, just return
	if len(msg_list) <= 1 {
		return
	}
	app_alias := msg_list[1]
	// get app_conf by app_alias
	app_conf, e := storage.GetAppByAlias(app_alias)
	if e != nil || app_conf == nil {
		if e != nil {
			// got error while get data from database, do logging
			t_msg := fmt.Sprintf("can't get stored app info by app_alias %v, due to incorrect APP alias or internal error.", app_alias)
			logger.Errorf("<bindAccountFromCMD> %v, failed with: %v\n", t_msg, e.Error())
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Incorrect App alias",
		})
		return
	}
	// bind account
	handleAccountAuth(ctx, b, chat_id, app_conf.ID, app_conf.ClientId, app_alias, 0, false)
}

// handle "/reAuth <app_alias> <user_alias>"
func reAuthUserByAppAliasUserAliasFromCMD(ctx context.Context, b *bot.Bot, chat_id int64, msg_list []string) {
	// no app_alias, just return
	if len(msg_list) <= 1 {
		return
	} else if len(msg_list) > 3 {
		logger.Errorf("<reAuthUserByAppAliasUserAliasFromCMD> msg_list is too long, len: %v\n", len(msg_list))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Incorrect cmd. Please input like \"/unbind <app_alias> <user_alias>\".",
		})
		return
	}
	// get app_alias and user_alias
	app_alias := msg_list[1]
	user_alias := msg_list[2]
	// get app_id
	appConf, e := storage.GetAppByAlias(app_alias)
	// handle database error
	if e != nil || appConf == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get APP ID by alias.",
		})
		if e != nil {
			// database error, logging
			logger.Errorf("<reAuthUserByAppAliasUserAliasFromCMD> failed to get app_id by alias %v, failed with: %v\n", app_alias, e.Error())
		}

		return
	}
	// get user's config by app_id, chat_id, user's alias,
	// this means we must make sure this <user_alias> - <user_alias> belong to this user specific by chat_id
	userConf, e := storage.GetUserByAppIDTgIDUserAlias(appConf.ID, chat_id, user_alias)
	// handle database error
	if e != nil || userConf == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get user ID by alias.",
		})
		if e != nil {
			// database error, logging
			logger.Errorf("<reAuthUserByAppAliasUserAliasFromCMD> failed to get user_id by alias %v, failed with: %v\n", user_alias, e.Error())
		}
	} else {
		// do re-auth account
		handleAccountAuth(ctx, b, chat_id, appConf.ID, appConf.ClientId, app_alias, userConf.ID, true)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// unbind account
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// handle "/unbind"
func unBindAccountHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// get usersconfig stat stat
	chat_id := update.Message.Chat.ID
	stat := storage.GetAppStatsForUser(chat_id)
	unbindAccountHandlerHelper(ctx, b, chat_id, stat)
}

// handle "/unbindOther"
// admin only
func unBindAccountHandlerOther(ctx context.Context, b *bot.Bot, update *models.Update) {
	// get usersconfig stat stat
	chat_id := update.Message.Chat.ID
	if config.AdminSet.Has(chat_id) {
		stat := storage.GetAppStatsOnlyBoundUser()
		unbindAccountHandlerHelper(ctx, b, chat_id, stat)
	}
}

// unbind account handler helper
func unbindAccountHandlerHelper(ctx context.Context, b *bot.Bot, chat_id int64, stat []*storage.AppConfigStat) {
	// no account bound, send message and return
	if len(stat) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "There is no account bound.",
		})
		return
	}
	// send app list for user to choose to unbind it's bounded account
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   "Choose an APP to unbind it's account.",
	})
	// send all app alias and offer an option to user to choose to unbind it's bounded account
	for _, app := range stat {
		time.Sleep(BotSendMsgInterval)
		t_msg := fmt.Sprintf("Unbind from: %v(%v)", app.Alias, app.Bound)
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Unbind", CallbackData: ActionYes},
				},
			},
		}
		m, e := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chat_id,
			ReplyMarkup: kb,
			Text:        t_msg,
		})
		if e != nil {
			logger.Errorf("<UnbindAccountHandler> failed to send app unbind option message to %v with error: %v\n", chat_id, e.Error())
			continue
		}
		msg_id := m.ID
		// cache to BindCached
		t_key := &MsgKey{m.Chat.ID, msg_id}
		t_value := &MsgValue{
			MsgType:   ReplyForUnbindAccountS1,
			ExpiredAt: GetExpiredTimeFromNowAfter(BindCacheTimeInSeconds),
			Extra: &ExtraData{
				ExtraData1String: app.ClientId,
				ExtraData1Uint:   app.ID,
				ExtraData2String: app.Alias,
			},
		}
		BindCachedObj.Add(t_key, t_value)
	}
}

// unbind account step1
// handle reply message which user choose an APP to unbind it's account
// reply to user all relevant account in that APP
func handleUnbindAccountS1(ctx context.Context, b *bot.Bot, key *MsgKey) {
	chat_id := key.ChatID
	if key.ChatID != chat_id {
		return
	}
	t_v := BindCachedObj.Get(key)
	client_id := t_v.Extra.ExtraData1String
	// get appid
	appConf, e := storage.GetAppByClientID(client_id)
	// handle database error
	if e != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get APP ID by client ID.",
		})
		return
	}
	// get all UsersConfig
	uc_list, e := storage.GetUsersByAppIDTgID(appConf.ID, chat_id)
	// handle database error
	if e != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get APP ID by client ID.",
		})
		return
	}
	// clean
	defer CleanTGMsgAndBindCached(ctx, b, key)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   "Choose an account to unbind.",
	})
	// loop through uc_list, and send option to user to choose to unbind it's account
	for _, uc := range uc_list {
		// sleep fro a while, to void potential issue
		time.Sleep(BotSendMsgInterval)
		handleUnbindAccountS1Helper(ctx, b, chat_id, uc, "handleUnbindAccountS1")
	}
}

// send unbind msg to TG, and cached in BindCachedObj
// <tag> is the identify the caller function which call this function, show main process
func handleUnbindAccountS1Helper(ctx context.Context, b *bot.Bot, tg_id int64, uc *storage.Users, tag string) {
	t_msg := fmt.Sprintf("Unbind: %v", uc.MsUsername)
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Unbind", CallbackData: ActionYes},
			},
		},
	}
	m, e := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      tg_id,
		ReplyMarkup: kb,
		Text:        t_msg,
	})
	if e != nil {
		logger.Errorf("<%v> failed to send account unbind option message to %v with error: %v\n", tag, tg_id, e.Error())
		return
	}
	msg_id := m.ID
	// add to BindCached
	t_key := &MsgKey{m.Chat.ID, msg_id}
	t_value := &MsgValue{
		MsgType:   ReplyForUnbindAccountS2,
		ExpiredAt: GetExpiredTimeFromNowAfter(BindCacheTimeInSeconds),
		Extra: &ExtraData{
			ExtraData1Uint:   uc.ID,
			ExtraData1String: uc.MsUsername,
			ExtraData2String: uc.Alias,
		},
	}
	BindCachedObj.Add(t_key, t_value)
}

// unbind account step2
// handle the unbound action according the user's choose a specific account
func handleUnbindAccountS2(ctx context.Context, b *bot.Bot, key *MsgKey) {
	t_v := BindCachedObj.Get(key)
	var user_id uint
	var e error
	user_id = t_v.Extra.ExtraData1Uint
	// get UsersConfig by user's ID
	uc, e := storage.GetUsersByID(user_id)
	if key.ChatID != uc.TgId {
		return
	}
	// handle database error
	if e != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: key.ChatID,
			Text:   "Failed to unbind account. Please contact admin.",
		})
		logger.Errorf("<handleUnbindAccountS2> failed to get usersConfig by id %v, failed with: %v\n", t_v.Extra, e.Error())
	}
	// uc is nil, means some error happened
	if uc == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: key.ChatID,
			Text:   "System error occurred. Please contact admin.",
		})
		return
	}

	// clean before return
	defer CleanTGMsgAndBindCached(ctx, b, key)
	// unbind, delete UsersConfig
	delUserConfig(ctx, b, key.ChatID, uc.ID, uc.MsUsername)
}

// delete user's config
func delUserConfig(ctx context.Context, b *bot.Bot, chat_id int64, user_id uint, username string) {
	// unbind, delete UsersConfig
	e := storage.DelUsers(user_id)
	// handle database error
	if e != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   fmt.Sprintf("Failed to unbind account %v. Please contact admin.", username),
		})
		logger.Errorf("<handleUnbindAccountS2> failed to delete users config by id %v, failed with: %v\n", user_id, e.Error())
		// just unlock
		// UsersConfigCacheObj.Unlock(uc.ID)
		return
	}
	// delete OpStats
	storage.DelStatsByUserID(user_id)

	// remove user cache
	UsersConfigCacheObj.DelCache(user_id)
	// send message to user
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   fmt.Sprintf("Account %v has been unbound.", username),
	})
}

// handle "/unbind <app_alias> <user_alias>", used by generic user
func unbindUserFromCMD(ctx context.Context, b *bot.Bot, chat_id int64, msg_list []string) {
	if len(msg_list) <= 1 {
		return
	} else if len(msg_list) > 3 {
		logger.Errorf("<unbindUserFromCMD> msg_list is too long, len: %v\n", len(msg_list))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Incorrect cmd. Please input like \"/unbind <app_alias> <user_alias>\".",
		})
		return
	}
	// get app_alias and user_alias
	app_alias := msg_list[1]
	user_alias := msg_list[2]
	// get app_id
	appConf, e := storage.GetAppByAlias(app_alias)
	// handle database error
	if e != nil || appConf == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get APP ID by alias.",
		})
		if e != nil {
			// database error, logging
			logger.Errorf("<unbindUserFromCMD> failed to get app_id by alias %v, failed with: %v\n", app_alias, e.Error())
		}

		return
	}
	// get user's config by app_id, chat_id, user's alias,
	// this means we must make sure this <user_alias> - <user_alias> belong to this user specific by chat_id
	userConf, e := storage.GetUserByAppIDTgIDUserAlias(appConf.ID, chat_id, user_alias)
	// handle database error
	if e != nil || userConf == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get user ID by alias.",
		})
		if e != nil {
			// database error, logging
			logger.Errorf("<unbindUserFromCMD> failed to get user_id by alias %v, failed with: %v\n", user_alias, e.Error())
		}
	}
	// delete
	delUserConfig(ctx, b, chat_id, userConf.ID, userConf.MsUsername)
}

// handle "/unbind <app_alias> <user_alias>", used by admin
func unbindUserOtherFromCMD(ctx context.Context, b *bot.Bot, chat_id int64, msg_list []string) {
	if len(msg_list) <= 1 {
		return
	} else if len(msg_list) > 3 {
		logger.Errorf("<UnbindUser> msg_list is too long, len: %v\n", len(msg_list))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Incorrect cmd. Please input like \"/unbind <app_alias> <user_alias>\".",
		})
		return
	}
	// get app_alias and user_alias
	app_alias := msg_list[1]
	user_alias := msg_list[2]
	// get app_id
	appConf, e := storage.GetAppByAlias(app_alias)
	// handle database error
	if e != nil || appConf == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get APP ID by alias.",
		})
		if e != nil {
			// database error, logging
			logger.Errorf("<UnbindUser> failed to get app_id by alias %v, failed with: %v\n", app_alias, e.Error())
		}

		return
	}
	// get user's config by app_id, user's alias, admin only
	userConf, e := storage.GetUsersByAppIDUserAlias(appConf.ID, user_alias)
	// handle database error
	if e != nil || userConf == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to get user ID by alias.",
		})
		if e != nil {
			// database error, logging
			logger.Errorf("<UnbindUser> failed to get user_id by alias %v, failed with: %v\n", user_alias, e.Error())
		}
	}
	// delete
	delUserConfig(ctx, b, chat_id, userConf.ID, userConf.MsUsername)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// delete APP
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// handle "/delApp"
// admin only
func delAppHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chat_id := update.Message.Chat.ID
	// no-admin will do nothing
	if !config.AdminSet.Has(chat_id) {
		return
	}
	// get all apps which no user bound
	s := storage.GetAppStatsOnlyNoUserBound()
	// no APP for deletion
	if len(s) == 0 {
		// send tips
		t_msg := "No APP for deletion. Please unbound their accounts before delete them."
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   t_msg,
		})
		return
	}
	// send tips
	t_msg := "Choose an APP to delete. Only APPs which have no account will be deleted."
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   t_msg,
	})
	// loop through s, and send option to user to choose to delete it's account
	for _, v := range s {
		// sleep a little moment, void potential failure
		time.Sleep(BotSendMsgInterval)
		t_msg := fmt.Sprintf("Delete: %v", v.Alias)
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Delete", CallbackData: ActionYes},
				},
			},
		}
		m, e := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chat_id,
			ReplyMarkup: kb,
			Text:        t_msg,
		})
		if e != nil {
			logger.Errorf("<DelAppHandler> failed to send delete option message %v to %v, failed with: %v\n", v.Alias, chat_id, e.Error())
		} else {
			msg_id := m.ID
			// add to BindCached
			t_key := &MsgKey{m.Chat.ID, msg_id}
			t_value := &MsgValue{
				MsgType:   ReplyForDeleteAPP,
				ExpiredAt: GetExpiredTimeFromNowAfter(BindCacheTimeInSeconds),
				Extra: &ExtraData{
					ExtraData1String: v.ClientId,
					ExtraData1Uint:   v.ID,
					ExtraData2String: v.Alias,
				},
			}
			BindCachedObj.Add(t_key, t_value)
		}
	}
}

// cached Extra like： ExtraData1String - client_id; ExtraData1Uint - app_id; ExtraData2String = app_alias
// handle APP deletion according user's choice
func handleAPPDeletion(ctx context.Context, b *bot.Bot, key *MsgKey) {
	value := BindCachedObj.Get(key)
	// client_id := value.Extra.ExtraData1String
	app_id := value.Extra.ExtraData1Uint
	app_alias := value.Extra.ExtraData2String
	// check if user is admin
	// non-admin will do nothing
	if !config.AdminSet.Has(key.ChatID) {
		return
	}
	// check this APP bound accounts or not
	app_stat := storage.GetAppStatByID(app_id)
	if app_stat == nil {
		logger.Errorf("<handleDeleteAPP> failed to get users config by app_id %v\n", app_id)
		return
	}
	// still have accounts, can't delete
	// send message to user
	if app_stat.Bound > 0 {
		t_msg := "Can't delete APP. Please unbind their accounts before delete them."
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: key.ChatID,
			Text:   t_msg,
		})
		return
	}
	// do deletion
	appDeletionHelper(ctx, b, key.ChatID, app_id, app_alias)
	// do clean all the time, no matter success or fail to delete APP
	CleanTGMsgAndBindCached(ctx, b, key)
}

// handle APP deletion according user's choice, return delete success or not
func appDeletionHelper(ctx context.Context, b *bot.Bot, chat_id int64, app_id uint, app_alias string) bool {
	// do deletion
	e := storage.DelAppConfig(app_id)
	if e != nil {
		logger.Errorf("<handleDeleteAPP> failed to delete app config by app id %v, failed with: %v\n", app_id, e.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   "Failed to delete APP. Please contact admin.",
		})
		return false
	} else {
		// send message to user about successfully delete
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   fmt.Sprintf("APP %v has been deleted successfully.", app_alias),
		})
		return true
	}
}

// handle "/delApp <app_alias>..."
// admin only
func delAppByFromCMD(ctx context.Context, b *bot.Bot, chat_id int64, cmd_list []string) {
	// no-admin will do nothing
	if !config.AdminSet.Has(chat_id) {
		return
	}
	if len(cmd_list) < 2 {
		// do nothing here
		return
	}
	for _, app_alias := range cmd_list[1:] {
		user_stat := storage.GetAppStatByAlias(app_alias)
		// no such APP
		if user_stat == nil {
			// logger.Errorf("<handleDeleteAPP> failed to get users config by app_alias %v\n", app_alias)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chat_id,
				Text:   fmt.Sprintf("No such APP with alias - %v.", app_alias),
			})
			return
		}
		// bound accounts, can't delete
		if user_stat.Bound > 0 {
			t_msg := "Can't delete APP. Please unbind their accounts before delete them."
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chat_id,
				Text:   t_msg,
			})
			return
		}
		// do deletion
		// sleep a little moment, void potential failure
		time.Sleep(BotSendMsgInterval)
		appDeletionHelper(ctx, b, chat_id, user_stat.ID, app_alias)
	}
}

// handle /listApps
func showAPPsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chat_id := update.Message.Chat.ID
	// get app config stats
	s := storage.GetAppsStat()
	if len(s) == 0 {
		// this is admin
		var t_msg string
		if config.AdminSet.Has(chat_id) {
			t_msg = "No APPs yet. Please create one."
		} else {
			t_msg = "No APPs yet. Please contact with Admins."
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   t_msg,
		})
		return
	}
	// loop through s, and send option to user to choose to delete it's account
	var t_msg string
	t_msg = "APP:"
	for _, v := range s {
		time.Sleep(BotSendMsgInterval)
		if config.AdminSet.Has(chat_id) {
			t_msg = fmt.Sprintf("%v\n    -  %v, Count: %v", t_msg, v.Alias, v.Bound)
		} else {
			t_msg = fmt.Sprintf("%v\n    -  %v", t_msg, v.Alias)
		}
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   t_msg,
	})
}

// hand "/listUsers"
func showBoundUsersHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chat_id := update.Message.Chat.ID
	// get user config stats
	var s []*storage.ResultStatsForNotify
	var e error
	// this is an admin
	if config.AdminSet.Has(chat_id) {
		s, e = storage.GetAllStatsWithAlias()
	} else {
		// non-admin
		s, e = storage.GetStatsByTgIDWithAlias(chat_id)
	}
	if len(s) == 0 || e != nil {
		t_msg := "No users yet. Please /bind one."
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat_id,
			Text:   t_msg,
		})
		return
	}
	// loop through s, make mapping
	var result_map map[string][]*storage.ResultStatsForNotify = make(map[string][]*storage.ResultStatsForNotify)
	for _, v := range s {
		if _, ok := result_map[v.AppAlias]; !ok {
			result_map[v.AppAlias] = make([]*storage.ResultStatsForNotify, 0)
		}
		result_map[v.AppAlias] = append(result_map[v.AppAlias], v)
	}
	var t_msg string
	t_msg = "APP:"
	for app_alias, v_1 := range result_map {
		t_msg = fmt.Sprintf("%v\n    -  %v", t_msg, app_alias)
		for _, v := range v_1 {
			t_msg = fmt.Sprintf("%v\n        -  %v", t_msg, v.UserAlias)
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   t_msg,
	})
}

// handle "/stat"
// for all users
func statHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	tg_id := update.Message.Chat.ID
	// get user's config
	results, e := storage.GetStatsByTgIDWithAlias(tg_id)
	if e != nil {
		logger.Errorf("<StatHandler> failed to get user's config by id %v, failed with: %v\n", tg_id, e.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: tg_id,
			Text:   "Failed to statistics. Please contact admin.",
		})
		return
	}
	handleSendStats(ctx, b, tg_id, results)
}

// handle "/statAll"
// admin only
func statAllHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	tg_id := update.Message.Chat.ID
	// check tg_id is not admin, just return
	if !config.AdminSet.Has(tg_id) {
		return
	}
	// get all OpStats
	results, e := storage.GetAllStatsWithAlias()
	if e != nil {
		logger.Errorf("<StatHandler> failed to get user's config by id %v, failed with: %v\n", tg_id, e.Error())
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: tg_id,
			Text:   "Failed to statistics. Please contact admin.",
		})
		return
	}
	handleSendStats(ctx, b, tg_id, results)
}

// handle send stats to user(tg_id)
func handleSendStats(ctx context.Context, b *bot.Bot, tg_id int64, results []*storage.ResultStatsForNotify) {
	if len(results) == 0 {
		logger.Debugf("<StatHandler> no account bound yet for %v\n", tg_id)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: tg_id,
			Text:   "no account bound yet.",
		})
		return
	}
	var results_map map[string]map[string][]*storage.ResultStatsForNotify = make(map[string]map[string][]*storage.ResultStatsForNotify)
	for _, v := range results {
		if results_map[v.AppAlias] == nil {
			results_map[v.AppAlias] = make(map[string][]*storage.ResultStatsForNotify, 0)
		}
		if results_map[v.AppAlias][v.UserAlias] == nil {
			results_map[v.AppAlias][v.UserAlias] = make([]*storage.ResultStatsForNotify, 0)
		}
		results_map[v.AppAlias][v.UserAlias] = append(results_map[v.AppAlias][v.UserAlias], v)
	}
	t_msg_builder := strings.Builder{}
	t_msg_builder.WriteString("Statistics:")
	for app_alias, v_1 := range results_map {
		t_t_msg := fmt.Sprintf("    -  %v", app_alias)
		t_msg_builder.WriteString(fmt.Sprintf("\n%v", t_t_msg))
		for user_alias, v_2 := range v_1 {
			t_t_msg := fmt.Sprintf("        -  %v", user_alias)
			t_msg_builder.WriteString(fmt.Sprintf("\n%v", t_t_msg))
			for _, v := range v_2 {
				t_time_str := utils.GetTimeString(v.LatestTime)
				t_t_msg := fmt.Sprintf("            -  %v: %v(s)/%v(f) - %v", v.OpAlias, v.Success, v.Failure, t_time_str)
				t_msg_builder.WriteString(fmt.Sprintf("\n%v", t_t_msg))
			}
		}
	}

	t_msg := t_msg_builder.String()
	t_msg_builder.Reset()
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: tg_id,
		Text:   t_msg,
	})
}

// hand "/help"
func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	tg_id := update.Message.Chat.ID
	var t_msg string
	if config.AdminSet.Has(tg_id) {
		t_msg = fmt.Sprintf("%s\n%s", config.Notice, HelpContentAdmin)
	} else {
		t_msg = fmt.Sprintf("%s\n%s", config.Notice, HelpContent)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: tg_id,
		Text:   t_msg,
	})
}
