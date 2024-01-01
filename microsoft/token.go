package microsoft

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// get device code by REST API
func GetDeviceCode(ctx context.Context, client_id string) (*TokenCache, error) {
	// make a cache to store token
	cache := &TokenCache{ClientID: client_id}

	b, e := getDeviceCode(Tenant, client_id, scope_new)
	if e != nil {
		return nil, e
	}
	if e_b := detect_error(b); e_b != nil {
		return nil, fmt.Errorf("failed to get device code, failed with: %v", string(e_b))
	}
	t_device_code_rest := DeviceCodeREST{}
	// parse DeviceCodeREST data in string which store in content into t_device_code_rest
	if e = json.Unmarshal(b, &t_device_code_rest); e != nil {
		return nil, e
	}
	t_device_code := t_device_code_rest.DeviceCode
	t_user_code := t_device_code_rest.UserCode
	t_verification_url := t_device_code_rest.VerificationUrl
	t_expires_in := t_device_code_rest.ExpiresIn
	// we set ExpireTime in 10 seconds before <time now> + <expireed_in>
	t_expires_time := time.Now().Unix() + int64(t_expires_in) - 10
	t_interval := t_device_code_rest.Interval
	t_next_req_time := time.Now().Unix() + int64(t_interval)
	t_message := t_device_code_rest.Message
	if len(t_device_code) == 0 || len(t_user_code) == 0 || len(t_verification_url) == 0 || t_expires_in == 0 || t_interval == 0 {
		return nil, fmt.Errorf("wrong device code, content: %s", string(b))
	}
	cache.DeviceCode = t_user_code
	cache.DeviceCodeDev = t_device_code
	cache.VerificationUrl = t_verification_url
	cache.ExpireIn = t_expires_in
	cache.ExpireTime = t_expires_time
	cache.Interval = t_interval
	cache.NextReqTime = t_next_req_time
	cache.Message = t_message
	// convert content into []byte, then store in cache.Content
	cache.Content = b
	return cache, e
}

// check authorization status of device code by REST API
// we get user info from https://graph.microsoft.com/v1.0/me, to get alias and mail address
func CheckAuthStatusOfDeviceCode(ctx context.Context, token_cache *TokenCache) error {
	for {
		time.Sleep(500 * time.Millisecond)
		b, e := authDeviceCode(Tenant, token_cache.ClientID, token_cache.DeviceCodeDev)
		if e != nil {
			continue
		}
		t_error_obj := ErrorREST{}
		err := json.Unmarshal(b, &t_error_obj)
		if err != nil || t_error_obj.Error == "" {
			// successful authorization
			t_auth_result := &AuthStatusOfDCFromRest{}
			if e = json.Unmarshal(b, &t_auth_result); e != nil {
				return fmt.Errorf("failed to parse authorization result")
			}
			token_cache.AccessToken = t_auth_result.AccessToken
			token_cache.RefreshToken = t_auth_result.RefreshToken
			token_cache.ExpireIn = t_auth_result.ExpiresIn
			token_cache.ExpireTime = time.Now().Add(time.Second * time.Duration(t_auth_result.ExpiresIn)).Unix()
			// get user
			t_alias, t_mail := GetMeInfo(token_cache.AccessToken, "")
			token_cache.Username = t_mail
			token_cache.Alias = t_alias
			break
		} else {
			// return a error json string
			t_err_str := t_error_obj.Error
			// return error with authorization_pending
			if t_err_str == "authorization_pending" {
				continue
			} else if t_err_str == "authorization_declined" || t_err_str == "bad_verification_code" || t_err_str == "expired_token" {
				// authorization declined or bad verification code or expired token
				return fmt.Errorf("got error while waiting for user to authenticate this device, failed with: %v", t_err_str)
			} else {
				return fmt.Errorf("got error while waiting for user to authenticate this device, failed with: %v, %v", t_err_str, t_error_obj.ErrorDescription)
			}
		}

	}
	return nil
}

// fresh access token by client_id and refresh token
// return access_token, refresh_token, expire_id, and error if there is any error.
func RefreshToken(client_id, refresh string) (string, string, int, error) {
	t_b, err := refreshToken(Tenant, client_id, refresh, scope_new)
	if err != nil {
		return "", "", -1, err
	}
	if err_bytes := detect_error(t_b); err_bytes != nil {
		return "", "", -1, fmt.Errorf("failed to refresh token, failed with: %v", string(t_b))
	}
	t_refresh_token := RefreshTokenREST{}
	err = json.Unmarshal(t_b, &t_refresh_token)
	if err != nil {
		return "", "", -1, fmt.Errorf("wrong token type, content: %s", string(t_b))
	}
	return t_refresh_token.AccessToken, t_refresh_token.RefreshToken, t_refresh_token.ExpiresIn, nil
}

// get device code
func getDeviceCode(tenant, clientID string, scopes []string) ([]byte, error) {
	// Create a URL-encoded form with the required parameters
	form := url.Values{}
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	form.Add("client_id", clientID)
	form.Add("scope", strings.Join(scopes, " "))

	// Create a new HTTP request with the form as the body
	req, err := http.NewRequest("POST", getDCReqUrl(tenant), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	return performHttpRequest(req)
}

// authenticate by device code
func authDeviceCode(tenant, clientID, deviceCode string) ([]byte, error) {
	// Create a URL-encoded form with the required parameters
	form := url.Values{}
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	form.Add("client_id", clientID)
	form.Add("device_code", deviceCode)

	// Create a new HTTP request with the form as the body
	req, err := http.NewRequest("POST", getTokenUrl(tenant), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	return performHttpRequest(req)
}

// refresh token
func refreshToken(tenant, clientID, refreshToken string, scopes []string) ([]byte, error) {
	// Create a URL-encoded form with the required parameters
	form := url.Values{}
	form.Add("client_id", clientID)
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	t_scopes_str := ""
	for _, v := range scopes {
		t_scopes_str = fmt.Sprintf("%v %v/%v", t_scopes_str, GraphUrl, v)
	}
	t_scopes_str = strings.TrimSpace(t_scopes_str)
	form.Add("scope", t_scopes_str)

	// Create a new HTTP request with the form as the body
	req, err := http.NewRequest("POST", getTokenUrl(tenant), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	return performHttpRequest(req)
}

func performHttpRequest(req *http.Request) ([]byte, error) {
	// Set the content type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Close the response body when the function returns
	defer resp.Body.Close()

	// Read the response body as a byte slice
	return io.ReadAll(resp.Body)
}
