package pastebin

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const endpoint = "https://pastebin.com/api/api_post.php"

type Client interface {
	Paste(title, body, private, expire string) string
}

type clientImpl struct {
	apiKey string
}

func NewClient(apiKey string) Client {
	return &clientImpl{apiKey: apiKey}
}

func (T *clientImpl) Paste(title, body, private, expire string) string {
	u := url.Values{}
	u.Set("api_option", "paste")
	u.Set("api_dev_key", T.apiKey)
	bodyTrimmed := body
	if len(body) > 100000 {
		bodyTrimmed = body[:100000]
	}
	u.Set("api_paste_code", bodyTrimmed)
	u.Set("api_paste_name", title)
	u.Set("api_paste_private", private)
	u.Set("api_paste_expire_date", expire)
	
	hc := http.Client{}
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(u.Encode()))
	if (err != nil) {
		return ""
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := hc.Do(req)
	if err != nil {
		return ""
	}

	if res.StatusCode == http.StatusOK {
		resp, _ := ioutil.ReadAll(res.Body)

		if strings.Index(string(resp), "Bad API") < 0 {
			return string(resp)
		}
	}

	return ""
}
