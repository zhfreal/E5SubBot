package bots

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/zhfreal/E5SubBot/microsoft"
	"github.com/zhfreal/E5SubBot/storage"
)

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

const (
	ActionBindLegacy       string = "BL"
	ActionBindDeviceCode   string = "DC"
	ActionYes              string = "YES"
	ActionNo               string = "NO"
	ActionCancel           string = "CANCEL"
	ActionBindAccount      string = "BA"
	BindCacheTimeInSeconds int    = 30 * 60
	RefreshTokenBefore     int    = 60 * 5 // time in seconds, we refresh token in this second before it actually expire
	CMDHelp                string = "/help"
	CMDBindApp             string = "/bindApp"
	CMDBind                string = "/bind"
	CMDUnbind              string = "/unbind"
	CMDListApps            string = "/listApps"
	CMDListUsers           string = "/listUsers"
	CMDStat                string = "/stat"
	CMDReAuth              string = "/reAuth"
	CMDUnbindOther         string = "/unbindOther"
	CMDStatAll             string = "/statAll"
	CMDDelApp              string = "/delApp"
)

const (
	ReplyForBindAPP = iota
	ReplyForBindAccount
	ReplyForReAuth
	ReplyForWaitingAuth
	ReplyForDeleteAPP
	ReplyForUnbindAccountS1
	ReplyForUnbindAccountS2
)

const (
	ReplyFromCallBack int = iota
	ReplyWithPureMsg
	InitialMsg
)

var (
	WelcomeContent    string = "欢迎使用E5SubBot!"
	HelpContentHeader string = `
    命令：
    /bind        绑定帐号
    /unbind      解绑账户
    /reAuth      重新授权已绑定账户
    /stat        查看统计信息
    /listApps    查看已绑定应用
    /listUsers   查看已绑定用户
`
	HelpContentTail string = `
    /help        帮助
    源码及使用方法: https://github.com/zhfreal/E5SubBot
`
	HelpContent      string = HelpContentHeader + HelpContentTail
	HelpContentAdmin string = HelpContentHeader + `
    /bindApp     绑定应用(管理员)
    /delApp      删除应用(管理员)
    /unbindOther 解给其他用户(管理员)
    /statAll     统计所有用户(管理员)
    ` + HelpContentTail
)

var (
	// instance of github.com/go-telegram/bot.Bot
	botTelegram *bot.Bot
	// device code caches
	AuthCachedObj *AuthCache
	// bind message caches
	BindCachedObj *BindCache
	// locker for user account, run task and delete account must accquire this locker
	UsersConfigCacheObj *UsersConfigCache
	// cache for refresh failure counter
)

type AuthCache struct {
	cachedData map[int64]*CachedData
}

type CachedData struct {
	Locker     *sync.Mutex
	Locked     bool
	DeviceCode *PendingDeviceCode
	Cancel     func()
}

type PendingDeviceCode struct {
	ClientID   string
	DeviceCode string
	Msg        string
}

func (a *AuthCache) Lock(chat_id int64, dc *PendingDeviceCode, cancel func()) bool {
	if !a.has(chat_id) {
		a.add(chat_id, dc, cancel)
		a.cachedData[chat_id].Locker.Lock()
		a.cachedData[chat_id].Locked = true
		return true
	}
	if a.cachedData[chat_id].Locked {
		return false
	} else {
		a.cachedData[chat_id].Locker.Lock()
		a.cachedData[chat_id].Locked = true
		return true
	}
}

func (a *AuthCache) Unlock(chat_id int64) {
	a.unlock(chat_id)
}

func (a *AuthCache) unlock(chat_id int64) {
	if !a.has(chat_id) {
		a.add(chat_id, &PendingDeviceCode{}, nil)
		return
	}
	if a.cachedData[chat_id].Locked {
		a.cachedData[chat_id].Locker.Unlock()
		a.cachedData[chat_id].Locked = false
		a.cachedData[chat_id].DeviceCode = &PendingDeviceCode{}
	}
}

func (a *AuthCache) Cancel(chat_id int64) {
	if !a.has(chat_id) {
		a.add(chat_id, &PendingDeviceCode{}, nil)
	} else {
		a.cancel(chat_id)
	}
}

func (a *AuthCache) IsLocked(chat_id int64) bool {
	if !a.has(chat_id) {
		a.add(chat_id, &PendingDeviceCode{}, nil)
		return false
	}
	return a.cachedData[chat_id].Locked
}

