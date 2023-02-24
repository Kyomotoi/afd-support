package afdian

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	queryOrderURL   = "https://afdian.net/api/open/query-order"
	querySponsorURL = "https://afdian.net/api/open/query-sponsor"
	getProfileURL   = "https://afdian.net/api/user/get-profile-by-slug?url_slug="
)

type AfdianAPIService struct {
	APIToken string
	UserID   string
}

func New(token string, userID string) *AfdianAPIService {
	return &AfdianAPIService{
		APIToken: token,
		UserID:   userID,
	}
}

func (a *AfdianAPIService) sign(params string, ts string, user_id string) string {
	content := fmt.Sprintf("%sparams%sts%suser_id%s", a.APIToken, params, ts, user_id)
	has := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", has)
}

func (a *AfdianAPIService) request(url string, params string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}

	now := strconv.FormatInt(time.Now().Unix(), 10)
	q := req.URL.Query()
	q.Add("user_id", a.UserID)
	q.Add("params", params)
	q.Add("ts", now)
	q.Add("sign", a.sign(params, now, a.UserID))
	req.URL.RawQuery = q.Encode()
	logrus.Debug(fmt.Sprintf("API 请求参数: %s", req.URL.RawQuery))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}

func (a *AfdianAPIService) QueryOrder(content string) (*AfdianQueryResponse, error) {
	var params string
	if len(content) < 20 {
		params = fmt.Sprintf(`{"page": %s}`, content)
	} else {
		params = fmt.Sprintf(`{"out_trade_no": %s}`, content)
	}

	resp, err := a.request(queryOrderURL, params)
	if err != nil {
		logrus.Error(fmt.Sprintf("API: 请求失败 %s", queryOrderURL))
		return &AfdianQueryResponse{}, err
	}
	defer resp.Body.Close()

	aqr := &AfdianQueryResponse{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, aqr)
	if err != nil {
		logrus.Error("API: 解析返回数据失败")
		return &AfdianQueryResponse{}, err
	}
	return aqr, nil
}

func (a *AfdianAPIService) QuerySponsorWithPage(page string) (*AfdianQueryResponse, error) {
	params := fmt.Sprintf(`{"page": %s}`, page)
	resp, err := a.request(querySponsorURL, params)
	if err != nil {
		logrus.Error("API: 请求失败 %s", querySponsorURL)
	}
	defer resp.Body.Close()

	aqr := &AfdianQueryResponse{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, aqr)
	if err != nil {
		logrus.Error("API: 解析返回数据失败")
		return &AfdianQueryResponse{}, err
	}
	return aqr, nil
}

func GetUserIDbyProfileURL(profileURL string) (*AfdianProfileResponse, error) {
	reg := regexp.MustCompile("afdian.net/a/(.*)")
	if matched := reg.FindStringSubmatch(profileURL); matched != nil {
		url := getProfileURL + matched[1]
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return &AfdianProfileResponse{}, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return &AfdianProfileResponse{}, err
		}
		defer resp.Body.Close()

		apr := &AfdianProfileResponse{}
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, apr)
		if err != nil {
			logrus.Error("API: 解析返回数据失败")
			return &AfdianProfileResponse{}, err
		}
		return apr, nil
	} else {
		return &AfdianProfileResponse{}, errors.New("afdian 昵称无匹配")
	}
}
