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
// return:
//
//		[4]interface{folder_id, msg_id, is_read, has_attachments}
//	 successful or not
//	 error
func readOneMail(access_token, folder_id, msg_id, proxy *string) ([4]interface{}, bool, error) {
	var s, f int = 0, 0
	var t_result [4]interface{}
	var t_url string
	var t_err error
	t_url, t_err = genGraphApiUrl(query_string, getMsgFoldersSubPath(folder_id, msg_id))
	if t_err != nil {
		return t_result, false, t_err
	}
	var t_status_code int
	var content string
	t_status_code, content, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		f += 1
		return t_result, false, t_err
	}
	// treat it as successful just return 200
	s++
	is_read := gjson.Get(content, "isRead").Bool()
	has_attachments := gjson.Get(content, "hasAttachments").Bool()
	t_result[0] = folder_id
	t_result[1] = msg_id
	t_result[2] = is_read
	t_result[3] = has_attachments
	return t_result, true, t_err
}

// list a mail's attachments
func listMailsAttachments(access_token, folder_id, msg_id, proxy *string) ([]string, bool, error) {
	var s, f int = 0, 0
	var t_url, content string
	var t_status_code int
	var t_err error
	var t_result []string
	// get attachments list
	t_url, t_err = genGraphApiUrl(query_string_full, getMsgFoldersSubPath(folder_id, msg_id), "/attachments")
	if t_err != nil {
		return t_result, false, t_err
	}
	// time.Sleep(APIInterval)
	t_status_code, content, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		f += 1
		return t_result, false, t_err
	}
	// treat it as successful just return 200
	s++
	t_results := gjson.GetMany(content, "value.#.id")
	t_id_list := t_results[0].Array()
	// get attachments content
	for _, t_a := range t_id_list {
		t_result = append(t_result, t_a.String())
	}
	return t_result, true, t_err
}

// download a mail's attachment
func downloadMailsAttachment(access_token, folder_id, msg_id, attachment_id, proxy *string) (bool, error) {
	var t_url, content string
	var t_status_code int
	var t_err error
	// get attachments list
	t_url, t_err = genGraphApiUrl(query_string_full, getMsgFoldersSubPath(folder_id, msg_id), "/attachments", *attachment_id)
	if t_err != nil {
		return false, t_err
	}
	// time.Sleep(APIInterval)
	t_status_code, content, t_err = performGraphApiGet(access_token, &t_url, proxy)
	if t_err != nil || t_status_code != 200 {
		return false, t_err
	}
	// treat it as successful just return 200
	t_attachment_id := gjson.Get(content, "id").String()
	ok := true
	if len(t_attachment_id) == 0 || t_attachment_id != *attachment_id {
		t_err = fmt.Errorf("failed to download attachment")
		ok = false
	}
	return ok, t_err
}

// mark a mail as read
func readMailMarkAsRead(access_token, folder_id, msg_id, proxy *string) (bool, error) {
	var t_url string
	var t_err error
	// get attachments list
	t_url, t_err = genGraphApiUrl(map[string]any{}, getMsgFoldersSubPath(folder_id, msg_id))
	if t_err != nil {
		return false, t_err
	}
	t_data := `{
    "isRead": true
}`
	// time.Sleep(APIInterval)
	return performGraphApiPatch(access_token, &t_url, &t_data, proxy)
}