func (a *AuthCache) CancelAll() {
	for k := range a.cachedData {
		a.cancel(k)
	}
}

func (a *AuthCache) AddCancelFunc(chat_id int64, cancel func()) {
	if !a.has(chat_id) {
		a.add(chat_id, &PendingDeviceCode{}, cancel)
	} else {
		a.cachedData[chat_id].Cancel = cancel
	}
}

func (a *AuthCache) Remove(chat_id int64) {
	a.remove(chat_id)
}

func (a *AuthCache) GetDeviceCode(chat_id int64) *PendingDeviceCode {
	if !a.has(chat_id) {
		return nil
	}
	return a.cachedData[chat_id].DeviceCode
}

func (a *AuthCache) has(chat_id int64) bool {
	_, ok := a.cachedData[chat_id]
	return ok
}

func (a *AuthCache) add(chat_id int64, dc *PendingDeviceCode, cancel func()) {
	if !a.has(chat_id) {
		a.cachedData[chat_id] = &CachedData{
			Locker:     &sync.Mutex{},
			Locked:     false,
			DeviceCode: dc,
			Cancel:     cancel,
		}
	}
}

func (a *AuthCache) cancel(chat_id int64) {
	if a.has(chat_id) && a.cachedData[chat_id].Locked {
		if a.cachedData[chat_id].Cancel != nil {
			a.cachedData[chat_id].Cancel()
		}
		a.unlock(chat_id)
	}
}

func (a *AuthCache) remove(chat_id int64) {
	if a.has(chat_id) {
		delete(a.cachedData, chat_id)
	}
}

func NewAuthCache() *AuthCache {
	return &AuthCache{
		cachedData: make(map[int64]*CachedData),
	}
}

type MsgKey struct {
	ChatID int64
	MsgID  int
}

func (m *MsgKey) String() string {
	return fmt.Sprintf("%d:%d", m.ChatID, m.MsgID)
}

type MsgValue struct {
	MsgType   int
	ExpiredAt time.Time
	Extra     *ExtraData
}

type ExtraData struct {
	ExtraData1String string
	ExtraData1Uint   uint
	ExtraData2String string
	ExtraData2Uint   uint
}

func (m *MsgValue) IsExpired() bool {
	return m.ExpiredAt.Before(time.Now())
}

type BindCache struct {
	MsgCached map[MsgKey]MsgValue
}

func NewBindCache() *BindCache {
	return &BindCache{
		MsgCached: make(map[MsgKey]MsgValue),
	}
}

func (b *BindCache) has(msg_cached *MsgKey) bool {
	_, ok := b.MsgCached[*msg_cached]
	return ok
}

func (b *BindCache) Has(msg_cached *MsgKey) bool {
	return b.has(msg_cached)
}

func (b *BindCache) Add(msg_cached *MsgKey, reply_type *MsgValue) bool {
	// check msg_cached exists or not
	if b.has(msg_cached) {
		return false
	}
	b.MsgCached[*msg_cached] = *reply_type
	return true
}

func (b *BindCache) FindMsgKeyByChatID(chat_id int64) []*MsgKey {
	// check msg_cached exists or not
	var keys []*MsgKey
	for k := range b.MsgCached {
		if k.ChatID == chat_id {
			keys = append(keys, &k)
		}
	}
	return keys
}

func (b *BindCache) FindMsgKeyByChatIDAndReplyType(chat_id int64, reply_type int) []*MsgKey {
	// check msg_cached exists or not
	var keys []*MsgKey
	for k, v := range b.MsgCached {
		if k.ChatID == chat_id && v.MsgType == reply_type {
			keys = append(keys, &k)
		}
	}
	return keys
}

func (b *BindCache) FindMsgKeyByChatIDAndReplyTypeExtraDataKey1(chat_id int64, reply_type int, extra_data_key_1 string) []*MsgKey {
	// check msg_cached exists or not
	var keys []*MsgKey
	for k, v := range b.MsgCached {
		if k.ChatID == chat_id && v.MsgType == reply_type && v.Extra.ExtraData1String == extra_data_key_1 {
			keys = append(keys, &k)
		}
	}
	return keys
}

func (b *BindCache) Del(msg_cached *MsgKey) bool {
	return b.del(msg_cached)
}

