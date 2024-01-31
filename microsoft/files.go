package microsoft

import (
	"time"

	"github.com/tidwall/gjson"
	"github.com/zhfreal/E5SubBot/config"
)

// GET /me/drive
func getMeDrive(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_url, t_err = genGraphApiUrl(map[string]any{}, "me/drive")
	if t_err != nil {
		return false
	}
	var content string
	t_status_code, content, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		return false
	}
	drive_id := gjson.Get(content, "id").String()
	return len(drive_id) > 0
}

// Get /me/drive/root/children
func getMeDriveRootChildren(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_url, t_err = genGraphApiUrl(map[string]any{}, "me/drive/root/children")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	return t_err == nil && t_status_code == 200
}

// GET /me/drive/recent
func getMeDriveRecent(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_url, t_err = genGraphApiUrl(map[string]any{}, "me/drive/recent")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	return t_err == nil && t_status_code == 200
}

// GET /me/drive/sharedWithMe
func getMeDriveSharedWithMe(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_url, t_err = genGraphApiUrl(map[string]any{}, "me/drive/sharedWithMe")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	return t_err == nil && t_status_code == 200
}

// GET /me/drive/following
func getMeDriveFollowing(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_url, t_err = genGraphApiUrl(map[string]any{}, "me/drive/following")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	return t_err == nil && t_status_code == 200
}

// Worker to call getMeDrive
func GetMeDrive(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if id > 0 && len(*access_token) > 0 {
		ok := getMeDrive(access_token, proxy)
		if ok {
			s++
		} else {
			f++
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeFileListFiles,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// Worker to call getMeDriveRootChildren
func GetMeDriveRootChildren(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if id > 0 && len(*access_token) > 0 {
		ok := getMeDriveRootChildren(access_token, proxy)
		if ok {
			s++
		} else {
			f++
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeFileListFiles,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// Worker to call getMeDriveRecent
func GetMeDriveRecent(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if id > 0 && len(*access_token) > 0 {
		ok := getMeDriveRecent(access_token, proxy)
		if ok {
			s++
		} else {
			f++
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeFileListFiles,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// worker to call getMeDriveSharedWithMe
func GetMeDriveSharedWithMe(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if id > 0 && len(*access_token) > 0 {
		ok := getMeDriveSharedWithMe(access_token, proxy)
		if ok {
			s++
		} else {
			f++
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeFileListFiles,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// worker to call getMeDriveFollowing
func GetMeDriveFollowing(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if id > 0 && len(*access_token) > 0 {
		ok := getMeDriveFollowing(access_token, proxy)
		if ok {
			s++
		} else {
			f++
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeFileListFiles,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}
