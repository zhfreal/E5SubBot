package microsoft

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
	Tenant         string        = "organizations"
	AuthBase       string        = "https://login.microsoftonline.com"
	OAuth          string        = "oauth2/v2.0"
	OAuthDC        string        = "devicecode"
	OAuthToken     string        = "token"
	GraphUrl       string        = "https://graph.microsoft.com"
	GraphVer       string        = "v1.0"
	OpGet          string        = "GET"
	OpPost         string        = "POST"
	OpDelete       string        = "DELETE"
	APIInterval    time.Duration = 50 * time.Millisecond // in milliseconds
	ReadMailsCount int           = 100
	ODataNextLink  string        = "@odata\\.nextLink" // using github.com/tidwall/gjson to search path, the "." must be escaped
	ODataContext   string        = "@odata\\.context"  // using github.com/tidwall/gjson to search path, the "." must be escaped
	// const_timeout_request_device_code        = 10 // timeout in seconds for request device code
)

const (
	OpTypeMail uint = iota + 1
)

var (
	scope_new = []string{
		"Calendars.ReadWrite",
		"Calendars.ReadWrite.Shared",
		"Contacts.ReadWrite",
		"Contacts.ReadWrite.Shared",
		"email",
		"Files.ReadWrite.All",
		"IMAP.AccessAsUser.All",
		"Mail.ReadWrite",
		"Mail.ReadWrite.Shared",
		"Mail.Send",
		"Mail.Send.Shared",
		"Notes.Create",
		"Notes.ReadWrite.All",
		"offline_access",
		"openid",
		"People.Read",
		"POP.AccessAsUser.All",
		"profile",
		"SMTP.Send",
		"Tasks.ReadWrite",
		"Tasks.ReadWrite.Shared",
		"User.ReadWrite",
	}
	myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	OPs    = map[uint]string{
		OpTypeMail: "Mails",
	}
	from = 0
	size = 20
)

type TokenCache struct {
	ClientID        string
	DeviceCode      string // device code in short format, used by user
	DeviceCodeDev   string // device code in long format, used by REST api
	AccountID       string
	Username        string
	Alias           string
	AccessToken     string
	RefreshToken    string
	ExpireTime      int64
	ExpireIn        int
	VerificationUrl string
	Message         string
	Interval        int   // interval for next device code
	NextReqTime     int64 // unix time for next request, equal to time.now().unix() + Interval
	Content         []byte
}

type ApiResult struct {
	ID        uint // UsersConfig's ID, indicate which user in storage
	OpID      uint
	S, F      int
	StartTime int64 // start time in unix time format
	Duration  int64 // Operation in millisecond
	EndTime   int64 // end time in unix time format
}

type Args struct {
	Func        func(id uint, access_token string, out chan ApiResult, proxy string)
	ID          uint
	AccessToken string
}

func detect_error(data []byte) []byte {
	e := ErrorREST{}
	err := json.Unmarshal(data, &e)
	if err != nil || len(e.Error) == 0 {
		return nil
	}
	return data
}

// Device Code json struct
type DeviceCodeREST struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUrl string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	Message         string `json:"message"`
}

// data json in:
//
//	{
//	    "token_type": "Bearer",
//	    "scope": "User.Read profile openid email",
//	    "expires_in": 3599,
//	    "access_token": ".......",
//	    "refresh_token": ".......",
//	}
type AuthStatusOfDCFromRest struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// define struct ErrorREST {} to store json data as:
//
//	{
//	   "error": "authorization_pending",
//	   "error_description": "........",
//	   "error_codes": [
//	       70016
//	   ],
//	   "timestamp": "2023-12-14 05:38:49Z",
//	   "trace_id": "447de602-50c7-42e9-94ab-d338e3c70300",
//	   "correlation_id": "3595c79e-1b4a-4a84-a1c9-bd8fd5d923e0",
//	   "error_uri": "https://login.microsoftonline.com/error?code=70016"
//	}
type ErrorREST struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCodes       []int  `json:"error_codes"`
	Timestamp        string `json:"timestamp"`
	TraceID          string `json:"trace_id"`
	CorrelationID    string `json:"correlation_id"`
	ErrorURI         string `json:"error_uri"`
}

// store json data as:
//
//	{
//	    "access_token": "........",
//	    "token_type": "Bearer",
//	    "expires_in": 3599,
//	    "scope": "https%3A%2F%2Fgraph.microsoft.com%2Fmail.read",
//	    "refresh_token": "........",
//	    "id_token": ".........",
//	}
type RefreshTokenREST struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

// ////////////////////////////////////////////////////////////
// store REST search json data
// define a struct that corresponds to the JSON structure
type ReqSearchData struct {
	EntityTypes []string `json:"entityTypes"`
	Query       struct {
		QueryString string `json:"queryString"`
	} `json:"query"`
	From int `json:"from"`
	Size int `json:"size"`
}

// define a struct that contains a slice of requests
type RequestSearchData struct {
	ReqDataList []ReqSearchData `json:"requests"`
}

/////////////////////////////////////////////////////////////

func NewRequestsData(keyword string, from, size int) *RequestSearchData {
	r := &RequestSearchData{}
	// assign values to its fields
	r.ReqDataList = []ReqSearchData{
		{
			EntityTypes: []string{"message"},
			Query: struct {
				QueryString string `json:"queryString"`
			}{
				QueryString: keyword,
			},
			From: from,
			Size: size,
		},
	}
	return r
}

func NewRequestsDataString(keyword string, from, size int) string {
	b_s, _ := json.Marshal(NewRequestsData(keyword, from, size))
	return string(b_s)
}