func (b *BindCache) del(msg_cached *MsgKey) bool {
	// check msg_cached exists or not
	if !b.has(msg_cached) {
		return false
	}
	delete(b.MsgCached, *msg_cached)
	return true
}

func (b *BindCache) Remove(msg_cached *MsgKey) bool {
	return b.del(msg_cached)
}

func (b *BindCache) Delete(msg_cached *MsgKey) bool {
	return b.del(msg_cached)
}

func (b *BindCache) Clear() {
	b.MsgCached = make(map[MsgKey]MsgValue)
}

func (b *BindCache) Get(msg_cached *MsgKey) *MsgValue {
	// check msg_cached exists or not
	if !b.has(msg_cached) {
		return nil
	}
	v := b.MsgCached[*msg_cached]
	return &v
}

// the key is the user's ID in UsersConfig table
type UsersConfigCache struct {
	lockerMappings map[uint]*usersData
}

type usersData struct {
	// locker           sync.Mutex
	// status           bool
	// clientId         string
	// msEmail          string
	// accessToken      string
	// refreshToken     string
	refreshFailCount int
	// expiresAt        time.Time
}

// new UsersConfigCache with all valid UsersConfigs stored in db
func NewUsersConfigCache() *UsersConfigCache {
	a, e := storage.GetAllUsers()
	if e != nil {
		fmt.Println("failed to get all users configs", e.Error())
		os.Exit(1)
	}
	u := &UsersConfigCache{
		lockerMappings: make(map[uint]*usersData, 0),
	}
	for _, v := range a {
		u.initFailCount(v.ID)
	}
	return u
}

func (u *UsersConfigCache) has(id uint) bool {
	_, ok := u.lockerMappings[id]
	return ok
}

// func (u *UsersConfigCache) remove(id uint) {
// 	if u.has(id) {
// 		u.unlock(id)
// 		delete(u.lockerMappings, id)
// 	}
// }

// func (u *UsersConfigCache) locked(id uint) bool {
// 	return u.has(id) && u.lockerMappings[id].status
// }

// func (u *UsersConfigCache) lock(id uint) {
// 	if !u.locked(id) {
// 		u.lockerMappings[id].locker.Lock()
// 		u.lockerMappings[id].status = true
// 	}
// }

// func (u *UsersConfigCache) unlock(id uint) {
// 	if u.locked(id) {
// 		u.lockerMappings[id].locker.Unlock()
// 		u.lockerMappings[id].status = false
// 	}
// }

func (u *UsersConfigCache) initFailCount(id uint) {
	u.lockerMappings[id] = &usersData{
		refreshFailCount: 0,
	}
}

func (u *UsersConfigCache) setFailCount(id uint) {
	u.lockerMappings[id].refreshFailCount++
}

func (u *UsersConfigCache) resetFailCount(id uint) {
	u.lockerMappings[id].refreshFailCount = 0
}

func (u *UsersConfigCache) delCache(id uint) {
	if !u.has(id) {
		return
	}
	delete(u.lockerMappings, id)
}

func (u *UsersConfigCache) getFailCount(id uint) int {
	return u.lockerMappings[id].refreshFailCount
}

func (u *UsersConfigCache) InitFailCount(id uint) {
	u.initFailCount(id)
}

func (u *UsersConfigCache) GetFailCount(id uint) int {
	return u.getFailCount(id)
}

func (u *UsersConfigCache) SetFailCount(id uint) {
	u.setFailCount(id)
}

func (u *UsersConfigCache) ResetFailCount(id uint) {
	u.resetFailCount(id)
}

func (u *UsersConfigCache) DelCache(id uint) {
	u.delCache(id)
}

// func (u *UsersConfigCache) GetEmail(id uint) string {
// 	return u.lockerMappings[id].msEmail
// }

// func (u *UsersConfigCache) Lock(id uint) {
// 	u.lock(id)
// }

// func (u *UsersConfigCache) Unlock(id uint) {
// 	u.unlock(id)
// }

// func (u *UsersConfigCache) Locked(id uint) bool {
// 	return u.locked(id)
// }

// func (u *UsersConfigCache) Add(uc *storage.UsersConfig) {
// 	u.addToken(uc)
// }

// func (u *UsersConfigCache) Remove(id uint) {
// 	u.remove(id)
// }