// get mails under folder
//
//	 support multiple keywords
//		return list - [(t_folder_id, t_msg_id, is_read)], count of successful call, count of failure call, f, error
func getMailsFromFolder(access_token, folder_id *string, count int, proxy *string, read_latest, get_unread bool, keywords []*string) ([][3]interface{}, int, int, error) {
	var content, t_url string
	var t_status_code int
	var err error
	var s, f int
	var t_results [][3]interface{}
	param := map[string]any{}
	// read unread mails, make sure we just do one call to finish this job
	var filter_rules string
	if get_unread {
		filter_rules = "isRead eq false"
		param["$filter"] = "isRead eq false"
	}
	if len(keywords) > 0 {
		var t_filter_rules string
		// filter via multiple keywords
		for _, t_k := range keywords {
			if len(t_filter_rules) > 0 {
				t_filter_rules += " or "
			}
			t_filter_rules += fmt.Sprintf("contains(body/content,'%v')", *t_k)
		}
		if len(filter_rules) > 0 {
			filter_rules = fmt.Sprintf("%v and (%v)", filter_rules, t_filter_rules)
		} else {
			filter_rules += t_filter_rules
		}
	}
	// order result by receivedDateTime reverse order
	if read_latest {
		param["$orderby"] = "receivedDateTime DESC"
		t_filter_rules := "receivedDateTime gt 1900-01-01T00:00:00Z"
		// if read_latest, we need order by receivedDateTime as descend order, and make sure receivedDateTime in $filter's first parameter
		if len(filter_rules) > 0 {
			filter_rules = fmt.Sprintf("%v and (%v)", t_filter_rules, filter_rules)
		} else {
			filter_rules = t_filter_rules
		}
	}
	if len(filter_rules) > 0 {
		param["$filter"] = filter_rules
	}
	// set $top to count, make sure we just do one call to finish this job
	t_top := count
	if t_top <= 0 {
		t_top = ReadMailsCount
	}
	param["$top"] = t_top
	t_url, err = genGraphApiUrl(param, getMsgFoldersSubPath(folder_id, nil))
	if err != nil {
		return t_results, s, f, err
	}
	t_fetched := 0
	for {
		t_status_code, content, err = performGraphApiGet(access_token, &t_url, proxy)
		if err != nil || t_status_code != 200 {
			f++
			return t_results, s, f, err
		}
		// invalid response
		if gjson.Get(content, ODataContext).String() == "" {
			f++
			break
			// return t_results, s, f, err
		}
		s++
		t_message_list := gjson.Get(content, "value").Array()
		for _, t_msg := range t_message_list {
			t_raw := t_msg.Raw
			t_msg_id := gjson.Get(t_raw, "id").String()
			t_folder_id := gjson.Get(t_raw, "parentFolderId").String()
			t_is_read := gjson.Get(t_raw, "isRead").Bool()
			// if we need to get unread mails and this mail does not read yet or we need read and unread mails.
			if !get_unread || (get_unread && !t_is_read) {
				// append to t_results
				t_results = append(t_results, [3]interface{}{t_folder_id, t_msg_id, t_is_read})
				// count the number we fetched
				t_fetched += 1
			}
		}
		t_next_url := gjson.Get(content, ODataNextLink).String()
		// no more results
		if len(t_next_url) == 0 || (count > 0 && t_fetched >= count) {
			break
		}
		t_url = t_next_url
		// time.Sleep(APIInterval)
	}
	return t_results, s, f, err
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
func searchEmailByKeywords(access_token *string, keywords []*string, from, size int, proxy *string) (string, error) {
	var content string
	var t_status_code int
	var err error
	t_url, err := genGraphApiUrl(map[string]any{}, "/search/query")
	if err != nil {
		return "", fmt.Errorf("fail to generate url, failed with %v", err.Error())
	}
	t_data := NewRequestsDataStringMultiple(keywords, from, size)
	t_status_code, content, err = performGraphApiPost(access_token, &t_url, &t_data, proxy)
	if err != nil {
		return "", fmt.Errorf("fail to search email, failed with %v", err.Error())
	}
	if t_status_code != 200 {
		return "", fmt.Errorf("fail to search email, status code is %v", t_status_code)
	}
	return content, nil
}

// get all mail folders
func getAllMailFolders(access_token, proxy *string) ([]gjson.Result, int, int, error) {
	var t_result_slice []gjson.Result
	var t_status_code int
	var s, f int = 0, 0
	var err error
	var content string
	t_url, _ := genGraphApiUrl(map[string]any{}, "/me/mailFolders")
	for {
		t_status_code, content, err = performGraphApiGet(access_token, &t_url, proxy)
		if err != nil || t_status_code != 200 {
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

// folder info like
//
//	{
//		  "id": "AQMkADYAAAIBDAAAAA==",
//		  "displayName": "Inbox",
//		  "parentFolderId": "AQMkADYAAAIBCAAAAA==",
//		  "childFolderCount": 1,
//		  "unreadItemCount": 70,
//		  "totalItemCount": 71,
//		  "isHidden": false
//	}
//
// get all mail folders from root, or get children folder list from folder_id, if isHidden is false
// return [][4]string as [(parent's id, it's own id, number of child folder, total Item Count in this folder)]
func getAllMailFoldersNew(folder_id, access_token, proxy *string) ([][4]interface{}, int, int, error) {
	var s, f int = 0, 0
	var err error
	var content string
	var t_url string
	var t_status_code int
	if folder_id == nil || *folder_id == "" {
		t_url, _ = genGraphApiUrl(map[string]any{}, "/me/mailFolders")
	} else {
		t_url, _ = genGraphApiUrl(map[string]any{}, "/me/mailFolders", *folder_id, "childFolders")
	}

	var t_result_slice [][4]interface{}
	for {
		t_status_code, content, err = performGraphApiGet(access_token, &t_url, proxy)
		if err != nil || t_status_code != 200 {
			f += 1
			break
		}
		s += 1
		t_values := gjson.Get(content, "value").Array()
		for _, r := range t_values {
			is_hidden := r.Get("isHidden").Bool()
			// just as part of return if it does not hidden
			if !is_hidden {
				t_result_slice = append(t_result_slice, [4]interface{}{
					r.Get("parentFolderId").String(),
					r.Get("id").String(),
					r.Get("childFolderCount").Int(),
					r.Get("totalItemCount").Int(),
				})
			}
		}
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

// use https://graph.microsoft.com/v1.0/search/query
//
//	return list - [(t_folder_id, t_msg_id, t_is_read)], count of successful call, count of failure call, f, error
func searchMailsByKeywords(access_token *string, keywords []*string, fetch_quantity int, proxy *string, read_unread bool) ([][3]interface{}, int, int, error) {
	var s, f int = 0, 0
	var fetched int
	t_from := 0
	var t_list [][3]interface{}
	for {
		content, err := searchEmailByKeywords(access_token, keywords, t_from, fetch_quantity, proxy)
		// bad request, or network issue or other error, rather than empty search result
		if err != nil {
			f += 1
			break
		}
		s++
		// get search results in path "content.value[0].hitsContainers[0]" use gjson
		t_search_content := gjson.Get(content, "value.0.hitsContainers.0").String()
		// no matched results, return as one successful try
		if t_search_content == "" {
			break
		}
		t_hits_string := gjson.Get(t_search_content, "hits").String()
		// t_total := gjson.Get(t_search_content, "total").Int()
		t_more_results_available := gjson.Get(t_search_content, "moreResultsAvailable").Bool()
		// empty search results, return as one successful try
		if len(t_hits_string) == 0 {
			break
		}
		t_msg_id_list := gjson.Get(t_hits_string, "#.hitId").Array()
		t_is_read_list := gjson.Get(t_hits_string, "#.resource.isRead").Array()
		t_folder_id_list := gjson.Get(t_hits_string, "#.resource.parentFolderId").Array()
		if len(t_msg_id_list) == 0 || len(t_folder_id_list) == 0 || len(t_is_read_list) == 0 {
			break
		}
		for i := 0; i < len(t_msg_id_list); i++ {
			t_msg_id := t_msg_id_list[i].String()
			t_is_read := t_is_read_list[i].Bool()
			t_folder_id := t_folder_id_list[i].String()
			if read_unread && !t_is_read {
				t_list = append(t_list, [3]interface{}{t_folder_id, t_msg_id, t_is_read})
			}
		}
		fetched += len(t_msg_id_list)
		// no more results available or fetched enough mails, break
		if !t_more_results_available || fetched >= fetch_quantity {
			break
		}
		t_from += len(t_msg_id_list)
	}
	return t_list, s, f, nil
}

// get mails delta
// TODO: read messages from delta
func getMailsDelta(access_token, folder_id, proxy *string) (int, int) {
	t_url, _ := genGraphApiUrl(map[string]any{}, "/me/mailFolders", *folder_id, "messages", "delta")
	t_status_code, _, err := performGraphApiGet(access_token, &t_url, proxy)
	if err != nil || t_status_code != 200 {
		return 0, 1
	}
	return 1, 0
}

// get mails folder delta
// TODO: read messages from delta
func getMailFoldersDelta(access_token, proxy *string) bool {
	t_url, _ := genGraphApiUrl(map[string]any{}, "/me/mailFolders", "delta")
	t_status_code, _, err := performGraphApiGet(access_token, &t_url, proxy)
	if err != nil || t_status_code != 200 {
		return false
	}
	return true
}

// send email
func sendEmail(access_token, proxy, subject, email_content, to *string) (bool, error) {
	t_url, _ := genGraphApiUrl(map[string]any{}, "/me/sendMail")
	t_json_content := NewEmailContentString(subject, &MailContentHtml, email_content, to, false)
	return performGraphApiPostSendMail(access_token, &t_url, &t_json_content, proxy)
}

// do list unread mails from root, or from a folder
func MailListMailsFrom(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var t_list [][3]interface{}
	var id uint
	var access_token *string
	var read_attachments bool
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	var folder_id *string
	if args[ArgFolderId] != nil {
		folder_id = args[ArgFolderId].(*string)
	}
	if _, ok := args[ArgReadAttachments]; ok {
		read_attachments = args[ArgReadAttachments].(bool)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		t_list, s, f, _ = getMailsFromFolder(access_token, folder_id, ms_config.Mail.ReadMails.Quantity, proxy, true, ms_config.Mail.ReadMails.ReadUnread, nil)
	}

	t_end_at := time.Now()
	t_task_list := []*Task{}
	// append task to MailReadMail only if ReadUnread is true
	if ms_config.Mail.ReadMails.ReadUnread {
		for _, t := range t_list {
			t_folder_id := t[0].(string)
			t_mail_id := t[1].(string)
			is_read := t[2].(bool)
			if !is_read {
				t_task := &Task{
					Func: MailReadMail,
					Args: MsArgs{
						ArgUserID:          id,
						ArgAccessToken:     access_token,
						ArgMailId:          &t_mail_id,
						ArgFolderId:        &t_folder_id,
						ArgReadAttachments: read_attachments,
					},
				}
				t_task_list = append(t_task_list, t_task)
			}
		}
	}
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailListRootUnreadMail,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     t_task_list,
	}
}

// read one mail
func MailReadMail(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var t_result [4]interface{}
	var ok, read_attachments bool
	var s, f int = 0, 0
	var folder_id *string
	var mail_id *string
	var id uint
	var access_token *string
	var t_task_list []*Task
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		folder_id = args[ArgFolderId].(*string)
	}
	if args[ArgMailId] != nil {
		mail_id = args[ArgMailId].(*string)
	}
	if _, ok := args[ArgReadAttachments]; ok {
		read_attachments = args[ArgReadAttachments].(bool)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		t_result, ok, _ = readOneMail(access_token, folder_id, mail_id, proxy)
		if ok {
			s += 1
			if len(t_result) > 0 {
				is_read := t_result[2].(bool)
				has_attachments := t_result[3].(bool)
				if !is_read {
					t_task_list = append(t_task_list, &Task{
						Func: MailReadMarkMailAsRead,
						Args: MsArgs{
							ArgUserID:      id,
							ArgAccessToken: access_token,
							ArgMailId:      mail_id,
							ArgFolderId:    folder_id,
						},
					})
				}
				if has_attachments && read_attachments {
					t_task_list = append(t_task_list, &Task{
						Func: MailReadListMailsAttachments,
						Args: MsArgs{
							ArgUserID:          id,
							ArgAccessToken:     access_token,
							ArgMailId:          mail_id,
							ArgFolderId:        folder_id,
							ArgReadAttachments: &read_attachments,
						},
					})
				}
			}
		} else {
			f += 1
		}
	}
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
		Tasks:     t_task_list,
	}
}

func MailReadListMailsAttachments(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var t_result []string
	var ok bool
	var s, f int = 0, 0
	var folder_id *string
	var mail_id *string
	var id uint
	var access_token *string
	var t_task_list []*Task
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		folder_id = args[ArgFolderId].(*string)
	}
	if args[ArgMailId] != nil {
		mail_id = args[ArgMailId].(*string)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		t_result, ok, _ = listMailsAttachments(access_token, folder_id, mail_id, proxy)
		if ok {
			s += 1
			for _, attachment_id := range t_result {
				t_task_list = append(t_task_list, &Task{
					Func: MailReadDownloadAttachment,
					Args: MsArgs{
						ArgUserID:       id,
						ArgAccessToken:  access_token,
						ArgMailId:       mail_id,
						ArgFolderId:     folder_id,
						ArgAttachmentId: &attachment_id,
					},
				})
			}
		} else {
			f += 1
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailReadMailsAttachments,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     t_task_list,
	}
}

func MailReadDownloadAttachment(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var ok bool
	var s, f int = 0, 0
	var folder_id, mail_id, attachment_id *string
	var id uint
	var access_token *string
	var t_task_list []*Task
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		folder_id = args[ArgFolderId].(*string)
	}
	if args[ArgMailId] != nil {
		mail_id = args[ArgMailId].(*string)
	}
	if args[ArgAttachmentId] != nil {
		attachment_id = args[ArgAttachmentId].(*string)
	}
	if id > 0 && access_token != nil && attachment_id != nil && len(*access_token) > 0 && len(*attachment_id) > 0 {
		ok, _ = downloadMailsAttachment(access_token, folder_id, mail_id, attachment_id, proxy)
		if ok {
			s += 1
		} else {
			f += 1
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailReadMailsAttachments,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     t_task_list,
	}
}

func MailReadMarkMailAsRead(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var ok bool
	var s, f int = 0, 0
	var folder_id, mail_id *string
	var id uint
	var access_token *string
	var t_task_list []*Task
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		folder_id = args[ArgFolderId].(*string)
	}
	if args[ArgMailId] != nil {
		mail_id = args[ArgMailId].(*string)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok, _ = readMailMarkAsRead(access_token, folder_id, mail_id, proxy)
		if ok {
			s += 1
		} else {
			f += 1
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailReadMarkMailAsRead,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     t_task_list,
	}
}

// list all mail folders from root or a folder, and make some tasks in chan out as []*Task
// these tasks include read mail's delta, read child folder and read mails from folder
func MailListMailFolders(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var t_list [][4]interface{}
	var folder_id *string
	var id uint
	var access_token *string
	var read_attachments bool
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		folder_id = args[ArgFolderId].(*string)
	}
	if _, ok := args[ArgReadAttachments]; ok {
		read_attachments = args[ArgReadAttachments].(bool)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		t_list, s, f, _ = getAllMailFoldersNew(folder_id, access_token, proxy)
	}
	t_end_at := time.Now()
	t_task_list := []*Task{}
	for _, t := range t_list {
		t_own_id := t[1].(string)
		// read mail's delta
		t_task_list = append(t_task_list, &Task{
			Func: MailReadMailsDelta,
			Args: MsArgs{
				ArgUserID:          id,
				ArgAccessToken:     access_token,
				ArgFolderId:        &t_own_id,
				ArgReadAttachments: read_attachments,
			},
		})
		// read child folder, if it has child folder
		t_child_folder_count := t[2].(int64)
		if t_child_folder_count > 0 {
			t_task_list = append(t_task_list, &Task{
				Func: MailListMailFolders,
				Args: MsArgs{
					ArgUserID:          id,
					ArgAccessToken:     access_token,
					ArgFolderId:        &t_own_id,
					ArgReadAttachments: read_attachments,
				},
			})
		}
		// read mails from folder
		t_total_item_count := t[3].(int64)
		if t_total_item_count > 0 {
			t_task_list = append(t_task_list, &Task{
				Func: MailListMailsFromFolder,
				Args: MsArgs{
					ArgUserID:          id,
					ArgAccessToken:     access_token,
					ArgFolderId:        &t_own_id,
					ArgReadAttachments: read_attachments,
				},
			})
		}
	}
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailListMailFolder,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     t_task_list,
	}
}

// read folder's delta
func MailReadMailFoldersDelta(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok := getMailFoldersDelta(access_token, proxy)
		if ok {
			s++
		} else {
			f++
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailReadMailFoldersDelta,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// read mail's delta
func MailReadMailsDelta(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var folder_id *string
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		folder_id = args[ArgFolderId].(*string)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		s, f = getMailsDelta(access_token, folder_id, proxy)
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailReadMailsDelta,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     []*Task{},
	}
}

// do list mails from folder,
func MailListMailsFromFolder(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	var s, f int = 0, 0
	t_start_at := time.Now()
	var folder_id *string
	var t_list [][3]interface{}
	var id uint
	var access_token *string
	var read_attachments bool
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		folder_id = args[ArgFolderId].(*string)
	}
	if _, ok := args[ArgReadAttachments]; ok {
		read_attachments = args[ArgReadAttachments].(bool)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		t_list, s, f, _ = getMailsFromFolder(access_token,
			folder_id,
			ms_config.Mail.ReadMailFolders.Quantity,
			proxy,
			false,
			ms_config.Mail.ReadMailFolders.ReadUnread, nil)
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	t_task_list := []*Task{}
	// read un-read mails, if ReadUnread is true
	if ms_config.Mail.ReadMailFolders.ReadUnread {
		for _, t := range t_list {
			t_folder_id := t[0].(string)
			t_mail_id := t[1].(string)
			is_read := t[2].(bool)
			if !is_read {
				t_task := &Task{
					Func: MailReadMail,
					Args: MsArgs{
						ArgUserID:          id,
						ArgAccessToken:     access_token,
						ArgMailId:          &t_mail_id,
						ArgFolderId:        &t_folder_id,
						ArgReadAttachments: read_attachments,
					},
				}
				t_task_list = append(t_task_list, t_task)
			}
		}
	}
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailListRootUnreadMail,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     t_task_list,
	}
}

// do mail search
func MailsSearch(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var t_list [][3]interface{}
	var id uint
	var access_token *string
	var read_attachments bool
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if _, ok := args[ArgReadAttachments]; ok {
		read_attachments = args[ArgReadAttachments].(bool)
	}
	// do search only if keyword is not empty
	if len(ms_config.Mail.SearchMails.Keywords) > 0 && id > 0 && access_token != nil && len(*access_token) > 0 {
		var t_keyword_slice []*string
		for i := range ms_config.Mail.SearchMails.Keywords {
			t_keyword_slice = append(t_keyword_slice, &ms_config.Mail.SearchMails.Keywords[i])
		}
		t_list, s, f, _ = searchMailsByKeywords(access_token, t_keyword_slice, ms_config.Mail.SearchMails.Quantity, proxy, ms_config.Mail.SearchMails.ReadUnread)
		t_end_at := time.Now()
		t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
		t_task_list := []*Task{}
		// read un-read mails, if ReadUnread is true
		if ms_config.Mail.SearchMails.ReadUnread {
			for _, t := range t_list {
				t_folder_id := t[0].(string)
				t_mail_id := t[1].(string)
				is_read := t[2].(bool)
				if !is_read {
					t_task := &Task{
						Func: MailReadMail,
						Args: MsArgs{
							ArgUserID:          id,
							ArgAccessToken:     access_token,
							ArgMailId:          &t_mail_id,
							ArgFolderId:        &t_folder_id,
							ArgReadAttachments: read_attachments,
						},
					}
					t_task_list = append(t_task_list, t_task)
				}
			}
		}
		out <- &ApiResult{
			ID:        id,
			OpID:      OpTypeMailSearch,
			S:         s,
			F:         f,
			StartTime: &t_start_at,
			Duration:  t_durations_milliseconds,
			EndTime:   &t_end_at,
			Tasks:     t_task_list,
		}
	}
}

// This is the beginning of mail deletion
func MailsDelListSpecificMailFolders(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	t_task_list := []*Task{}
	var target_folder_list []string
	var id uint
	var access_token *string
	var read_attachments bool
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if _, ok := args[ArgReadAttachments]; ok {
		read_attachments = args[ArgReadAttachments].(bool)
	}
	// get all folders, get the folder_id of "Inbox", "Sent Items", "Drafts"
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		var folder_list []*string
		for _, folder_name := range ms_config.Mail.AutoDeleteMails.FolderName {
			folder_list = append(folder_list, &folder_name)
		}
		if len(folder_list) == 0 {
			folder_list = append(folder_list, &MailBoxFolderInBox)
		}
		target_folder_list, s, f, _ = getSpecificFolderId(access_token, proxy, folder_list)
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	for _, t := range target_folder_list {
		t_len := len(ms_config.Mail.AutoDeleteMails.Keywords)
		if t_len > 0 {
			t_task_list = append(t_task_list, &Task{
				Func: MailsDelListFilteredMails,
				Args: MsArgs{
					ArgUserID:          id,
					ArgAccessToken:     access_token,
					ArgFolderId:        &t,
					ArgReadAttachments: read_attachments,
				},
			})
		}
		// for _, t_keyword := range ms_config.Mail.AutoDeleteMails.Keywords {
		// 	// add task to ListFilteredMails, Loop
		// 	t_task_list = append(t_task_list, &Task{
		// 		Func: MailsDelListFilteredMails,
		// 		Args: MsArgs{
		// 			ArgUserID:          id,
		// 			ArgAccessToken:     access_token,
		// 			ArgFolderId:        &t,
		// 			ArgKeyword:         &t_keyword,
		// 			ArgReadAttachments: read_attachments,
		// 		},
		// 	})
		// }
		// add task to ListFilteredMails, Loop
	}
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailListMailFolder,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     t_task_list,
	}
}

