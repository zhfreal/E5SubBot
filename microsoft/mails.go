package microsoft

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
	"github.com/zhfreal/E5SubBot/config"
)

var query_string = map[string]string{
	"$select": "sender,subject,from,toRecipients,ccRecipients,hasAttachments,isRead,isDraft",
}

var query_string_full = map[string]string{
	"$select": "sender,subject,from,body,flag,importance,toRecipients,ccRecipients,hasAttachments,isRead,isDraft",
}

// read one mail
func readOneMail(access_token, folder_id, msg_id, proxy string) (int, int) {
	var s, f int = 0, 0
	t_url, t_err := getGraphApiUrl(query_string, getMsgFoldersSubPath(folder_id, msg_id))
	if t_err != nil {
		return s, f
	}
	var content string
	content, t_err = performGraphApiGet(access_token, t_url, proxy)
	if t_err != nil {
		return s, f + 1
	}
	s += 1
	is_read := gjson.Get(content, "isRead").Bool()
	// get attachments
	if !is_read {
		// get attachments list
		t_url, t_err = getGraphApiUrl(query_string_full, getMsgFoldersSubPath(folder_id, msg_id), "/attachments")
		if t_err != nil {
			return s, f
		}
		time.Sleep(APIInterval)
		content, t_err = performGraphApiGet(access_token, t_url, proxy)
		if t_err != nil {
			return s, f + 1
		}
		s += 1
		t_results := gjson.GetMany(content, "value.#.id")
		t_id_list := t_results[0].Array()
		// get attachments content
		for _, t_a := range t_id_list {
			t_url, t_err = getGraphApiUrl(map[string]string{}, getMsgFoldersSubPath(folder_id, msg_id), "/attachments/", t_a.String())
			if t_err != nil {
				continue
			}
			time.Sleep(APIInterval)
			_, t_err = performGraphApiGet(access_token, t_url, proxy)
			if t_err != nil {
				f += 1
			} else {
				s += 1
			}
		}
	}
	return s, f
}

// get mails under folder, we get message from two APIs:
//
//	-- https://graph.microsoft.com/v1.0/me/messages/{id}; /me/messages/{id};
//	-- https://graph.microsoft.com/v1.0/me/mailFolders/{id}/messages/{msg_id}; /me/mailFolders/{id}/messages/{id};
//
//	"access_token" can't be empty
//	"folder_id" and "proxy" can be empty
//
//	count: specific count of mails to read, when count <=0, means read all mails
func readMailsFromFolder(access_token, folder_id string, count int, proxy string, read_latest bool) (int, int) {
	var content string
	var err error
	var s, f int = 0, 0
	param := map[string]string{}
	if read_latest {
		param["$orderby"] = "sentDateTime DESC"
	}
	t_url, t_err := getGraphApiUrl(param, getMsgFoldersSubPath(folder_id, ""))
	if t_err != nil {
		return s, f
	}
	t_fetched := 0
	for count <= 0 || t_fetched < count {
		content, err = performGraphApiGet(access_token, t_url, proxy)
		if err != nil {
			f += 1
			return s, f
		}
		s += 1
		// invalid response
		if gjson.Get(content, ODataContext).String() == "" {
			return s, f
		}
		results := gjson.GetMany(content, "value.#.id", "value.#.isRead")
		// if len(results) != 2 {
		// 	return s, f
		// }
		id_slice := results[0].Array()
		read_slice := results[1].Array()
		for i := 0; i < len(id_slice); i++ {
			t_msg_id := id_slice[i].String()
			t_is_read := read_slice[i].Bool()
			// does not read yet
			if !t_is_read {
				// read message after 100 milliseconds
				time.Sleep(APIInterval)
				t_s, t_f := readOneMail(access_token, folder_id, t_msg_id, proxy)
				s += t_s
				f += t_f
			}
		}
		t_fetched += len(id_slice)
		t_next_url := gjson.Get(content, ODataNextLink).String()
		// no more results
		if len(t_next_url) == 0 {
			break
		}
		t_url = t_next_url
		time.Sleep(APIInterval)
	}
	return s, f
}

