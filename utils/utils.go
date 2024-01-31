package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetURLValue(Url, key string) string {
	u, _ := url.Parse(Url)
	query := u.Query()
	return query.Get(key)
}

func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}

func getDateTimeByTime(t *time.Time) *time.Time {
	t_n := t.Truncate(24 * time.Hour)
	// print t_n as human readable date time
	// fmt.Print(t_n.Format("2006-01-02 15:04:05 -0700"))
	return &t_n
}

// get date string "2006-01-02" based on given time
func getTimeDateString(t *time.Time) string {
	return t.Format("2006-01-02")
}

// get hour:min:sec string "15:04:05" based on given time
func getTimeHMSString(t *time.Time) string {
	return t.Format("15:04:05")
}

// return date string "2006-01-02", if given unix time before today in localtime;
// return hour:min:sec string "15:04:05", if given unix time is today in localtime;
func GetTimeString(t int64) string {
	t_t := time.Unix(t, 0)
	t_t_d := *getDateTimeByTime(&t_t)
	t_s := time.Now()
	t_s_d := *getDateTimeByTime(&t_s)
	if t_t_d.Before(t_s_d) {
		return getTimeDateString(&t_t)
	} else {
		return getTimeHMSString(&t_t)
	}
}

// parse proxy string format, valid proxy like:
//
//	   "socks5://127.0.0.1:1080" or "http://127.0.0.1:8080" or "https://<domain or ip>:<port>"
//	   the leading schema must be socks5 or http or https, and must contain <domain or ip>:<port> in host port
//	use url.Parse to parse URI
//
// return bool, scheme, username, password, host, port, insecure
func ParseProxy(proxy string) (*url.URL, error) {
	// parse URI using url.Parse
	proxy = strings.ToLower(proxy)
	u, e := url.Parse(proxy)
	return u, e
	// // invalid URI, return false
	// if e != nil {
	//     return false, "", "", "", "", 0, false
	// }
	// // invalid schema, return false
	// if u.Scheme != "socks5" && u.Scheme != "http" && u.Scheme != "https" {
	//     return false, "", "", "", "", 0, false
	// }
	// // check host port, must contain <domain or ip>:<port> in host port
	// if u.Host == "" {
	//     return false, "", "", "", "", 0, false
	// }
	// t_r, t_p := parsePort(u.Host)
	// if !t_r {
	//     return false, "", "", "", "", 0, false
	// }
	// // check username and password
	// t_username, t_password := "", ""
	// if u.User == nil {
	//     t_username = u.User.Username()
	//     var t_r bool
	//     t_password, t_r = u.User.Password()
	//     if !t_r || len(t_password) == 0 {
	//         return false, "", "", "", "", 0, false
	//     }
	// }
	// insecure := false
	// if u.Scheme == "https" {
	//     q := u.Query()
	//     ins := q.Get("insecure")
	//     if ins != "" {
	//         insecure = checkBoolString(ins)
	//     }
	// }
	// return true, u.Scheme, t_username, t_password, u.Host, t_p, insecure
}

// give a port string is invalid port or not
// valid port must be in range of [1, 65535]
func ParsePort(port string) (bool, uint16) {
	p, e := strconv.Atoi(port)
	if e != nil {
		return false, 0
	}
	if p < 1 || p > 65535 {
		return false, 0
	}
	return true, uint16(p)
}

// check bool string to bool:
//
//	"Yes", "true", "1", to bool true
//	"No", "false", "0", to bool false
//	all string ignore upper or lower case
func CheckBoolString(b string) bool {
	b = strings.ToLower(b)
	if b == "yes" || b == "true" || b == "1" {
		return true
	}
	return false
}

// define a function that takes two values of any type and returns an int
func CompareNumbers(a, b interface{}) int {
	// use the reflect package to get the value and kind of the arguments
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	ka := va.Kind()
	kb := vb.Kind()

	// check if the arguments are numbers
	if ka >= reflect.Int && ka <= reflect.Float64 && kb >= reflect.Int && kb <= reflect.Float64 {
		// convert the arguments to float64 for comparison
		fa := va.Convert(reflect.TypeOf(float64(0))).Float()
		fb := vb.Convert(reflect.TypeOf(float64(0))).Float()

		// compare the numbers and return the result
		if fa > fb {
			return 1
		} else if fa == fb {
			return 0
		} else {
			return -1
		}
	} else {
		// return an error if the arguments are not numbers
		return -2
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// make folder, and make it's permission to 775
func MakeFolderAndSetPermission(path string) error {
	// check path exists or not
	// if not exists, make it and set permission to 775
	//    return error if any error shown
	// if exists, check it's a directory or not
	//    if it's a directory, set permission to 775
	//    if it's not a directory, return error
	//    return error if any error shown
	// path is not exists, make it
	folder, _ := filepath.Split(path)
	if folder == "" {
		folder = "./"
	}
	info, err := os.Stat(folder)
	if err == nil {
		if info.IsDir() {
			// path is a directory, set permission to 775
			err = os.Chmod(folder, 0o775)
			if err != nil {
				err = fmt.Errorf("%v is exists, but we failed to set it's permission to 0755, failed with: %v", folder, err.Error())
			}
			return err
		}
		// path is not a directory, return error
		return fmt.Errorf("%v is exists, but it's not a directory", folder)
	} else {
		err = os.MkdirAll(folder, 0o775)
		if err != nil {
			err = fmt.Errorf("%v is not exists, but we failed to make it, failed with: %v", folder, err.Error())
		}
		return err
	}
}

// split string whit whitespace indentation: \s+
// func SplitStringByWhiteSpaces(s string) []string {
// 	s = strings.TrimSpace(s)
// 	r := regexp.MustCompile(`\s+`)
// 	t_list := r.Split(s, -1)
// 	return t_list
// }

// split string whit whitespace indentation: \s+
func SplitStringByWhiteSpaces(s string) []string {
	return SplitStringByRegStr(s, "\\s+")
}

func SplitStringByRegStr(s, reg_str string) []string {
	s = strings.TrimSpace(s)
	r := regexp.MustCompile(reg_str)
	t_list := r.Split(s, -1)
	return t_list
}

// shuffle slice
func ShuffleSlice(s []*string) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i] // swap the elements
	})
}

func GetTimeWithDelta(day, hour int, negative bool) time.Time {
	t_day := day
	t_hour := hour
	if t_day <= 1 {
		t_day = 1
	}
	if t_day >= 30 {
		t_day = 30
	}
	if t_hour <= 1 {
		t_hour = 1
	}
	if t_hour >= 24 {
		t_hour = 24
	}
	t_rand_days := rand.Intn(day)
	t_rand_hours := rand.Intn(hour)
	if negative {
		t_rand_days = -t_rand_days
		t_rand_hours = -t_rand_hours
	}
	return time.Now().AddDate(0, 0, t_rand_days).Add(time.Duration(t_rand_hours) * time.Hour)
}
