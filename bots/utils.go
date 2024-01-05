package bots

import (
	"context"
	"fmt"
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
		access_token, refresh_token, expire_in, err = microsoft.RefreshToken(context.Background(), uc.MsId, uc.RefreshToken)
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
			access_token, refresh_token, expire_in, err = microsoft.RefreshToken(context.Background(), uc.MsId, uc.RefreshToken)
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