// get all mailFolders and their child folders, and their mails
func readMailsFromAllFolders(access_token, proxy string) (int, int) {
	var content string
	var err error
	var s, f int = 0, 0
	t_folder_result, t_s, t_f, err := getAllMailFolders(access_token, proxy)
	s += t_s
	f += t_f
	if err != nil {
		return s, f
	}
	// has child folders
	folders := make(map[string]map[string]gjson.Result)
	t_id_list := make([]string, 0)
	for _, r := range t_folder_result {
		key := r.Get("id").String()
		v := r.Map()
		folders[key] = v
		t_id_list = append(t_id_list, key)
	}
	// loop to get all mails and sub folders
	for i := 0; i < len(t_id_list); i++ {
		time.Sleep(APIInterval)
		t_id := t_id_list[i]
		t_s, t_f := readMailsFromFolder(access_token, t_id, ReadMailsCount, proxy, true)
		s += t_s
		f += t_f
		// has child folders
		if folders[t_id]["childFolders"].Int() > 0 {
			// get child folders
			time.Sleep(APIInterval)
			url, _ := getGraphApiUrl(map[string]string{}, "/me/mailFolders", t_id, "/childFolders")
			content, err = performGraphApiGet(access_token, url, proxy)
			if err != nil {
				f += 1
				return s, f
			}
			s += 1
			result := gjson.Get(content, "value")
			if result.IsArray() {
				t_result_slice := result.Array()
				for _, r := range t_result_slice {
					key := r.Get("id").String()
					// child folders didn't add into to folders
					if _, ok := folders[key]; !ok {
						v := r.Map()
						folders[key] = v
						t_id_list = append(t_id_list, key)
					}
				}
			}
		}
	}

	return s, f
}

func deleteOneEmail(access_token, folder_id, msg_id, proxy string) (bool, error) {
	ok := false
	t_url, err := getGraphApiUrl(map[string]string{}, getMsgFoldersSubPath(folder_id, msg_id))
	if err != nil {
		return ok, fmt.Errorf("fail to generate url, failed with %v", err.Error())
	}
	ok, err = performGraphApiDelete(access_token, t_url, proxy)
	return ok, err
}

// graph REST API for search: https://graph.microsoft.com/v1.0/search/query
func searchEmailByKeyword(access_token, keyword string, from, size int, proxy string) (string, error) {
	var content string
	var err error
	t_url, err := getGraphApiUrl(map[string]string{}, "/search/query")
	if err != nil {
		return "", fmt.Errorf("fail to generate url, failed with %v", err.Error())
	}
	content, err = performGraphApiPost(access_token, t_url, NewRequestsDataString(keyword, from, size), proxy)
	if err != nil {
		return "", fmt.Errorf("fail to search email, failed with %v", err.Error())
	}
	return content, nil
}

// read filtered mails by a keyword in specific folder with it's folder_id
// count: specific count of mails to read, when count <=0, means read all mails
func readFilteredMails(folder_id, access_token, keyword string, count int, proxy string) ([]gjson.Result, int, int, error) {
	var content string
	var err error
	var s, f int = 0, 0
	var t_result_slice []gjson.Result
	t_filter_string := fmt.Sprintf("contains(body/content,'%v')", keyword)
	t_url, err := getGraphApiUrl(map[string]string{"filter": t_filter_string}, "me/mailFolders", folder_id, "messages")
	if err != nil {
		return t_result_slice, s, f, err
	}
	t_fetched := 0
	for count <= 0 || t_fetched < count {
		content, err = performGraphApiGet(access_token, t_url, proxy)
		if err != nil {
			f += 1
			break
		}
		s += 1
		t_values := gjson.Get(content, "value")
		if !t_values.IsArray() {
			err = fmt.Errorf("invalid response")
			break
		}
		t_values_slice := t_values.Array()
		t_fetched += len(t_values_slice)
		t_result_slice = append(t_result_slice, t_values_slice...)
		t_next_url := gjson.Get(content, ODataNextLink).String()
		// no more results
		if len(t_next_url) == 0 {
			break
		}
		t_url = t_next_url
		time.Sleep(APIInterval)
	}
	return t_result_slice, s, f, err
}

// get all mail folders
func getAllMailFolders(access_token, proxy string) ([]gjson.Result, int, int, error) {
	var t_result_slice []gjson.Result
	var s, f int = 0, 0
	var err error
	var content string
	t_url, _ := getGraphApiUrl(map[string]string{}, "/me/mailFolders")
	for {
		content, err = performGraphApiGet(access_token, t_url, proxy)
		if err != nil {
			f += 1
			break
		}
		s += 1
		t_values := gjson.Get(content, "value")
		if !t_values.IsArray() {
			err = fmt.Errorf("invalid response")
			break
		}
		t_result_slice = append(t_result_slice, t_values.Array()...)
		t_next_url := gjson.Get(content, ODataNextLink).String()
		// no more results
		if len(t_next_url) == 0 {
			break
		}
		t_url = t_next_url
	}
	return t_result_slice, s, f, err
}

