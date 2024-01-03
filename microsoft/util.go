package microsoft

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func GetRegURLNew() string {
	return "https://aka.ms/appregistrations"
}

func getDCReqUrl(tenant string) string {
	return fmt.Sprintf("%v/%v", getAuthUrl(tenant), OAuthDC)
}

func getTokenUrl(tenant string) string {
	return fmt.Sprintf("%v/%v", getAuthUrl(tenant), OAuthToken)
}

func getAuthUrl(tenant string) string {
	return fmt.Sprintf("%v/%v/%v", AuthBase, tenant, OAuth)
}

func getMsgFoldersSubPath(folder_id, msg_id string) string {
	var sub_path string = "/me"
	if len(folder_id) > 0 {
		sub_path = fmt.Sprintf("%v/mailFolders/%v", sub_path, folder_id)
	}
	sub_path = fmt.Sprintf("%v/messages", sub_path)
	if len(msg_id) > 0 {
		sub_path = fmt.Sprintf("%v/%v", sub_path, msg_id)
	}
	return sub_path
}

func getGraphApiUrl(query map[string]string, paths ...string) (string, error) {
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
		queryParams.Add(k, v)
	}

	// Update the URL with the modified query parameters
	u.RawQuery = queryParams.Encode()
	return u.String(), nil
}

func Rand_Choice(choices []string) string {
	return choices[myRand.Intn(len(choices))]
}

func performGraphApiGet(access_token, url_str, proxy string) (string, error) {
	resp, err := performGraphApi(OpGet, access_token, url_str, "", proxy)
	if err != nil {
		return "", err
	}
	t_b, t_err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if t_err != nil {
		return "", fmt.Errorf("fail to fetch return content, failed with: %v", t_err.Error())
	}
	return string(t_b), nil
}

func performGraphApiPost(access_token, url_str, data, proxy string) (string, error) {
	resp, err := performGraphApi(OpPost, access_token, url_str, data, proxy)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		t_b, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("%v, %v", resp.Status, string(t_b))
		return string(t_b), err
	}
	t_b, t_err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if t_err != nil {
		return "", fmt.Errorf("fail to fetch return content, failed with: %v", t_err.Error())
	}
	return string(t_b), nil
}

// perform graph delete
func performGraphApiDelete(access_token, url_str, proxy string) (bool, error) {
	ok := false
	resp, err := performGraphApi(OpDelete, access_token, url_str, "", proxy)
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
func performGraphApi(action, access_token, url_str, data, proxy string) (*http.Response, error) {
	var req *http.Request
	var err error
	if action == OpPost {
		b := bytes.NewBuffer([]byte(data))
		req, err = http.NewRequest(action, url_str, b)
	} else {
		req, err = http.NewRequest(action, url_str, nil)
	}
	// fail to create request
	if err != nil {
		return &http.Response{}, err
	}
	// set authorization
	req.Header.Set("Authorization", access_token)
	req.Header.Set("Accept", "application/json")
	// set Content-Type if we post data
	if action == OpPost {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	// Create a new HTTP client and send the request
	client := &http.Client{}
	if proxy != "" {
		t_url_obj, _ := url.Parse(proxy)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(t_url_obj),
		}
	}
	return client.Do(req)
}

func WorkingOnMsFromChan(in chan Args, out chan ApiResult, done chan bool, wg *sync.WaitGroup, proxy string) {
	for {
		select {
		case args := <-in:
			args.Func(args.ID, args.AccessToken, out, proxy)
		case ok := <-done:
			if ok {
				wg.Done()
				return
			}
		default:
			time.Sleep(APIInterval)
		}
	}
}
