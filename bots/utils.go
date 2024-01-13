package bots

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/zhfreal/E5SubBot/microsoft"
	"github.com/zhfreal/E5SubBot/storage"
)

// get token
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
		access_token, refresh_token, expire_in, err = microsoft.RefreshToken(context.Background(), &uc.MsId, &uc.RefreshToken)
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
	// print title
	fmt.Printf("%s\t%s\t%s\n", "Email", "ExpiresAt", "AccessToken")
	// check access token expired or not, if expired, then refresh it;
	for _, uc := range uc_list {
		var access_token, refresh_token string
		var expire_in int
		var err error = nil
		t_expire_at := uc.ExpiresAt
		expiresAt := time.Unix(t_expire_at, 0)
		t_time_now := time.Now()
		if expiresAt.Before(t_time_now) || len(uc.RefreshToken) == 0 {
			access_token, refresh_token, expire_in, err = microsoft.RefreshToken(context.Background(), &uc.MsId, &uc.RefreshToken)
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

// summarized: to sum level-2 test into level-1, eg. add MailsRead stats into Mails
func GetAppsStats(s []*storage.StatsForShow, summarized bool) []*storage.AppsStats {
	t_apps := make([]*storage.AppsStats, 0)
	result_map := make(map[string]map[string]map[string][]*storage.BasicOpsStats)
	for _, result := range s {
		app_alias := result.AppAlias
		user_alias := result.UserAlias
		op_alias := result.OpAlias
		op_id := result.OpId
		if result_map[app_alias] == nil {
			result_map[app_alias] = make(map[string]map[string][]*storage.BasicOpsStats)
		}
		if result_map[app_alias][user_alias] == nil {
			result_map[app_alias][user_alias] = make(map[string][]*storage.BasicOpsStats)
		}
		// reset op_id, op_alias to level-1's
		// current we just have mails options
		if summarized {
			if op_id >= microsoft.OpTypeMail && op_id < microsoft.OpTypeFile {
				op_id = microsoft.OpTypeMail
				op_alias = microsoft.Ops[op_id]
			}
		}
		if result_map[app_alias][user_alias][op_alias] == nil {
			result_map[app_alias][user_alias][op_alias] = []*storage.BasicOpsStats{}
		}
		result_map[app_alias][user_alias][op_alias] = append(result_map[app_alias][user_alias][op_alias], &storage.BasicOpsStats{
			OpAlias:    op_alias,
			OpId:       op_id,
			TgId:       result.TgId,
			Success:    result.Success,
			Failure:    result.Failure,
			LatestTime: result.LatestTime,
		})
	}
	for k1, v1 := range result_map {
		t_app := &storage.AppsStats{
			AppAlias: k1,
		}
		for k2, v2 := range v1 {
			t_user := &storage.UsersStatsData{
				UserAlias: k2,
				OpsStats:  []*storage.BasicOpsStats{},
			}
			for _, v3 := range v2 {
				tmp := &storage.BasicOpsStats{
					OpId:       v3[0].OpId,
					OpAlias:    v3[0].OpAlias,
					TgId:       v3[0].TgId,
					Success:    v3[0].Success,
					Failure:    v3[0].Failure,
					LatestTime: v3[0].LatestTime,
				}
				// summarize level-2 test into level-1
				if len(v3) > 1 {
					for _, v4 := range v3[1:] {
						tmp.Success += v4.Success
						tmp.Failure += v4.Failure
						if v4.LatestTime > tmp.LatestTime {
							tmp.LatestTime = v4.LatestTime
						}
					}
				}
				t_user.OpsStats = append(t_user.OpsStats, tmp)
				sort.Slice(t_user.OpsStats, func(i, j int) bool {
					if t_user.OpsStats[i].OpAlias == t_user.OpsStats[j].OpAlias {
						return t_user.OpsStats[i].TgId < t_user.OpsStats[j].TgId
					}
					return t_user.OpsStats[i].OpAlias < t_user.OpsStats[j].OpAlias
				})
			}
			t_app.UsersStats = append(t_app.UsersStats, t_user)
			sort.Slice(t_app.UsersStats, func(i, j int) bool {
				return t_app.UsersStats[i].UserAlias < t_app.UsersStats[j].UserAlias
			})
		}
		t_apps = append(t_apps, t_app)
	}
	sort.Slice(t_apps, func(i, j int) bool {
		return t_apps[i].AppAlias < t_apps[j].AppAlias
	})
	return t_apps
}

func GetAppsStatsByTgId(tg_id int64) ([]*storage.AppsStats, error) {
	t, e := storage.GetStatsForShowByTgId(tg_id)
	if e != nil {
		return nil, e
	}
	return GetAppsStats(t, true), nil
}

func GetAllAppsStats() ([]*storage.AppsStats, error) {
	t, e := storage.GetAllStatsForShow()
	if e != nil {
		return nil, e
	}
	return GetAppsStats(t, true), nil
}