// get the folder_id of "Inbox", "Sent Items", "Drafts"
func getFolderId(access_token, proxy string) ([]string, int, int, error) {
	var folder_id_list []string
	var s, f int = 0, 0
	t_folder_result, t_s, t_f, err := getAllMailFolders(access_token, proxy)
	s += t_s
	f += t_f
	if err != nil {
		return folder_id_list, s, f, err
	}
	for _, r := range t_folder_result {
		folder_name := r.Get("displayName").String()
		folder_id := r.Get("id").String()
		if folder_name == "Inbox" || folder_name == "Sent Items" || folder_name == "Drafts" {
			folder_id_list = append(folder_id_list, folder_id)
		}
	}
	return folder_id_list, s, f, nil
}

// delete specific mails by keywords under specific folders, each string in keywords will be a condition to query mails to delete
// we delete mails in "Inbox", "Sent Items", "Drafts"
func deleteOutlookMails(access_token string, keywords []string, quantity_for_delete int, proxy string) (int, int) {
	var ok bool
	var s, f int = 0, 0
	t_deleted := 0
	// get all folders, get the folder_id of "Inbox", "Sent Items", "Drafts"
	target_folder_list, t_s, t_f, err := getFolderId(access_token, proxy)
	s += t_s
	f += t_f
	if err != nil || len(target_folder_list) == 0 {
		return s, f
	}
	for _, folder_id := range target_folder_list {
		// loop keywords, search mails and delete mails
		for _, t_keywords := range keywords {
		OUTER_LOOP:
			for {
				t_result_list, t_s, t_f, err := readFilteredMails(folder_id, access_token, t_keywords, quantity_for_delete, proxy)
				s += t_s
				f += t_f
				// bad request, or network issue or other error, rather than empty search result
				if err != nil {
					break
				}
				if len(t_result_list) == 0 {
					break
				}

				for _, r := range t_result_list {
					t_msg_id := r.Get("id").String()
					t_is_read := r.Get("isRead").Bool()
					t_folder_id := r.Get("parentFolderId").String()
					if !t_is_read {
						// read this message
						time.Sleep(APIInterval)
						t_s, t_f := readOneMail(access_token, t_folder_id, t_msg_id, proxy)
						s += t_s
						f += t_f
					}
					// do deleteOneEmail
					time.Sleep(APIInterval)
					ok, err = deleteOneEmail(access_token, t_folder_id, t_msg_id, proxy)
					if err != nil || !ok {
						// fail to delete
						f += 1
					} else {
						// successfully delete
						s += 1
					}
					t_deleted += 1 // add 1 into t_deleted, mean we try to delete one mail
					if t_deleted >= quantity_for_delete {
						break OUTER_LOOP
					}
				}
			}
		}
	}

	return s, f
}

// use https://graph.microsoft.com/v1.0/search/query
func searchAndLoopMails(access_token string, keywords []string, fetch_quantity int, proxy string) (int, int) {
	var s, f int = 0, 0
	var fetched int
	t_from := from
	// loop keywords, search mails and delete mails
	for _, t_keywords := range keywords {
		for {
			content, err := searchEmailByKeyword(access_token, t_keywords, t_from, size, proxy)
			// bad request, or network issue or other error, rather than empty search result
			if err != nil {
				f += 1
				break
			}
			// get search results in path "content.value[0].hitsContainers[0]" use gjson
			t_search_content := gjson.Get(content, "value.0.hitsContainers.0").String()
			// no matched results, return as one successful try
			if t_search_content == "" {
				s += 1
				break
			}
			t_hits_string := gjson.Get(t_search_content, "hits").String()
			// t_total := gjson.Get(t_search_content, "total").Int()
			t_more_results_available := gjson.Get(t_search_content, "moreResultsAvailable").Bool()
			// empty search results, return as one successful try
			if len(t_hits_string) == 0 {
				s += 1
				break
			}
			t_msg_id_list := gjson.Get(t_hits_string, "#.hitId").Array()
			t_is_read_list := gjson.Get(t_hits_string, "#.resource.isRead").Array()
			t_folder_id_list := gjson.Get(t_hits_string, "#.resource.parentFolderId").Array()
			if len(t_msg_id_list) == 0 || len(t_folder_id_list) == 0 || len(t_is_read_list) == 0 {
				s += 1
				break
			}
			for i := 0; i < len(t_msg_id_list); i++ {
				t_msg_id := t_msg_id_list[i].String()
				t_is_read := t_is_read_list[i].Bool()
				t_folder_id := t_folder_id_list[i].String()
				if !t_is_read {
					// read this message if it is not read yet
					time.Sleep(APIInterval)
					t_s, t_f := readOneMail(access_token, t_folder_id, t_msg_id, proxy)
					s += t_s
					f += t_f
				}
			}
			fetched += len(t_msg_id_list)
			// no more results available or fetched enough mails, break
			if !t_more_results_available || fetched >= fetch_quantity {
				break
			}
			from += len(t_msg_id_list)
		}
	}
	return s, f
}

