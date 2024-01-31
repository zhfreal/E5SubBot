package microsoft

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetRegURLNew() string {
	return "https://aka.ms/appregistrations"
}

func getDCReqUrl(tenant *string) string {
	return fmt.Sprintf("%v/%v", getAuthUrl(tenant), OAuthDC)
}

func getTokenUrl(tenant *string) string {
	return fmt.Sprintf("%v/%v", getAuthUrl(tenant), OAuthToken)
}

func getAuthUrl(tenant *string) string {
	return fmt.Sprintf("%v/%v/%v", AuthBase, *tenant, OAuth)
}

func getMsgFoldersSubPath(folder_id, msg_id *string) string {
	var sub_path string = "/me"
	if folder_id != nil && len(*folder_id) > 0 {
		sub_path = fmt.Sprintf("%v/mailFolders/%v", sub_path, *folder_id)
	}
	sub_path = fmt.Sprintf("%v/messages", sub_path)
	if msg_id != nil && len(*msg_id) > 0 {
		sub_path = fmt.Sprintf("%v/%v", sub_path, *msg_id)
	}
	return sub_path
}

func genGraphApiUrl(query map[string]interface{}, paths ...string) (string, error) {
	// Construct the URL
	var t_url string
	var e error
	if len(paths) == 0 {
		return "", fmt.Errorf("paths is empty")
	}
	if strings.HasPrefix(paths[0], "http") {
		if len(paths) == 1 {
			t_url = paths[0]
		} else {
			t_url, e = url.JoinPath(paths[0], paths[1:]...)
		}
	} else {
		t_url, e = url.JoinPath(fmt.Sprintf("%s/%s", GraphUrl, GraphVer), paths...)
	}

	if e != nil {
		return "", e
	}

	u, err := url.Parse(t_url)
	if err != nil {
		return "", e
	}

	// Set query parameters using url.Values
	queryParams := u.Query()
	for k, v := range query {
		queryParams.Add(k, fmt.Sprintf("%v", v))
	}

	// Update the URL with the modified query parameters
	u.RawQuery = queryParams.Encode()
	// param := ""
	// for k, v := range query {
	// 	t_k := url.QueryEscape(fmt.Sprint(k))
	// 	t_v := url.QueryEscape(fmt.Sprint(v))
	// 	t_v = strings.Replace(t_v, "+", "%20", -1)
	// 	param = fmt.Sprintf("%v&%v=%v", param, t_k, t_v)
	// }

	// if len(param) > 0 {
	// 	param = strings.TrimLeft(param, "&")
	// 	t_url = fmt.Sprintf("%v?%v", t_url, param)
	// }
	return u.String(), nil
}

func Rand_Choice(choices []string) string {
	return choices[myRand.Intn(len(choices))]
}

// return status_code, body, error
func performGraphApiGet(access_token, url_str, proxy *string) (int, string, error) {
	resp, err := performGraphApi(&OpGet, access_token, url_str, nil, proxy)
	if err != nil {
		return -1, "", err
	}
	t_status_code := resp.StatusCode
	t_b, t_err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if t_err != nil {
		return t_status_code, "", fmt.Errorf("fail to fetch return content, failed with: %v", t_err.Error())
	}
	return t_status_code, string(t_b), nil
}

// return status_code, body, error
func performGraphApiPost(access_token, url_str, data, proxy *string) (int, string, error) {
	resp, err := performGraphApi(&OpPost, access_token, url_str, data, proxy)
	if err != nil {
		return -1, "", err
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		t_b, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("%v, %v", resp.Status, string(t_b))
		return -1, string(t_b), err
	}
	t_status_code := resp.StatusCode
	t_b, t_err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if t_err != nil {
		return t_status_code, "", fmt.Errorf("fail to fetch return content, failed with: %v", t_err.Error())
	}
	return t_status_code, string(t_b), nil
}

func performGraphApiPatch(access_token, url_str, data, proxy *string) (bool, error) {
	resp, err := performGraphApi(&OpPatch, access_token, url_str, data, proxy)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		t_b, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("%v, %v", resp.Status, string(t_b))
		return false, err
	}
	return true, nil
}

func performGraphApiPostSendMail(access_token, url_str, data, proxy *string) (bool, error) {
	resp, err := performGraphApi(&OpPost, access_token, url_str, data, proxy)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 202 {
		defer resp.Body.Close()
		t_b, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("%v, %v", resp.Status, string(t_b))
		return false, err
	}
	return true, nil
}

// perform graph delete
func performGraphApiDelete(access_token, url_str, proxy *string) (bool, error) {
	ok := false
	resp, err := performGraphApi(&OpDelete, access_token, url_str, nil, proxy)
	if err != nil {
		return ok, err
	}
	if resp.StatusCode != 204 && resp.StatusCode != 404 {
		defer resp.Body.Close()
		t_b, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("%v, %v", resp.Status, string(t_b))
	} else {
		// return 204, represent deletion successful
		ok = true
	}
	return ok, err
}

// support http GET, POST and DELETE, POST action support add post data
// POST data must be json string
// to void api failed we need set interval between two api requests.
func performGraphApi(action, access_token, url_str, data, proxy *string) (*http.Response, error) {
	var req *http.Request
	var err error
	if action != nil && (*action == OpPost || *action == OpPatch) {
		b := bytes.NewBuffer([]byte(*data))
		req, err = http.NewRequest(*action, *url_str, b)
	} else {
		req, err = http.NewRequest(*action, *url_str, nil)
	}
	// fail to create request
	if err != nil {
		return &http.Response{}, err
	}
	// set authorization
	req.Header.Set("Authorization", *access_token)
	req.Header.Set("Accept", "application/json")
	// set Content-Type if we post data
	if action != nil && (*action == OpPost || *action == OpPatch) {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	// Create a new HTTP client and send the request
	client := &http.Client{}
	if proxy != nil && *proxy != "" {
		t_url_obj, _ := url.Parse(*proxy)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(t_url_obj),
		}
	}
	return client.Do(req)
}
