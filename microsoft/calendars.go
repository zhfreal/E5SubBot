package microsoft

import (
	"fmt"
	"time"

	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/utils"
)

// GET /me/calendars
func listCalendars(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_url, t_err = genGraphApiUrl(map[string]any{}, "/me/calendars")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		return false
	}
	return true
}

// GET /me/calendar
func getMeCalendar(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_url, t_err = genGraphApiUrl(map[string]any{}, "/me/calendar")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		return false
	}
	return true
}

// GET /me/calendar/calendarView?startDateTime={start_datetime}&endDateTime={end_datetime}
func getCalendarView(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_start_at := utils.GetTimeWithDelta(30, 24, true)
	t_end_at := utils.GetTimeWithDelta(7, 24, false)
	t_query := map[string]any{
		"startDateTime": t_start_at.Format(time.RFC3339),
		"endDateTime":   t_end_at.Format(time.RFC3339),
	}
	t_url, t_err = genGraphApiUrl(t_query, "/me/calendar/calendarView")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		return false
	}
	return true
}

// GET /me/calendar/events
func listEvents(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_url, t_err = genGraphApiUrl(map[string]any{}, "/me/calendar/events")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		return false
	}
	return true
}

// GET https://graph.microsoft.com/v1.0/me/reminderView(startDateTime='2017-06-05T10:00:00.0000000',endDateTime='2017-06-11T11:00:00.0000000')
func listReminder(access_token, proxy *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_start_at := utils.GetTimeWithDelta(30, 24, true)
	t_end_at := utils.GetTimeWithDelta(7, 24, false)
	t_sub_url := fmt.Sprintf("/me/reminderView(startDateTime='%v',endDateTime='%v')",
		t_start_at.Format(time.RFC3339), t_end_at.Format(time.RFC3339))
	t_url, t_err = genGraphApiUrl(map[string]any{}, t_sub_url)
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		return false
	}
	return true
}

// POST https://graph.microsoft.com/v1.0/me/calendar/getSchedule
func getSchedule(access_token, proxy, mail *string) bool {
	var t_url string
	var t_err error
	var t_status_code int
	t_start_at := utils.GetTimeWithDelta(30, 24, true)
	t_end_at := utils.GetTimeWithDelta(7, 24, false)
	t_req_data := NewScheduleRequestString(mail, &t_start_at, &t_end_at)
	t_url, t_err = genGraphApiUrl(map[string]any{}, "/me/calendar/getSchedule")
	if t_err != nil {
		return false
	}
	t_status_code, _, t_err = performGraphApiPost(access_token, &t_url, t_req_data, proxy)
	if t_err != nil || t_status_code != 200 {
		return false
	}
	return true
}

// Worker to call listCalendars
func ListCalendars(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
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
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok := listCalendars(access_token, proxy)
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
		OpID:      OpTypeCalendarListCalendars,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// Worker to call getMeCalendar
func GetMeCalendar(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
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
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok := getMeCalendar(access_token, proxy)
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
		OpID:      OpTypeCalendarGetMeCalendar,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// worker to call getCalendarView
func GetCalendarView(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
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
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok := getCalendarView(access_token, proxy)
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
		OpID:      OpTypeCalendarGetCalendarView,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// worker to call listEvents
func ListEvents(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
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
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok := listEvents(access_token, proxy)
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
		OpID:      OpTypeCalendarListEvents,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// worker to call listReminder
func ListReminder(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
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
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok := listReminder(access_token, proxy)
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
		OpID:      OpTypeCalendarListReminders,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// worker to call getSchedule
func GetSchedule(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var id uint
	var access_token, mail *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgMailAddr] != nil {
		mail = args[ArgMailAddr].(*string)
	}
	if id > 0 && access_token != nil && mail != nil && len(*access_token) > 0 && len(*mail) > 0 {
		ok := getSchedule(access_token, proxy, mail)
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
		OpID:      OpTypeCalendarGetSchedule,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}