// get mailFolders delta
// TODO: read messages from delta
func getMailFoldersDelta(access_token, proxy string) (int, int) {
	t_url, _ := getGraphApiUrl(map[string]string{}, "/me/mailFolders/delta")
	_, err := performGraphApiGet(access_token, t_url, proxy)
	if err != nil {
		return 0, 1
	}
	return 1, 0
}

// list mails
func listMails(access_token, proxy string, count int) (int, int) {
	t_s, t_f := readMailsFromFolder(access_token, "", ReadMailsCount, proxy, true)
	return t_s, t_f
}

func WorkingOnMails(id uint, access_token string, out chan ApiResult, proxy string) {
	var s, f int = 0, 0
	t_start_at := time.Now()
	t_s, t_f := listMails(access_token, proxy, ReadMailsCount)
	s += t_s
	f += t_f
	time.Sleep(APIInterval)
	t_s, t_f = readMailsFromAllFolders(access_token, proxy)
	s += t_s
	f += t_f
	time.Sleep(APIInterval)
	t_s, t_f = getMailFoldersDelta(access_token, proxy)
	s += t_s
	f += t_f
	t_s, t_f = searchAndLoopMails(access_token, config.MailAutoDeleteKeyWords, config.MailAutoDeleteQuantity, proxy)
	s += t_s
	f += t_f
	// do deleteOutlookMails just according the config file
	if config.MailAutoDeleteEnabled {
		time.Sleep(APIInterval)
		t_s, t_f = deleteOutlookMails(access_token, config.MailAutoDeleteKeyWords, config.MailAutoDeleteQuantity, proxy)
		s += t_s
		f += t_f
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- ApiResult{
		ID:        id,
		OpID:      OpTypeMail,
		S:         s,
		F:         f,
		StartTime: t_start_at.Unix(),
		Duration:  t_durations_milliseconds,
		EndTime:   t_end_at.Unix(),
	}
}

func DoListAllMails(id uint, access_token string, out chan ApiResult, proxy string) {
	t_start_at := time.Now()
	t_s, t_f := readMailsFromFolder(access_token, "", ReadMailsCount, proxy, true)
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- ApiResult{
		ID:        id,
		OpID:      OpTypeMail,
		S:         t_s,
		F:         t_f,
		StartTime: t_start_at.Unix(),
		Duration:  t_durations_milliseconds,
		EndTime:   t_end_at.Unix(),
	}
}

func DoListAllMailFolders(id uint, access_token string, out chan ApiResult, proxy string) {
	t_start_at := time.Now()
	t_s, t_f := readMailsFromAllFolders(access_token, proxy)
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- ApiResult{
		ID:        id,
		OpID:      OpTypeMail,
		S:         t_s,
		F:         t_f,
		StartTime: t_start_at.Unix(),
		Duration:  t_durations_milliseconds,
		EndTime:   t_end_at.Unix(),
	}
}

func DoMailDeletion(id uint, access_token string, out chan ApiResult, proxy string) {
	t_start_at := time.Now()
	t_s, t_f := deleteOutlookMails(access_token, config.MailAutoDeleteKeyWords, config.MailAutoDeleteQuantity, proxy)
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- ApiResult{
		ID:        id,
		OpID:      OpTypeMail,
		S:         t_s,
		F:         t_f,
		StartTime: t_start_at.Unix(),
		Duration:  t_durations_milliseconds,
		EndTime:   t_end_at.Unix(),
	}
}