// list filtered mails by a keyword in specific folder with it's folder_id
// count: specific count of mails to read, when count <=0, means read all mails
// return list - [(t_folder_id, t_msg_id, is_read)], count of successful call, count of failure call, f, error
func MailsDelListFilteredMails(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	// escape keyword
	var t_folder_id *string
	var t_list [][3]interface{}
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		t_folder_id = args[ArgFolderId].(*string)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		var t_keyword_slice []*string
		for i := range ms_config.Mail.AutoDeleteMails.Keywords {
			t_keyword_slice = append(t_keyword_slice, &ms_config.Mail.AutoDeleteMails.Keywords[i])
		}
		t_list, s, f, _ = getMailsFromFolder(access_token, t_folder_id, ms_config.Mail.AutoDeleteMails.Quantity, proxy, true, false, t_keyword_slice)
	}
	t_task_list := []*Task{}
	// put task to task list, for each mail, read it, if ms_config.Mail.AutoDeleteMails.ReadUnread and mail is unread.
	if ms_config.Mail.AutoDeleteMails.ReadUnread {
		for _, t := range t_list {
			t_folder_id := t[0].(string)
			t_mail_id := t[1].(string)
			is_read := t[2].(bool)
			if !is_read {
				t_task := &Task{
					Func: MailReadMarkMailAsRead,
					Args: MsArgs{
						ArgUserID:      id,
						ArgAccessToken: access_token,
						ArgMailId:      &t_mail_id,
						ArgFolderId:    &t_folder_id,
					},
				}
				t_task_list = append(t_task_list, t_task)
			}
		}
	}
	// put task to task list, for each mail, delete it.
	for _, t := range t_list {
		t_folder_id := t[0].(string)
		t_mail_id := t[1].(string)
		t_task := &Task{
			Func: MailsDelDeleteOneMail,
			Args: MsArgs{
				ArgUserID:      id,
				ArgAccessToken: access_token,
				ArgMailId:      &t_mail_id,
				ArgFolderId:    &t_folder_id,
			},
		}
		t_task_list = append(t_task_list, t_task)
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	//
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailReadFilteredMails,
		S:         s,
		F:         f,
		StartTime: &t_end_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     t_task_list,
	}
}

