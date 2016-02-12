package pastebin

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const endpoint = "http://pastebin.com/api/api_post.php"

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
	u.Set("api_paste_code", body)
	u.Set("api_paste_name", title)
	u.Set("api_paste_private", private)
	u.Set("api_paste_expire_date", expire)

	res, err := http.DefaultClient.PostForm(endpoint, u)
	if err != nil {
		log.Printf("error pasting: %s", err.Error())
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