// func (u *UsersConfigCache) addToken(uc *storage.UsersConfig) {
// 	if !u.has(uc.ID) {
// 		u.lockerMappings[uc.ID] = &usersData{
// 			locker:           sync.Mutex{},
// 			status:           false,
// 			refreshFailCount: 0,
// 		}
// 	}
// 	u.lockerMappings[uc.ID].clientId = uc.MsId
// 	u.lockerMappings[uc.ID].msEmail = uc.MsUsername
// 	u.lockerMappings[uc.ID].accessToken = uc.AccessToken
// 	u.lockerMappings[uc.ID].refreshToken = uc.RefreshToken
// 	u.lockerMappings[uc.ID].expiresAt = time.Unix(uc.ExpiresAt, 0)
// }

// func (u *UsersConfigCache) updateToken(id uint, uc *storage.UsersConfig) {
// 	if !u.has(id) {
// 		return
// 	}
// 	u.lockerMappings[id].accessToken = uc.AccessToken
// 	u.lockerMappings[id].refreshToken = uc.RefreshToken
// 	u.lockerMappings[id].expiresAt = time.Unix(uc.ExpiresAt, 0)
// }

// func (u *UsersConfigCache) AddToken(uc *storage.UsersConfig) {
// 	u.addToken(uc)
// }

// func (u *UsersConfigCache) UpdateToken(id uint, uc *storage.UsersConfig) {
// 	u.updateToken(id, uc)
// }

// func (u *UsersConfigCache) refreshToken(id uint) (string, error) {
// 	t_access_token, t_refresh_token, t_expires_in, e := microsoft.RefreshToken(u.lockerMappings[id].clientId, u.lockerMappings[id].refreshToken)
// 	if e != nil {
// 		// add refreshFailed
// 		u.setFailCount(id)
// 		return "", e
// 	}
// 	// reset refreshFailed
// 	u.resetFailCount(id)
// 	expire_at := GetExpiredTimeFromNowAfter(t_expires_in)
// 	uc := storage.UsersConfig{
// 		AccessToken:  t_access_token,
// 		RefreshToken: t_refresh_token,
// 		ExpiresAt:    expire_at.Unix(),
// 	}
// 	u.updateToken(id, &uc)
// 	storage.UpdateUsersConfigTokens(id, &uc)
// 	return t_access_token, nil
// }

// // get token from cache, if it expired, then refresh it;
// //
// //	return access_token, refresh_failed_or_not, error
// func (u *UsersConfigCache) GetToken(id uint) (string, bool, error) {
// 	var access_token string
// 	is_refreshed_err := false
// 	var err error = nil
// 	// we don't cache this id
// 	if !u.has(id) {
// 		err = fmt.Errorf("user is not in cache")
// 	} else {
// 		// check access token expired or not, if expired, then refresh it;
// 		// we set expiresAt is sooner than actual expired value about <RefreshTokeBefore> seconds
// 		expiresAt := u.lockerMappings[id].expiresAt.Add(-time.Duration(RefreshTokenBefore) * time.Second)
// 		t_time_now := time.Now()
// 		if expiresAt.Before(t_time_now) || len(u.lockerMappings[id].refreshToken) == 0 {
// 			access_token, err = u.refreshToken(id)
// 			if err != nil {
// 				is_refreshed_err = true
// 			}
// 		} else {
// 			access_token = u.lockerMappings[id].accessToken
// 		}
// 	}
// 	return access_token, is_refreshed_err, err
// }

// // show token
// func (u *UsersConfigCache) ShowToken(email string) {
// 	var access_token string
// 	var err error = nil
// 	var uc_list []*storage.UsersConfig
// 	var e error
// 	if len(email) == 0 {
// 		uc_list, e = storage.GetAllUsersConfigs()
// 	} else {
// 		uc_list, e = storage.GetUsersConfigIDByEmail(email)
// 	}