// just do mail deletion
func MailsDelDeleteOneMail(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	t_start_at := time.Now()
	var s, f int = 0, 0
	var t_folder_id *string
	var t_msg_id *string
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgFolderId] != nil {
		t_folder_id = args[ArgFolderId].(*string)
	}
	if args[ArgMailId] != nil {
		t_msg_id = args[ArgMailId].(*string)
	}
	// delete one mail, return true if success, false if failure.
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok, _ := deleteOneEmail(access_token, t_folder_id, t_msg_id, proxy)
		if ok {
			s += 1
		} else {
			f += 1
		}
	}
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
		Tasks:     nil, // no task in this case.,
	}
}

// do mail send
func MailsSend(out chan *ApiResult, proxy *string, ms_config *config.ConfigMs, args MsArgs) {
	var s, f int = 0, 0
	t_start_at := time.Now()
	var to *string
	var id uint
	var access_token *string
	if args[ArgUserID] != nil {
		id = args[ArgUserID].(uint)
	}
	if args[ArgAccessToken] != nil {
		access_token = args[ArgAccessToken].(*string)
	}
	if args[ArgTo] != nil {
		to = args[ArgTo].(*string)
	}
	if id > 0 && access_token != nil && len(*access_token) > 0 {
		ok, _ := sendEmail(access_token, proxy, &ms_config.Mail.AutoSendMails.Subject, &ms_config.Mail.AutoSendMails.TemplateContent, to)
		if ok {
			s += 1
		} else {
			f += 1
		}
	}
	t_end_at := time.Now()
	t_durations_milliseconds := t_end_at.Sub(t_start_at).Milliseconds()
	out <- &ApiResult{
		ID:        id,
		OpID:      OpTypeMailSend,
		S:         s,
		F:         f,
		StartTime: &t_start_at,
		Duration:  t_durations_milliseconds,
		EndTime:   &t_end_at,
		Tasks:     nil, // no task in this case.,
	}
}
