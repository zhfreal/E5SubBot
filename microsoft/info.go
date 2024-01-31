package microsoft

import (
	"github.com/tidwall/gjson"
)

// get user info from https://graph.microsoft.com/v1.0/me
// return displayName, mail
func GetMeInfo(access_token, proxy *string) (string, string) {
	t_url, _ := genGraphApiUrl(map[string]any{}, "/me")
	t_status_code, t_content, err := performGraphApiGet(access_token, &t_url, proxy)
	if err != nil || t_status_code != 200 {
		return "", ""
	}
	displayName := gjson.Get(t_content, "displayName").String()
	mail := gjson.Get(t_content, "mail").String()
	return displayName, mail
}