// 	if e != nil {
// 		fmt.Print("Failed to get data from database. Please try again!")
// 		return
// 	}
// 	if len(uc_list) == 0 {
// 		fmt.Print("No data found in database!")
// 		return
// 	}
// 	for _, uc := range uc_list {
// 		id := uc.ID
// 		mail := uc.MsUsername
// 		// we don't cache this id
// 		if !u.has(id) {
// 			continue
// 		} else {
// 			// check access token expired or not, if expired, then refresh it;
// 			expiresAt := u.lockerMappings[id].expiresAt
// 			t_time_now := time.Now()
// 			if expiresAt.Before(t_time_now) || len(u.lockerMappings[id].refreshToken) == 0 {
// 				access_token, err = u.refreshToken(id)
// 				if err != nil {
// 					continue
// 				}
// 			} else {
// 				access_token = u.lockerMappings[id].accessToken
// 			}
// 			expiresAt_human := u.lockerMappings[id].expiresAt.Format("2006-01-02 T15:04:05 -0700")
// 			fmt.Printf("%s\t%s\t%s\n", mail, expiresAt_human, access_token)
// 		}
// 	}
// }

// independent function to get token from *storage.UsersConfig,
//
//	if it expired, then refresh it
//	    if it succeed to refresh , then store new access token,
//	       refresh token, and expiresAt in *storage.UsersConfig; and reset
//	       failure counter in UsersConfigCacheObj;
//	    if the refresh failed, then add failure counter into
//	        UsersConfigCacheObj, and return with error.
//	 if it's not expired, just return it.
//	return access_token, refresh_failed_or_not, error
func GetToken(uc *storage.Users) (string, bool, error) {
	var access_token, refresh_token string
	var expire_in int
	failed_to_refreshed := false
	var err error = nil
	// check access token expired or not, if expired, then refresh it;
	// we set expiresAt is sooner than actual expired time about <RefreshTokeBefore> seconds
	t_expire_at := uc.ExpiresAt
	expiresAt := time.Unix(t_expire_at, 0).Add(-time.Duration(RefreshTokenBefore) * time.Second)
	t_time_now := time.Now()
	if expiresAt.Before(t_time_now) || len(uc.RefreshToken) == 0 {
		access_token, refresh_token, expire_in, err = microsoft.RefreshToken(uc.MsId, uc.RefreshToken)
		if err != nil {
			// failed to refresh
			failed_to_refreshed = true
			// add refreshFailCount
			UsersConfigCacheObj.SetFailCount(uc.ID)
		} else {
			// succeed to refresh, then store new access token,
			// reset refreshFailCount after
			UsersConfigCacheObj.ResetFailCount(uc.ID)
			expire_at := GetExpiredTimeFromNowAfter(expire_in)
			uc.AccessToken = access_token
			uc.RefreshToken = refresh_token
			uc.ExpiresAt = expire_at.Unix()
			storage.UpdateUsersTokens(uc.ID, uc)
		}
	} else {
		access_token = uc.AccessToken
	}
	return access_token, failed_to_refreshed, err
}

// show token from *storage.UsersConfig, if it expired, then refresh it;
func ShowToken(account string) {
	var uc_list []*storage.Users
	var e error
	if len(account) == 0 {
		uc_list, e = storage.GetAllUsers()
	} else {
		uc_list, e = storage.GetUsersConfigIDByEmail(account)
	}

	if e != nil {
		fmt.Print("Failed to get data from database. Please try again!")
		return
	}
	if len(uc_list) == 0 {
		fmt.Print("No data found in database!")
		return
	}
	// check access token expired or not, if expired, then refresh it;
	for _, uc := range uc_list {
		var access_token, refresh_token string
		var expire_in int
		var err error = nil
		t_expire_at := uc.ExpiresAt
		expiresAt := time.Unix(t_expire_at, 0)
		t_time_now := time.Now()
		if expiresAt.Before(t_time_now) || len(uc.RefreshToken) == 0 {
			access_token, refresh_token, expire_in, err = microsoft.RefreshToken(uc.MsId, uc.RefreshToken)
			if err != nil {
				continue
			} else {
				// succeed to refresh, then store new access token,
				expire_at := GetExpiredTimeFromNowAfter(expire_in)
				uc.AccessToken = access_token
				uc.RefreshToken = refresh_token
				uc.ExpiresAt = expire_at.Unix()
				storage.UpdateUsersTokens(uc.ID, uc)
			}
		}
		access_token = uc.AccessToken
		email := uc.MsUsername
		expiresAt_human := time.Unix(uc.ExpiresAt, 0).Format("2006-01-02 T15:04:05 -0700")
		fmt.Printf("%s\t%s\t%s\n", email, expiresAt_human, access_token)
	}
}

func Init() {
	AuthCachedObj = NewAuthCache()
	BindCachedObj = NewBindCache()
	UsersConfigCacheObj = NewUsersConfigCache()
}
