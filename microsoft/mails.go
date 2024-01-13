package microsoft

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
	"github.com/zhfreal/E5SubBot/config"
)

const (
	MailTemplate string = `<!DOCTYPE html>
<html lang="en">

<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">

  <title>ieSoft Best quote</title>
  <style media="all" type="text/css">
    /* -------------------------------------
    GLOBAL RESETS
------------------------------------- */

    body {
      font-family: Helvetica, sans-serif;
      -webkit-font-smoothing: antialiased;
      font-size: 16px;
      line-height: 1.3;
      -ms-text-size-adjust: 100%;
      -webkit-text-size-adjust: 100%;
    }

    table {
      border-collapse: separate;
      mso-table-lspace: 0pt;
      mso-table-rspace: 0pt;
      width: 100%;
    }

    table td {
      font-family: Helvetica, sans-serif;
      font-size: 16px;
      vertical-align: top;
    }

    /* -------------------------------------
    BODY & CONTAINER
------------------------------------- */

    body {
      background-color: #f4f5f6;
      margin: 0;
      padding: 0;
    }

    .body {
      background-color: #f4f5f6;
      width: 100%;
    }

    .container {
      margin: 0 auto !important;
      max-width: 600px;
      padding: 0;
      padding-top: 24px;
      width: 600px;
    }

    .content {
      box-sizing: border-box;
      display: block;
      margin: 0 auto;
      max-width: 600px;
      padding: 0;
    }

    /* -------------------------------------
    HEADER, FOOTER, MAIN
------------------------------------- */

    .main {
      background: #ffffff;
      border: 1px solid #eaebed;
      border-radius: 16px;
      width: 100%;
    }

    .wrapper {
      box-sizing: border-box;
      padding: 24px;
    }

    .footer {
      clear: both;
      padding-top: 24px;
      text-align: center;
      width: 100%;
    }

    .footer td,
    .footer p,
    .footer span,
    .footer a {
      color: #9a9ea6;
      font-size: 16px;
      text-align: center;
    }

    /* -------------------------------------
    TYPOGRAPHY
------------------------------------- */

    p {
      font-family: Helvetica, sans-serif;
      font-size: 16px;
      font-weight: normal;
      margin: 0;
      margin-bottom: 16px;
    }

    a {
      color: #0867ec;
      text-decoration: underline;
    }

    /* -------------------------------------
    BUTTONS
------------------------------------- */

    .btn {
      box-sizing: border-box;
      min-width: 100% !important;
      width: 100%;
    }

    .btn>tbody>tr>td {
      padding-bottom: 16px;
    }

    .btn table {
      width: auto;
    }

    .btn table td {
      background-color: #ffffff;
      border-radius: 4px;
      text-align: center;
    }

    .btn a {
      background-color: #ffffff;
      border: solid 2px #0867ec;
      border-radius: 4px;
      box-sizing: border-box;
      color: #0867ec;
      cursor: pointer;
      display: inline-block;
      font-size: 16px;
      font-weight: bold;
      margin: 0;
      padding: 12px 24px;
      text-decoration: none;
      text-transform: capitalize;
    }

    .btn-primary table td {
      background-color: #0867ec;
    }

    .btn-primary a {
      background-color: #0867ec;
      border-color: #0867ec;
      color: #ffffff;
    }

    @media all {
      .btn-primary table td:hover {
        background-color: #ec0867 !important;
      }

      .btn-primary a:hover {
        background-color: #ec0867 !important;
        border-color: #ec0867 !important;
      }
    }

    /* -------------------------------------
    OTHER STYLES THAT MIGHT BE USEFUL
------------------------------------- */

    .last {
      margin-bottom: 0;
    }

    .first {
      margin-top: 0;
    }

    .align-center {
      text-align: center;
    }

    .align-right {
      text-align: right;
    }

    .align-left {
      text-align: left;
    }

    .text-link {
      color: #0867ec !important;
      text-decoration: underline !important;
    }

    .clear {
      clear: both;
    }

    .mt0 {
      margin-top: 0;
    }

    .mb0 {
      margin-bottom: 0;
    }

    .preheader {
      color: transparent;
      display: none;
      height: 0;
      max-height: 0;
      max-width: 0;
      opacity: 0;
      overflow: hidden;
      mso-hide: all;
      visibility: hidden;
      width: 0;
    }

    .powered-by a {
      text-decoration: none;
    }

    /* -------------------------------------
    RESPONSIVE AND MOBILE FRIENDLY STYLES
------------------------------------- */

    @media only screen and (max-width: 640px) {

      .main p,
      .main td,
      .main span {
        font-size: 16px !important;
      }

      .wrapper {
        padding: 8px !important;
      }

      .content {
        padding: 0 !important;
      }

      .container {
        padding: 0 !important;
        padding-top: 8px !important;
        width: 100% !important;
      }

      .main {
        border-left-width: 0 !important;
        border-radius: 0 !important;
        border-right-width: 0 !important;
      }

      .btn table {
        max-width: 100% !important;
        width: 100% !important;
      }

      .btn a {
        font-size: 16px !important;
        max-width: 100% !important;
        width: 100% !important;
      }
    }

    /* -------------------------------------
    PRESERVE THESE STYLES IN THE HEAD
------------------------------------- */

    @media all {
      .ExternalClass {
        width: 100%;
      }

      .ExternalClass,
      .ExternalClass p,
      .ExternalClass span,
      .ExternalClass font,
      .ExternalClass td,
      .ExternalClass div {
        line-height: 100%;
      }

      .apple-link a {
        color: inherit !important;
        font-family: inherit !important;
        font-size: inherit !important;
        font-weight: inherit !important;
        line-height: inherit !important;
        text-decoration: none !important;
      }

      #MessageViewBody a {
        color: inherit;
        text-decoration: none;
        font-size: inherit;
        font-family: inherit;
        font-weight: inherit;
        line-height: inherit;
      }
    }
  </style>
</head>

<body>
  <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body">
    <tbody>
      <tr>
        <td>&nbsp;</td>
        <td class="container">
          <div class="content">
            <!-- START CENTERED WHITE CONTAINER -->
            <span class="preheader">ieSoft Best quote</span>
            <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="main">

              <!-- START MAIN CONTENT AREA -->
              <tbody>
                <tr>
                  <td class="wrapper">
                    <p>Hi there</p>
                    <p>We have just celebrated a series of holidays, and now returning to everyday work might seem
                      challenging. But don't worry, we have a couple of tips to help you! </p>
                    <p>&nbsp;&nbsp;-&nbsp;Start with easy tasks and gradually move to more complex ones!</p>
                    <p>&nbsp;&nbsp;-&nbsp;Remember what makes your work so special and try to focus on the tasks you
                      enjoy the most! </p>
                    <p>&nbsp;&nbsp;-&nbsp;And don't forget to add a little festive atmosphere to your workday!</p>
                    <p>&nbsp;&nbsp;-&nbsp;After the holidays, pay attention to proper sorting of accumulated emails to
                      avoid losing important messages in spam or among deferred ones. </p>
                    <p>We, the team at ieSoft Group, want to wish you an easy adjustment after the holiday break.
                      Remember, we are always here to make your experience working with us as comfortable as possible.
                    </p>
                    <p>================<br><a href="https://www.zhfreal.top"
                        target="_blank">www.zhfreal.top</a><br>ieSoft LTD.<br>================</p>
                  </td>
                </tr>
                <tr>
                  <td style="text-align: center;" valign="top" class="m_-1042650705938460172footerContent">
                    <a href="https://cloud.zhfreal.top" target="_blank">visit our website</a>
                    <span class="m_-1042650705938460172hide-mobile"> | </span>
                    <a href="https://cloud.zhfreal.top/" target="_blank">log in to your account</a>
                    <span class="m_-1042650705938460172hide-mobile"> | </span>
                    <a href="https://cloud.zhfreal.top/submitticket.php" target="_blank">get support</a> <br>
                    Copyright © <a href="http://zhfreal.top" target="_blank">ieSoft LTD.</a>, All rights reserved.
                    <br> <br>
                  </td>
                </tr>
                <!-- END MAIN CONTENT AREA -->
              </tbody>
            </table>
            <!-- END CENTERED WHITE CONTAINER -->
          </div>
        </td>
        <td>&nbsp;</td>
      </tr>
    </tbody>
  </table>
</body>

</html>
`
)

var query_string = map[string]any{
	"$select": "sender,subject,from,toRecipients,ccRecipients,hasAttachments,isRead,isDraft",
}

var query_string_full = map[string]any{
	"$select": "sender,subject,from,body,flag,importance,toRecipients,ccRecipients,hasAttachments,isRead,isDraft",
}

// read one mail
func readOneMail(access_token, folder_id, msg_id, proxy *string) (int, int) {
	var s, f int = 0, 0
	t_url, t_err := genGraphApiUrl(query_string, getMsgFoldersSubPath(folder_id, msg_id))
	if t_err != nil {
		return s, f
	}
	var content string
	content, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil {
		return s, f + 1
	}
	s += 1
	is_read := gjson.Get(content, "isRead").Bool()
	// get attachments
	if !is_read {
		// get attachments list
		t_url, t_err = genGraphApiUrl(query_string_full, getMsgFoldersSubPath(folder_id, msg_id), "/attachments")
		if t_err != nil {
			return s, f
		}
		// time.Sleep(APIInterval)
		content, t_err = performGraphApiGet(access_token, &t_url, proxy)
		if t_err != nil {
			return s, f + 1
		}
		s += 1
		t_results := gjson.GetMany(content, "value.#.id")
		t_id_list := t_results[0].Array()
		// get attachments content
		for _, t_a := range t_id_list {
			t_url, t_err = genGraphApiUrl(map[string]any{}, getMsgFoldersSubPath(folder_id, msg_id), "/attachments/", t_a.String())
			if t_err != nil {
				continue
			}
			// time.Sleep(APIInterval)
			_, t_err = performGraphApiGet(access_token, &t_url, proxy)
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
//			-- https://graph.microsoft.com/v1.0/me/messages/{id}; /me/messages/{id};
//			-- https://graph.microsoft.com/v1.0/me/mailFolders/{id}/messages/{msg_id}; /me/mailFolders/{id}/messages/{id};
//
//		 params:
//			"access_token": api access token;
//			"folder_id": folder's id;
//			"count": specific count of mails to read, when count <=0, means read all mails;
//	        "proxy": specific the proxy when we perform REST API;
//	        "filter_unread": filter unread mails;
//	        "read_latest": read latest mails;
//	        "read_unread": read unread mails;
func readMailsFromFolder(access_token, folder_id *string, count int, proxy *string, filter_unread, read_latest, read_unread bool) (int, int) {
	var content string
	var err error
	var s, f int = 0, 0
	param := map[string]any{}
	// read unread mails
	if filter_unread {
		param["$filter"] = "isRead eq false"
	}
	// order result by receivedDateTime reverse order
	if read_latest {
		param["$orderby"] = "receivedDateTime DESC"
	}
	// set $top
	t_top := count
	if t_top <= 0 {
		t_top = ReadMailsCount
	}
	param["$top"] = t_top
	t_url, t_err := genGraphApiUrl(param, getMsgFoldersSubPath(folder_id, nil))
	if t_err != nil {
		return s, f
	}
	t_fetched := 0
	for {
		content, err = performGraphApiGet(access_token, &t_url, proxy)
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
			if read_unread && !t_is_read {
				// read message after 100 milliseconds
				// time.Sleep(APIInterval)
				t_s, t_f := readOneMail(access_token, folder_id, &t_msg_id, proxy)
				s += t_s
				f += t_f
			}
		}
		t_fetched += len(id_slice)
		t_next_url := gjson.Get(content, ODataNextLink).String()
		// no more results
		if len(t_next_url) == 0 || (count > 0 && t_fetched >= count) {
			break
		}
		t_url = t_next_url
		// time.Sleep(APIInterval)
	}
	return s, f
}

// get all mailFolders and their child folders, and their mails
func readMailsFromAllFolders(access_token, proxy *string, read_unread bool) (int, int) {
	var content string
	var err error
	var s, f int = 0, 0
	t_folder_result, t_s, t_f, err := getAllMailFolders(access_token, proxy)
	s += t_s
	f += t_f
	if err != nil {
		return s, f
	}
	// TODO： change folders into map[*string]map[*string]gjson.Result, to reduce string replicas
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
		// time.Sleep(APIInterval)
		t_id := t_id_list[i]
		// mails delta
		t_s, t_f := getMailFoldersDelta(access_token, &t_id, proxy)
		s += t_s
		f += t_f
		t_s, t_f = readMailsFromFolder(access_token, &t_id, ReadMailsCount, proxy, false, false, read_unread)
		s += t_s
		f += t_f
		// has child folders
		if folders[t_id]["childFolders"].Int() > 0 {
			// get child folders
			// time.Sleep(APIInterval)
			url, _ := genGraphApiUrl(map[string]any{}, "/me/mailFolders", t_id, "/childFolders")
			content, err = performGraphApiGet(access_token, &url, proxy)
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

func deleteOneEmail(access_token, folder_id, msg_id, proxy *string) (bool, error) {
	ok := false
	t_url, err := genGraphApiUrl(map[string]any{}, getMsgFoldersSubPath(folder_id, msg_id))
	if err != nil {
		return ok, fmt.Errorf("fail to generate url, failed with %v", err.Error())
	}
	ok, err = performGraphApiDelete(access_token, &t_url, proxy)
	return ok, err
}

// graph REST API for search: https://graph.microsoft.com/v1.0/search/query
func searchEmailByKeyword(access_token, keyword *string, from, size int, proxy *string) (string, error) {
	var content string
	var err error
	t_url, err := genGraphApiUrl(map[string]any{}, "/search/query")
	if err != nil {
		return "", fmt.Errorf("fail to generate url, failed with %v", err.Error())
	}
	t_data := NewRequestsDataString(keyword, from, size)
	content, err = performGraphApiPost(access_token, &t_url, &t_data, proxy)
	if err != nil {
		return "", fmt.Errorf("fail to search email, failed with %v", err.Error())
	}
	return content, nil
}

// get filtered mails by a keyword in specific folder with it's folder_id
// count: specific count of mails to read, when count <=0, means read all mails
// use $search instead of $filter for more convenient call
// return: the mails in []gjson.Result, the success operation count, the failure operation account, and the error
func getFilteredMails(folder_id, access_token, keyword *string, count int, proxy *string) ([]gjson.Result, int, int, error) {
	var content string
	var err error
	var s, f int = 0, 0
	var t_result_slice []gjson.Result
	// quote the keyword
	var params map[string]any = map[string]any{"$search": fmt.Sprintf("\"%v\"", *keyword)}
	t_top := count
	if t_top <= 0 {
		t_top = ReadMailsCount
	}
	params["$top"] = t_top
	t_url, err := genGraphApiUrl(params, "me/mailFolders", *folder_id, "messages")
	if err != nil {
		return t_result_slice, s, f, err
	}
	t_fetched := 0
	for {
		content, err = performGraphApiGet(access_token, &t_url, proxy)
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
		if len(t_next_url) == 0 || (count > 0 && t_fetched >= count) {
			break
		}
		t_url = t_next_url
		// time.Sleep(APIInterval)
	}
	return t_result_slice, s, f, err
}

// get all mail folders
func getAllMailFolders(access_token, proxy *string) ([]gjson.Result, int, int, error) {
	var t_result_slice []gjson.Result
	var s, f int = 0, 0
	var err error
	var content string
	t_url, _ := genGraphApiUrl(map[string]any{}, "/me/mailFolders")
	for {
		content, err = performGraphApiGet(access_token, &t_url, proxy)
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
func getSpecificFolderId(access_token, proxy *string, folder_name_list []*string) ([]string, int, int, error) {
	var folder_id_list []string
	var s, f int = 0, 0
	t_folder_result, t_s, t_f, err := getAllMailFolders(access_token, proxy)
	s += t_s
	f += t_f
	if err != nil {
		return folder_id_list, s, f, err
	}
	// make temp folder name map for quick search
	t_folder_name_map := make(map[string]bool)
	for _, r := range folder_name_list {
		t_folder_name_map[*r] = true
	}
	for _, r := range t_folder_result {
		folder_name := r.Get("displayName").String()
		folder_id := r.Get("id").String()
		if t_folder_name_map[folder_name] {
			folder_id_list = append(folder_id_list, folder_id)
		}
	}
	return folder_id_list, s, f, nil
}

// delete specific mails by keywords under specific folders, each string in keywords will be a condition to query mails to delete
// we delete mails in "Inbox", "Sent Items", "Drafts"
func deleteOutlookMails(access_token *string, keywords []string, quantity_for_delete int, proxy *string) (int, int) {
	var ok bool
	var s, f int = 0, 0
	t_deleted := 0
	// get all folders, get the folder_id of "Inbox", "Sent Items", "Drafts"
	target_folder_list, t_s, t_f, err := getSpecificFolderId(access_token, proxy, []*string{&MailBoxFolderInBox, &MailBoxFolderSent, &MailBoxFolderDrafts})
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
				t_result_list, t_s, t_f, err := getFilteredMails(&folder_id, access_token, &t_keywords, quantity_for_delete, proxy)
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
						// time.Sleep(APIInterval)
						t_s, t_f := readOneMail(access_token, &t_folder_id, &t_msg_id, proxy)
						s += t_s
						f += t_f
					}
					// do deleteOneEmail
					// time.Sleep(APIInterval)
					ok, err = deleteOneEmail(access_token, &t_folder_id, &t_msg_id, proxy)
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
func searchAndLoopMails(access_token *string, keywords []string, fetch_quantity int, proxy *string, read_unread bool) (int, int) {
	var s, f int = 0, 0
	var fetched int
	t_from := 0
	// loop keywords, search mails and delete mails
	for _, t_keywords := range keywords {
		for {
			content, err := searchEmailByKeyword(access_token, &t_keywords, t_from, ReadMailsCount, proxy)
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
				if read_unread && !t_is_read {
					// read this message if it is not read yet
					// time.Sleep(APIInterval)
					t_s, t_f := readOneMail(access_token, &t_folder_id, &t_msg_id, proxy)
					s += t_s
					f += t_f
				}
			}
			fetched += len(t_msg_id_list)
			// no more results available or fetched enough mails, break
			if !t_more_results_available || fetched >= fetch_quantity {
				break
			}
			t_from += len(t_msg_id_list)
		}
	}
	return s, f
}

// get mailFolders delta
// TODO: read messages from delta
func getMailFoldersDelta(access_token, folder_id, proxy *string) (int, int) {
	t_url, _ := genGraphApiUrl(map[string]any{}, "/me/mailFolders", *folder_id, "messages", "delta")
	_, err := performGraphApiGet(access_token, &t_url, proxy)
	if err != nil {
		return 0, 1
	}
	return 1, 0
}

// list unread mails
func listUnreadMails(access_token, proxy *string, count int, read_unread bool) (int, int) {
	t_s, t_f := readMailsFromFolder(access_token, nil, ReadMailsCount, proxy, true, true, read_unread)
	return t_s, t_f
}

// send email
func sendEmail(access_token, proxy, subject, email_content, to *string) (int, int) {
	t_url, _ := genGraphApiUrl(map[string]any{}, "/me/sendMail")
	t_json_content := NewEmailContentString(subject, &MailContentHtml, email_content, to, false)
	ok, _ := performGraphApiPostSendMail(access_token, &t_url, &t_json_content, proxy)
	if !ok {
		return 0, 1
	}
	return 1, 0
}

func WorkingOnMails(id uint, access_token *string, out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, to *string) {
	var s, f int = 0, 0
	t_start_at := time.Now()
	t_s, t_f := listUnreadMails(access_token, proxy, ms_config.Mail.ReadMails.Quantity, ms_config.Mail.ReadMails.Enabled)
	s += t_s
	f += t_f
	// time.Sleep(APIInterval)
	t_s, t_f = readMailsFromAllFolders(access_token, proxy, ms_config.Mail.ReadMails.Enabled)
	s += t_s
	f += t_f
	// time.Sleep(APIInterval)
	t_s, t_f = searchAndLoopMails(access_token, ms_config.Mail.SearchMails.Keywords, ms_config.Mail.SearchMails.Quantity, proxy, ms_config.Mail.SearchMails.ReadUnread)
	s += t_s
	f += t_f
	// do deleteOutlookMails just according the config file
	if ms_config.Mail.AutoDeleteMails.Enabled {
		// time.Sleep(APIInterval)
		t_s, t_f = deleteOutlookMails(access_token, ms_config.Mail.AutoDeleteMails.Keywords, ms_config.Mail.AutoDeleteMails.Quantity, proxy)
		s += t_s
		f += t_f
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMail,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
	}
}

// do mail read, include listUnreadMails() and readMailsFromAllFolders()
func WorkingOnMailsRead(id uint, access_token *string, out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, to *string) {
	var s, f int = 0, 0
	t_start_at := time.Now()
	t_s, t_f := listUnreadMails(access_token, proxy, ms_config.Mail.ReadMails.Quantity, ms_config.Mail.ReadMails.Enabled)
	s += t_s
	f += t_f
	// time.Sleep(APIInterval)
	t_s, t_f = readMailsFromAllFolders(access_token, proxy, ms_config.Mail.ReadMails.Enabled)
	s += t_s
	f += t_f
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailRead,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
	}
}

// do mail search
func WorkingOnMailsSearch(id uint, access_token *string, out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, to *string) {
	var s, f int = 0, 0
	t_start_at := time.Now()
	t_s, t_f := searchAndLoopMails(access_token, ms_config.Mail.SearchMails.Keywords, ms_config.Mail.SearchMails.Quantity, proxy, ms_config.Mail.SearchMails.ReadUnread)
	s += t_s
	f += t_f
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailSearch,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
	}
}

// do mail deletion
func WorkingOnMailsDelete(id uint, access_token *string, out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, to *string) {
	var s, f int = 0, 0
	t_start_at := time.Now()
	t_s, t_f := deleteOutlookMails(access_token, ms_config.Mail.AutoDeleteMails.Keywords, ms_config.Mail.AutoDeleteMails.Quantity, proxy)
	s += t_s
	f += t_f
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailDelete,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
	}
}

// do mail send
func WorkingOnMailsSend(id uint, access_token *string, out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, to *string) {
	var s, f int = 0, 0
	t_start_at := time.Now()
	t_s, t_f := sendEmail(access_token, proxy, &ms_config.Mail.AutoSendMails.Subject, &ms_config.Mail.AutoSendMails.TemplateContent, to)
	s += t_s
	f += t_f
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailDelete,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
	}
}

// func DoListAllMails(id uint, access_token string, out chan ApiResult, proxy string) {
// 	t_start_at := time.Now()
// 	t_s, t_f := readMailsFromFolder(access_token, "", ReadMailsCount, proxy, false, false)
// 	t_end_at := time.Now()
// 	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
// 	out <- ApiResult{
// 		ID:        id,
// 		OpID:      OpTypeMail,
// 		S:         t_s,
// 		F:         t_f,
// 		StartTime: t_start_at.Unix(),
// 		Duration:  t_durations_milliseconds,
// 		EndTime:   t_end_at.Unix(),
// 	}
// }

// func DoListAllMailFolders(id uint, access_token string, out chan ApiResult, proxy string) {
// 	t_start_at := time.Now()
// 	t_s, t_f := readMailsFromAllFolders(access_token, proxy)
// 	t_end_at := time.Now()
// 	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
// 	out <- ApiResult{
// 		ID:        id,
// 		OpID:      OpTypeMail,
// 		S:         t_s,
// 		F:         t_f,
// 		StartTime: t_start_at.Unix(),
// 		Duration:  t_durations_milliseconds,
// 		EndTime:   t_end_at.Unix(),
// 	}
// }

// func DoMailDeletion(id uint, access_token string, out chan ApiResult, proxy string) {
// 	t_start_at := time.Now()
// 	t_s, t_f := deleteOutlookMails(access_token, config.MailAutoDeleteKeyWords, config.MailAutoDeleteQuantity, proxy)
// 	t_end_at := time.Now()
// 	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
// 	out <- ApiResult{
// 		ID:        id,
// 		OpID:      OpTypeMail,
// 		S:         t_s,
// 		F:         t_f,
// 		StartTime: t_start_at.Unix(),
// 		Duration:  t_durations_milliseconds,
// 		EndTime:   t_end_at.Unix(),
// 	}
// }
